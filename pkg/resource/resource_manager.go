// package resource manages http resources that you serve based on a static path to the resource.
package resource

import (
	"bytes"
	"context"
	"github.com/goradd/goradd/pkg/pool"
	"net/http"
)

var resourceManager = newResourceManager() // Create a new singleton page manager.


type ResourcePathHandler func(ctx context.Context, buf *bytes.Buffer) (headers map[string]string, err error)

type ResourceManagerI interface {
	RegisterPath(path string, handler ResourcePathHandler)
}

// The ResourceManager is a singleton global that manages the registration and deployment of rest paths. It acts like
// a URL router, returning the ResourcePathHandler that corresponds to a particular URL path.
// init() functions should be created for each path that associates a function to create a rest path,
// with the URL that corresponds to the path.
type ResourceManager struct {
	pathRegistry   map[string]ResourcePathHandler // maps paths to functions that create forms
}

// ResourcePathPrefix is a prefix you can use in front of all goradd rest paths, like a directory path, to indicate that
// this is a goradd rest path.
var ResourcePathPrefix = ""

// GetResourceManager returns the current page manager.
func GetResourceManager() *ResourceManager {
	return resourceManager
}

func newResourceManager() *ResourceManager {
	return &ResourceManager{pathRegistry: make(map[string]ResourcePathHandler)}
}

// RegisterPath associates the given URL path with the given creation function.
// Call this from an init() function. Afterwards, whenever a user navigates to the given path, the
// result of the query will be presented to the user.
func RegisterPath(path string, handler ResourcePathHandler) {
	if _, ok := resourceManager.pathRegistry[path]; ok {
		panic("Page is already registered: " + path)
	}
	resourceManager.pathRegistry[path] = handler
}


func HandleRequest(w http.ResponseWriter, r *http.Request) bool {

	p := r.URL.Path

	handler,ok := resourceManager.getHandler(p)
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
		_,_ = w.Write(buf.Bytes())
	}
	return true
}


func (m *ResourceManager) getHandler(path string) (f ResourcePathHandler, ok bool) {
	f, ok = m.pathRegistry[path]
	return
}


// runHandler processes the resource and writes the response into the buffer. Any special response headers are returned.
func runHandler(ctx context.Context, handler ResourcePathHandler, buf *bytes.Buffer) (headers map[string]string, httpErrCode int) {
	defer func() {
		if r := recover(); r != nil {
			switch v := r.(type) {
			case error:
				httpErrCode = 500
				buf.WriteString(v.Error())
			case string:
				httpErrCode = 500
				buf.WriteString(v)
			case *HttpError:	// A kind of http panic that just returns a response code and headers
				headers = v.headers
				httpErrCode = v.errCode
			default:
				httpErrCode = 500
			}
		}
	}()


	headers, err := handler(ctx, buf)

	if err != nil {
		httpErrCode = 500
		buf.WriteString(err.Error())
		return
	}
	return
}

// HttpError represents an error response to a http request.
type HttpError struct {
	headers map[string] string
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