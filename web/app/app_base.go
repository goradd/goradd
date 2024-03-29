// Package app contains the default web application server.
package app

import (
	"context"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/alexedwards/scs/v2/memstore"
	"github.com/goradd/goradd/pkg/base"
	"github.com/goradd/goradd/pkg/config"
	"github.com/goradd/goradd/pkg/goradd"
	http2 "github.com/goradd/goradd/pkg/http"
	grlog "github.com/goradd/goradd/pkg/log"
	"github.com/goradd/goradd/pkg/messageServer"
	"github.com/goradd/goradd/pkg/messageServer/ws"
	"github.com/goradd/goradd/pkg/orm/broadcast"
	"github.com/goradd/goradd/pkg/orm/db"
	"github.com/goradd/goradd/pkg/session"
	"github.com/goradd/goradd/pkg/watcher"
	"github.com/goradd/html5tag"
	"net/http/pprof"
	"time"

	"github.com/goradd/goradd/pkg/page"
	"net/http"
	"os"

	_ "github.com/goradd/goradd/web/assets"
)

// ApplicationI allows you to override default functionality in the ApplicationBase.
// The main routine offers a way of creating mock applications, and alternate versions of the application from the default
type ApplicationI interface {
	Init()
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	PutContext(*http.Request) *http.Request
	SetupErrorHandling()
	SetupPagestateCaching()
	SetupSessionManager()
	SetupMessenger()
	SetupDatabaseWatcher()
	SetupPaths()
	SessionHandler(next http.Handler) http.Handler
	HSTSHandler(next http.Handler) http.Handler
	AccessLogHandler(next http.Handler) http.Handler
	PutDbContextHandler(next http.Handler) http.Handler
	ServeAppMux(next http.Handler) http.Handler
	ServeRequestHandler() http.Handler
	ServePatternMux(next http.Handler) http.Handler
}

// Application is the application base, to be embedded in your application
type Application struct {
	base.Base
	httpErrorReporter http2.ErrorReporter
}

// Init initializes the application base.
func (a *Application) Init(self ApplicationI) {
	a.Base.Init(self)

	self.SetupErrorHandling()
	self.SetupPagestateCaching()
	self.SetupSessionManager()
	self.SetupMessenger()
	self.SetupPaths()
	self.SetupDatabaseWatcher()

	page.DefaultCheckboxLabelDrawingMode = html5tag.LabelAfter
}

func (a *Application) this() ApplicationI {
	return a.Self().(ApplicationI)
}

// SetupErrorHandling prepares the application for various types of error handling
func (a *Application) SetupErrorHandling() {

	// Create the top level http error reporter that will catch panics throughout the application
	// The default will intercept anything unexpected and set it to StdErr. Override this to do something else.
	a.httpErrorReporter = http2.ErrorReporter{}

}

// SetupPagestateCaching sets up the service that saves pagestate information that reflects the state of a goradd form to
// our go code. The default sets up a one server-one process cache that does not scale, which works great for development, testing, and
// for moderate amounts of traffic. Override and replace the page cache with one that serializes the page state and saves
// it to a database to make it scalable.
func (a *Application) SetupPagestateCaching() {
	// Controls how pages are cached. This will vary depending on whether you are using multiple machines to run your app,
	// and whether you are in development mode, etc. This default is for an in-memory store on one server and only one
	// process on that server. It basically does not serialize anything and leaves the entire pagestate intact in memory.
	// This makes for a very fast server, but one that takes up quite a bit of RAM if you have a lot of simultaneous users.
	if page.GetPagestateCache() == nil {
		page.SetPagestateCache(page.NewFastPageCache(1000, 60*60*24))
	}

	// Controls how pages are serialized if a serialization cache is being used. This version uses the gob encoder.
	// You likely will not need to change this, but you might if your database cannot handle binary data.
	page.SetPageEncoder(page.GobPageEncoder{})
}

// SetupSessionManager sets up the global session manager. The session can be used to save data that is specific to a user
// and specific to the user's time on a browser. Sessions are often used to save login credentials so that you know
// the current user is logged in.
//
// The default uses a 3rd party session manager, stores the session in memory, and tracks sessions using cookies.
// This setup is useful for development, testing, debugging, and for moderately used websites.
// However, this default does not scale, so if you are launching multiple copies of the app in production,
// you should override this with a scalable storage mechanism.
func (a *Application) SetupSessionManager() {
	s := scs.New()
	store := memstore.NewWithCleanupInterval(24 * time.Hour) // replace this with a different store if desired
	s.Store = store
	if config.ProxyPath != "" {
		s.Cookie.Path = config.ProxyPath
	}
	sm := session.NewScsManager(s)
	sm.(session.ScsManager).SessionManager.IdleTimeout = 6 * time.Hour

	session.SetSessionManager(sm)
}

