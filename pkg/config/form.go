package config

import (
	"path"
)

// AssetPrefix is the path prefix for all goradd assets. It indicates to the program to look for the given file in the assets collection of files
// which in development mode is wherever the file is on the disk, and in release mode, the central asset directory where
// all assets get copied
var AssetPrefix = "/assets/"

// Minify controls whether we try to strip out unnecessary whitespace from our HTML output
var Minify bool = !Debug

var assetDirectory string
var htmlDirectory string

var DefaultDateFormat = "January 2, 2006"
var DefaultTimeFormat = "3:04 am"
var DefaultDateTimeFormat = "January 2, 2006 3:04am"

var jQueryPath string
var jQueryPathParams map[string]string


// JQueryPath returns the url, and possibly other parameters needed to get JQuery
func JQueryPath() (string, map[string]string) {
	if jQueryPath == "" {
		return path.Join(GoraddAssets(),"/js/jquery3.js"), nil
	} else {
		return jQueryPath, jQueryPathParams
	}
}

func SetJqueryPath(path string, params map[string]string) {
	jQueryPath = path
	jQueryPathParams = params
}

var jQueryUIPath string
var jQueryUIPathParams map[string]string

// JQueryUIPath returns the url, and possible other parameters to get JQueryUI.
func JQueryUIPath() (string, map[string]string) {
	if jQueryUIPath == "" {
		return path.Join(GoraddAssets(),"/js/jquery-ui.js"), nil
	} else {
		return jQueryUIPath, jQueryUIPathParams
	}
}

func SetJqueryUIPath(path string, params map[string]string) {
	jQueryUIPath = path
	jQueryUIPathParams = params
}


func SetAssetDirectory(assetDir string) {
	if Release && assetDir == "" {
		panic("The -assetDir flag is required when running the release build")
	}
	assetDirectory = assetDir
}

func AssetDirectory() string {
	return assetDirectory
}

func SetHtmlDirectory(d string) {
	htmlDirectory = d
}

func HtmlDirectory() string {
	return htmlDirectory
}



