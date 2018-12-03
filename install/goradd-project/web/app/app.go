// The app package contains your local application object. It uses an object oriented model to implement a default
// application, and provides hooks for you to customize its behavior. Web applications can grow in complicated ways,
// and this is the main place you will customize how the server itself behaves.
package app

import (
	"flag"
	"fmt"
	"github.com/spekary/goradd/pkg/messageServer"
	"github.com/spekary/goradd/pkg/page"
	"github.com/spekary/goradd/web/app"
	"goradd-project/config"
	"net/http"
	"net/http/fcgi"
	"net/http/pprof"
	"path/filepath"
)

type Application struct {
	app.Application

	// Your own vars, methods and overrides
}

// Create the application object and related objects
// You can potentially read command line params and make other versions of the app for testing purposes
func MakeApplication(assetDir string) *Application {
	flag.Parse() // Parse the flags so we can read them

	config.Init(assetDir)
	config.InitDatabases()

	a := new(Application)
	a.Init()
	return a
}

func (a *Application) Init() {
	a.Application.Init(a)

}

// Uncomment and edit to change the error page. You can call the embedded Application version first, and then alter it too.
/*
func (a *Application) SetupErrorPageTemplate() {
	if config.Debug {
		page.ErrorPageFunc = page.DebugErrorPageTmpl
	} else {
		page.ErrorPageFunc = page.ReleaseErrorPageTmpl
	}
}
*/

// Uncomment and edit to change the page cache. You can call the embedded Application version first, and then alter it too.
/*
func (a *Application) SetupPageCaching() {
	// Control how pages are cached. This will vary depending on whether you are using multiple machines to run your app,
	// and whether you are in development mode, etc. This default is for an in-memory store on one server and only one
	// process on that server. It basically does not serialize anything and leaves the entire formstate intact in memory.
	// This makes for a very fast server, but one that takes up quite a bit of RAM if you have a lot of simultaneous users.
	page.SetPageCache(page.NewFastPageCache())

	// Control how pages are serialized if a serialization cache is being used. This version uses the gob encoder.
	// You likely will not need to change this, but you might if your database cannot handle binary data.
	page.SetPageEncoder(page.GobPageEncoder{})
}
*/

// InitializeLoggers allows you to set up the various types of logs for various types of builds. By default, the DebugLog
// and FrameworkDebugLogs will be deactivated when the config.Debug variables are false. Otherwise, configure how you
// want, and simply remove a log if you don't want it to log anything. To activate it, just uncomment the function below.
/*
func (a *Application) InitializeLoggers() {
	a.Application.InitializeLoggers()

	// This example turns the error log into an email logger in release mode so we will be notified when errors
	// occur on our production server
	if config.Release {
		log.Loggers[log.ErrorLog] = log.EmailLogger{log.New(os.Stdout,
		"Error: ", log.Ldate|log.Lmicroseconds|log.Lshortfile), []string{"errors@mybusiness.com", "supervisor@mybusiness.com"}}
	}
}
*/

// SetupAssetDirectories sets up directories that will serve assets. Its best to put your assets in your project/assets
// directory, but if you need to serve assets from another directory too, you can uncomment the code below to add
// whatever assets you need.
/*
func (a *Application) SetupAssetDirectories() {
	a.Application.SetupAssetDirectories()
	page.RegisterAssetDirectory(location, config.AssetPrefix + name)

}
*/

// RunWebServer launches the main webserver.

