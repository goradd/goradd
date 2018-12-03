package config

const GoraddWebSocketPort		= 8101
const GoraddWebSocketTLSPort 	= 8102 // This will require ssl certificates

// You will need to put in the path to your certfile and keyfile below.
// The default implementation only uses these for the release build.
// You can use the same ones that you use for normal SSL communication over http.
const GoraddWebSocketTLSCertFile = ""
const GoraddWebSocketTLSKeyFile = ""

