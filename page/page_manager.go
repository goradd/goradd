package page

import (
	"context"
	"goradd/config"
	"strings"
	"net/http"
	"fmt"
	"github.com/shiyanhui/hero"
	"html/template"
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
	gCtx := GetContext(ctx)

	pageStateId := gCtx.pageStateId

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
		stateId := m.cache.NewPageId()
		page.GetPageBase().stateId = pageStateId
		m.cache.Set(stateId, page)
		isNew = true
	}
	return
}


func (m *PageManager) RunPage(ctx context.Context, w http.ResponseWriter) {
	defer func() {
		if r := recover(); r != nil {
			switch v := r.(type) {
			case error:
				err := newRunError(ctx, v)
				m.makeErrorResponse(ctx, err, "", w)
			case string:
				err := newRunError(ctx, v)
				m.makeErrorResponse(ctx, err, "", w)
			default:
				err := newRunError(ctx,fmt.Errorf("Unknown package error: %v", r))
				m.makeErrorResponse(ctx, err, "", w)
			}
		}
	}()

	//gCtx := GetContext(ctx)
	pageI, isNew := m.getPage(ctx)

	page := pageI.GetPageBase()

	defer m.cleanup(pageI)

	page.renderStatus = UNRENDERED
	pageI.Run()

	if !isNew {
		// TODO: handle actions
	}

	//if server {
	// TODO: Make our own version of this so we are not dependent on hero
	buf := hero.GetBuffer()
	defer hero.PutBuffer(buf)

	err := pageI.Draw(ctx, buf)


	if err != nil {
		m.makeErrorResponse(ctx, newRunError(ctx, err), buf.String(), w)
	} else {
		w.Write(buf.Bytes())
	}
	//}

}

func (m *PageManager) cleanup (p PageI) {


}

func (m *PageManager) makeErrorResponse(ctx context.Context,
	err *Error,
	html string,
	w http.ResponseWriter) {

	if ErrorPageFunc == nil {
		panic("No error page template function is defined")
	}

	buf := hero.GetBuffer()
	defer hero.PutBuffer(buf)

	ErrorPageFunc(ctx, html, err, buf)
	w.Write(buf.Bytes())
}

func (m *PageManager) IsAsset (ctx context.Context) bool {
	return assetIsRegistered(GetContext(ctx).HttpContext.URL.Path)
}

func (m *PageManager) ServeAsset (ctx context.Context, w http.ResponseWriter, r *http.Request) {
	localpath := GetAssetFilePath(GetContext(ctx).HttpContext.URL.Path)
	if localpath == "" {
		panic("Invalid asset")
	}
	http.ServeFile(w, r, localpath)
}


