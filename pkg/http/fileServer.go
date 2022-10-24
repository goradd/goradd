package http

import (
	"compress/gzip"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/andybalholm/brotli"
	"github.com/goradd/goradd/pkg/config"
	"github.com/goradd/goradd/pkg/log"
	strings2 "github.com/goradd/goradd/pkg/strings"
)

// FileSystemServer serves a file system as an http.Handler.
//
// The file system contained can point to compressed versions of http resources, and those
// compressed version will be served when possible. If you only store compressed versions,
// and if the browser does not support compression, the compressed file will be decompressed here
// before serving the file. This lets you save space by only storing a compressed file, at the cost
// of some speed. Since most browsers support compression, this should not be a big deal.
//
// The files in Fsys must implement the io.ReaderSeeker interface. Both embed and
// traditional OS file systems do this.
//
// Files found will be sent to any registered file processors if the extension matches one of the
// processors (see RegisterFileProcessor).
type FileSystemServer struct {
	// Fsys is the file system being served.
	Fsys fs.FS

	// SendModTime will send the modification time of the file when it is served. Generally, you want
	// to do this for files that can be bookmarked, like html files, since there really is no other way
	// to try to get the server to reload the file when it is changed. However, for asset files that are
	// using the cache buster, you should not do this, since cache busting will take care of notifying
	// the user the file is changed.
	SendModTime bool

	// UseCacheBuster will look for cache buster paths and fix them.
	UseCacheBuster bool

	// Hide is a slice of file endings that will be blocked from being served. These endings do not have to just
	// be file extensions, but any string. So if you specify an ending of "_abc.txt", any file ending in the
	// string will NOT be shown.
	Hide []string
}

// ServeHTTP will serve the file system.
func (f FileSystemServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !f.serveStaticFile(w, r) {
		http.NotFound(w, r)
	}
}

// serveStaticFile serves up files found in the file system of f.
// If the file is not found or cannot be opened, it will return false.
// File errors will be logged.
func (f FileSystemServer) serveStaticFile(w http.ResponseWriter, r *http.Request) bool {
	p := r.URL.Path

	if f.UseCacheBuster {
		p = StripCacheBusterPath(p)
	}

	if len(p) > 0 && p[0] == '/' {
		p = p[1:]
	}

	if p == "" {
		p = "index.html"
	} else if p[len(p)-1] == '/' {
		p = path.Join(p, "index.html")
	} else if s, err := fs.Stat(f.Fsys, p); err == nil && s.IsDir() { // a directory after all
		p = path.Join(p, "/index.html")
	}

	if !fs.ValidPath(p) {
		return false
	}

	for _, bl := range f.Hide {
		if strings2.EndsWith(p, bl) {
			return false // cannot show this kind of file
		}
	}

	var acceptsGzip bool
	var acceptsBr bool

	if values, ok := r.Header["Accept-Encoding"]; ok {
		for _, value1 := range values {
			for _, value := range strings.Split(value1, ",") {
				if value == "gzip" {
					acceptsGzip = true
				} else if value == "br" {
					acceptsBr = true
				}
			}
		}
	}

	// Check for compressed versions
	if foundPath := p + ".br"; acceptsBr && f.pathExists(foundPath) {
		if err := f.servePath(w, r, p, foundPath, "br"); err != nil {
			log.Error(err)
			return false
		}
		return true
	}
	if foundPath := p + ".gz"; acceptsGzip && f.pathExists(foundPath) {
		if err := f.servePath(w, r, p, foundPath, "gzip"); err != nil {
			log.Error(err)
			return false
		}
		return true
	}

	// Check for uncompressed version
	if f.pathExists(p) {
		if err := f.servePath(w, r, p, p, ""); err != nil {
			panic(err)
		}
		return true
	}

	if foundPath := p + ".br"; f.pathExists(foundPath) {
		if err := f.serveDecompressedBrotli(w, r, p, foundPath); err != nil {
			log.Error(err)
			return false
		}
		return true
	}
	if foundPath := p + ".gz"; f.pathExists(foundPath) {
		if err := f.serveDecompressedGzip(w, r, p, foundPath); err != nil {
			log.Error(err)
			return false
		}
		return true
	}
	return false
}

// serveDecompressedBrotli will decompress a found brotli file and serve it up
// as its decompressed counterpart.
func (f FileSystemServer) serveDecompressedBrotli(
	w http.ResponseWriter,
	r *http.Request,
	name string,
	path string) error {
	tempFile, err := os.CreateTemp("", "goradd")
	if err != nil {
		return err
	}
	defer func() {
		// close and delete temp file when done
		p2 := tempFile.Name()
		_ = tempFile.Close()
		_ = os.Remove(p2)
	}()

	var file fs.File
	if file, err = f.Fsys.Open(path); err != nil {
		return err
	}
	defer func(file fs.File) {
		_ = file.Close()
	}(file)

	br := brotli.NewReader(file)

	if _, err = io.Copy(tempFile, br); err != nil {
		return err
	}

	_, _ = tempFile.Seek(0, 0)
	return f.serveFile(w, r, name, tempFile)
}

// serveDecompressedGzip will decompress a found gzip file and serve it up
// as its decompressed counterpart.
func (f FileSystemServer) serveDecompressedGzip(
	w http.ResponseWriter,
	r *http.Request,
	name string,
	path string) error {
	tempFile, err := os.CreateTemp("", "goradd")
	if err != nil {
		return err
	}
	defer func() {
		// close and delete temp file when done
		p2 := tempFile.Name()
		_ = tempFile.Close()
		_ = os.Remove(p2)
	}()

	var file fs.File
	if file, err = f.Fsys.Open(path); err != nil {
		return err
	}
	defer func(file fs.File) {
		_ = file.Close()
	}(file)

	var gz *gzip.Reader
	if gz, err = gzip.NewReader(file); err != nil {
		return err
	}
	if _, err = io.Copy(tempFile, gz); err != nil {
		return err
	}

	_, _ = tempFile.Seek(0, 0)
	return f.serveFile(w, r, name, tempFile)
}

func (f FileSystemServer) pathExists(path string) bool {
	_, err := fs.Stat(f.Fsys, path)
	return err == nil
}

// servePath opens the pathInFS file and serves it.
func (f FileSystemServer) servePath(w http.ResponseWriter,
	r *http.Request,
	name string,
	pathInFS string,
	encoding string) error {
	file, err := f.Fsys.Open(pathInFS)
	if err != nil {
		return err
	}
	defer func(file fs.File) {
		_ = file.Close()
	}(file)

	ext := path.Ext(name)
	if p := fileProcessors[ext]; p != nil {
		return p(file, w, r)
	}
	if c := contentTypes[ext]; c != "" {
		w.Header().Set("Content-Type", c)
	}
	if encoding != "" {
		w.Header().Set("Content-Encoding", encoding)
	}
	return f.serveFile(w, r, name, file)
}

// serveFile serves an open file.
//
// The file must be an io.ReadSeeker
func (f FileSystemServer) serveFile(w http.ResponseWriter,
	r *http.Request,
	name string,
	file fs.File) error {
	var modTime time.Time
	if f.SendModTime {
		if stat, err := file.Stat(); err != nil {
			return err
		} else {
			modTime = stat.ModTime()
		}
	}

	if !config.Release {
		// During development, tell the browser not to cache our static files and assets so that if we change an asset,
		// we don't have to deal with trying to get the browser to refresh
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate, max-age=1")
	}
	http.ServeContent(w, r, name, modTime, file.(io.ReadSeeker))
	return nil
}
