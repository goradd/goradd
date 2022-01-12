package app

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/alexedwards/scs/v2/memstore"
	"github.com/goradd/gengen/pkg/maps"
	"github.com/goradd/goradd/pkg/base"
	"github.com/goradd/goradd/pkg/config"
	"github.com/goradd/goradd/pkg/goradd"
	"github.com/goradd/goradd/pkg/html"
	http2 "github.com/goradd/goradd/pkg/http"
	grlog "github.com/goradd/goradd/pkg/log"
	"github.com/goradd/goradd/pkg/messageServer"
	"github.com/goradd/goradd/pkg/messageServer/ws"
	"github.com/goradd/goradd/pkg/orm/broadcast"
	"github.com/goradd/goradd/pkg/orm/db"
	"github.com/goradd/goradd/pkg/session"
	strings2 "github.com/goradd/goradd/pkg/strings"
	"github.com/goradd/goradd/pkg/sys"
	"github.com/goradd/goradd/pkg/watcher"
	"hash/crc64"
	"io/ioutil"
	"log"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/goradd/goradd/pkg/page"
	"net/http"
	"os"
)

// StaticDirectoryPaths is a map of patterns to directory locations to serve statically.
// These can be registered at the command line or in the application
var StaticDirectoryPaths *maps.StringSliceMap

// StaticBlacklist is the list of file terminators that specify what files we always want to hide from view
// when a static file directory is searched. The default will always hide .go files. Add to it if you have
// other kinds of files in your static directories that you do not want to show. Do this only at startup.
var StaticBlacklist = []string{".go"}

type staticFileProcessor struct {
	ending    string
	processor StaticFileProcessorFunc
}

type StaticFileProcessorFunc func(file string, w http.ResponseWriter, r *http.Request)

// StaticFileProcessors is a map that connects file endings to processors that will process the content and return it
// to the output stream, bypassing other means of processing static files.
var staticFileProcessors []staticFileProcessor

// The application interface. A minimal set of commands that the main routine will ask the application to do.
// The main routine offers a way of creating mock applications, and alternate versions of the application from the default
type ApplicationI interface {
	Init()
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	PutContext(*http.Request) *http.Request
	SetupErrorPageTemplate()
	SetupPagestateCaching()
	InitializeLoggers()
	SetupAssetDirectories()
	SetupSessionManager()
	SetupMessenger()
	SetupDatabaseWatcher()
	SetupCacheBuster()
	SessionHandler(next http.Handler) http.Handler
	HSTSHandler(next http.Handler) http.Handler
	ServeStaticFile(w http.ResponseWriter, r *http.Request) bool
	AccessLogHandler(next http.Handler) http.Handler
	PutDbContextHandler(next http.Handler) http.Handler
	ServeAppMuxHandler(next http.Handler) http.Handler
}

// The application base, to be embedded in your application
type Application struct {
	base.Base
}

