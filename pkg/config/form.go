package config

// AssetPrefix is the path prefix for all goradd assets. It indicates to the program to look for the given file in the assets collection of files
// which in development mode is wherever the file is on the disk, and in release mode, the central asset directory where
// all assets get copied
var AssetPrefix = "/assets/"
var ApiPrefix = "/api/"

// Minify controls whether we try to strip out unnecessary whitespace from our HTML output
var Minify bool = !Debug

var assetDirectory string
var htmlDirectory string

var DefaultDateFormat = "January 2, 2006"
var DefaultTimeFormat = "3:04 am"
var DefaultDateTimeFormat = "January 2, 2006 3:04am"

// SelectOneString is used in selection lists as the default item to indicate that a selection is required but has not yet been made
var SelectOneString = "- Select One -"
// NoSelectionString is used in selection lists as the item that indicates no selection when a selection is not required
var NoSelectionString = "-"

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
