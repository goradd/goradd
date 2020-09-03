package page

import (
	"github.com/goradd/goradd/pkg/config"
	"log"
	"mime"
	"net/http"
	"path"
	"path/filepath"
	"strings"

	strings2 "github.com/goradd/goradd/pkg/strings"
	"github.com/goradd/goradd/pkg/sys"
)

var assetDirectories = map[string]string{}

// GetAssetLocation returns the disk location of the asset file indicated by the given url.
// Asset directories must be registered with the RegisterAssetDirectory function. In debug mode, the
// file is taken from the registered location, but in release mode, the file will have been copied to
// a location on the server, and we will serve the file from there.
func GetAssetLocation(url string) string {
	// If we have an AssetDirectory, either we are in release mode, or we are locally testing the release process
	if config.AssetDirectory() != "" {
		if !strings2.StartsWith(url, config.AssetPrefix) {
			panic("Assets must start with the asset prefix.")
		}
		fPath := StripCacheBusterPath(url)
		fPath = strings.TrimPrefix(fPath, config.AssetPrefix)
		return filepath.Join(config.AssetDirectory(), filepath.Clean(fPath))
	}
	for dirUrl, dir := range assetDirectories {
		if strings2.StartsWith(url, dirUrl) {
			fPath := strings.TrimPrefix(url, dirUrl)
			return filepath.Join(dir, filepath.Clean(fPath))
		}
	}
	return ""
}

// GetAssetUrl returns the url that corresponds to the asset at the given location. Its the reverse of
// GetAssetLocation.
func GetAssetUrl(location string) string {
	var outPath string
	if config.AssetDirectory() != "" {
		if config.Release {
			if !strings2.StartsWith(location, config.AssetPrefix) {
				panic("In the release build, asset locations should be the same as asset paths. " + location + " does not start with " + config.AssetPrefix)
			}
		} else {
			// debug build, but a given asset directory.
			for url, dir := range assetDirectories {
				if strings2.StartsWith(location, dir) {
					fPath := strings.TrimPrefix(location, dir)
					location = path.Join(url, filepath.ToSlash(fPath))
					break
				}
			}
		}
		outPath = CacheBustedPath(location)
	} else {
		for url, dir := range assetDirectories {
			if strings2.StartsWith(location, dir) {
				fPath := strings.TrimPrefix(location, dir)
				outPath = path.Join(url, filepath.ToSlash(fPath))
				break
			}
		}
	}
	if outPath == "" {
		return ""
	}
	return config.MakeLocalPath(outPath)
}

// ServeAsset is the default server for files in asset directories.
//
// Some things you get for free are:
// - It will look for and use zipped versions of assets in the release version if available. See the goradd-project/build directory
//   for examples on how to do that.
// - It will use Go's default file server, which automatically gives you Last-Modified support in the header of these assets in
//   the release version of the application. However, in order to actually get the client to validate, you must specify a
//   Cache-Control policy of no-cache, or a max-age, etc. with the file. Also, remember that this will check the modification
//   date of the file on the server, which will change each time you copy the file to the server.
//
// Much more complicated servers are possible. You could, for example, package all your static files with the application itself
// and use an ETag to control caching. But, that is not available here, you would need to implement that yourself.
func ServeAsset(w http.ResponseWriter, r *http.Request) {
	localpath := GetAssetLocation(r.URL.Path)
	if localpath == "" {
		log.Printf("Invalid asset %s", r.URL.Path)
		return
	}
	//log.Printf("Served %s", localpath)

	if !config.Release && config.AssetDirectory() == "" {
		if ext := filepath.Ext(localpath); ext == "" {
			panic("Asset file does not have an extension: " + localpath)
		}

		// During development, tell the browser not to cache our assets so that if we change an asset, we don't have to deal with
		// trying to get the browser to refresh
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate, max-age=1")
		http.ServeFile(w, r, localpath)
	} else {
		var ext = filepath.Ext(localpath)

		var minFileName string

		if !strings2.EndsWith(localpath, ".min"+ext) {
			minFileName = localpath[:len(localpath)-len(ext)] + ".min" + ext
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

		type compType struct {
			file string
			typ  string
		}
		// build search file list in the order we want to use them
		var files []compType
		if acceptsBr {
			if minFileName != "" {
				files = append(files, compType{minFileName + ".br", "br"})
			}
			files = append(files, compType{localpath + ".br", "br"})
		}
		if acceptsGzip {
			files = append(files, compType{localpath + ".gz", "gzip"})
			if minFileName != "" {
				files = append(files, compType{minFileName + ".gz", "gzip"})
			}
		}

		if minFileName != "" {
			files = append(files, compType{minFileName, ""})
		}

		if ext != "" {
			for _, comp := range files {
				if sys.PathExists(comp.file) {
					if comp.typ != "" {
						w.Header().Set("Content-Encoding", comp.typ)
					}
					ctype := mime.TypeByExtension(ext)
					w.Header().Set("Content-Type", ctype)
					http.ServeFile(w, r, comp.file)
					return
				}
			}
		}
		http.ServeFile(w, r, localpath)
	}
}

// RegisterAssetDirectory registers the given directory as a static file server. The files are served by the normal
// go static file process. RegisterAssetDirectory must be called during application initialization, as the static file directories
// are added to the MUX at startup time.
func RegisterAssetDirectory(dir string, pattern string) {
	if d, ok := assetDirectories[pattern]; ok && d != dir {
		panic(pattern + " is already registered as an asset directory. ")
	}
	assetDirectories[pattern] = dir
}


// CacheBustedPath returns a path to an asset that was previously registered with the CacheBuster. The new path
// will contain a hash of the file that will change whenever the file changes, and cause the browser to reload the file.
// Since we are in control of serving these files, we will later remove the hash before serving it.
func CacheBustedPath(url string) string {
	if p,ok := config.CacheBuster[url]; ok {
		// inject the crc as a part of the path.
		dir,file := path.Split(url)
		url = path.Join(dir, config.CacheBusterPrefix + p, file)
	}
	return url
}

// StripCacheBusterPath removes the hash of the asset file from the path to the asset.
func StripCacheBusterPath(fPath string) string {
	dir,f := path.Split(fPath)
	dir2,f2 := path.Split(dir[:len(dir) - 1]) // ignore trailing slash
	// ignore the cache buster
	if strings2.StartsWith(f2, config.CacheBusterPrefix) {
		fPath = path.Join(dir2, f)
	}
	return fPath
}

