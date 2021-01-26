package main

// The spaces in this import block serve to prevent goimports from rearranging the order of the files.
import (
	_ "goradd-project/config" // Initialize required variables. This MUST come first.

	// _ "goradd-project/api" // Uncomment this if you are implementing an API (i.e. REST api).

	"flag"
	"fmt"
	"github.com/goradd/goradd/pkg/config"
	"github.com/goradd/goradd/pkg/sys"
	app2 "github.com/goradd/goradd/web/app"
	"goradd-project/web/app"
	"log"
	"strings"

	// Below is where you import packages that register forms
	_ "goradd-project/web/form" // Your  forms.
	// Custom paths, including additional form directories
	// _ "mysite"
)

var port = flag.Int("port", -1, "Serve as a webserver from the given port. Default = 8000.")
var tlsPort = flag.Int("tlsPort", -1, "Serve securely from the given port.")
var tlsKeyFile = flag.String("keyFile", "", "Path to key file for tls.")
var tlsCertFile = flag.String("certFile", "", "Path to cert file for tls.")

var assetDir = flag.String("assetDir", "", "The centralized asset directory. Required to run the release version of the app.")
var htmlDir = flag.String("htmlDir", "", "The centralized html directory. Required to run the release version of the app if you are serving static files.")
var proxyPath = flag.String("proxyPath", "", "The url path to the application.")
var staticPaths arrayFlags // See below. Holds the static paths to be registered with the app.
// Create other flags you might care about here

func main() {
	var err error

	useFlags()
	a := app.MakeApplication()
	fmt.Println("\nLaunching server on " + sys.GetIpAddress())
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

	if *proxyPath != "" {
		config.ProxyPath = *proxyPath
	}

	if *port != -1 {
		config.Port = *port
	}

	if *tlsPort != -1 {
		config.TLSPort = *tlsPort
	}

	if *tlsKeyFile != "" {
		config.TLSKeyFile = *tlsKeyFile
	}
	if *tlsCertFile != "" {
		config.TLSCertFile = *tlsCertFile
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
	return fmt.Sprint(*i)
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func init() {
	flag.Var(&staticPaths, "map", "Map a physical directory or file to a url path. The path comes first, followed by the directory or file, separated by a comma or semicolon. Directory paths should end with a forward slash (/). You can specify this flag multiple times.")
	flag.Parse() // Parse the command line flags
}
