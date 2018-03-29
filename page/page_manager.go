package page

import (
	"context"
	"goradd/config"
	"strings"
	"fmt"
	"html/template"
	"bytes"
	"reflect"
)

var pageManager *PageManager
var errorTemplate *template.Template

type PageCreationFunc func() PageI


type PageManagerI interface {
	RegisterPage(path string, creationFunction PageCreationFunc)
}

type PageManager struct {
	pathRegistry map[string] PageCreationFunc // maps paths to functions that create pages
	typeRegistry map[string] PageCreationFunc // maps object types to functions that create pages
}


func GetPageManager() *PageManager {
	return pageManager
}

func NewPageManager() *PageManager {
	return &PageManager{pathRegistry:make(map[string] PageCreationFunc), typeRegistry:make(map[string] PageCreationFunc)}
}

func RegisterPage(path string, creationFunction PageCreationFunc) {
	if pageManager == nil {
		pageManager = NewPageManager()	// Create a new singleton page manager
	}
	if _,ok := pageManager.pathRegistry[path]; ok {
		panic ("Page is already registered: " + path)
	}
	pageManager.pathRegistry[path] = creationFunction

	p := creationFunction()
	pageManager.typeRegistry[reflect.TypeOf(p).String()] = creationFunction
}

func (m *PageManager) getNewPageFunc (ctx context.Context) (f PageCreationFunc, path string, ok bool) {
	path = GetContext(ctx).URL.Path
	prefix := config.PAGE_PATH_PREFIX
	if prefix != "" {
		if strings.Index(path, prefix) == 0 { // starts with prefix
			path = path[len(prefix):] // remove path
		} else {
			return // not found in path
		}
	}
	f,ok = m.pathRegistry[path]
	return
}

func (m *PageManager) IsPage (ctx context.Context) bool {
	_,_,ok := m.getNewPageFunc(ctx)
	return ok
}

func (m *PageManager) getPage(ctx context.Context) (page PageI, isNew bool) {
	var pageStateId string

	gCtx := GetContext(ctx)

	pageStateId = gCtx.pageStateId

	if pageStateId != "" {
		page = pageCache.Get(pageStateId)
	}

	if page == nil {
		// page was not found, so make a new one
		f,path, _ := m.getNewPageFunc(ctx)
		if f == nil {
			panic("Could not find the page creation function")
		}
		page = f()	// call the page create function
		page.GetPageBase().Self = page
		page.Init(ctx, path)
		pageStateId = pageCache.NewPageId()
		page.GetPageBase().stateId = pageStateId
		//pageCache.Set(pageStateId, page)
		isNew = true
	} else {
		page.GetPageBase().Self = page
		page.Restore()
	}
	return
}

// RunPage processes the page and writes the response into the buffer. Any special response headers are returned.
func (m *PageManager) RunPage(ctx context.Context, buf *bytes.Buffer) map[string]string {
	defer func() {
		if r := recover(); r != nil {
			switch v := r.(type) {
			case error:
				err := newRunError(ctx, v)
				m.makeErrorResponse(ctx, err, "", buf)
			case string:
				err := newRunError(ctx, v)
				m.makeErrorResponse(ctx, err, "", buf)
			default:
				err := newRunError(ctx,fmt.Errorf("Unknown package error: %v", r))
				m.makeErrorResponse(ctx, err, "", buf)
			}
		}
	}()

	pageI, isNew := m.getPage(ctx)
	page := pageI.GetPageBase()
	defer m.cleanup(pageI)
	err := page.runPage(ctx, buf, isNew)

	if err != nil {
		var html = buf.String() // copy current html
		buf.Reset()
		m.makeErrorResponse(ctx, newRunError(ctx, err), html, buf)
		return nil
	}
	return page.responseHeader
}

func (m *PageManager) cleanup (p PageI) {


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

func (m *PageManager) IsAsset (ctx context.Context) bool {
	return assetIsRegistered(GetContext(ctx).HttpContext.URL.Path)
}



