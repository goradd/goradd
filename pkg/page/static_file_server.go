package page

import (
	"github.com/spekary/goradd/pkg/config"
	"log"
	"mime"
	"net/http"
	"path"
	"path/filepath"
	"strings"

	strings2 "github.com/spekary/goradd/pkg/strings"
	"github.com/spekary/goradd/pkg/sys"
)

var assetDirectories = map[string]string{}

// RenderAssetTag will render a tag that points to a static file asset that should be served by the MUX. filePath points
// to the file on the development server, and it will be copied to the appropriate subdirectory in the local assets directory
// for easy deployment. tag is the tag label to put in the tag, and attributes are additional attributes to include in the tag.
// The copied location of the file and structure of the tag will be deduced from the type of tag and the label of the file.
// The type of tag will also be used to automatically insert the location of the file into the correct tag attribute.
/*
func RenderAssetTag(filePath string, tag string, attributes *html.Attributes, content string) string {
	var typ string

	_,fileName := filepath.Split(filePath)
	ext := path.Ext(fileName)


	switch ext {
	case "js": fallthrough
	case "javascript":
		typ = "js"
	case "css":
		typ = "css"
	case "jpg": fallthrough
	case "jpeg":fallthrough
	case "png": fallthrough
	case "gif": fallthrough
	case "bmp": fallthrough
	case "ico":
		typ = "img"
	default:
		switch tag {
		case "script":
			typ = "js"
		case "a":
			typ = "file" // a download file type likely, if we haven't already recognized it as something else.
		default:
			panic ("Unknown file type")
		}
	}

	url := "/" + typ + "/" + fileName

	url = RegisterAssetFile(url, filePath)

	switch tag {
	case script
	}
}
*/

func GetAssetLocation(url string) string {
	// If we have an AssetDirectory, either we are in release mode, or we are locally testing the release process
	if config.AssetDirectory != "" {
		if !strings2.StartsWith(url, config.AssetPrefix) {
			panic("Assets must start with the asset prefix.")
		}
		fPath := strings.TrimPrefix(url, config.AssetPrefix)
		return filepath.Join(config.AssetDirectory, filepath.FromSlash(fPath))
	}
	for dirUrl, dir := range assetDirectories {
		if strings2.StartsWith(url, dirUrl) {
			fPath := strings.TrimPrefix(url, dirUrl)
			return filepath.Join(dir, filepath.FromSlash(fPath))
		}
	}
	return ""
}

func GetAssetUrl(location string) string {
	if config.Release {
		if !strings2.StartsWith(location, config.AssetPrefix) {
			panic("In the release build, asset locations should be the same as asset paths.")
		}
		return location
	}
	for url, dir := range assetDirectories {
		if strings2.StartsWith(location, dir) {
			fPath := strings.TrimPrefix(location, dir)
			return path.Join(url, filepath.ToSlash(fPath))
		}
	}
	return ""
}

func ServeAsset(w http.ResponseWriter, r *http.Request) {
	localpath := GetAssetLocation(r.URL.Path)
	if localpath == "" {
		log.Printf("Invalid asset %s", r.URL.Path)
		return
	}
	//log.Printf("Served %s", localpath)

	if !config.Release && config.AssetDirectory == "" {
		// TODO: Set up per file cache control

		if ext := filepath.Ext(localpath); ext == "" {
			panic("Asset file does not have an extension: " + localpath)
		}

		// During development, tell the browser not to cache our assets so that if we change an asset, we don't have to deal with
		// trying to get the browser to refresh
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		http.ServeFile(w, r, localpath)
	} else {
		// TODO: Set up a validating cache control
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
// go static file process. This must happen during application initialization, as the static file directories
// are added to the MUX at startup time.
func RegisterAssetDirectory(dir string, pattern string) {
	if _, ok := assetDirectories[pattern]; ok {
		panic(pattern + " is already registered as an asset directory.")
	}
	assetDirectories[pattern] = dir
}