// SetupMessenger injects the global messenger that permits pub/sub communication between the server and client.
//
// You can use this mechanism to set up your own messaging system for application use too.
func (a *Application) SetupMessenger() {
	// The default sets up a websocket based messenger appropriate for development and single-server applications
	messenger := new(ws.WsMessenger)
	messageServer.Messenger = messenger
	messenger.Start()
}

// SetupPaths sets up default path handlers for the application
func (a *Application) SetupPaths() {
	if config.Debug {
		// standard go profiling support
		http2.RegisterPrefixHandler("/debug/pprof/", http.HandlerFunc(pprof.Index))
		http2.RegisterPrefixHandler("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
		http2.RegisterPrefixHandler("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
		http2.RegisterPrefixHandler("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
		http2.RegisterPrefixHandler("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))
	}

	if config.WebsocketMessengerPrefix != "" {
		http2.RegisterPrefixHandler(config.WebsocketMessengerPrefix, http.HandlerFunc(WebsocketMessengerHandler))
	}
}

// SetupDatabaseWatcher injects the global database watcher
// and the database broadcaster which together detect database changes and
// then draws controls that are watching for those
func (a *Application) SetupDatabaseWatcher() {
	watcher.Watcher = &watcher.DefaultWatcher{}
	broadcast.Broadcaster = &broadcast.DefaultBroadcaster{}
}

// PutContext puts application data into the context.
func (a *Application) PutContext(r *http.Request) *http.Request {
	return page.PutContext(r, os.Args[1:])
}

// ServeHTTP serves a Goradd form.
func (a *Application) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pm := page.GetPageManager()
	if pm.IsPage(r.URL.Path) {
		ctx := r.Context()
		pm.RunPage(ctx, w, r)
	}
}

// MakeAppServer creates the handler chain that will handle http requests. There are a ton of ways to do this, 3rd party
// libraries to help with this, and middlewares you can use. This is a working example, and not a declaration of any
// "right" way to do this, since it can be very application specific. Generally you must make sure that
// PutAppContextHandler is called before ServePageHandler in the chain.
func (a *Application) MakeAppServer() http.Handler {
	// the handler chain gets built in the reverse order of getting called

	// These handlers are called in reverse order
	h := a.this().ServeRequestHandler() // Should go at the end of the chain to catch whatever is missed above
	h = a.this().ServeAppMux(h)         // Serves other dynamic files, and possibly the api
	h = a.ServePageHandler(h)           // Serves the Goradd dynamic pages
	h = a.PutAppContextHandler(h)
	h = a.this().SessionHandler(h)
	h = a.BufferedOutputHandler(h) // Must be in front of the session handler
	h = a.StatsHandler(h)
	h = a.this().ServePatternMux(h) // Serves most static files and websocket requests.
	// Must be after the error handler so panics are intercepted by the error reporter
	// and must be in front of the buffered output handler because of websocket server
	h = a.this().PutDbContextHandler(h) // This is here so that the PatternMux handlers can use the ORM
	h = a.validateHttpHandler(h)
	h = a.httpErrorReporter.Use(h) // Default http error handler to intercept panics.
	h = a.this().HSTSHandler(h)
	h = a.this().AccessLogHandler(h)

	return h
}

// SessionHandler initializes the global session handler. This default version uses the injected global session handler. Feel
// free to replace it with the session handler of your choice.
func (a *Application) SessionHandler(next http.Handler) http.Handler {
	return session.Use(next)
}

// HSTSHandler sets the browser to HSTS mode using the given timeout. HSTS will force a browser to accept only
// HTTPS connections for everything coming from your domain, if the initial page was served over HTTPS. Many browsers
// already do this. What this additionally does is prevent the user from overriding this. Also, if your
// certificate is bad or expired, it will NOT allow the user the option of using your website anyways.
// This should be safe to send in development mode if your local server is not using HTTPS, since the header
// is ignored if a page is served over HTTP.
//
// Once the HSTS policy has been sent to the browser, it will remember it for the amount of time
// specified, even if the header is not sent again. However, you can override it by sending another header, and
// clear it by setting the timeout to 0. Set the timeout to -1 to turn it off. You can also completely override this by
// implementing this function in your app.go file.
func (a *Application) HSTSHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if config.HSTSTimeout >= 0 {
			w.Header().Set("Strict-Transport-Security", fmt.Sprintf("max-age=%d; includeSubDomains", config.HSTSTimeout))
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

// ServeRequestHandler is the last handler on the default call chain.
// It returns a simple not found error.
// By default, this handler is never reached, because of the html root handler registered in
// goradd-project/web/embedder.go. You will need to modify or delete that handler
// to reach this handler. You can override this handler by duplicating it in your app object.
func (a *Application) ServeRequestHandler() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	}
	return http.HandlerFunc(fn)
}

// ServePageHandler processes requests for automated goradd pages.
func (a *Application) ServePageHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		pm := page.GetPageManager()
		if pm == nil {
			panic("No page manager defined")
		}
		if pm.IsPage(r.URL.Path) {
			a.ServeHTTP(w, r)
		} else {
			next.ServeHTTP(w, r)
		}
	}
	return http.HandlerFunc(fn)
}

