package config

import "github.com/goradd/goradd/pkg/config"

// The variables below are initialized from the main() function. Delete whatever you don't need.

// Port is the http port the app will run on. If not set, it will use the default of 80.
var Port int = 0

// The following variables are needed only if you are serving secure communications directly from the application.
// You can alternatively put the application behind a reverse proxy provided by apache or nginx and let those
// applications handle the encryption layer. See the build/docker/docker-compose.yml files for more info.

// TLSPort is the port to serve SSL from. If this is unset, it will not serve SSL. If set, you must also
// provide a TLSCertFile and TLSKeyFile below.
var TLSPort int = 0

// TLSCertFile is the path for the https certificate.
var TLSCertFile = ""

// TLSKeyFile is the path to the https key file.
var TLSKeyFile = ""

// This is space to put your apps own config constants and variables, or modify globals as needed.

func initApp() {
	if !config.Release {
		// Put initializations here just for the development build
	} else {
		// Put initializations here for the release build
	}

	// Put initializations here for all builds
}
