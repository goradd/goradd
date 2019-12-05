package config

import "path"

// AssetPrefix is the path prefix for all goradd assets. It indicates to the program to look for the given file in the assets collection of files
// which in development mode is wherever the file is on the disk, and in release mode, the central asset directory where
// all assets get copied
var AssetPrefix = "/assets/"
var ApiPrefix = "/api/"

// Minify controls whether we try to strip out unnecessary whitespace from our HTML output
var Minify bool = !Debug

var assetDirectory string
var htmlDirectory string

// AliasPath is the url path to the application. By default, this is the root, but you can set it
// to any path. This is particularly useful to making the application appear as if it is running in a subdirectory
// of the root path. This is great for putting behind an Apache server, and using ProxyPass and ProxyPassReverse to direct
// traffic from a particular path to the application. This gets stripped off incoming urls automatically by the server,
// but needs to be added to all links to resources on the server.
var ProxyPath string

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

type LocalPathMaker func(string)string

var localPathMaker LocalPathMaker = defaultLocalPathMaker

func defaultLocalPathMaker(p string) string {
	var hasSlash bool
	if p == "" {
		panic(`cannot make a local path to an empty path. If you are trying to refer to the root, use '/'.`)
	}
	if p[len(p)-1:] == "/" {
		hasSlash = true
	}
	if p[0:1] == "/" {
		p = path.Join(ProxyPath, p) // will strip trailing slashes
		if hasSlash {
			p = p + "/"
		}
	}
	return p

}

// MakeLocalPath turns a path that points to a resource on this computer into a path that will reach
// that resource. It takes into account a variety of settings that may affect the path and that will
// depend on how the app is deployed.
// You can inject your own local path maker using SetLocalPathMaker
func MakeLocalPath(p string) string {
	return localPathMaker(p)
}

func SetLocalPathMaker(f LocalPathMaker) {
	localPathMaker = f
}