func (a *Application) RunWebServer(port string, useFcgi bool) (err error) {
	// The message server communicates to the browser UI changes caused by database changes. If you are simply redrawing
	// everything manually, and are not concerned about multi-user scenarios, you can comment it out.
	messageServer.Start(a.MakeWebsocketMux())

	mux := a.MakeServerMux()

	// If you are directly responding to encrypted requests, launch a server here. Note that you CAN put the app behind
	// a web server, like Nginx or Apache and let the web server handle the certificate issues.

	/*
	if config.Release { // Depends on whether you need encryption during local development
		go func() {
			log.Fatal(http.ListenAndServeTLS(fmt.Sprintf(":%d", config.AppTLSPort), config.AppTLSCertFile, config.AppTLSKeyFile, mux))
		}()
	}*/

	// The  "Serve" functions below will launch go routines for each request, so that multiple requests can be
	// processed in parallel.
	if useFcgi { // Run as FCGI via standard I/O
		err = fcgi.Serve(nil, mux)
	} else {
		addr := ":" + port
		err = http.ListenAndServe(addr, mux)
	}

	return
}

func (a *Application) MakeServerMux() *http.ServeMux {
	mux := http.NewServeMux()

	// Add handlers for your straight html files and anything you want to simply be served without processing and
	// that you can put in a specific directory. Note that you can also serve files from the ServeRequest handler below.
	if !config.Release {
		// This registers the goradd-project/web/html directory to serve files with urls that start with "/html". Feel
		// free to change this.
		mux.Handle("/html/", http.StripPrefix("/html/", http.FileServer(http.Dir(filepath.Join(config.ProjectDir(), "web", "html")))))
	}

	if config.Debug {
		// standard go profiling support
		mux.HandleFunc("/debug/pprof/", pprof.Index)
		mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
		mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	}

	// Handle the favicon request.
	mux.Handle("/favicon.ico", http.HandlerFunc(faviconHandler))

	// serve up static asset files
	mux.Handle(config.AssetPrefix, http.HandlerFunc(page.ServeAsset)) // serve up application assets
	mux.Handle("/", a.MakeAppServer())   // send anything you don't explicitly handle above to the goradd page server

	return mux
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, page.GetAssetLocation("/assets/project/image/favicon.ico"))
}

// Uncomment the code below and add your own code if you need additional authentication mechansims for websockets
// beyond the default.
/*
func (a *Application) WebSocketAuthHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pagestate := r.FormValue("id")

		if !page.GetPageManager().HasPage(pagestate) {
			// The page manager has no record of the pagestate, so either it is expired or never existed
			return // TODO: return error?
		}

		next.ServeHTTP(w, r)
	})
}
*/

// ServeRequest is the place to serve up any files that have not been handled in any other way, either by a previously
// declared handler, or by the goradd app server. ServeRequest is only called when all the other methods have failed.
// This is a good place to handle serving up static html files, pdfs,
// or any other kind of custom request.
func (a *Application) ServeRequest (w http.ResponseWriter, r *http.Request) {
	url := r.URL.EscapedPath()

	if !config.Release {
		// serve up the /form/index.html file in development mode so that we can get to the code-generated forms.
		if url == "/form/index.html" || url == "/form" || url == "/form/" {
			http.ServeFile(w, r, filepath.Join(config.ProjectDir(), "gen", "index.html"))
			return
		}
	}

	// If the url simply points to nothing, then serve up an appropriate error page.
	w.WriteHeader(404)
	fmt.Fprint(w, "<!DOCTYPE html><html><h1>Page Not Found</h1><p>The page you are looking for is not here.</p></html>")
}


// SessionHandler initializes the global session handler. The default version uses the scs session handler.
// To replace it with the session handler of your choice, uncommend the code below and implement your session handler here.
/*
func (a *Application) SessionHandler(next http.Handler) http.Handler {
	return session.Use(next)
}
*/

// PutContext allocates a blank context object for our application specific context data, to be populated later.
// Activate it by uncommenting the function below, and then edit the accompanying context.go file to add your
// application specific context data.
func (a *Application) PutContext(r *http.Request) *http.Request {
	ctx := r.Context()
	ctx = PutContext(ctx)
	r = r.WithContext(ctx)
	// be sure to call the superclass version so the goradd framework can operate
	return a.Application.PutContext(r)
}

