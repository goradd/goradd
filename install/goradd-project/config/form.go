package config

import "path"

const MultiPartFormMax = 10000000 // 10 MB max size for uploaded forms. Change how you want.

const PageCacheMaxSize = 1000     // How many items are allowed in the override state before they are thrown out
const PageCacheTTL = 60 * 60 * 24 // Time that a override can go untouched before it is thrown from the cache (in nanos). This value is a day.

const MaxBufferSize = 10000 // Maximum size we will allow in our memory buffer pool. If too small, we will thrash memory. If too big, we will waste memory.

const PageCacheVersion = int32(1) // Change this every time we might change how pages are cached. Only needed if we are serializing pages.
const DefaultPageSize = 10 // Default number of items in Pager controlled controls
const MaxPageButtons = 10 // Maximum number of override buttons in a Pager control

const PagePathPrefix = ""           // A path added to all goradd pages
const AssetPrefix = "/assets/" 		// A prefix to indicate we are serving up an asset

// Minify controls whether we try to strip out unnecessary whitespace from our HTML output
var Minify bool = !Debug

var AssetDirectory string


const (
	DefaultDateFormat = "January 2, 2006"
	DefaultTimeFormat = "3:04 am"
	DefaultDateTimeFormat = "January 2, 2006 3:04am"
)

// JQueryPath returns either the local physical location of jquery, or a URL where to get jQuery. The default below
// uses a local path for development, but gets jquery from a public location in the release version. Change how you want.
func JQueryPath() (string, map[string]string) {
	if Release {
		return "https://code.jquery.com/jquery-3.3.1.min.js", map[string]string{"integrity": "sha256-FgpCb/KJQlLNfOu91ta32o/NMZxltwRo8QtmkMRdAu8=", "crossorigin": "anonymous"}
	} else {
		return path.Join(GoraddAssets(),"/js/jquery3.js"), nil
	}
}

// JQueryUIPath returns either the local physical location of jquery, or a URL where to get jQuery. The default below
// uses a local path for development, but gets jquery from a public location in the release version. Change how you want.
func JQueryUIPath() (string, map[string]string) {
	if Release {
		return "https://code.jquery.com/ui/1.12.1/jquery-ui.min.js", map[string]string{"integrity": "sha256-VazP97ZCwtekAsvgPBSUwPFKdrwD3unUfSGVYrahUqU=", "crossorigin": "anonymous"}
	} else {
		return path.Join(GoraddAssets(),"/js/jquery-ui.js"), nil
	}
}

func Init(assetDir string) {
	if Release && assetDir == "" {
		panic("The -assetDir flag is required when running the release build")
	}
	AssetDirectory = assetDir
}