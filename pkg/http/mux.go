package http

import (
	"github.com/goradd/goradd/pkg/config"
	"net/http"
	"path"
)

// Muxer represents the typical functions available in a mux and allows you
// to replace the default muxer here with a 3rd party mux, like the Gorilla mux.
//
// However, beware. The default Go muxer will do redirects. If this goradd application
// is behind a reverse proxy that is rewriting the url, the Go muxer will not correctly
// do rewrites because it will not include the reverse proxy path in the rewrite
// rule, and things will break.
//
// If you create your own mux and you want to do redirects, use MakeLocalPath to
// create the redirect url. See also maps.SafeMap for a map you can use if you
// are modifying paths while using the mux.
type Muxer interface {
	// Handle associates a handler with the given pattern in the url path
	Handle(pattern string, handler http.Handler)

	// Handler returns the handler associate with the request, if one exists. It
	// also returns the actual path registered to the handler
	Handler(r *http.Request) (h http.Handler, pattern string)

	// ServeHTTP sends a request to the MUX, to be forwarded on to the registered handler,
	// or responded with an unknown resource error.
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

// Mux is the default muxer for Goradd.
//
// Mux was written to fix the following problems with the Go default muxer:
// - It cannot work behind a reverse proxy because of its rewrite rules
// - Its algorithm for finding a handler gets slow the more handlers you give it
// - It does read locks, which is unnecessary in our case since goradd does not modify the mux after startup
//
// If you register a path with a slash at the end of it, that handler will
// be redirected to the path with the slash at the end. To change this behavior,
// register for both paths, one with the slash, and one without.
//
// Once you use the mux, you cannot add entries to it. This prevents
// race condition problems if the mux is being changed and accessed at the
// same time from two different go routines.
type Mux struct {
	m           map[string]http.Handler
	wasAccessed bool
}

func NewMux() *Mux { return new(Mux) }

// Handle registers the handler for the given pattern.
// If a handler already exists for pattern, Handle panics.
func (mux *Mux) Handle(pattern string, handler http.Handler) {
	if pattern == "" {
		panic("http: cannot register an empty pattern")
	}
	if handler == nil {
		panic("http: nil handler")
	}
	if _, exist := mux.m[pattern]; exist {
		panic("http: multiple registrations for " + pattern)
	}

	if mux.wasAccessed {
		panic("http: attempted to modify the mux after using it")
	}

	if mux.m == nil {
		mux.m = make(map[string]http.Handler)
	}
	mux.m[pattern] = handler
}

// HandleFunc registers the handler function for the given pattern.
func (mux *Mux) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	if handler == nil {
		panic("http: nil handler")
	}
	mux.Handle(pattern, http.HandlerFunc(handler))
}

// Handler returns the handler to use for the given request r.URL.Path.
// The incoming path must be cleaned before being sent here.
//
// If the path cannot be found, but its equivalent path with an ending
// slash is found,  the
// handler for the slashed path will be returned.
//
// The path and host are used unchanged for CONNECT requests.
//
// Handler also returns the registered pattern that matches the
// request or, in the case of internally-generated redirects,
// the pattern that will match after following the redirect.
//
// If there is no registered handler that applies to the request,
// Handler returns a “page not found” handler and an empty pattern.
func (mux *Mux) Handler(r *http.Request) (h http.Handler, pattern string) {
	var pathChanged bool

	mux.wasAccessed = true
	path := r.URL.Path

	// CONNECT requests are not canonicalized.
	if r.Method != "CONNECT" {
		path = cleanPath(path)
		if path != r.URL.Path {
			pathChanged = true
		}
	}
	var ok bool

	// If the given path is /tree/ find /tree/{{config.DefaultPage}}. As in /tree/index.html.
	if path[len(path)-1] == '/' && config.DefaultPage != "" {
		np := path + config.DefaultPage
		if h, ok = mux.m[np]; ok {
			if pathChanged {
				pattern = path
				h = http.RedirectHandler(MakeLocalPath(path), http.StatusMovedPermanently)
			}
			pattern = path
			return
		}
	}

	// find exact match
	if h, ok = mux.m[path]; ok {
		if pathChanged {
			pattern = path
			h = http.RedirectHandler(MakeLocalPath(path), http.StatusMovedPermanently)
		}
		pattern = path
		return
	}

	// If the given path is /tree and its handler is not registered,
	// find /tree/ and redirect if found.
	if path[len(path)-1] != '/' {
		np := path + "/"
		if h, ok = mux.m[np]; ok {
			pattern = np
			h = http.RedirectHandler(MakeLocalPath(np), http.StatusMovedPermanently)
			return
		}
	}

	// If the given path is /tree/ find /tree/{{config.DefaultPage}}. As in /tree/index.html.
	if path[len(path)-1] == '/' && config.DefaultPage != "" {
		np := path + config.DefaultPage
		if h, ok = mux.m[np]; ok {
			if pathChanged {
				pattern = path
				h = http.RedirectHandler(MakeLocalPath(path), http.StatusMovedPermanently)
			}
			pattern = path
			return
		}
	}

	// walk backwards in the path looking for slash ending paths to match
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '/' {
			if h, ok = mux.m[path[:i+1]]; ok {
				pattern = path[:i+1]

				if pathChanged {
					h = http.RedirectHandler(MakeLocalPath(path), http.StatusMovedPermanently)
				}

				return
			}
		}
	}

	// not found
	h = http.NotFoundHandler()
	pattern = ""
	return
}

// cleanPath returns the canonical path for p, eliminating . and .. elements.
func cleanPath(p string) string {
	if p == "" || p == "/" {
		return "/"
	}
	if p[0] != '/' {
		p = "/" + p
	}
	hasTrailingSlash := p[len(p)-1] == '/'
	p = path.Clean(p)
	if hasTrailingSlash {
		p += "/"
	}
	return p
}

func (mux *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI == "*" {
		if r.ProtoAtLeast(1, 1) {
			w.Header().Set("Connection", "close")
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	h, _ := mux.Handler(r)
	h.ServeHTTP(w, r)
}

// UseMuxer serves a muxer such that if a handler cannot be found, or the found handler does not respond,
// control is past to the next handler.
//
// Note that the default Go Muxer is NOT recommended, as it improperly handles
// redirects if this is behind a reverse proxy.
func UseMuxer(mux Muxer, next http.Handler) http.Handler {
	if next == nil {
		panic("next may not be nil. Pass a http.NotFoundHandler if this is the end of the handler chain")
	}
	if mux == nil {
		panic("mux may not be nil")
	}
	fn := func(w http.ResponseWriter, r *http.Request) {

		var h http.Handler
		var p string

		h, p = mux.Handler(r)
		if p == "" {
			// not found, so go to next handler
			next.ServeHTTP(w, r) // skip
		} else {
			// match, so serve normally
			h.ServeHTTP(w, r)
		}
	}
	return http.HandlerFunc(fn)
}
