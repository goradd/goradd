package page

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/goradd/goradd/pkg/goradd"
	"github.com/goradd/goradd/pkg/log"
	"github.com/goradd/goradd/pkg/orm/db"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
)

// RequestMode tracks what kind of request we are processing.
type RequestMode int

const (
	// Server indicates we are calling back to a previously sent form using a standard form post
	Server     RequestMode = iota
	// Http indicates this is a first-time request for a page
	Http
	// Ajax indicates we are calling back in to a currently showing form using an ajax request
	Ajax
	// CustomAjax indicates we are calling an entry point from ajax, but not through our js file. This could be used to
	// implement a Rest API at a custom location.
	CustomAjax
	// Cli indicates we are being called from the command line and not through the http server.
	Cli
)

const HtmlVarAction = "Goradd_Action"
const htmlVarPagestate = "Goradd__PageState"
const htmlVarParams  = "Goradd__Params"


// MultipartFormMax is the maximum size of a mult-part form that we will allow.
var MultipartFormMax int64 = 10000000 // 10MB max in memory file

// String satisfies the Stringer interface and returns a description of the RequestMode
func (m RequestMode) String() string {
	switch m {
	case Server:
		return "Server"
	case Http:
		return "Http"
	case Ajax:
		return "Ajax"
	case CustomAjax:
		return "Custom Ajax"
	case Cli:
		return "Command-line"
	}
	return "Unknown"
}

/*
type ContextI interface {
	Http() *HttpContext
	App() *AppContext
}*/

// Context is the page context that we embed in the context.Context object that is passed throughout the application,
// and contains the per-request information that needs to be sent to various parts of the program. It primarily
// consists of items that we unpack from the http request. To get to it, simply call GetContext(ctx), where
// ctx is the context taken from the http request. The framework will take care of setting this up when
// a request is received.
type Context struct {
	HttpContext
	AppContext
}

// HttpContext contains typical things we can extract from an http request.
type HttpContext struct {
	// Req is the original http.Request object
	Req        *http.Request
	// URL is the url being queried
	URL        *url.URL
	// formVars is a private version of the form variables. Use the FormValue and FormValues functions to get these
	formVars   url.Values
	// Host is the host value extracted from the request
	Host       string
	// RemoteAddr is the ip address of the client
	RemoteAddr string
	// Referrer is the referring url, if there is one and it is included in the request. In other words, if a link was
	// clicked to get here, it would be the URL of the page that had the link
	Referrer   string
	// Cookies are the cookies coming from the client, mapped by name
	Cookies    map[string]*http.Cookie
	// Files are the files being uploaded, if this is a file upload. This currently only works with Server calls
	// in response to a file upload control.
	Files      map[string][]*multipart.FileHeader
	// Header is the http header coming from the client.
	Header     http.Header
}

// AppContext has Goradd application specific information.
type AppContext struct {
	err                 error // An error that occurred during the unpacking of the context. We save this for later so we can let the override manager display it if we get that far.
	requestMode         RequestMode
	cliArgs             []string // All arguments from the command line, whether from the command line call, or the ones that started the daemon
	pageStateId         string
	customControlValues map[string]map[string]interface{} // map of new control values keyed by control id. This supplements what comes through in the formVars as regular post variables. Numbers are preserved as json.Number types.
	checkableValues     map[string]interface{}            // map of checkable control values, keyed by id. Values could be a true/false, an id from a radio group, or an array of ids from a checkbox group
	actionControlID     string                            // If an action, the control sending the action
	eventID             EventID                           // The event to send to the control
	actionValues        actionValues
	// OutBuf is the output buffer being used to respond to the request. At the end of processing, it will be written to the response.
	OutBuf              *bytes.Buffer
	// NoJavaScript indicates javascript is turned off by the browser
	NoJavaScript        bool
}


// String is a string representation of all the information in the context, and should primarily be used for debugging.
func (c *Context) String() string {
	b, _ := json.Marshal(c.actionValues)
	actionValues := string(b[:])
	s := fmt.Sprintf("URL: %s, Mode: %s, FormBase Values: %v, Control ID: %s, Event ID: %d, Action Values: %s, Page State: %s", c.URL, c.requestMode, c.formVars, c.actionControlID, c.eventID, actionValues, c.pageStateId)

	if c.err != nil {
		s += fmt.Sprintf(", Error: %s", c.err.Error())
	}
	return s
}

// PutContext is used by the framework to insert the goradd context as a value in the standard GO context.
// You should not normally call this, unless you are customizing how your http server works.
func PutContext(r *http.Request, cliArgs []string) *http.Request {
	ctx := r.Context()
	grctx := &Context{}

	err := grctx.fillHttp(r)
	if err != nil {
		log.Error("Error creating http context: " + err.Error())
	}
	grctx.fillApp(cliArgs)
	ctx = context.WithValue(ctx, goradd.PageContext, grctx)

	// Create a context that the orm can use
	ctx = context.WithValue(ctx, goradd.SqlContext, &db.SqlContext{})

	return r.WithContext(ctx)
}

func (ctx *Context) fillHttp(r *http.Request) (err error) {
	if contentType := r.Header.Get("content-type"); contentType != "" {
		// Per comments in the ResponseWriter, we need to read and processs the entire request before attempting to write.
		if strings.Contains(contentType, "multipart") {
			// TODO: The Go doc is vague about how it handles file uploads larger than this value. Some doc suggests it
			// will return an error, and other doc suggests it will just split it into multiple partial files.
			// Nothing explains how to prevent malicious code from attempting to upload a gigantic file.
			// Likely we need to check the header for a size before attempting to parse. We will need to experiment to try to prevent this.
			err = r.ParseMultipartForm(MultipartFormMax)
		} else {
			err = r.ParseForm()
		}
	} else {
		err = r.ParseForm()
	}

	ctx.Req = r
	ctx.URL = r.URL
	ctx.formVars = r.Form
	ctx.Host = r.Host
	ctx.RemoteAddr = r.RemoteAddr
	ctx.Referrer = r.Referer()
	ctx.Header = r.Header

	ctx.Cookies = make(map[string]*http.Cookie)
	for _, c := range r.Cookies() {
		ctx.Cookies[c.Name] = c
	}

	if r.MultipartForm != nil {
		ctx.Files = r.MultipartForm.File
	}
	return
}

