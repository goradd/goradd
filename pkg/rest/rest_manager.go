package rest

import (
	"bytes"
	"context"
	"github.com/goradd/goradd/pkg/pool"
	"net/http"
	"strings"
)

var restManager = newRestManager() // Create a new singleton page manager.

type RestPathHandler func(ctx context.Context, buf *bytes.Buffer) error

type RestManagerI interface {
	RegisterPath(path string, creationFunction RestPathHandler)
}

// The RestManager is a singleton global that manages the registration and deployment of rest paths. It acts like
// a URL router, returning the RestPathHandler that corresponds to a particular URL path.
// init() functions should be created for each path that associates a function to create a rest path,
// with the URL that corresponds to the path.
type RestManager struct {
	pathRegistry map[string]RestPathHandler // maps paths to functions that create forms
}

// RestPathPrefix is a prefix you can use in front of all goradd rest paths, like a directory path, to indicate that
// this is a goradd rest path.
var RestPathPrefix = ""

// GetRestManager returns the current page manager.
func GetRestManager() *RestManager {
	return restManager
}

func newRestManager() *RestManager {
	return &RestManager{pathRegistry: make(map[string]RestPathHandler)}
}

// RegisterPath associates the given URL path with the given creation function.
// Call this from an init() function. Afterwards, whenever a user navigates to the given path, the
// result of the query will be presented to the user.
func RegisterPath(path string, handler RestPathHandler) {
	if _, ok := restManager.pathRegistry[path]; ok {
		panic("Page is already registered: " + path)
	}
	restManager.pathRegistry[path] = handler
}

func HandleRequest(w http.ResponseWriter, r *http.Request) bool {

	p := r.URL.Path
	pathItems := strings.Split(p, "/")
	if len(pathItems) == 0 {
		return false
	}

	handler, ok := restManager.getHandler(pathItems[0])
	if !ok {
		return false
	}

	ctx := r.Context()
	buf := pool.GetBuffer()
	defer pool.PutBuffer(buf)
	headers, errCode := runHandler(ctx, handler, buf)
	if headers != nil {
		for k, v := range headers {
			// Multi-value headers can simply be separated with commas I believe
			w.Header().Set(k, v)
		}
	}
	if errCode != 0 {
		w.WriteHeader(errCode)
	} else {
		_, _ = w.Write(buf.Bytes())
	}
	return true
}

func (m *RestManager) getHandler(path string) (f RestPathHandler, ok bool) {
	if RestPathPrefix != "" {
		if strings.Index(path, RestPathPrefix) == 0 { // starts with prefix
			path = path[len(RestPathPrefix):] // remove prefix from path
		} else {
			return // not found in path
		}
	}
	f, ok = m.pathRegistry[path]
	return
}

// RunPage processes the page and writes the response into the buffer. Any special response headers are returned.
func runHandler(ctx context.Context, handler RestPathHandler, buf *bytes.Buffer) (headers map[string]string, httpErrCode int) {
	defer func() {
		if r := recover(); r != nil {
			switch v := r.(type) {
			case error:
				httpErrCode = 500
				buf.WriteString(v.Error())
			case string:
				httpErrCode = 500
				buf.WriteString(v)
			case *HttpError: // A kind of http panic that just returns a response code and headers
				headers = v.headers
				httpErrCode = v.errCode
			default:
				httpErrCode = 500
			}
		}
	}()

	err := handler(ctx, buf)

	if err != nil {
		httpErrCode = 500
		buf.WriteString(err.Error())
		return
	}
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

// Redirect aborts the current page load and tells the browser to load a different url.
func Redirect(url string) {
	e := HttpError{}
	e.SetResponseHeader("Location", url)
	e.Send(307)
}