// PutAppContextHandler is an HTTP handler that adds the application context to the current context.
func (a *Application) PutAppContextHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		r = a.this().PutContext(r)
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

// PutDbContextHandler is an HTTP handler that adds the database contexts to the current context.
func (a *Application) PutDbContextHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = db.PutContext(ctx)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

// BufferedOutputHandler manages the buffering of http output.
// It will save all output in a buffer, and make sure any and all Header sets can happen before
// writing the buffer out to the stream.
func (a *Application) BufferedOutputHandler(next http.Handler) http.Handler {
	return http2.BufferedOutputManager().Use(next)
}

// RegisterStaticPath registers the given url path such that it points to the given directory in the host file system.
//
// For example, passing "/test", "/my/test/dir" will serve documents out of /my/test/dir whenever a URL
// has /test in front of it. Only call this during application startup.
//
// hide is a list of file endings for files that should not be served. If a file is searched for, but is not
// found, a NotFound error will be sent to the http server.
func RegisterStaticPath(
	path string,
	directory string,
	useCacheBuster bool,
	hide []string,
) {
	fileSystem := os.DirFS(directory)
	fs := http2.FileSystemServer{
		Fsys:           fileSystem,
		SendModTime:    !useCacheBuster,
		UseCacheBuster: useCacheBuster,
		Hide:           hide}
	http2.RegisterPrefixHandler(path, fs)
	grlog.FrameworkInfof("Registering static path %s to %s", path, directory)
}

// ServeAppMux serves up the http2.AppMuxer, which handles REST calls,
// and dynamically created files.
//
// To use your own AppMuxer, override this function in app.go and create your own.
func (a *Application) ServeAppMux(next http.Handler) http.Handler {
	return http2.UseAppMuxer(http2.NewMux(), next)
}

// ServePatternMux serves up the http2.PatternMuxer.
//
// The pattern muxer routes patterns early in the handler stack. It is primarily for handlers that
// do not need the session manager or buffered output, things like static files.
//
// The default version injects a standard http muxer. Override to use your own muxer.
func (a *Application) ServePatternMux(next http.Handler) http.Handler {
	return http2.UsePatternMuxer(http2.NewMux(), next)
}

// ServeWebsocketMessengerHandler is the authorization handler for watcher requests.
// It uses the pagestate as the client id, verifying the page state is valid
func (a *Application) ServeWebsocketMessengerHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pagestate := r.FormValue("id")

		if !page.HasPage(pagestate) {
			// The page manager has no record of the pagestate, so either it is expired or never existed
			return // TODO: return error?
		}

		// Inject the pagestate as the client ID so the next handler down can read it
		ctx := context.WithValue(r.Context(), goradd.WebSocketContext, pagestate)
		messageServer.Messenger.(*ws.WsMessenger).WebSocketHandler().ServeHTTP(w, r.WithContext(ctx))
	})
}

// WebsocketMessengerHandler is a handler for Websockets.
func WebsocketMessengerHandler(w http.ResponseWriter, r *http.Request) {
	pagestate := r.FormValue("id")

	if !page.HasPage(pagestate) {
		// The page manager has no record of the pagestate, so either it is expired or never existed
		return // TODO: return error?
	}

	// Inject the pagestate as the client ID so the next handler down can read it
	ctx := context.WithValue(r.Context(), goradd.WebSocketContext, pagestate)
	messageServer.Messenger.(*ws.WsMessenger).WebSocketHandler().ServeHTTP(w, r.WithContext(ctx))
}

// AccessLogHandler simply logs requests.
func (a *Application) AccessLogHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		grlog.FrameworkInfo("Serving: ", r.RequestURI)
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

// validateHttpHandler performs OWASP style validation on a request.
func (a *Application) validateHttpHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if !http2.ValidateHeader(r.Header) {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			next.ServeHTTP(w, r)
		}
	}
	return http.HandlerFunc(fn)
}
