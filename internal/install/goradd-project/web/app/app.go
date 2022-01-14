// The app package contains your local application object. It uses an object oriented model to implement a default
// application, and provides hooks for you to customize its behavior. Web applications can grow in complicated ways,
// and this is the main place you will customize how the server itself behaves.
package app

import (
	"fmt"
	"github.com/goradd/goradd/pkg/config"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/web/app"
	"log"
	"net/http"
	"net/http/pprof"
)

type Application struct {
	app.Application

	// Your own vars, methods and overrides
}

// Create the application object and related objects
// You can potentially read command line params and make other versions of the app for testing purposes
func MakeApplication() *Application {
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
func (a *Application) SetupPagestateCaching() {
	// Controls how pages are cached. This will vary depending on whether you are using multiple machines to run your app,
	// and whether you are in development mode, etc.

	// This default is for an in-memory store on one server and only one
	// process on that server. It basically does not serialize anything and leaves the entire formstate intact in memory.
	// This makes for a very fast server, but one that takes up quite a bit of RAM if you have a lot of simultaneous users.
	page.SetPagestateCache(page.NewFastPageCache())

	// comment out the above, and uncomment below to change to a serialized page cache. It is still
	// in memory, but can be used to test whether the page cache could be stored in a database instead.
	//page.SetPagestateCache(page.NewSerializedPageCache(100, 60*60*24))

	// Controls how pages are serialized if a serialization cache is being used. This version uses the gob encoder.
	// You likely will not need to change this, but you might if your database cannot handle binary data.
	page.SetPageEncoder(page.GobPageEncoder{})
}
*/

// InitializeLoggers allows you to set up the various types of logs for various types of builds. By default, the DebugLog
// and FrameworkDebugLogs will be deactivated when the config.Debug variables are false. Otherwise, configure how you
// want, and simply remove a log if you don't want it to log anything.
/*
func (a *Application) InitializeLoggers() {
	a.Application.InitializeLoggers()

	// This example turns the error log into an email logger in release mode so we will be notified when errors
	// occur on our production server
	if config.Release {
		log2.Loggers[log2.ErrorLog] = log2.EmailLogger{log.New(os.Stdout,
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

// SetupMessenger injects the global messenger that permits pub/sub communication between the server and client.
// Uncomment the following if you need to change parameters on the hub.
// If you don't need this at all, you can uncomment below and simply make it an empty function.
// Or, you can setup a different pub/sub messaging service here.
/*
func (a *Application) SetupMessenger() {
	// The default sets up a websocket based messenger appropriate for development and single-server applications
	messenger := new (ws.WsMessenger)
	messageServer.Messenger = messenger
	hub := messenger.Start()
	hub.WriteWait = 10 * time.Second // for example
}

*/

/*
This is an example of how you can setup your own custom handler for websockets.
Implement the functions you need.

// This is an example of a websocket auth handler for a custom websocket based messenger.
// At a minimum you must identify the user and set the client ID so that messages go to that client.
func (a *Application) myWebsocketAuthHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.FormValue("id")

		// confirm that the client is authorized. Substitute your own way of doing this here.
		parts := strings.Split(id, "-")

		challenge := parts[0] + "mySalt"

		sum := sha256.Sum256([]byte(challenge))
		pSum := fmt.Sprintf("%x", sum)
		if pSum != parts[1] {
			return
		}

		// Put the client ID in the context so that the framework's websocket handler can be used it to identify the client
		ctx := context.WithValue(r.Context(), goradd.WebSocketContext, parts[0])
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
*/

// SetupDatabaseWatcher injects the global database watcher and broadcaster
// which detects database changes and then draws controls that are watching for those changes.
//
// The default uses the provided goradd websocket message server to broadcast changes to the database, which is sufficient
// for a single-server application. If you need a multi-server scalable version, change the watcher here to something that
// uses a distributed pub/sub mechanism.
//
// Changing the broadcaster will let you do additional things on the server side when specific
// database items change.
/*
func (a *Application) SetupDatabaseWatcher() {
	watcher.Watcher = &watcher.DefaultWatcher{}
	broadcast.Broadcaster = &broadcast.DefaultBroadcaster{}
}
*/

// RunWebServer launches the main webserver.
func (a *Application) RunWebServer() (err error) {
	mux := a.MakeServerMux()

	// If you are directly responding to encrypted requests, launch a server here. Note that you CAN put the app behind
	// a web server, like Nginx or Apache and let the web server handle the certificate issues.

	if config.TLSPort != 0 {
		go func() {
			var addr string
			if config.TLSPort != 0 {
				addr = fmt.Sprint(":", config.TLSPort)
			}

			log.Fatal(app.ListenAndServeTLSWithTimeouts(addr, config.TLSCertFile, config.TLSKeyFile, mux))
		}()
	}

	// TODO: Make a way so that we will automatically redirect to https if specified to do so
	// I think its a simple matter of providing a mux just for this purpose
	var addr string
	if config.Port != 0 {
		addr = fmt.Sprint(":", config.Port)
	}
	err = app.ListenAndServeWithTimeouts(addr, mux)

	return
}

func (a *Application) MakeServerMux() *http.ServeMux {
	mux := http.NewServeMux()

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

	// Below is an example of how you can add your own handler that piggybacks
	// on the framework's websocket messenger.
	//mux.Handle("/ws/myws", a.myWebsocketAuthHandler(messageServer.Messenger.(*ws.WsMessenger).WebSocketHandler()))

	
	// Serve up the websocket messenger
	mux.Handle(config.WebsocketMessengerPrefix, http.HandlerFunc(app.WebsocketMessengerHandler))

	// serve up static application asset files
	mux.Handle(config.AssetPrefix, http.HandlerFunc(page.ServeAsset))

	// send anything you don't explicitly handle above to the goradd app server
	// note that the app server can serve up static html too. See ServeStaticFile.
	mux.Handle("/", a.MakeAppServer())

	// Uncomment this and implement ServeData to serve up custom generated
	// files like PDFs, CSVs, images, etc.
	//http2.RegisterPathHandler("/data", http2.ErrorHandler(http.HandlerFunc(a.ServeData)))

	return mux
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	//http.ServeFile(w, r, page.GetAssetLocation("/assets/project/image/favicon.ico"))
}


// ServeRequest is the place to serve up any files that have not been handled in any other way, either by a previously
// declared handler, or by the goradd app server, or the static file server. ServeRequest is only called when all
// the other methods have failed. The default serves up a 404 not found error, and you can customize whatever error
// message you want to present to the user when a bad url is entered into the browser.
/*
func (a *Application) ServeRequest (w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)

	// If you want to log these errors to detect a potentially bad link somewhere on your site, uncomment below.
	//log.Error("Bad url entered: " + r.URL.Path)
}
*/

// SetupSessionManager sets up the global session manager. The session manager can be used to save data that is
// specific to a user and specific to the user's time on a browser. Sessions are often used to save
// login credentials so that you know the current user is logged in.
//
// The default uses a 3rd party session manager, stores the session in memory, and tracks sessions using cookies.
// This setup is useful for development, testing, debugging, and for moderately used websites.
// However, this default does not scale, so if you are launching multiple copies of the app in production,
// you should override this with a scalable storage mechanism.
//
// The code here calls the default, and then gives you the option to further refine how the session works.
/*
func (a *Application) SetupSessionManager() {
	a.Application.SetupSessionManager()
	sm := session.SessionManager()

	// If you are only serving your application over https, you should do this too for added security.
	if (config.Release) {
		sm.(session.ScsManager).SessionManager.Cookie.Secure = true
	}
}
*/

// SessionHandler initializes the global session handler. The default version uses the global session handler, which is
// highly configurable. However, if you want to use a completely different session handler, you can do so here.
/*
func (a *Application) SessionHandler(next http.Handler) http.Handler {
	return session.Use(next)
}
*/

// PutContext allocates a blank context object for our application specific context data, to be populated later.
// Activate it by uncommenting the function below, and then edit the accompanying context.go file to add your
// application specific context data.
/*
func (a *Application) PutContext(r *http.Request) *http.Request {
	ctx := r.Context()
	// ctx = auth.PutContext(ctx) uncomment this to use the auth package for authentication
	ctx = PutLocalContext(ctx) // puts your local context if you need that
	r = r.WithContext(ctx)
	// be sure to call the superclass version so the goradd framework can operate
	return a.Application.PutContext(r)
}
*/


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

// ServeData serves up custom generated files from the /data directory.
/*
  This is an example of how to create a custom data server. See above for
  how this is called.
func (a *Application) ServeData (w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "data.csv" {
		serveData(w,r)
		return
	}
	http.NotFound(w, r)
}
*/