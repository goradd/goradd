package http

import (
	"net/http"
)

// Muxer represents the typical functions available in a mux and allows you to replace the default
// Golang muxer here with a 3rd party mux, like the Gorilla mux.
type Muxer interface {
	// Handle associates a handler with the given pattern in the url path
	Handle(pattern string, handler http.Handler)
	// ServeHTTP sends a request to the MUX, to be forwarded on to the registered handler,
	// or responded with an unknown resource error.
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

// AppMuxer is the application muxer that lets you do traditional http handling
// from behind the application facilities of session management, output buffering,
// etc.
//
//It is called from the default MakeAppServer implementation.
var AppMuxer = http.NewServeMux()

// UseAppMuxer is called by the framework at application startup to place the
// application muxer in the handler stack.
func UseAppMuxer(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		h,p := AppMuxer.Handler(r)
		if p != "" {
			// handler was found
			h.ServeHTTP(w,r)
		} else {
			next.ServeHTTP(w,r) // skip
		}
	}
	return http.HandlerFunc(fn)
}

// RegisterAppMuxerHandler registers a handler for the given directory prefix.
//
// The handler will be called with the prefix stripped away. When the prefix is
// stripped, a rooted path will be passed along. In other words, if the path
// is /api/file, the called handler will receive /file.
//
// Note that you CAN register a handler for the root directory.
func RegisterAppMuxerHandler(prefix string, handler http.Handler) {
	if prefix == "" {
		prefix = "/"
	} else {
		if prefix[0] != '/' {
			prefix = "/" + prefix
		}
		// Make sure the registered prefix ends with a /
		if prefix[len(prefix) - 1] != '/' {
			prefix = prefix + "/"
		}
	}
	// Here we register the handler with a closing / so that similar names will not be confused,
	// but we do not strip the last / from the file name passed on.
	AppMuxer.Handle(prefix, http.StripPrefix(prefix[0:len(prefix) - 1],handler))
}

// ErrorHandler wraps the given handler in a default HTTP error handler that
// will respond appropriately to any panics that happen within the given handler.
//
// Panic with an http.Error value to get a specific kind of http error to
// be output.
func ErrorHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				switch v := r.(type) {
				case Error: // A kind of http panic that just returns a response code and headers
					header := w.Header()
					for k,v := range v.headers {
						header.Set(k,v)
					}
					w.WriteHeader(v.errCode)
					if v.message != "" {
						_,_ = w.Write([]byte(v.message))
					}
				case error:
					w.WriteHeader(http.StatusInternalServerError)
					_,_ =w.Write([]byte(v.Error()))
				case string:
					w.WriteHeader(http.StatusInternalServerError)
					_,_ =w.Write([]byte(v))
				case int:
					w.WriteHeader(v)
				default:
					w.WriteHeader(http.StatusInternalServerError)
				}
			}
		}()
		h.ServeHTTP(w, r)
	})
}