func (a *Application) Init(self ApplicationI) {
	a.Base.Init(self)

	self.SetupErrorPageTemplate()
	self.SetupPagestateCaching()
	self.InitializeLoggers()
	self.SetupAssetDirectories()
	self.SetupSessionManager()
	self.SetupMessenger()
	self.SetupDatabaseWatcher()
	self.SetupCacheBuster()

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

// SetupPagestateCaching sets up the service that saves pagestate information that reflects the state of a goradd form to
// our go code. The default sets up a one server-one process cache that does not scale, which works great for development, testing, and
// for moderate amounts of traffic. Override and replace the page cache with one that serializes the page state and saves
// it to a database to make it scalable.
func (a *Application) SetupPagestateCaching() {
	// Controls how pages are cached. This will vary depending on whether you are using multiple machines to run your app,
	// and whether you are in development mode, etc. This default is for an in-memory store on one server and only one
	// process on that server. It basically does not serialize anything and leaves the entire pagestate intact in memory.
	// This makes for a very fast server, but one that takes up quite a bit of RAM if you have a lot of simultaneous users.
	page.SetPagestateCache(page.NewFastPageCache(1000, 60*60*24))

	// Controls how pages are serialized if a serialization cache is being used. This version uses the gob encoder.
	// You likely will not need to change this, but you might if your database cannot handle binary data.
	page.SetPageEncoder(page.GobPageEncoder{})
}

// SetupCacheBuster sets up the cache busting strategy. Cache busting permits the browser to cache files, but notifies
// the server when a file has changed without requiring the browser to check the server. The default cache busting strategy
// is to add a CRC value to the path on all the files in the assets directory. That will cause the file to change its name
// whenever the file changes, forcing the browser to reload the file.
func (a *Application) SetupCacheBuster() {
	t := crc64.MakeTable(crc64.ECMA)
	config.CacheBuster = make (map[string]string)
	assetDir := config.AssetDirectory()
	if assetDir == "" {return}
	if err := filepath.Walk(assetDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err // stop
		}
		if info.IsDir() {
			return nil // keep going
		}
		if filepath.Ext(path) == ".gz" {
			// skip encrypted files, since they will be handled automatically
			return nil
		}
		if data,err := ioutil.ReadFile(path); err != nil {
			return err
		} else {
			// CRC it
			c := crc64.Checksum(data, t)
			e := strconv.FormatInt(int64(c), 36)
			s := strings.TrimPrefix(path, assetDir)
			s = filepath.ToSlash(s)
			s = filepath.Join(config.AssetPrefix, s)
			config.CacheBuster[s] = e
			return nil
		}
	}); err != nil {
		panic("failed walking the asset directory " + err.Error())
	}
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
	page.RegisterAssetDirectory(config.GoraddAssets(), config.AssetPrefix+"goradd")
	page.RegisterAssetDirectory(config.ProjectAssets(), config.AssetPrefix+"project")

	// If serving static html out of the root path, this will point it to the HtmlDirectory
	if dir := config.HtmlDirectory(); dir != "" {
		RegisterStaticPath("/", dir)
	}
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
	session.SetSessionManager(sm)
}

// SetupMessenger injects the global messenger that permits pub/sub communication between the server and client.
//
// You can use this mechanism to setup your own messaging system for application use too.
func (a *Application) SetupMessenger() {
	// The default sets up a websocket based messenger appropriate for development and single-server applications
	messenger := new (ws.WsMessenger)
	messageServer.Messenger = messenger
	messenger.Start()
}


// SetupDatabaseWatcher injects the global database watcher
// and the database broadcaster which together detect database changes and
// then draws controls that are watching for those
func (a *Application) SetupDatabaseWatcher() {
	watcher.Watcher = &watcher.DefaultWatcher{}
	broadcast.Broadcaster = &broadcast.DefaultBroadcaster{}
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
	buf := page.OutputBuffer(ctx)
	if pm.IsPage(r.URL.Path) {
		headers, errCode := pm.RunPage(ctx, buf)
		if errCode == 0 {
			if headers != nil {
				for k, v := range headers {
					// Multi-value headers can simply be separated with commas I believe
					w.Header().Set(k, v)
				}
			}
		} else {
			log.Print(errCode)
			grctx := page.GetContext(ctx)
			if grctx.RequestMode() == page.Ajax {
				js := []interface{}{errCode, headers}
				s,err := json.Marshal(js)
				if err == nil {
					w.Write(s)
				}
				w.WriteHeader(400)
			} else {
				if headers != nil {
					for k, v := range headers {
						// Multi-value headers can simply be separated with commas I believe
						w.Header().Set(k, v)
					}
				}
				w.WriteHeader(errCode)
			}

		}
	}
}


