package api

import (
	"github.com/goradd/goradd/pkg/config"
	"github.com/goradd/goradd/pkg/page"
	"net/http"
	"net/url"
	"strings"
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

func init() {
	// Create a default ApiManager if one does not already exist
	// The default can be over-ridden in the project's config package, since those
	// init functions will be called before this one
	if config.ApiManager == nil && config.ApiPrefix != "" {
		config.ApiManager = NewDefaultApiManager()
	}
}

// The DefaultApiManager type represents the default api manager that is created at bootup.
// The Mux variable can be replaced with the Muxer of your choice.
type DefaultApiManager struct {
	Mux Muxer
}

// Creates a new api manager that uses config.ApiPrefix as a path prefix to determine
// if the given call is a REST call. It also uses Go's standard ServeMux as a muxer to
// direct which handler should handle the request. The Mux is public, so you can replace
// it with your own Muxer if desired.
func NewDefaultApiManager() DefaultApiManager {
	return DefaultApiManager{
		Mux: http.NewServeMux(),
	}
}

func (a DefaultApiManager) RegisterPattern(pattern string, handler http.Handler) {
	a.Mux.Handle(pattern, handler)
}

func (a DefaultApiManager) HandleRequest(w http.ResponseWriter, r *http.Request) (isApiCall bool) {
	// Strip out the api prefix so that handlers don't need to know it,
	// and we can use it to determine if this is an API request.
	if strings.HasPrefix(r.URL.Path, config.ApiPrefix) {
		isApiCall = true
		p := strings.TrimPrefix(r.URL.Path, config.ApiPrefix)
		r2 := new(http.Request)
		*r2 = *r // shallow copy
		r2.URL = new(url.URL)
		*r2.URL = *r.URL
		r2.URL.Path = p

		// Trap all panics and serve up an http error when that happens so that the application
		// keeps running.
		defer func() {
			if r := recover(); r != nil {
				switch v := r.(type) {
				case page.FrameworkError:
					w.WriteHeader(v.HttpError())
					w.Write([]byte(v.Error()))
				case *HttpError: // A kind of http panic that just returns a response code and headers
					header := w.Header()
					for k,v := range v.headers {
						header.Set(k,v)
					}
					w.WriteHeader(v.errCode)
				case error:
					w.WriteHeader(500)
					w.Write([]byte(v.Error()))
				case string:
					w.WriteHeader(500)
					w.Write([]byte(v))
				case int:
					w.WriteHeader(v)
				default:
					w.WriteHeader(500)
				}
				isApiCall = true
			}
		}()
		a.Mux.ServeHTTP(w, r2)
	}
	return
}

// RegisterPattern associates the given URL path with the given handler.
// Call this from an init() function. Whenever a user navigates to the given pattern, the
// handler will be called. Note that the pattern should NOT include the RestPathPrefix.
func RegisterPattern(pattern string, handler http.HandlerFunc) {
	if config.ApiManager == nil {
		panic ("the ApiManager has not been initialized")
	}
	config.ApiManager.RegisterPattern(pattern, handler)
}

// HttpError represents an error response to a http request.
type HttpError struct {
	headers map[string]string
	errCode int
}

// SetResponseHeader sets a key-value in the header response.
func (e *HttpError) SetResponseHeader(key, value string) {
	if e.headers == nil {
		e.headers = map[string]string{key: value}
	} else {
		e.headers[key] = value
	}
}

// Send will cause the page to error with the given http error code.
func (e *HttpError) Send(errCode int) {
	e.errCode = errCode
	panic(e)
}

