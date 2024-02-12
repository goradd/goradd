// Package app contains your local application object. It uses an object oriented model to implement a default
// application, and provides hooks for you to customize its behavior. Web applications can grow in complicated ways,
// and this is the main place you will customize how the server itself behaves.
package app

import (
	"fmt"
	"github.com/goradd/goradd/web/app"
	"goradd-project/config"
	"log"
	"net/http"
)

type Application struct {
	app.Application

	// Your own vars, methods and overrides
}

// MakeApplication creates the application object and related objects.
//
// You can potentially read command line params and make other versions of the app for testing purposes.
func MakeApplication() *Application {
	a := new(Application)
	a.Init()
	return a
}

// Init initializes the application object.
func (a *Application) Init() {
	a.Application.Init(a)
}

// MakeAppServer creates the handler that serves the application.
//
// This is typically where you create the middleware stack that divides the server
// into small pieces that each do one job.
//
// The default use the base Application's middleware stack, which itself is quite flexible
// and has hooks where you can override pieces. Or, you can just replace the whole thing
// and reimplement it here.
//
// See also the Init function where can assign additional handlers to specific paths via
// the application muxers.
func (a *Application) MakeAppServer() http.Handler {
	return a.Application.MakeAppServer()
}

// RunWebServer launches the main webserver.
func (a *Application) RunWebServer() (err error) {
	handler := a.MakeAppServer()

	// If you are directly responding to encrypted requests, launch a server here. Note that you CAN put the app behind
	// a web server, like Nginx or Apache and let the web server handle the certificate issues.

	if config.TLSPort != 0 {
		go func() {
			var addr string
			if config.TLSPort != 0 {
				addr = fmt.Sprint(":", config.TLSPort)
			}

			log.Fatal(app.ListenAndServeTLSWithTimeouts(addr, config.TLSCertFile, config.TLSKeyFile, handler))
		}()
	}

	var addr string
	if config.Port != 0 {
		addr = fmt.Sprint(":", config.Port)
	}
	err = app.ListenAndServeWithTimeouts(addr, handler)

	return
}

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

	// Set your idle and session lifetimes as appropriate
	sm.(session.ScsManager).SessionManager.IdleTimeout = 6 * time.Hour
	sm.(session.ScsManager).SessionManager.Lifetime = 24 * time.Hour

	if config.Release {
		// If you are only serving your application over https, you should do this too for added security.
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

// ServeRequestHandler is the last handler on the default call chain.
// The default below returns a simple not found error.
// By default, this handler is never reached, because of the html root handler registered in
// goradd-project/web/embedder.go. You will need to modify or delete that handler
// to reach this handler.
func (a *Application) ServeRequestHandler() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	}
	return http.HandlerFunc(fn)
}
*/
