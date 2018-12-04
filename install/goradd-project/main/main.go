package main

import (
	"flag"
	"fmt"
	"github.com/spekary/goradd/pkg/config"
	"goradd-project/web/app"
	// Below is where you import packages that register forms
	_ "goradd-project/web/form" // Your  forms.

	// Custom paths, including additional form directories
	// _ "mysite"
)

var port = flag.Int("port", 0, "Serve as webserver from given port. Default = 8000.")
var tlsPort = flag.Int("tlsPort", 0, "Serve securely from the given port.")
var wsPort = flag.Int("wsPort", 0, "Serve the websocket from given port.")
var wsTlsPort = flag.Int("wsTlsPort", 0, "Serve securely the websocket from the given port.")
var tlsKeyFile = flag.String("keyFile", "", "Path to key file for tls.")
var tlsCertFile = flag.String("certFile", "", "Path to cert file for tls.")
var wsKeyFile = flag.String("wsKeyFile", "", "Path to key file for websocket.")
var wsCertFile = flag.String("wsCertFile", "", "Path to cert file for websocket.")

var useFcgi = flag.Bool("fcgi", false, "Serve as fcgi.")
var assetDir = flag.String("assetDir", "", "The centralized asset directory. Required to run the release version of the app.")
var htmlDir = flag.String("htmlDir", "", "The centralized html directory. Required to run the release version of the app if you are serving static files.")
// Create other flags you might care about here

func main() {
	var err error

	useFlags()
	a := app.MakeApplication()
	err = a.RunWebServer()

	if err != nil {
		fmt.Println(err)
	}
}

func useFlags() {
	if *assetDir != "" {
		config.SetAssetDirectory(*assetDir)
	}

	if *htmlDir != "" {
		config.SetHtmlDirectory(*htmlDir)
	}

	if *useFcgi {
		config.UseFCGI = true
	}

	if *port != 0 {
		config.Port = *port
	}

	if *tlsPort != 0 {
		config.TLSPort = *tlsPort
	}


	if *wsPort != 0 {
		config.WebSocketPort = *wsPort
	}

	if *wsTlsPort != 0 {
		config.WebSocketTLSPort = *wsTlsPort
	}

	if *tlsKeyFile != "" {
		config.TLSKeyFile = *tlsKeyFile
	}
	if *tlsCertFile != "" {
		config.TLSCertFile = *tlsCertFile
	}
	if *wsKeyFile != "" {
		config.WebSocketTLSKeyFile = *wsKeyFile
	}
	if *wsCertFile != "" {
		config.WebSocketTLSCertFile = *wsCertFile
	}
}

func init() {
	flag.Parse() // Parse the command line flags
}
