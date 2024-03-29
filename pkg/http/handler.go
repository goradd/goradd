package http

import (
	"context"
	"io"
	"mime"
	"net/http"
	"path"
	"path/filepath"
)

// The following two maps collect handler registration during Go's init process. These are
// then registered to the application muxers when the application starts up. This makes it
// possible for parts of the app to turn themselves on just by being imported.

type handlerMap map[string]http.Handler

// patternHandlers are the handlers that are processed immediately based on the static path.
var patternHandlers = make(handlerMap)

// appHandlers are the handlers processed at the end of the application handler stack,
// so these handlers go through session management, authentication, etc. They are also
// associated with a path.
var appHandlers = make(handlerMap)

// PatternMuxer is the muxer that immediately routes handlers based on the path without
// going through the application handlers. It is automatically loaded during app startup.
var PatternMuxer Muxer

// AppMuxer is the application muxer that lets you do traditional http handling
// from behind the application facilities of session management, output buffering,
// etc. It is automatically loaded during app startup.
var AppMuxer Muxer

// UsePatternMuxer is called by the framework at application startup to place the
// pattern mux in the handler stack.
//
// All previously registered pattern handlers will be put in the given muxer. The muxer will
// be remembered so that future registrations will go to that muxer.
//
// next specifies a handler that will be used if the muxer processes a URL that
// it does not recognize.
func UsePatternMuxer(mux Muxer, next http.Handler) http.Handler {
	if patternHandlers == nil {
		panic("the PatternMuxer has already been registered")
	}
	PatternMuxer = mux
	return useMuxer(mux, next, &patternHandlers)
}

// UseAppMuxer is called by the framework at application startup to place the
// application muxer at the end of the application handler stack
//
// next specifies a handler that will be used if the AppMuxer is presented a URL that
// it does not recognize.
func UseAppMuxer(mux Muxer, next http.Handler) http.Handler {
	if appHandlers == nil {
		panic("the AppMuxer has already been registered")
	}
	AppMuxer = mux
	return useMuxer(mux, next, &appHandlers)
}

// RegisterHandler registers a handler for the given pattern.
//
// Use this when registering a handler to a specific path. Use RegisterPrefixHandler if registering
// a handler for a whole subdirectory of a path.
//
// The given handler is served immediately by the application without going through the application
// handler stack. If you need session management, HSTS protection, authentication, etc., use
// RegisterAppHandler.
//
// You may call this from an init() function.
func RegisterHandler(pattern string, handler http.Handler) {
	registerHandler(pattern, handler, patternHandlers, PatternMuxer)
}

// RegisterAppHandler registers a handler for the given pattern.
//
// Use this when registering a handler to a specific path. Use RegisterAppPrefixHandler if registering
// a handler for a whole subdirectory of a path.
//
// The given handler is served near the end of the application handler stack, so
// you will have access to session management and any other middleware handlers
// in the application stack.
//
// You may call this from an init() function.
func RegisterAppHandler(pattern string, handler http.Handler) {
	registerHandler(pattern, handler, appHandlers, AppMuxer)
}

// RegisterPrefixHandler registers a handler for the given directory prefix.
//
// The handler will be called immediately based on the path and will not be sent
// through the application handler middleware stack. Use RegisterAppPrefixHandler
// for the equivalent function processed at the end of the application handler stack.
//
// The handler will be called with the prefix stripped away. When the prefix is
// stripped, a rooted path will be passed along. In other words, if the
// prefix is /api, and the path being served
// is /api/file, the called handler will receive /file.
//
// Note that you CAN register a handler for the root directory.
//
// If the handler is presented a URL that it does not recognize, it should
// return an HTTP error to the ResponseWriter.
//
// You may call this from an init() function.
func RegisterPrefixHandler(prefix string, handler http.Handler) {
	registerPrefixHandler(prefix, handler, patternHandlers, PatternMuxer)
}

// RegisterAppPrefixHandler registers a handler for the given directory prefix.
//
// The handler will be called at the end of the application handler middleware stack.
//
// The handler will be called with the prefix stripped away. When the prefix is
// stripped, a rooted path will be passed along. In other words, if the
// prefix is /api, and the path being served
// is /api/file, the called handler will receive /file.
//
// Note that you CAN register a handler for the root directory.
//
// If the handler is presented a URL that it does not recognize, it should
// return an HTTP error to the ResponseWriter.
//
// You may call this from an init() function.
func RegisterAppPrefixHandler(prefix string, handler http.Handler) {
	registerPrefixHandler(prefix, handler, appHandlers, AppMuxer)
}

// A DrawFunc sends output to the Writer. goradd uses this signature in its template functions.
type DrawFunc func(ctx context.Context, w io.Writer) (err error)

// RegisterDrawFunc registers an output function for the given pattern.
//
// This could be used to register template output with a path, for example. See the renderResource
// template macro and the configure.tpl.got file in the welcome application for an example.
//
// The file name extension will be used first to determine the Content-Type. If that fails, then
// the content will be inspected to determine the Content-Type.
//
// Registered handlers are served by the AppMuxer.
func RegisterDrawFunc(pattern string, f DrawFunc) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		name := path.Base(r.URL.Path)

		// set the content type by extension first
		ctype := mime.TypeByExtension(filepath.Ext(name))
		if ctype != "" {
			w.Header().Set("Content-Type", ctype)
		}

		// the write process below will fill in the content-type if not set
		err := f(ctx, w)
		if err != nil {
			panic(err)
		}
	}
	h := http.HandlerFunc(fn)
	RegisterAppHandler(pattern, h)
}

func registerHandler(pattern string, handler http.Handler, m handlerMap, mux Muxer) {
	if m == nil {
		// The muxer has already been recorded, so register the pattern directly to the muxer
		mux.Handle(pattern, handler)
	} else {
		// the muxer has not yet been used, so cache the pattern in anticipation of the muxer being recorded

		if _, ok := m[pattern]; ok {
			panic("the handler for " + pattern + " is already registered")
		}
		m[pattern] = handler
	}
}

func registerPrefixHandler(prefix string, handler http.Handler, m handlerMap, mux Muxer) {
	if prefix == "" {
		prefix = "/"
	}

	if prefix != "/" {
		if prefix[0] != '/' {
			prefix = "/" + prefix
		}
		// Add a trailing slash if its not there
		if prefix[len(prefix)-1] != '/' {
			prefix = prefix + "/"
		}
	}
	registerHandler(prefix, http.StripPrefix(prefix[:len(prefix)-1], handler), m, mux)
}

func useMuxer(mux Muxer, next http.Handler, m *handlerMap) http.Handler {
	for p, h := range *m {
		mux.Handle(p, h)
	}
	*m = nil
	return UseMuxer(mux, next)
}
