package rest2

import (
	"github.com/goradd/goradd/pkg/page"
	"net/http"
	"strings"
)

var restManager = newRestManager() // Create a new singleton REST manager.

type RestManagerI interface {
	RegisterPath(path string, handler http.HandlerFunc)
}

// The RestManager is a singleton global that manages the registration and deployment of rest paths. It acts like
// a URL router, returning the RestPathHandler that corresponds to a particular URL path.
// init() functions should be created for each path that associates a function to create a rest path,
// with the URL that corresponds to the path.
type RestManager struct {
	pathRegistry map[string]http.HandlerFunc // maps paths to REST functions
}

// RestPathPrefix is a prefix you can use in front of all goradd rest paths, like a directory path, to indicate that
// this is a goradd REST path. You can normally leave this blank if you are only implementing a REST API.
// If you are implementing different APIs, specify the path here to indicate that the user is making
// a REST call.
var RestPathPrefix = ""

// GetRestManager returns the current page manager.
func GetRestManager() *RestManager {
	return restManager
}

func newRestManager() *RestManager {
	return &RestManager{pathRegistry: make(map[string]http.HandlerFunc)}
}

// RegisterPath associates the given URL path with the given handler.
// Call this from an init() function. Afterwards, whenever a user navigates to the given path, the
// result of the query will be presented to the user.
// TODO: Substitute a standard MUX for this, and that way other MUX's, like the Gorilla Mux, could
// be used to enhance REST management
func RegisterPath(path string, handler http.HandlerFunc) {
	if path != "" && path[0] == '/' {
		path = path[1:]
	}
	if _, ok := restManager.pathRegistry[path]; ok {
		panic("Page is already registered: " + path)
	}
	restManager.pathRegistry[path] = handler
}

// HandleRequest is the application request handler for a REST request. It returns true if the
// request was for a recognized resource, and false if not.
func HandleRequest(w http.ResponseWriter, r *http.Request) bool {
	p := r.URL.Path
	pathItems := strings.Split(p, "/")
	if len(pathItems) == 0 {
		return false
	}

	handler := restManager.getHandler(pathItems[0])
	if handler == nil {
		return false
	}

	runHandler(w, r, handler)
	return true
}

func (m *RestManager) getHandler(path string) (f http.HandlerFunc) {
	if RestPathPrefix != "" {
		if strings.Index(path, RestPathPrefix) == 0 { // starts with prefix
			path = path[len(RestPathPrefix):] // remove prefix from path
		} else {
			return // not found in path
		}
	}

	f = m.pathRegistry[path]
	return
}

// runHandler processes the REST handler and writes the response into the writer.
// Handlers are responsible for setting the headers correctly and writing into the response writer in
// general, though this routine will also trap any panics and return an appropriate error code.
func runHandler(w http.ResponseWriter, r *http.Request, handler http.HandlerFunc) {
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
			default:
				w.WriteHeader(500)
			}
		}
	}()

	handler(w,r)
	return
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

