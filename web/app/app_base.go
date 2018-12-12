package app

import (
	"bytes"
	"github.com/alexedwards/scs"
	"github.com/alexedwards/scs/stores/memstore"
	"github.com/spekary/goradd/pkg/base"
	"github.com/spekary/goradd/pkg/config"
	"github.com/spekary/goradd/pkg/html"
	grlog "github.com/spekary/goradd/pkg/log"
	"github.com/spekary/goradd/pkg/session"
	"github.com/spekary/goradd/pkg/sys"
	"path/filepath"
	"time"

	"github.com/spekary/goradd/pkg/messageServer"
	"github.com/spekary/goradd/pkg/page"
	"net/http"
	"os"
)

// The application interface. A minimal set of commands that the main routine will ask the application to do.
// The main routine offers a way of creating mock applications, and alternate versions of the application from the default
type ApplicationI interface {
	Init()
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	PutContext(*http.Request) *http.Request
	SetupErrorPageTemplate()
	SetupPageCaching()
	InitializeLoggers()
	SetupAssetDirectories()
	SetupSessionManager()
	WebSocketAuthHandler(next http.Handler) http.Handler
	SessionHandler(next http.Handler) http.Handler
	ServeRequest (w http.ResponseWriter, r *http.Request)
	ServeStaticFile (w http.ResponseWriter, r *http.Request) bool
}

// The application base, to be embedded in your application
type Application struct {
	base.Base
}

func (a *Application) Init(self ApplicationI) {
	a.Base.Init(self)

	self.SetupErrorPageTemplate()
	self.SetupPageCaching()
	self.InitializeLoggers()
	self.SetupAssetDirectories()
	self.SetupSessionManager()

	page.DefaultCheckboxLabelDrawingMode = html.LabelAfter
}

func (a *Application) this() ApplicationI {
	return a.Self.(ApplicationI)
}

// SetupErrorPageTemplate sets the template that controls the output when an error happens during the processing of a
// page request, including any code that panics. By default, in debug mode, it will popup an error message in the browser with debug
// information when an error occurs. And in release mode it will popup a simple message that an error occurred and will log the
// error to the error log. You can implement this function in your local application object to override it and do something different.
func (a *Application) SetupErrorPageTemplate() {
	if config.Debug {
		page.ErrorPageFunc = page.DebugErrorPageTmpl
	} else {
		page.ErrorPageFunc = page.ReleaseErrorPageTmpl
	}
}

// SetupPageCaching sets up the service that saves pagestate information that reflects the state of a goradd form to
// our go code. The default sets up a one server-one process cache that does not scale, which works great for development, testing, and
// for moderate amounts of traffic. Override and replace the page cache with one that serializes the page state and saves
// it to a database to make it scalable.
func (a *Application) SetupPageCaching() {
	// Control how pages are cached. This will vary depending on whether you are using multiple machines to run your app,
	// and whether you are in development mode, etc. This default is for an in-memory store on one server and only one
	// process on that server. It basically does not serialize anything and leaves the entire formstate intact in memory.
	// This makes for a very fast server, but one that takes up quite a bit of RAM if you have a lot of simultaneous users.
	page.SetPageCache(page.NewFastPageCache(1000 ,60 * 60 * 24))

	// Control how pages are serialized if a serialization cache is being used. This version uses the gob encoder.
	// You likely will not need to change this, but you might if your database cannot handle binary data.
	page.SetPageEncoder(page.GobPageEncoder{})
}

// InitializeLoggers sets up the various types of logs for various types of builds. By default, the DebugLog
// and FrameworkDebugLogs will be deactivated when the config.Debug variables are false. Otherwise, configure how you
// want, and simply remove a log if you don't want it to log anything.
func (a *Application) InitializeLoggers() {
	grlog.CreateDefaultLoggers()
}

// SetupAssetDirectories registers default directories that will contain web assets. These assets are served up in
// place in development mode, and served from a specified asset directory in release mode. This means the assets will
// need to be copied to a central location and moved to the release server. See the build directory for info.
func (a *Application) SetupAssetDirectories() {
	page.RegisterAssetDirectory(config.GoraddAssets(), config.AssetPrefix + "goradd")
	page.RegisterAssetDirectory(config.ProjectAssets(), config.AssetPrefix + "project")
}


// SetupSessionManager sets up the session manager. The session can be used to save data that is specific to a user
// and specific to the user's time on a browser. Sessions are often used to save login credentials so that you know
// the current user is logged in.
//
// The default uses a 3rd party session manager, and stores the session in memory, which is useful for development,
// testing, debugging, and for moderately used websites. The default does not scale, so replace it with a different
// storage mechanism is you are launching multiple copies of the app.
func (a *Application) SetupSessionManager() {
	// create the session manager. The default uses an in-memory storage engine. Change as you see fit.
	interval, _ := time.ParseDuration("24h")
	session.SetSessionManager(session.NewSCSManager(scs.NewManager(memstore.New(interval))))
}

func (a *Application) PutContext(r *http.Request) *http.Request {
	return page.PutContext(r, os.Args[1:])
}

