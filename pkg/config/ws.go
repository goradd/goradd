package config

var		UseFCGI bool

var 	Port int = 8000
var 	TLSPort int = 0 // This will require ssl certificates. The default has this turned off.

// You will need to put in the path to your certfile and keyfile below.
// The default implementation only uses these for the release build.
var  	TLSCertFile = ""
var  	TLSKeyFile = ""

var 	WebSocketPort int = 8101 // Default can be reset later, or via command line, but before the application starts. Set to zero to turn it off.
var 	WebSocketTLSPort int = 0 // This will require ssl certificates. The default has this turned off.

// You will need to put in the path to your certfile and keyfile below.
// The default implementation only uses these for the release build.
// You can use the same ones that you use for normal SSL communication over http.
var  WebSocketTLSCertFile = ""
var  WebSocketTLSKeyFile = ""

