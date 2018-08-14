package config

import (
	"path"
	"runtime"
)

type AppMode int

const (
	AppModeDevelopment    AppMode = iota
	AppModeDeploymentPrep    // preparing for deployment
	AppModeDebug
	AppModeRelease
)

var Mode AppMode // filled in by Goradd. The application mode.

var LocalDir string
var GoraddDir string // filled in by Goradd
var ProjectDir string
var Minify bool

const JQuery = ""

const PagePathPrefix = ""     // A path added to all goradd pages
const AssetPrefix = "/assets" // A prefix added to all asset paths (js, css, etc. is added after this)

const MultiPartFormMax = 10000000 // 10 MB max size for uploaded forms. Change how you want.

const PageCacheMaxSize = 1000     // How many items are allowed in the override state before they are thrown out
const PageCacheTTL = 60 * 60 * 24 // Time that a override can go untouched before it is thrown from the cache (in nanos). This value is a day.

const MaxBufferSize = 10000 // Maximum size we will allow in our memory buffer pool. If too small, we will thrash memory. If too big, we will waste memory.

const PageCacheVersion = int32(1) // Change this every time we might change how pages are cached. Only needed if we are serializing pages.
const DefaultPageSize = 10 // Default number of items in Pager controlled controls
const MaxPageButtons = 10 // Maximum number of override buttons in a Pager control


func init() {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("No caller information")
	}
	LocalDir = path.Dir(path.Dir(filename))
	ProjectDir = path.Dir(LocalDir) + "/project"
}

func GoraddAssets() string {
	return GoraddDir + "/assets"
}

func LocalAssets() string {
	return LocalDir + "/assets"
}

func ProjectAssets() string {
	return ProjectDir + "/assets"
}

const (
	DefaultDateFormat = "January 2, 2006"
	DefaultTimeFormat = "3:04 am"
	DefaultDateTimeFormat = "January 2, 2006 3:04am"
)
