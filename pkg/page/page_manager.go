package page

import (
	"context"
	"fmt"
	"github.com/goradd/goradd/pkg/config"
	"github.com/goradd/goradd/pkg/log"
	"io"
	"reflect"
)

var pageManager PageManagerI = newPageManager() // Create a new singleton page manager.

type FormCreationFunction func(context.Context) FormI
type formInfo struct {
	formID string
	typ reflect.Type
}

type PageManagerI interface {
	RegisterForm(path string, form FormI, formID string)
	RunPage(ctx context.Context, w io.Writer) (headers map[string]string, httpErrCode int)
	IsPage(path string) bool
}

// The PageManager is a singleton global that manages the registration and deployment of pages. It acts like
// a URL router, returning the page that corresponds to a particular URL path. init() functions should be
// created for each page that associate a function to create a page, with the URL that corresponds to the page,
// and the ID of the page.
type PageManager struct {
	forms map[string]formInfo                      // maps paths to form info
}

// PagePathPrefix is a prefix you can use in front of all goradd pages, like a directory path, to indicate that
// this is a goradd path.
var PagePathPrefix = ""

// GetPageManager returns the current page manager.
func GetPageManager() PageManagerI {
	return pageManager
}

// SetPageManger injects a new page manager. Call this at init time.
func SetPageManager(p PageManagerI)  {
	pageManager = p
}

// newPageManager creates a default page manager that associates url paths with forms.
func newPageManager() *PageManager {
	return &PageManager{forms: make(map[string]formInfo)}
}

// RegisterForm associates the given URL path with the given form and form id and registers it with page manager.
// Call this from an init() function. Afterwards, whenever a user navigates to the given path, the form will be
// created and presented to the user.
func RegisterForm(path string, form FormI, id string) {
	pageManager.RegisterForm(path, form, id)
}

func (m *PageManager)RegisterForm(path string, f FormI, id string) {
	if path == "" {
		panic(`you cannot register the empty path. If you want a default, register just a slash. ("/")`)
	}
	if _, ok := m.forms[path]; ok {
		panic("Form is already registered: " + path)
	}
	if !controlIsRegistered(f) {
		RegisterControl(f) // a form is a control, and needs to be registered for the serializer
	}
	typ := reflect.Indirect(reflect.ValueOf(f)).Type()

	m.forms[path] = formInfo{typ: typ, formID: id}
}

// IsPage returns true if the given path has been registered with the page manager.
func (m *PageManager) IsPage(path string) bool {
	if path == "" {
		path = "/"
	}
	_, ok := m.forms[path]
	return ok
}

// getPage returns a page from the cache, whether it is previously allocated, or
// is a new page. If there is an error creating a new page, it should panic. If this is
// an ajax call and the previous page could not be found in the cache, then return nil in page.
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
			return
		}
		// page was not found, so make a new one
		var form FormI
		path := gCtx.URL.Path
		if path == "" {
			path = "/"
		}
		if info, ok := m.forms[path]; !ok {
			panic("form not found for path: " + gCtx.URL.Path)
		} else {
			form = reflect.New(info.typ).Interface().(FormI)
			form.control().Self = form
			form.Init(ctx, info.formID)
		}
		page = form.Page()
		pageStateId = pageCache.NewPageID()
		page.stateId = pageStateId
		log.FrameworkDebugf("Created page %s", pageStateId)
		isNew = true
	}
	return
}

// RunPage processes the page and writes the response into the buffer. Any special response headers are returned.
func (m *PageManager) RunPage(ctx context.Context, w io.Writer) (headers map[string]string, httpErrCode int) {
	defer func() {
		if r := recover(); r != nil {
			switch v := r.(type) {
			case error:
				err := newRunError(ctx, v)
				m.makeErrorResponse(ctx, err, "", w)
			case string:
				err := newRunError(ctx, v)
				m.makeErrorResponse(ctx, err, "", w)
			case *HttpError: // A kind of http panic that just returns a response code and headers
				headers = v.headers
				httpErrCode = v.errCode
			default:
				err := newRunError(ctx, fmt.Errorf("unknown package error: %v", r))
				m.makeErrorResponse(ctx, err, "", w)
			}
		}
	}()

	page, isNew := m.getPage(ctx)

	if page == nil {
		// An ajax call, but we could not deserialize the old page. Refresh the entire page to get server access.
		io.WriteString(w, `{"loc":"reload"}`) // the refresh will be handled in javascript
		return map[string]string{"Content-Type": "application/json"}, 0
	}

	defer m.cleanup(page)
	page.renderStatus = PageIsRendering
	log.FrameworkDebugf("Page started rendering %s, %s", page.stateId, GetContext(ctx))

	err := page.runPage(ctx, w, isNew)

	if e,ok := err.(FrameworkError); ok {
		if e.Err == FrameworkErrNotAuthorized {
			io.WriteString(w, page.form.GT(e.Error()))
			return nil, 403
		} else if e.Err == FrameworkErrRedirect {
			page.SetResponseHeader("Location", e.Location)
			return page.responseHeader, 303
		}
	}

	if err != nil {
		log.Error(err)
		buf := OutputBuffer(ctx)
		var html = buf.String() // copy current html
		buf.Reset()
		m.makeErrorResponse(ctx, newRunError(ctx, err), html, w)
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
	w io.Writer) {

	if ErrorPageFunc == nil {
		panic("No error page template function is defined")
	}

	// TODO: After GoT changes to writing to io.Writer, change this
	ErrorPageFunc(ctx, html, err, OutputBuffer(ctx))
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

// Redirect aborts the current page load and tells the browser to load a different url. Usually you should
// use Form.ChangeLocation, but you can use this in extreme situations where you really do not want to return
// a goradd page at all and just change locations. For example, if you detect some kind of attempt to hack
// your website, you can use this to redirect to a login page or an error page.
func Redirect(url string) {
	e := HttpError{}
	e.SetResponseHeader("Location", config.MakeLocalPath(url))
	e.Send(303)
}
