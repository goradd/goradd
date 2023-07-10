package config

// AssetPrefix is the path prefix for all goradd assets. It indicates to the program to
// look for the given file in the assets collection of files
// which in development mode is wherever the file is on the disk, and in release mode, the central asset directory where
// all assets get copied. set to blank to turn off the default asset management.
var AssetPrefix = "/assets/"

// WebsocketMessengerPrefix is the url prefix that indicates this is a Websocket call to our messenger service.
//
// The default turns on Websockets and uses this to implement the Watcher and Messenger
// mechanisms. Override this in the goradd_project/config/goradd.go file to set to the value of your choice,
// or set to blank to turn off handling of websockets.
var WebsocketMessengerPrefix = "/ws/"

// Minify controls whether we try to strip out unnecessary whitespace from our HTML output
var Minify bool = !Debug

// ProxyPath is the url path to the application. By default, this is the root, but you can set it
// to any path. This is particularly useful to making the application appear as if it is running in a subdirectory
// of the root path. This is great for putting behind an Apache server, and using ProxyPass and ProxyPassReverse to direct
// traffic from a particular path to the application. This gets stripped off incoming urls automatically by the server,
// but needs to be added to all links to resources on the server, and to cookies.
var ProxyPath string

var DefaultDateFormat = "January 2, 2006"
var DefaultTimeFormat = "3:04 pm"
var DefaultDateTimeFormat = "January 2, 2006 3:04 pm"

var DefaultDateEntryFormat = "1/2/06"
var DefaultTimeEntryFormat = "3:04 pm"
var DefaultDateTimeEntryFormat = "1/2/06 3:04 pm"

// SelectOneString is used in selection lists as the default item to indicate that a selection is required but has not yet been made
var SelectOneString = "- Select One -"

// NoSelectionString is used in selection lists as the item that indicates no selection when a selection is not required
var NoSelectionString = "-"

// DefaultFormFieldWrapperIdSuffix is the default suffix added to form control wrapper ids
var DefaultFormFieldWrapperIdSuffix = "-ff"
