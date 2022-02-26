package http

import (
	"hash/crc64"
	"io/fs"
	"path"
	"strconv"

	"github.com/goradd/goradd/pkg/config"
	strings2 "github.com/goradd/goradd/pkg/strings"
)

var crcTable = crc64.MakeTable(crc64.ECMA)

// cacheBuster maps paths to checksums that are used to tell the browser when it is time to reload a resource.
var cacheBuster = make(map[string]string)

// RegisterAssetDirectory maps a file system to a URL path in the application.
//
// The files in the path are registered with the cache buster, so that when you edit the file
// a new URL will be generated forcing the browser to reload the asset. This is much better than
// using a cache control header.
//
// If the browser attempts to access a file in the file system that does not exist, a 404 NotFound
// error will be sent back to the browser.
func RegisterAssetDirectory(prefix string, fsys fs.FS) {
	prefix = path.Clean(prefix)
	serv := FileSystemServer{Fsys: fsys, SendModTime: false, UseCacheBuster: true}
	RegisterPrefixHandler(prefix, serv)

	// Walk the entire file system provided to register each file with the cache buster cache so
	// that we know what hash to provide for that file.
	if err := fs.WalkDir(fsys, ".", func(p2 string, d fs.DirEntry, err error) error {
		if err != nil {
			return err // stop
		}
		if d.IsDir() {
			return nil // keep going
		}
		var data []byte
		if data, err = fs.ReadFile(fsys, p2); err != nil {
			return err
		} else {
			// CRC it
			c := crc64.Checksum(data, crcTable)
			e := strconv.FormatInt(int64(c), 36)
			s := path.Join(prefix, p2)
			cacheBuster[s] = e
			return nil
		}
	}); err != nil {
		panic("failed walking the asset directory " + prefix + ": " + err.Error())
	}
}

// GetAssetUrl returns the url that corresponds to the asset at the given path.
//
// This will add the cache-buster path, and the proxy path if there is one.
func GetAssetUrl(location string) string {
	location = CacheBustedPath(location)
	return MakeLocalPath(location)
}

// StripCacheBusterPath removes the hash of the asset file from the path to the asset.
func StripCacheBusterPath(fPath string) string {
	dir, f := path.Split(fPath)
	dir2, f2 := path.Split(dir[:len(dir)-1]) // ignore trailing slash
	if strings2.StartsWith(f2, config.CacheBusterPrefix) {
		fPath = path.Join(dir2, f)
	}
	return fPath
}

// CacheBustedPath returns a path to an asset that was previously registered with the CacheBuster. The new path
// will contain a hash of the file that will change whenever the file changes, and cause the browser to reload the file.
// Since we are in control of serving these files, we will later remove the hash before serving it.
func CacheBustedPath(url string) string {
	if p, ok := cacheBuster[url]; ok {
		// inject the crc as a part of the path.
		dir, file := path.Split(url)
		url = path.Join(dir, config.CacheBusterPrefix+p, file)
	}
	return url
}
