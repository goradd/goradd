package page

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/spekary/goradd"
	"github.com/spekary/goradd/orm/db"
	"goradd/config"
	"math"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
)

type RequestMode int

const (
	Server     RequestMode = iota // calling back in to currently showing page using a standard form post
	Http                          // new page request
	Ajax                          // calling back in to a currently showing page using an ajax request
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
	err                 error // An error that occurred during the unpacking of the context. We save This for later so we can let the page manager display it if we get that far.
	requestMode         RequestMode
	cliArgs             []string // All arguments from the command line, whether from the command line call, or the ones that started the daemon
	pageStateId         string
	customControlValues map[string]interface{} // map of new control values keyed by control id. This supplements what comes through in the formVars as regular post variables. Numbers are preserved as json.Number types.
	checkableValues     map[string]interface{} // map of checkable control values, keyed by id. Values could be a true/false, an id from a radio group, or an array of ids from a checkbox group
	actionControlID     string                 // If an action, the control sending the action
	eventID             EventID                // The event to send to the control
	actionValues        ActionValues
	WasHandled          bool
	// TODO: Session object
}

type Context struct {
	HttpContext
	AppContext
}

func (c *Context) String() string {
	b, _ := json.Marshal(c.actionValues)
	actionValues := string(b[:])
	return fmt.Sprintf("URL: %s, Mode: %s, Form Values: %v, Control ID: %s, Event ID: %d, Action Values: %s, Page State: %s", c.URL, c.requestMode, c.formVars, c.actionControlID, c.eventID, actionValues, c.pageStateId)
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

// Returns the corresponding form value as a string slice. Use This when you are expecting more than one value in the
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

// FillApp fills the app structure with app specific information from the request
// Do not panic here!
func (ctx *Context) FillApp(cliArgs []string) {
	var ok bool
	var v string = ""
	//var i interface{}
	var err error

	if ctx.URL != nil {
		if v, ok = ctx.FormValue(htmlVarParams); ok {
			if h := ctx.Header.Get("X-Requested-With"); strings.ToLower(h) == "xmlhttprequest" {
				ctx.requestMode = Ajax
			} else {
				ctx.requestMode = Server
			}

			var params struct {
				ControlValues   map[string]interface{} `json:"controlValues"`
				CheckableValues map[string]interface{} `json:"checkableValues"`
				ControlID       string                 `json:"controlID"`
				EventID         int                    `json:"eventID"`
				Values          ActionValues           `json:"actionValues"`
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

				/*
						// We have posted back from our form. Unpack our params values
					var params map[string]interface{}
					if err = json.Unmarshal([]byte(v), &params); err == nil {
						if i, ok = params["controlValues"]; !ok {
							ctx.customControlValues = make(map[string]interface{}) // empty map so we don't have to check for nil
						} else {
							if ctx.customControlValues, ok = i.(map[string]interface{}); !ok {
								ctx.err = fmt.Errorf("customControlValues wrong type")
							}
						}

						if i, ok = params["checkableValues"]; !ok {
							ctx.checkableValues = make(map[string]interface{})
						} else {
							if ctx.checkableValues, ok = i.(map[string]interface{}); !ok {
								ctx.err = fmt.Errorf("checkableValues wrong type")
							}
						}

						if i, ok = params["controlID"]; ok {
							if ctx.actionControlID, ok = i.(string); !ok {
								ctx.err = fmt.Errorf("controlID wrong type")
							}
						}

						if i, ok = params["eventID"]; ok {
							if v, ok2 := i.(int); !ok2 {
								ctx.err = fmt.Errorf("eventID wrong type")
							} else {
								ctx.eventID = EventID(v)
							}
						}

						if i, ok = params["values"]; ok {
							if m,ok2 := i.(map[string]interface{}); !ok2 {
								ctx.err = fmt.Errorf("values wrong type")
							} else {
								ctx.actionValues = m
							}
						}
				*/
				if ctx.pageStateId, ok = ctx.FormValue("Goradd__FormState"); !ok {
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

func (ctx *Context) RequestMode() RequestMode {
	return ctx.requestMode
}

func GetContext(ctx context.Context) *Context {
	return ctx.Value(goradd.PageContext).(*Context)
}

func fixActionValues(values ActionValues) ActionValues {
	values.Control = fixActionValue(values.Control)
	values.Event = fixActionValue(values.Event)
	values.Action = fixActionValue(values.Action)
	return values
}

func fixActionValue(val interface{}) interface{} {
	switch v := val.(type) {
	case float64:
		if i := math.Trunc(v); v == i {
			return int(i)
		}
		return val
	}
	return val
}
