package page

import (
	"bytes"
	"context"
	"fmt"
	"github.com/spekary/goradd/log"
	"goradd/config"
	"strings"
)

var pageManager *PageManager

type FormCreationFunction func(context.Context) FormI

type PageManagerI interface {
	RegisterPage(path string, creationFunction FormCreationFunction)
}

type PageManager struct {
	pathRegistry   map[string]FormCreationFunction // maps paths to functions that create forms
	formIdRegistry map[string]FormCreationFunction // maps form ids to functions that create forms
}

func GetPageManager() *PageManager {
	return pageManager
}

func NewPageManager() *PageManager {
	return &PageManager{pathRegistry: make(map[string]FormCreationFunction), formIdRegistry: make(map[string]FormCreationFunction)}
}

func RegisterPage(path string, creationFunction FormCreationFunction, formId string) {
	if pageManager == nil {
		pageManager = NewPageManager() // Create a new singleton page manager
	}
	if _, ok := pageManager.pathRegistry[path]; ok {
		panic("Page is already registered: " + path)
	}
	pageManager.pathRegistry[path] = creationFunction
	pageManager.formIdRegistry[formId] = creationFunction
}

func (m *PageManager) getNewPageFunc(ctx context.Context) (f FormCreationFunction, path string, ok bool) {
	path = GetContext(ctx).URL.Path
	prefix := config.PAGE_PATH_PREFIX
	if prefix != "" {
		if strings.Index(path, prefix) == 0 { // starts with prefix
			path = path[len(prefix):] // remove path
		} else {
			return // not found in path
		}
	}
	f, ok = m.pathRegistry[path]
	return
}

func (m *PageManager) IsPage(ctx context.Context) bool {
	_, _, ok := m.getNewPageFunc(ctx)
	return ok
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
			// TODO: If this happens, we need to reload the whole page, because we lost the formstate completely
			log.FrameworkDebug("Ajax lost the page state") // generally this should only happen if the page state drops out of the cash, which might happen after a long time
		}
		// page was not found, so make a new one
		f, _, _ := m.getNewPageFunc(ctx)
		if f == nil {
			panic("Could not find the page creation function")
		}
		page = f(ctx).Page() // call the page create function and get the page
		pageStateId = pageCache.NewPageID()
		page.GetPageBase().stateId = pageStateId
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
				grCtx := GetContext(ctx)
				grCtx.WasHandled = true // Notify listeners that the app handled the page
				headers = v.headers
				httpErrCode = v.errCode
			default:
				err := newRunError(ctx, fmt.Errorf("Unknown package error: %v", r))
				m.makeErrorResponse(ctx, err, "", buf)
			}
		}
	}()

	pageI, isNew := m.getPage(ctx)
	page := pageI.GetPageBase()
	defer m.cleanup(pageI)
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

func (m *PageManager) IsAsset(ctx context.Context) bool {
	return assetIsRegistered(GetContext(ctx).HttpContext.URL.Path)
}

type HttpError struct {
	headers map[string] string
	errCode int
}

func (e *HttpError) SetResponseHeader(key, value string) {
	if e.headers == nil {
		e.headers = map[string]string{key: value}
	} else {
		e.headers[key] = value
	}
}

func (e *HttpError) Send(errCode int) {
	e.errCode = errCode
	panic(e)
}

// Redirect aborts the current page load and tells the browser to load a different page
func Redirect(url string) {
	e := HttpError{}
	e.SetResponseHeader("Location", url)
	e.Send(307)
}