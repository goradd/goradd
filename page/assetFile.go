package page

import (
	"path/filepath"
	"goradd/config"
	"github.com/spekary/goradd/util"
	"github.com/spekary/goradd/util/types"
	"github.com/spekary/goradd/html"
	"os"
	"strings"
	"net/http"
	"log"
)

// A css, js or other file we want the browser to add
type AssetFile struct {
	url string
	filePath string
	attributes *html.Attributes
	localPath string	// where this file is located locally, in case we copied it to a more accessible location
}


var assetFiles = types.NewSafeOrderedMap()

// RenderAssetTag will render a tag that points to a static file asset that should be served by the MUX. filePath points
// to the file on the development server, and it will be copied to the appropriate subdirectory in the local assets directory
// for easy deployment. tag is the tag name to put in the tag, and attributes are additional attributes to include in the tag.
// The copied location of the file and structure of the tag will be deduced from the type of tag and the name of the file.
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

// RegisterAssetFile registers an asset file with the global asset file manager. url is the path from docroot that will
// appear in the browser, and by convention it is of the form /dir/filename.
// filePath is the path on the development system where the file is located. This file will be copied to the url path
// under the config.LocalAssets() directory if the app is in development mode.
// Returns the url. Panics if the url is already associated with a different filePath.
func RegisterAssetFile(url string, filePath string) string {
	if !assetFiles.Has(url) {
		var dir,fileName string = filepath.Split(url)

		dir = strings.TrimPrefix(dir, config.ASSET_PREFIX)

		var localDir = config.LocalAssets() + dir
		var localPath = localDir + fileName

		// if we are in the correct mode, copy the file to the local assets directory. Otherwise, we will trust its already there.
		if config.Mode == config.DeploymentPrep {
			os.MkdirAll(localDir, 0777)
			err := util.FileCopyIfNewer(filePath, localPath)
			if err != nil {
				panic(err)
			}
		}

		a := AssetFile{url:url, filePath: filePath, localPath: localPath}
		assetFiles.Set(url, a)
		return url
	} else {
		if !assetIsRegistered(url) {
			panic ("No file for " + url + " has been registered.")
		}
		a := assetFiles.Get(url).(AssetFile)
		if config.Mode <= config.Debug {
			if a.filePath != filePath {
				panic("Attempting to register two different files to the same url:" + filePath)
			}
		}
		return url
	}
}

func RegisterCssFile(urlPath string, filePath string) string {
	return RegisterAssetFile(config.ASSET_PREFIX + "/css/" + urlPath, filePath)
}

func RegisterJsFile(urlPath string, filePath string) string {
	return RegisterAssetFile(config.ASSET_PREFIX + "/js/" + urlPath, filePath)
}



func GetAssetFilePath(url string) string {
	if asset := assetFiles.Get(url); asset == nil {
		return ""
	} else if config.Mode == config.Dev {
		return asset.(AssetFile).filePath
	} else {
		return asset.(AssetFile).localPath
	}

}

func assetIsRegistered(url string) bool {
	return assetFiles.Has(url)
}

func ServeAsset (w http.ResponseWriter, r *http.Request) {
	localpath := GetAssetFilePath(r.URL.Path)
	if localpath == "" {
		log.Printf("Invalid asset %s", r.URL.Path)
		return
	}
	//log.Printf("Served %s", localpath)

	if config.Mode == config.Dev {
		// TODO: Set up per file cache control
		w.Header().Set("Cache-Control",  "no-cache, no-store, must-revalidate")
	} else {
		// TODO: Set up a validating cache control
	}

	http.ServeFile(w, r, localpath)
}
