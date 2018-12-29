// The app package contains your local application object. It uses an object oriented model to implement a default
// application, and provides hooks for you to customize its behavior. Web applications can grow in complicated ways,
// and this is the main place you will customize how the server itself behaves.
package app

import (
	"fmt"
	"github.com/goradd/goradd/pkg/config"
	"github.com/goradd/goradd/pkg/messageServer"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/web/app"
	"log"
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
func MakeApplication() *Application {

	configure() // see config.go

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

func (a *Application) RunWebServer() (err error) {
	// The message server communicates to the browser UI changes caused by database changes. If you are simply redrawing
	// everything manually, and are not concerned about multi-user scenarios, you can comment it out.
	messageServer.Start(a.MakeWebsocketMux())

	mux := a.MakeServerMux()

	// If you are directly responding to encrypted requests, launch a server here. Note that you CAN put the app behind
	// a web server, like Nginx or Apache and let the web server handle the certificate issues.

	if config.TLSPort != 0 {
		go func() {
			log.Fatal(http.ListenAndServeTLS(fmt.Sprintf(":%d", config.TLSPort), config.TLSCertFile, config.TLSKeyFile, mux))
		}()
	}

	// The  "Serve" functions below will launch go routines for each request, so that multiple requests can be
	// processed in parallel.
	if config.UseFCGI { // Run as FCGI via standard I/O
		// FCGI requires a serialized pagestate server, outside session server and care with the websocket server
		err = fcgi.Serve(nil, mux)
	} else {
		addr := fmt.Sprintf(":%d", config.Port)
		err = http.ListenAndServe(addr, mux)
	}

	return
}

func (a *Application) MakeServerMux() *http.ServeMux {
	mux := http.NewServeMux()

	// Add handlers for your straight html files and anything you want to simply be served without processing and
	// that you can put in a specific directory. Note that you can also serve files from the ServeRequest handler below.

	// This registers the goradd-project/web/html directory to serve files with urls that start with "/html". .
	if !config.Release {
		mux.Handle("/html/", http.StripPrefix("/html/", http.FileServer(http.Dir(filepath.Join(config.ProjectDir(), "web", "html")))))
	} else {
		// TODO: Do the equivalent of ServeAsset with html files. Perhaps just expand ServeAsset to do html too.
		//mux.Handle("/html/", http.StripPrefix("/html/", http.FileServer(http.Dir(filepath.Join(config.ProjectDir(), "web", "html")))))
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
// declared handler, or by the goradd app server, or the static file server. ServeRequest is only called when all
// the other methods have failed. The default serves up a 404 not found error, and you can customize whatever error
// message you want to present to the user when a bad url is entered into the browser.
func (a *Application) ServeRequest (w http.ResponseWriter, r *http.Request) {
	s := fmt.Sprintf("<!DOCTYPE html><html><h1>Page Not Found</h1><p>The page you are looking for is not here.</p></html>")
	http.Error(w, s, http.StatusNotFound)

	// If you want to log these errors to detect a potentially bad link somewhere on your site, uncomment below.
	//log.Error("Bad url entered: " + r.URL.Path)
}


// SessionHandler initializes the global session handler. The default version uses the scs session handler.
// To replace it with the session handler of your choice, uncomment the code below and implement your session handler here.
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


// ServeStaticFile serves up static html and other files. The default will serve up the generated form index
// and any files you put in the HTML directory. If you want to serve up files from other directories, uncomment
// the line below, but remember you will have to put those files on your release server and point your custom
// static file server there.
/*
func (a *Application) ServeStaticFile (w http.ResponseWriter, r *http.Request) bool {

	// If you do not want the default behavior, remove the following lines
	if a.Application.ServeStaticFile(w,r) {
		return true
	}

	// Serve files from other directories here

	return false	// indicates no static file was found
}
*/
