package page

import (
	"bytes"
	"context"
	"fmt"
	"github.com/goradd/goradd/pkg/log"
	"strings"
)

var pageManager = newPageManager() // Create a new singleton page manager.


type FormCreationFunction func(context.Context) FormI

type PageManagerI interface {
	RegisterPage(path string, creationFunction FormCreationFunction)
}

// The PageManager is a singleton global that manages the registration and deployment of pages. It acts like
// a URL router, returning the page that corresponds to a particular URL path. init() functions should be
// created for each page that associate a function to create a page, with the URL that corresponds to the page,
// and the ID of the page.
type PageManager struct {
	pathRegistry   map[string]FormCreationFunction // maps paths to functions that create forms
	formIdRegistry map[string]FormCreationFunction // maps form ids to functions that create forms
}

// PagePathPrefix is a prefix you can use in front of all goradd pages, like a directory path, to indicate that
// this is a goradd path.
var PagePathPrefix = ""

// GetPageManager returns the current page manager.
func GetPageManager() *PageManager {
	return pageManager
}

func newPageManager() *PageManager {
	return &PageManager{pathRegistry: make(map[string]FormCreationFunction), formIdRegistry: make(map[string]FormCreationFunction)}
}

// RegisterPath associates the given URL path with the given form creation function and form id and registers it with page manager.
// Call this from an init() function. Afterwards, whenever a user navigates to the given path, the form will be
// created and presented to the user.
func RegisterPage(path string, creationFunction FormCreationFunction, formId string) {
	if _, ok := pageManager.pathRegistry[path]; ok {
		panic("Page is already registered: " + path)
	}
	pageManager.pathRegistry[path] = creationFunction
	pageManager.formIdRegistry[formId] = creationFunction
}

func (m *PageManager) getNewPageFunc(path string) (f FormCreationFunction, ok bool) {
	if PagePathPrefix != "" {
		if strings.Index(path, PagePathPrefix) == 0 { // starts with prefix
			path = path[len(PagePathPrefix):] // remove prefix from path
		} else {
			return // not found in path
		}
	}
	f, ok = m.pathRegistry[path]
	return
}

// IsPage returns true if the given path has been registered with the page manager.
func (m *PageManager) IsPage(path string) bool {
	_, ok := m.getNewPageFunc(path)
	return ok
}

// HasPage returns true if the given page state is currently in the page cache. This indicates that a user
// has recently accessed a page with the given state id. You can use this to validate client interactions.
func (m *PageManager) HasPage(pageStateId string) bool {
	return pageCache.Has(pageStateId)
}


func (m *PageManager) getPage(ctx context.Context) (page *Page, isNew bool) {
	var pageStateId string

	gCtx := GetContext(ctx)

	pageStateId = gCtx.pageStateId

	if pageStateId != "" {
		page = pageCache.Get(pageStateId)
	}

	if page == nil {
		if gCtx.requestMode == Ajax {
			// TODO: If this happens, we need to reload the whole page, because we lost the pagestate completely
			log.FrameworkDebug("Ajax lost the page state") // generally this should only happen if the page state drops out of the cache, which might happen after a long time
		}
		// page was not found, so make a new one
		f, _ := m.getNewPageFunc(gCtx.URL.Path)
		if f == nil {
			panic("Could not find the page creation function")
		}
		page = f(ctx).Page() // call the page create function and get the page
		pageStateId = pageCache.NewPageID()
		page.stateId = pageStateId
		log.FrameworkDebugf("Created page %s", pageStateId)
		//pageCache.Set(pageStateId, page)
		isNew = true
	} else {
		//page.Restore() // TODO: Only restore if we were deserealized. Tricky to detect.
	}
	return
}

// RunPage processes the page and writes the response into the buffer. Any special response headers are returned.
func (m *PageManager) RunPage(ctx context.Context, buf *bytes.Buffer) (headers map[string]string, httpErrCode int) {
	defer func() {
		if r := recover(); r != nil {
			switch v := r.(type) {
			case error:
				err := newRunError(ctx, v)
				m.makeErrorResponse(ctx, err, "", buf)
			case string:
				err := newRunError(ctx, v)
				m.makeErrorResponse(ctx, err, "", buf)
			case *HttpError:	// A kind of http panic that just returns a response code and headers
				headers = v.headers
				httpErrCode = v.errCode
			default:
				err := newRunError(ctx, fmt.Errorf("unknown package error: %v", r))
				m.makeErrorResponse(ctx, err, "", buf)
			}
		}
	}()

	page, isNew := m.getPage(ctx)

	defer m.cleanup(page)
	page.renderStatus = PageIsRendering
	log.FrameworkDebugf("Page started rendering %s, %s", page.stateId, GetContext(ctx))

	err := page.runPage(ctx, buf, isNew)

	if err != nil {
		log.Error(err)
		var html = buf.String() // copy current html
		buf.Reset()
		m.makeErrorResponse(ctx, newRunError(ctx, err), html, buf)
		return
	}
	return page.responseHeader, page.responseError
}

func (m *PageManager) cleanup(p *Page) {
	p.renderStatus = PageIsNotRendering
	log.FrameworkDebugf("Page stopped rendering %s", p.stateId)
}

func (m *PageManager) makeErrorResponse(ctx context.Context,
	err *Error,
	html string,
	buf *bytes.Buffer) {

	if ErrorPageFunc == nil {
		panic("No error page template function is defined")
	}

	ErrorPageFunc(ctx, html, err, buf)
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

// Redirect aborts the current page load and tells the browser to load a different url. Usually you should
// use Form.ChangeLocation, but you can use this in extreme situations where you really do not want to return
// a goradd page at all and just change locations. For example, if you detect some kind of attempt to hack
// your website, you can use this to redirect to a login page or an error page.
func Redirect(url string) {
	e := HttpError{}
	e.SetResponseHeader("Location", url)
	e.Send(307)
}