// FormValue returns the given form variable value, either from post or get variables.
// If the value does not exist, or is a multi-part value, returns false in ok.
// Use FormValues for multipart values.
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

// FormValues returns the corresponding form value as a string slice. Use this when you are expecting more than
// one value in the given form variable
func (ctx *Context) FormValues(key string) (value []string, ok bool) {
	if ctx.formVars == nil {
		return
	}
	value, ok = ctx.formVars[key]
	return
}

// CheckableValue returns the value of the named checkable value. This would be something coming from a
// checkbox or radio button. You do not normally call this unless you are implementing a checkable control widget.
func (ctx *Context) CheckableValue(key string) (value interface{}, ok bool) {
	if ctx.NoJavaScript {
		// checkable values do not exist, and we are POSTing, so we have to trust that everything is on screen.
		value,ok = ctx.FormValue(key)
		if !ok {
			// In a POST, checkable values only exist if they are checked.
			// This requires great care when using a parent control that is paging a lot of child controls.
			value = false
		}
		return
	}
	if ctx.checkableValues == nil {
		return
	}
	value, ok = ctx.checkableValues[key]
	return
}

// CheckableValues returns multiple checkable values. You do not normally call this unless you are implementing
// a widget that would have multiple checkable values, like a checklist.
func (ctx *Context) CheckableValues() map[string]interface{} {
	return ctx.checkableValues
}

// CustomControlValue returns the value of a control that is using the custom control mechanism to report
// its values. You would only call this if your are implementing a control that has custom javascript to
// operate its UI.
func (ctx *Context) CustomControlValue(id string, key string) interface{} {
	if m,ok := ctx.customControlValues[id]; ok {
		if v,ok2 := m[key]; ok2 {
			return v
		}
	}
	return nil
}

// fillApp fills the app structure with app specific information from the request
// Do not panic here!
func (ctx *Context) fillApp(cliArgs []string) {
	var ok bool
	var v string = ""
	//var i interface{}
	var err error

	if ctx.URL != nil {
		if v, ok = ctx.FormValue(htmlVarParams); ok  {
			if v == "" {
				// javascript is turned off, meaning we are forced to use server submits
				// we are in a minimalist environment, where only buttons submit forms
				ctx.NoJavaScript = true
				ctx.requestMode = Server
				ctx.actionControlID, _ = ctx.FormValue(HtmlVarAction)
				if ctx.pageStateId, ok = ctx.FormValue(htmlVarPagestate); !ok {
					ctx.err = fmt.Errorf("No pagestate found in response")
					return
				}
				return
			}
			if h := ctx.Header.Get("X-Requested-With"); strings.ToLower(h) == "xmlhttprequest" {
				ctx.requestMode = Ajax
			} else {
				ctx.requestMode = Server
			}

			var params struct {
				ControlValues   map[string]map[string]interface{} `json:"controlValues"`
				CheckableValues map[string]interface{} `json:"checkableValues"`
				ControlID       string                 `json:"controlID"`
				EventID         int                    `json:"eventID"`
				Values          actionValues           `json:"actionValues"`
			}

			dec := json.NewDecoder(strings.NewReader(v))
			dec.UseNumber()
			if err = dec.Decode(&params); err == nil {
				ctx.customControlValues = params.ControlValues
				ctx.checkableValues = params.CheckableValues
				ctx.actionControlID = params.ControlID
				if params.EventID != 0 {
					ctx.eventID = EventID(params.EventID)
				}
				ctx.actionValues = params.Values

				if ctx.pageStateId, ok = ctx.FormValue(htmlVarPagestate); !ok {
					ctx.err = fmt.Errorf("No pagestate found in response")
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
				// A new call to our web override
				ctx.requestMode = Http
			}
		}
	} else {
		ctx.requestMode = Cli
	}
	ctx.cliArgs = cliArgs
}

// RequestMode returns the request mode of the current request.
func (ctx *Context) RequestMode() RequestMode {
	return ctx.requestMode
}

// GetContext returns the page context from the GO context.
func GetContext(ctx context.Context) *Context {
	return ctx.Value(goradd.PageContext).(*Context)
}

// ConvertToBool is a helper function that can convert Put or Get values and other possible kinds of values into
// a bool value.
func ConvertToBool(v interface{}) bool {
	var val bool
	switch s := v.(type) {
	case string:
		sLower := strings.ToLower(s)
		if sLower == "true" || sLower == "on" || sLower == "1" {
			val = true
		} else if sLower == "false" || sLower == "off" || sLower == "" || sLower == "0" {
			val = false
		} else {
			panic(fmt.Errorf("unknown checkbox string value: %s", s))
		}
	case int:
		if s == 0 {
			val = false
		} else {
			val = true
		}
	case bool:
		val = s
	default:
		panic(fmt.Errorf("unknown checkbox value: %v", v))
	}

	return val
}

// NewMockContext creates a context for testing.
func NewMockContext() (ctx context.Context) {
	r := httptest.NewRequest("", "/", nil)
	r = PutContext(r, nil)
	ctx = r.Context()
	return
}