// MakeAppServer creates the handler chain that will handle http requests. There are a ton of ways to do this, 3rd party
// libraries to help with this, and middlewares you can use. This is a working example, and not a declaration of any
// "right" way to do this, since it can be very application specific. Generally you must make sure that
// PutAppContextHandler is called before ServeAppHandler in the chain.
func (a *Application) MakeAppServer() http.Handler {
	// the handler chain gets built in the reverse order of getting called

	// These handlers are called in reverse order
	h := a.ServeRequestHandler()
	h = a.this().ServeAppMuxHandler(h)
	h = a.ServeStaticFileHandler(h)
	h = a.ServeAppHandler(h)
	h = a.PutAppContextHandler(h)
	h = a.this().PutDbContextHandler(h)
	h = a.this().SessionHandler(h)
	h = a.this().HSTSHandler(h)
	h = a.BufferedOutputHandler(h)
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
// It returns a simple not found error by default.
func (a *Application) ServeRequestHandler() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	}
	return http.HandlerFunc(fn)
}

// ServeStaticFileHandler serves up static files by calling ServeStaticFile.
//
// The difference between this and registering a handler with a muxer is that a muxer
// will return a 404 error if the file is not found, whereas the below method will pass
// control to the next handler if the file is not found.
// This lets you serve static files and dynamically generated files from the same logical web path.
func (a *Application) ServeStaticFileHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		if !a.this().ServeStaticFile(w, r) && next != nil {
			next.ServeHTTP(w, r)
		}
	}
	return http.HandlerFunc(fn)
}

// ServeAppHandler processes requests for goradd forms
func (a *Application) ServeAppHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		a.ServeHTTP(w, r)
		head := w.Header()
		if next != nil &&
			page.OutputBuffer(r.Context()).Len() == 0 &&
			len(head) <= 1 { // This could be a hack. Not sure of any other way to tell if we have responded.
			next.ServeHTTP(w, r)
		}
	}
	return http.HandlerFunc(fn)
}

// PutAppContextHandler is an http handler that adds the application context to the current context.
func (a *Application) PutAppContextHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		r = a.this().PutContext(r)
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

// PutDbContextHandler is an http handler that adds the database context to the current context.
//
// This allows the context to be used by various pieces of the app further down the chain. The default
// assumes a SQL database. Override it to change it.
func (a *Application) PutDbContextHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		// Create a context that the ORM can use
		ctx = context.WithValue(ctx, goradd.SqlContext, &db.SqlContext{})
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

// ServeStaticFile serves up static html and other files found in registered directories.
// If the file is not found, it will return false.
func (a *Application) ServeStaticFile(w http.ResponseWriter, r *http.Request) bool {
	url := r.URL.Path
	var path string

	// StaticDirectoryPaths should be sorted longest to shortest at this point
	StaticDirectoryPaths.Range(func(pattern string, dir string) bool {
		if strings2.StartsWith(url, pattern) {
			fPath := strings.TrimPrefix(url, pattern)
			if fPath != "" && fPath[0:1] != "/" {
				// We only matched part of a directory, so not a match
				return true // go to next iteration
			}
			cleaned := strings.TrimPrefix(fPath, "..") // This prevents someone from hacking by using .. to refer to files outside of the directory
			cleaned = filepath.Clean(cleaned)
			path = filepath.Join(dir, cleaned)
			return false // stop iterating
		}
		return true
	})

	if path == "" {
		return false // the directory was not found in our list of static file directories
	}

	for _, bl := range StaticBlacklist {
		if strings2.EndsWith(path, bl) {
			return false // cannot show this kind of file
		}
	}

	if sys.IsDir(path) {
		path = filepath.Join(path, "index.html")
	}

	if sys.PathExists(path) {
		for _, p := range staticFileProcessors {
			if strings2.EndsWith(path, p.ending) {
				p.processor(path, w, r)
				return true
			}
		}

		http.ServeFile(w, r, path)
		return true
	}

	return false // indicates no static file was found
}


