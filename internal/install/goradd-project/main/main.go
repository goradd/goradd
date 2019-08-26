package main

import (
	_ "goradd-project/config" // Initialize required variables
	"goradd-project/web/app"
	"log"
	"strings"

	"flag"
	"fmt"
	"github.com/goradd/goradd/pkg/config"
	app2 "github.com/goradd/goradd/web/app"
	// Below is where you import packages that register forms
	_ "goradd-project/web/form" // Your  forms.
	// Custom paths, including additional form directories
	// _ "mysite"
)

var port = flag.Int("port", 0, "Serve as a webserver from the given port. Default = 8000.")
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
var staticPaths arrayFlags // See below. Holds the static paths to be registered with the app.
// Create other flags you might care about here

func main() {
	var err error

	useFlags()
	a := app.MakeApplication()
	fmt.Println("\nLaunching Server...")
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

	for _, s := range staticPaths {
		i := strings.IndexAny(s, ":;")
		if i == -1 {
			log.Fatal("map must specify a path and directory, separated by a colon or semicolon")
		}
		app2.RegisterStaticPath(s[:i], s[i+1:])
	}

}

// Setup the ability to add multiple static paths
type arrayFlags []string

func (i *arrayFlags) String() string {
	return fmt.Sprintf("%v", *i)
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func init() {
	flag.Var(&staticPaths, "map", "Map a physical directory or file to a url path. The path comes first, followed by the directory or file, separated by a comma or semicolon. Directory paths should end with a forward slash (/). You can specify this flag multiple times.")
	flag.Parse() // Parse the command line flags
}
