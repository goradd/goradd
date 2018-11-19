package page

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/spekary/goradd/internal/goradd"
	"github.com/spekary/goradd/pkg/orm/db"
	"goradd-project/config"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
)

type RequestMode int

const (
	Server     RequestMode = iota // calling back in to currently showing override using a standard form post
	Http                          // new override request
	Ajax                          // calling back in to a currently showing override using an ajax request
	CustomAjax                    // calling an entry point from ajax, but not through our js file. REST API perhaps?
	Cli                           // From command line
)

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

type ContextI interface {
	Http() *HttpContext
	App() *AppContext
}

// Typical things we can extract from http
type HttpContext struct {
	Req        *http.Request
	URL        *url.URL
	formVars   url.Values
	Host       string
	RemoteAddr string
	Referrer   string
	Cookies    map[string]*http.Cookie
	Files      map[string][]*multipart.FileHeader
	Header     http.Header
}

// AppContext has Goradd application specific nodes
type AppContext struct {
	err                 error // An error that occurred during the unpacking of the context. We save This for later so we can let the override manager display it if we get that far.
	requestMode         RequestMode
	cliArgs             []string // All arguments from the command line, whether from the command line call, or the ones that started the daemon
	pageStateId         string
	customControlValues map[string]map[string]interface{} // map of new control values keyed by control id. This supplements what comes through in the formVars as regular post variables. Numbers are preserved as json.Number types.
	checkableValues     map[string]interface{}            // map of checkable control values, keyed by id. Values could be a true/false, an id from a radio group, or an array of ids from a checkbox group
	actionControlID     string                            // If an action, the control sending the action
	eventID             EventID                           // The event to send to the control
	actionValues        actionValues
	OutBuf              *bytes.Buffer
}

type Context struct {
	HttpContext
	AppContext
}

func (c *Context) String() string {
	b, _ := json.Marshal(c.actionValues)
	actionValues := string(b[:])
	s := fmt.Sprintf("URL: %s, Mode: %s, FormBase Values: %v, Control ID: %s, Event ID: %d, Action Values: %s, Page State: %s", c.URL, c.requestMode, c.formVars, c.actionControlID, c.eventID, actionValues, c.pageStateId)

	if c.err != nil {
		s += fmt.Sprintf(", Error: %s", c.err.Error())
	}
	return s
}

func PutContext(r *http.Request, cliArgs []string) *http.Request {
	ctx := r.Context()
	grctx := &Context{}

	grctx.FillHttp(r)
	grctx.FillApp(cliArgs)
	ctx = context.WithValue(ctx, goradd.PageContext, grctx)

	// Create a context that the orm can use
	ctx = context.WithValue(ctx, goradd.SqlContext, &db.SqlContext{})

	return r.WithContext(ctx)
}

func (ctx *Context) FillHttp(r *http.Request) {
	if _, ok := r.Header["content-type"]; ok {
		contentType := r.Header["content-type"][0]
		// Per comments in the ResponseWriter, we need to read and processs the entire request before attempting to write.
		if strings.Contains(contentType, "multipart") {
			r.ParseMultipartForm(config.MultiPartFormMax)
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

	ctx.Cookies = make(map[string]*http.Cookie)
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

func (ctx *Context) CheckableValue(key string) (value interface{}, ok bool) {
	if ctx.checkableValues == nil {
		return
	}
	value, ok = ctx.checkableValues[key]
	return
}

func (ctx *Context) CheckableValues() map[string]interface{} {
	return ctx.checkableValues
}

func (ctx *Context) CustomControlValue(id string, key string) interface{} {
	if m,ok := ctx.customControlValues[id]; ok {
		if v,ok2 := m[key]; ok2 {
			return v
		}
	}
	return nil
}

// FillApp fills the app structure with app specific information from the request
// Do not panic here!
func (ctx *Context) FillApp(cliArgs []string) {
	var ok bool
	var v string = ""
	//var i interface{}
	var err error

	if ctx.URL != nil {
		if v, ok = ctx.FormValue(htmlVarParams); ok  {
			if v == "" {
				// This is a javascript error. We unsuccessfully tried to gather form parameters
				ctx.err = fmt.Errorf("Javascript was not able to gather the goradd parameters.")
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

				if ctx.pageStateId, ok = ctx.FormValue("Goradd__FormState"); !ok {
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

func (ctx *Context) RequestMode() RequestMode {
	return ctx.requestMode
}

func GetContext(ctx context.Context) *Context {
	return ctx.Value(goradd.PageContext).(*Context)
}

// ConvertToBool is a helper function that can convert Put or Get values and other possible kinds of values into
// a bool value.
func ConvertToBool(v interface{}) bool {
	var val bool
	switch s := v.(type) {
	case string:
		slower := strings.ToLower(s)
		if slower == "true" || slower == "on" || slower == "1" {
			val = true
		} else if slower == "false" || slower == "off" || slower == "" || slower == "0" {
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