// RegisterStaticPath registers the given url path such that it points to the given directory. For example, passing
// "/test", "/my/test/dir" will statically serve everything out of /my/test/dir whenever a url has /test in front of it.
// You can only call this during application startup.
func RegisterStaticPath(path string, directory string) {
	if path[0:1] != "/" {
		log.Fatal("path " + path + " must begin with a slash (must be a rooted path)")
	}

	if !sys.IsDir(directory) {
		log.Fatal("path " + directory + " is not a valid directory")
	}

	var err error
	directory,err = filepath.Abs(directory)
	if err != nil {
		log.Fatal("could not get absolute path of " + directory + ": " + err.Error())
	}

	if path[len(path)-1:] == "/" {
		// Strip ending slash so that we can handle both /a/b/ and /a/b urls as directories and treat them the same.
		path = path[:len(path)-1]
	}

	if StaticDirectoryPaths == nil {
		StaticDirectoryPaths = maps.NewStringSliceMap()
		// sort the directory paths longest to shortest so that when we iterate them, we won't short circuit
		// longer paths with shorter versions of the same path.
		StaticDirectoryPaths.SetSortFunc(func(key1,key2 string, val1, val2 string) bool {
			// order longest to shortest keys
			return len(key1) > len(key2)
		})
	}
	StaticDirectoryPaths.Set(path, directory)
	grlog.Infof("Registering static path %s to %s", path, directory)
}

// ServeAppMuxHandler serves up the AppMuxHandler, which handles REST calls,
// and dynamically created files.
//
// To use your own AppMuxer, simply set a new http.AppMuxer.
// To register additional handlers, override this.
func (a *Application) ServeAppMuxHandler(next http.Handler) http.Handler {
	if config.ApiManager != nil {
		config.ApiManager.Use()
	}

	return http2.UseAppMuxer(next)
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
		grlog.Info("Serving: ", r.RequestURI)
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

// RegisterStaticFileProcessor registers a processor function for static files that have a particular suffix.
// Do this at init time.
func RegisterStaticFileProcessor(ending string, processorFunc StaticFileProcessorFunc) {
	staticFileProcessors = append(staticFileProcessors, staticFileProcessor{ending, processorFunc})
}

// ListenAndServeTLSWithTimeouts starts a secure web server with timeouts. The default http server does
// not have timeouts by default, which leaves the server open to certain attacks that would start
// a connection, but then very slowly read or write. Timeout values are taken from global variables
// defined in config, which you can set at init time.
func ListenAndServeTLSWithTimeouts(addr, certFile, keyFile string, handler http.Handler) error {
	// Here we confirm that the CertFile and KeyFile exist. If they don't, ListenAndServeTLS just exit with an open error
	// and you will not know why.
	if !sys.PathExists(certFile) {
		log.Fatalf("TLSCertFile does not exist: %s", config.TLSCertFile)
	}

	if !sys.PathExists(keyFile) {
		log.Fatalf("TLSKeyFile does not exist: %s", config.TLSKeyFile)
	}

	// TODO: https://blog.gopheracademy.com/advent-2016/exposing-go-on-the-internet/ recommends keeping track
	// of open connections using the ConnState hook for debugging purposes.

	srv := &http.Server{
		Addr: addr,
		ReadTimeout:  config.ReadTimeout,
		ReadHeaderTimeout: config.ReadHeaderTimeout,
		WriteTimeout: config.WriteTimeout,
		IdleTimeout:  config.IdleTimeout,
		Handler:      handler,
	}
	return srv.ListenAndServeTLS(certFile, keyFile)
}

// ListenAndServeWithTimeouts starts a web server with timeouts. The default http server does
// not have timeouts, which leaves the server open to certain attacks that would start
// a connection, but then very slowly read or write. Timeout values are taken from global variables
// defined in config, which you can set at init time. This non-secure version is appropriate
// if you are serving behind another server, like apache or nginx.
func ListenAndServeWithTimeouts(addr string, handler http.Handler) error {

	// TODO: https://blog.gopheracademy.com/advent-2016/exposing-go-on-the-internet/ recommends keeping track
	// of open connections using the ConnState hook for debugging purposes.

	srv := &http.Server{
		Addr: addr,
		ReadTimeout:  config.ReadTimeout,
		ReadHeaderTimeout: config.ReadHeaderTimeout,
		WriteTimeout: config.WriteTimeout,
		IdleTimeout:  config.IdleTimeout,
		Handler:      handler,
	}
	return srv.ListenAndServe()
}