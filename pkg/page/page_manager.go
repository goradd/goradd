package page

import (
	"context"
	"fmt"
	"github.com/goradd/goradd/pkg/http"
	"github.com/goradd/goradd/pkg/log"
	"io"
	http2 "net/http"
	"reflect"
)

var pageManager PageManagerI = newPageManager() // Create a new singleton page manager.

var HtmlErrorMessage =
	`<h1 id="err-title">Error</h1>
<p>
An unexpected error has occurred and your request could not be processed. The error has been logged and we will
attempt to fix the problem as soon as possible.
</p>`

type FormCreationFunction func(context.Context) FormI
type formInfo struct {
	formID string
	typ reflect.Type
}

type PageManagerI interface {
	RegisterForm(path string, form FormI, formID string)
	RunPage(ctx context.Context, w http2.ResponseWriter, r *http2.Request) ()
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
func (m *PageManager) RunPage(ctx context.Context, w http2.ResponseWriter, req *http2.Request) () {
	defer func() {
		if r := recover(); r != nil {
			var msg string
			switch v := r.(type) {
			case error:
				msg = v.Error()
			case string:
				msg = v
			case http.Error: // A kind of http panic that just returns a response code and headers
				panic(v) // send it up the panic chain
			case int: // Just an http error code
				panic (v)
			default:
				msg = fmt.Sprintf("%v", v)
			}
			// TODO: Need to translate the error message, but we don't have a page context to do it.

			// pass the error on to the error handler above
			panic (http.NewServerError(msg, req, 2, HtmlErrorMessage))
		}
	}()

	page, isNew := m.getPage(ctx)

	if page == nil {
		// An ajax call, but we could not deserialize the old page. Refresh the entire page to get server access.
		w.Header().Set("Content-Type", "application/json")
		_,_ = io.WriteString(w, `{"loc":"reload"}`) // the refresh will be handled in javascript
		return
	}

	defer m.cleanup(page)
	page.renderStatus = PageIsRendering
	log.FrameworkDebugf("Page started rendering %s, %s", page.stateId, GetContext(ctx))

	err := page.runPage(ctx, w, isNew)
	if err != nil {
		// TODO: remove this. All errors should panic in place so we can know where the problem is
		panic (err)
	}
}

func (m *PageManager) cleanup(p *Page) {
	p.renderStatus = PageIsNotRendering
	log.FrameworkDebugf("Page stopped rendering %s", p.stateId)
}

func (m *PageManager) makeErrorResponse(ctx context.Context,
	err *Error,
	html []byte,
	w io.Writer) {

	if ErrorPageFunc == nil {
		panic("No error page template function is defined")
	}

	ErrorPageFunc(ctx, html, err, w)
}

