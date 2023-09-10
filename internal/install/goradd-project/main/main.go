package main

// The spaces in this import block serve to prevent goimports from rearranging the order of the files.
import (
	"github.com/goradd/goradd/pkg/log"
	_ "goradd-project/config" // Initialize required variables. This MUST come first.
	config2 "goradd-project/config"

	// _ "goradd-project/api" // Uncomment this if you are implementing an API (i.e. REST api).

	"flag"
	"fmt"
	"github.com/goradd/goradd/pkg/config"
	"github.com/goradd/goradd/pkg/sys"
	"goradd-project/web/app"
	_ "goradd-project/web/form" // Registers forms through init calls.
	// _ "goradd-project/api" // Uncomment this if you are implementing an API (i.e. REST api).
	// Custom paths, including additional form directories
	// _ "mysite"
)

var port = flag.Int("port", -1, "Serve as a webserver from the given port. Default = 80.")
var tlsPort = flag.Int("tlsPort", -1, "Serve securely from the given port.")
var tlsKeyFile = flag.String("keyFile", "", "Path to key file for tls.")
var tlsCertFile = flag.String("certFile", "", "Path to cert file for tls.")

var proxyPath = flag.String("proxyPath", "", "The url path to the application.")
var logLevel = flag.Int("logLevel", 0, "The logging level. See the log package for details.")

// Create other flags you might care about here

// dbConfigFile is actually read and used in config/db.go, but we define it here so it can be part of the usage message
var _ = flag.String("dbConfigFile", "", "The path to the database configuration file.")

func main() {
	var err error

	flag.Parse()
	useFlags()
	a := app.MakeApplication()
	fmt.Println("\nLaunching server on " + sys.GetIpAddress())
	err = a.RunWebServer()

	if err != nil {
		fmt.Println(err)
	}
}

func useFlags() {
	if *proxyPath != "" {
		config.ProxyPath = *proxyPath
	}

	if *port != -1 {
		config2.Port = *port
	}

	if *tlsPort != -1 {
		config2.TLSPort = *tlsPort
	}

	if *tlsKeyFile != "" {
		config2.TLSKeyFile = *tlsKeyFile
	}
	if *tlsCertFile != "" {
		config2.TLSCertFile = *tlsCertFile
	}
	log.SetLoggingLevel(*logLevel)
}
