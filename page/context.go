package page

import (
	"net/url"
	"net/http"
	"mime/multipart"
	"strings"
	"goradd/config"
	"context"
)

type RequestMode int

const (
	Server RequestMode = iota 	// calling back in to currently showing page using a standard form post
	Http						// new page request
	Ajax						// calling back in to a currently showing page using an ajax request
	CustomAjax					// calling an entry point from ajax, but not through our js file. REST API perhaps?
	Cli							// From command line
)

type ContextI interface {
	Http() *HttpContext
	App() *AppContext
}

// Typical things we can extract from http
type HttpContext struct {
	Req *http.Request
	URL *url.URL
	formVars url.Values
	Host string
	RemoteAddr string
	Referrer string
	Cookies map[string]*http.Cookie
	Files map[string][]*multipart.FileHeader
	Header http.Header
}

// Goradd application specific nodes
type AppContext struct {
	requestMode RequestMode
	cliArgs []string			// All arguments from the command line, whether from the command line call, or the ones that started the daemon
	pageStateId string
	// TODO: Session object
}

type Context struct {
	HttpContext
	AppContext
	Response
}

func (ctx *Context) FillFromRequest(cliArgs []string, r *http.Request)() {
	ctx.FillHttp(r);
	ctx.FillApp(cliArgs);
}

func (ctx *Context) FillHttp(r *http.Request) {
	if _,ok := r.Header["content-type"]; ok {
		contentType := r.Header["content-type"][0]
		// Per comments in the ResponseWriter, we need to read and processs the entire request before attempting to write.
		if strings.Contains(contentType, "multipart") {
			r.ParseMultipartForm(config.MULTI_PART_FORM_MAX)
		} else {
			r.ParseForm()
		}
	} else {
		r.ParseForm()
	}

	ctx.Req = r
	ctx.URL = r.URL
	ctx.formVars = r.Form
	ctx.Host = r.Host
	ctx.RemoteAddr = r.RemoteAddr
	ctx.Referrer = r.Referer()
	ctx.Header = r.Header

	ctx.Cookies = make (map[string]*http.Cookie)
	for _, c := range r.Cookies() {
		ctx.Cookies[c.Name] = c
	}

	if r.MultipartForm != nil {
		ctx.Files = r.MultipartForm.File
	}
	return
}

// Returns the given form variable value, either from post or get variables
// If the value does not exist, or is a multi-part value, returns false in ok
// Use FormValues for multipart values
func (ctx *Context) FormValue(key string) (value string, ok bool) {
	if ctx.formVars == nil {
		return
	}
	var v []string
	if v, ok = ctx.formVars[key]; ok && len(v) == 1 {
		value = v[0];
	}
	return
}

// Returns the corresponding form value as a string slice. Use this when you are expecting more than one value in the
// given form variable
func (ctx *Context) FormValues(key string) (value []string, ok bool) {
	if ctx.formVars == nil {
		return
	}
	value, ok = ctx.formVars[key]
	return
}


func (ctx *Context) FillApp(cliArgs []string) {
	if ctx.URL != nil {
		if v,ok := ctx.FormValue("Qform__FormCallType"); ok {
			if v == "Ajax" {
				ctx.requestMode = Ajax
			} else {
				ctx.requestMode = Server
			}
		} else {
			if h := ctx.Header.Get("HTTP_X_REQUESTED_WITH"); strings.ToLower(h) == "xmlhttprequest" {
				ctx.requestMode = CustomAjax
			} else {
				ctx.requestMode = Http
			}
		}
	} else {
		ctx.requestMode = Cli
	}
	ctx.cliArgs = cliArgs
	ctx.pageStateId,_ = ctx.FormValue("Qform__FormState")

}

func GetContext(ctx context.Context) *Context {
	return ctx.Value("goradd").(*Context) // TODO: Must replace the context key with something that is not a basic string. See https://medium.com/@matryer/context-keys-in-go-5312346a868d.
}
