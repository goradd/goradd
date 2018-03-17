package page

import (
	"context"
	"goradd/config"
	"strings"
	"fmt"
	"html/template"
	"bytes"
)

var pageManager *PageManager
var errorTemplate *template.Template

type PageCreationFunc func(ctx context.Context) PageI


type PageManagerI interface {
	RegisterPage(path string, creationFunction PageCreationFunc)
}

type PageManager struct {
	cache *pageCache
	registry map[string] PageCreationFunc // maps paths to functions that create pages
	renderState int
	ErrorTemplate string
}


func GetPageManager() *PageManager {
	return pageManager
}

func NewPageManager() *PageManager {

	return &PageManager{cache: NewPageCache(), registry:make(map[string] PageCreationFunc), ErrorTemplate:"error.tpl.html"}
}

func RegisterPage(path string, creationFunction PageCreationFunc) {
	if pageManager == nil {
		pageManager = NewPageManager()	// Create a new singleton page manager
	}
	if _,ok := pageManager.registry[path]; ok {
		panic ("Page is already registered: " + path)
	}
	pageManager.registry[path] = creationFunction
}

func (m *PageManager) getNewPageFunc (ctx context.Context) (f PageCreationFunc, ok bool) {
	path := GetContext(ctx).URL.Path
	prefix := config.PAGE_PATH_PREFIX
	if prefix != "" {
		if strings.Index(path, prefix) == 0 { // starts with prefix
			path = path[len(prefix):] // remove path
		} else {
			return // not found in path
		}
	}
	f,ok = m.registry[path]
	return
}

func (m *PageManager) IsPage (ctx context.Context) bool {
	_,ok := m.getNewPageFunc(ctx)
	return ok
}

func (m *PageManager) getPage(ctx context.Context) (page PageI, isNew bool) {
	var pageStateId string

	gCtx := GetContext(ctx)

	pageStateId = gCtx.pageStateId

	if pageStateId != "" {
		page = m.cache.Get(pageStateId)
	}

	if page == nil {
		// page was not found, so make a new one
		f,_ := m.getNewPageFunc(ctx)
		if f == nil {
			panic("Could not find the page creation function")
		}
		page = f(ctx)	// call the page create function
		pageStateId = m.cache.NewPageId()
		page.GetPageBase().stateId = pageStateId
		m.cache.Set(pageStateId, page)
		isNew = true
	}
	return
}


func (m *PageManager) RunPage(ctx context.Context, buf *bytes.Buffer) {
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

	grCtx := GetContext(ctx)

	if grCtx.err != nil {
		panic(grCtx.err)	// If we received an error during the unpacking process, let the deferred code above handle the error.
	}
	pageI, isNew := m.getPage(ctx)

	page := pageI.GetPageBase()

	defer m.cleanup(pageI)

	page.renderStatus = UNRENDERED
	pageI.Run()

	if !isNew {
		pageI.Form().control().updateValues(grCtx)	// Tell all the controls to update their values.
		// if this is an event response, do the actions associated with the event
		if c := pageI.GetControl(grCtx.actionControlId); c != nil {
			c.control().doAction(ctx)
		}
	}

	//if server {

	err := pageI.Draw(ctx, buf)

	if err != nil {
		var html = buf.String() // copy current html
		buf.Reset()
		m.makeErrorResponse(ctx, newRunError(ctx, err), html, buf)
	}

	//}

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



