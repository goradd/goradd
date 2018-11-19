package main

import (
	"flag"
	"fmt"
	"github.com/spekary/goradd/web/app"
	"github.com/spekary/goradd/global"
	"net/http"
	fcgiserver "net/http/fcgi"
	"os"
	"path/filepath"

	// local imports
	localapp "goradd-project/app"

	"github.com/alexedwards/scs"
	"github.com/alexedwards/scs/stores/memstore"
	"github.com/spekary/goradd/pkg/session"
	"time"
	"bytes"
	"github.com/spekary/goradd/pkg/page"
	"goradd-project/config"
	_ "goradd-project/form" // Your pre-built goradd forms. Move these to another package as needed.

	// These are the packages that contain your actual goradd forms. init() code should register the forms
	_ "github.com/spekary/goradd/pkg/bootstrap/examples"

	// Custom paths, including additional form directories
	_ "site"

)

var local = flag.String("local", "", "serve as webserver from given port, example: -local 8000")
var fcgi = flag.Bool("fcgi", false, "serve as fcgi, example: -fcgi")
var assetDir = flag.String("assetDir", "", "The centralized asset directory. Required to run the release version of the app.")

// Create other flags you might care about here

func main() {
	var err error

	if *local != "" || *fcgi { // Run as a local web server
		err = runWebServer()
	} else {
		// Run in command line mode, and pass the flags on to make a context.
		a := goraddApp
		a.ProcessCommand(os.Args[1:])
		//log.Fatal(err)
	}

	if err != nil {
		fmt.Println(err)
	}
}

func runWebServer() (err error) {

	mux := http.NewServeMux()

	// Add handlers for your straight html files and anything not processed by goradd
	if !config.Release {
		mux.Handle("/form/", http.StripPrefix("/form/", http.FileServer(http.Dir(filepath.Join(config.ProjectDir() , "form")))))
	}

	// serve up local static asset files
	mux.Handle(config.AssetPrefix, http.HandlerFunc(page.ServeAsset)) // serve up application assets
	mux.Handle("/", makeAppServer())   // send anything you don't explicitly handle to goradd

	// The two "Serve" functions below will launch go routines for each request, so that multiple requests can be
	// processed in parallel. This may mean multiple requests for the same override, depending on the structure of the override.
	if *local != "" { // Run as a local web server
		addr := ":" + *local
		err = http.ListenAndServe(addr, mux)
	} else if *fcgi { // Run as FCGI via standard I/O
		err = fcgiserver.Serve(nil, mux)
	}

	return err
}

// makeAppServer creates the handler chain that will handle http requests. There are a ton of ways to do this, 3rd party
// libraries to help with this, and middlewares you can use. This is a working example, and not a declaration of any
// "right" way to do this, since it can be very application specific. There are some general requirements though:
// 1) You must call the putContextHandler before calling the serveAppHandler in the chain.
func makeAppServer() http.Handler {
	// the handler chain gets built in the reverse order of getting called
	buf := page.GetBuffer()
	defer page.PutBuffer(buf)

	// These handlers are called in reverse order
	h := serveFileHandler(buf)
	h = serveAppHandler(buf, h)
	h = sessionHandler(h)
	h = putContextHandler(h)

	return h
}

// putContextHandler is an http handler that adds the application context to the current context
func putContextHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		r = goraddApp.PutContext(r)

		// Setup the output buffer
		grctx := page.GetContext(r.Context())
		grctx.AppContext.OutBuf = page.GetBuffer()
		defer page.PutBuffer(grctx.AppContext.OutBuf)
		next.ServeHTTP(w, r)
		w.Write(grctx.AppContext.OutBuf.Bytes())
	}
	return http.HandlerFunc(fn)
}

// serveAppHandler is the main handler that processes the current request
func serveAppHandler(buf *bytes.Buffer, next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		goraddApp.ServeHTTP(w, r)
		if next != nil {
			next.ServeHTTP(w, r)
		}
	}
	return http.HandlerFunc(fn)
}

// serveFileHandler is where you would serve static html files
func serveFileHandler(buf *bytes.Buffer) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		grctx := page.GetContext(r.Context())

		if grctx.AppContext.OutBuf.Len() == 0 {
			found := serveStaticFile(w,r)
			// if not handled by your static html output, return not found error
			if !found {
				w.WriteHeader(404)
				fmt.Fprint(w, "<!DOCTYPE html><html><h1>Not Found</h1></html>")
			}

		}
	}
	return http.HandlerFunc(fn)

}


// Customize this function to serve your static files that you can't serve using the mux
func serveStaticFile(w http.ResponseWriter, r *http.Request) bool {
	return false
}

// sessionHandler initializes the global session handler. This default version uses the scs session handler. Feel
// free to replace it with the session handler of your choice.
func sessionHandler(next http.Handler) http.Handler {
	return session.Use(next)
}

// Global app object. Use GetApplication to get it.
// The application object will initialize BEFORE any init() package functions, so they can call GetApplication and get the
// global application object

var goraddApp app.ApplicationI = makeApplication()

// Create the application object and related objects
// You can potentially read command line params and make other versions of the app for testing purposes
func makeApplication() app.ApplicationI {
	flag.Parse() // Parse the flags so we can read them

	config.InitDatabases()
	config.Init(*assetDir)

	// create various flavors of application here
	a := &localapp.Application{}
	a.Init()
	global.App = a // inject the created app into the global space

	// create the session manager. The default uses an in-memory storage engine. Change as you see fit.
	interval, _ := time.ParseDuration("24h")
	session.SetSessionManager(session.NewSCSManager(scs.NewManager(memstore.New(interval))))

	return a
}
