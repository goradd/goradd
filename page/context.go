package page

import (
	"net/url"
	"net/http"
	"mime/multipart"
	"strings"
	"goradd/config"
	"context"
	"encoding/json"
	"fmt"
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

// AppContext has Goradd application specific nodes
type AppContext struct {
	err error					// An error that occurred during the unpacking of the context. We save this for later so we can let the page manager display it if we get that far.
	requestMode RequestMode
	cliArgs []string			// All arguments from the command line, whether from the command line call, or the ones that started the daemon
	pageStateId string
	customControlValues map[string]interface{} // map of new control values keyed by control id
	checkableValues map[string]interface{} // map of checkable control values, keyed by id. Values could be a true/false, an id from a radio group, or an array of ids from a checkbox group
	// TODO: Session object
}

type Context struct {
	HttpContext
	AppContext
}

func (ctx *Context) FillFromRequest(cliArgs []string, r *http.Request)() {
	ctx.FillHttp(r)
	ctx.FillApp(cliArgs)
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
		value = v[0]
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

// FillApp fills the app structure with app specific information from the request
// Do not panic here!
func (ctx *Context) FillApp(cliArgs []string) {
	var ok bool
	var v string = ""
	var i interface{}
	var err error


	if ctx.URL != nil {
		if v,ok = ctx.FormValue("Goradd__Params"); ok {
			if h := ctx.Header.Get("X-Requested-With"); strings.ToLower(h) == "xmlhttprequest" {
				ctx.requestMode = Ajax
			} else {
				ctx.requestMode = Server
			}

			// We have posted back from our form. Unpack our params values
			var params map[string]interface{}
			if err = json.Unmarshal([]byte(v), &params); err == nil {
				if i, ok = params["controlValues"]; !ok {
					ctx.customControlValues = make(map[string]interface{}) // empty map so we don't have to check for nil
				} else {
					ctx.customControlValues = i.(map[string]interface{})
				}

				if i, ok = params["checkableValues"]; !ok {
					ctx.checkableValues = make(map[string]interface{})
				} else {
					ctx.checkableValues = i.(map[string]interface{})
				}

				if ctx.pageStateId,ok = ctx.FormValue("Goradd__FormState"); !ok {
					ctx.err = fmt.Errorf("No formstate found in response")
					return
				}
			} else {
				ctx.err = err
				return
			}

		} else {
			// Scenarios where we are not posting the form

			if h := ctx.Header.Get("X-Requested-With"); strings.ToLower(h) == "xmlhttprequest" {
				// A custom ajax call
				ctx.requestMode = CustomAjax
			} else {
				// A new call to our web page
				ctx.requestMode = Http
			}
		}
	} else {
		ctx.requestMode = Cli
	}
	ctx.cliArgs = cliArgs
}

func GetContext(ctx context.Context) *Context {
	return ctx.Value("goradd").(*Context) // TODO: Must replace the context key with something that is not a basic string. See https://medium.com/@matryer/context-keys-in-go-5312346a868d.
}