func (a *Application) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	pm := page.GetPageManager()
	if pm == nil {
		panic("No page manager defined")
	}

	ctx := r.Context()
	grctx := page.GetContext(ctx)
	buf := grctx.AppContext.OutBuf
	if pm.IsPage(r.URL.Path) {
		headers, errCode := pm.RunPage(ctx, buf)
		if headers != nil {
			for k, v := range headers {
				// Multi-value headers can simply be separated with commas I believe
				w.Header().Set(k, v)
			}
		}
		if errCode != 0 {
			w.WriteHeader(errCode)
		}
	}
}

// MakeWebsocketMux creates the mux for the default websocket handler. The default handler provides session data to
// the web socket handler below, since its very common to need to get to session data to authenticate to user before
// responding to the request.
func (a *Application) MakeWebsocketMux() (*http.ServeMux) {
	mux := http.NewServeMux()

	mux.Handle("/ws", a.this().SessionHandler(a.this().WebSocketAuthHandler(messageServer.WebsocketHandler())))

	return mux
}


// WebSocketAuthHandler is the default authenticator of the web socket. This version simply makes sure the form
// has a pagestate, since if it doesn't, we should not be handling a request. If you want to authenticate using
// information out of the session, like to see whether the user is logged in, you should override this in your
// application instance.
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

// MakeAppServer creates the handler chain that will handle http requests. There are a ton of ways to do this, 3rd party
// libraries to help with this, and middlewares you can use. This is a working example, and not a declaration of any
// "right" way to do this, since it can be very application specific. Generally you must make sure that
// PutContextHandler is called before ServeAppHandler in the chain.
func (a *Application) MakeAppServer() http.Handler {
	// the handler chain gets built in the reverse order of getting called
	buf := page.GetBuffer()
	defer page.PutBuffer(buf)

	// These handlers are called in reverse order
	h := a.ServeRequestHandler(buf)
	h = a.ServeStaticFileHandler(buf, h)	// TODO: Speed this handler up by checking to see if the url is a goradd form before deciding to get context and session
	h = a.ServeAppHandler(buf, h)
	h = a.this().SessionHandler(h)
	h = a.PutContextHandler(h)

	return h
}

// SessionHandler initializes the global session handler. This default version uses the scs session handler. Feel
// free to replace it with the session handler of your choice.
func (a *Application) SessionHandler(next http.Handler) http.Handler {
	return session.Use(next)
}

// ServeRequestHandler is the last handler on the default call chain. It calls ServeRequest so the sub-class can handle it.
func (a *Application) ServeRequestHandler(buf *bytes.Buffer) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		grctx := page.GetContext(r.Context())

		if grctx.AppContext.OutBuf.Len() == 0 {
			a.this().ServeRequest(w,r)
		}
	}
	return http.HandlerFunc(fn)
}

// ServeStaticFileHandler serves up static files by calling ServeStaticFile.
func (a *Application) ServeStaticFileHandler(buf *bytes.Buffer, next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		if !a.this().ServeStaticFile(w,r) && next != nil {
			next.ServeHTTP(w, r)
		}
	}
	return http.HandlerFunc(fn)
}


// ServeAppHandler is the main handler that processes the current request
func (a *Application) ServeAppHandler(buf *bytes.Buffer, next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		a.ServeHTTP(w, r)

		grctx := page.GetContext(r.Context())

		if next != nil && grctx.AppContext.OutBuf.Len() == 0 {
			next.ServeHTTP(w, r)
		}
	}
	return http.HandlerFunc(fn)
}

// PutContextHandler is an http handler that adds the application context to the current context.
func (a *Application) PutContextHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		r = a.this().PutContext(r)

		// Setup the output buffer
		grctx := page.GetContext(r.Context())
		grctx.AppContext.OutBuf = page.GetBuffer()
		defer page.PutBuffer(grctx.AppContext.OutBuf)
		next.ServeHTTP(w, r)
		_,_ = w.Write(grctx.AppContext.OutBuf.Bytes())
	}
	return http.HandlerFunc(fn)
}


// ServeStaticFile serves up static html and other files. The default will serve up the generated form index
// and any files you put in the HTML directory. It is overridable by creating a ServeStaticFile in your local
// app.go file.
func (a *Application) ServeStaticFile (w http.ResponseWriter, r *http.Request) bool {
	url := r.URL.Path

	if !config.Release {
		// serve up the /form/index.html file in development mode so that we can get to the code-generated forms.
		if url == "/form/index.html" || url == "/form" || url == "/form/" {
			http.ServeFile(w, r, filepath.Join(config.ProjectDir(), "gen", "index.html"))
			return true
		}
	}

	// Attempt to serve the file out of the html directory
	if dir := config.HtmlDirectory(); dir != "" {
		if url[len(url)-1:] == "/" {
			url += "index.html"
		}

		// This prevents someone from hacking by using .. to refer to files outside of the html directory
		fp := filepath.Clean(url)

		file := filepath.Join(dir, fp)
		if sys.PathExists(file) {
			http.ServeFile(w, r, file)
			return true
		}
	}

	return false	// indicates no static file was found
}


// ServeRequest is the place to serve up any files that have not been handled in any other way, either by a previously
// declared handler, or by the goradd app server, or the static file server. ServeRequest is only called when all
// the other methods have failed. The default serves up a 404 not found error. Override it to handle other files,
// or to change the messaging when a bad url is attempted.
func (a *Application) ServeRequest (w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}

