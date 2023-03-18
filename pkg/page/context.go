package page

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/goradd/goradd/pkg/crypt"
	"github.com/goradd/goradd/pkg/goradd"
	http2 "github.com/goradd/goradd/pkg/http"
	"github.com/goradd/goradd/pkg/log"
	"github.com/goradd/goradd/pkg/page/action"
	"github.com/goradd/goradd/pkg/page/event"
	"github.com/goradd/goradd/pkg/session"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
)

// RequestMode tracks what kind of request we are processing.
type RequestMode int

const (
	// Server indicates we are calling back to a previously sent form using a standard form post
	Server RequestMode = iota
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
const HtmlVarPagestate = "Goradd__PageState"
const HtmlVarApistate = "__ApiState"
const htmlVarParams = "Goradd__Params"
const htmlCsrfToken = "Goradd__Csrf"

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
		return "NewEvent Ajax"
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
	// Request is the original http.Request object
	Request *http.Request
	// URL is the url being queried
	URL *url.URL
	// formVars is a private version of the form variables. Use the FormValue and FormValues functions to get these
	formVars url.Values
	// Host is the host value extracted from the request
	Host string
	// RemoteAddr is the ip address of the client
	RemoteAddr string
	// Referrer is the referring url, if there is one and it is included in the request. In other words, if a link was
	// clicked to get here, it would be the URL of the page that had the link
	Referrer string
	// Cookies are the cookies coming from the client, mapped by name
	Cookies map[string]*http.Cookie
	// Files are the files being uploaded, if this is a file upload. This currently only works with Server calls
	// in response to a file upload control.
	Files map[string][]*multipart.FileHeader
	// Header is the http header coming from the client.
	Header http.Header
}

// AppContext has Goradd application specific information.
type AppContext struct {
	err                  error // An error that occurred during the unpacking of the context. We save this for later so we can let the override manager display it if we get that far.
	requestMode          RequestMode
	cliArgs              []string // All arguments from the command line, whether from the command line call, or the ones that started the daemon
	pageStateId          string
	customControlValues  map[string]map[string]interface{} // map of new control values keyed by control id. This supplements what comes through in the formVars as regular post variables. Numbers are preserved as json.Number types.
	actionControlID      string                            // If an action, the control sending the action
	eventID              event.EventID                     // The event to send to the control
	actionValues         action.RawActionValues
	refreshIDs           []string
	hasTimezoneInfo      bool
	clientTimezoneOffset int
	clientTimezone       string

	// NoJavaScript indicates javascript is turned off by the browser
	NoJavaScript bool
}

// String is a string representation of all the information in the context, and should primarily be used for debugging.
func (ctx *Context) String() string {
	b, _ := json.Marshal(ctx.actionValues)
	actionValues := string(b[:])
	s := fmt.Sprintf("URL: %s, Mode: %s, FormBase Values: %v, ControlBase ID: %s, Event ID: %d, DoAction Values: %s, Page State: %s", ctx.URL, ctx.requestMode, ctx.formVars, ctx.actionControlID, ctx.eventID, actionValues, ctx.pageStateId)

	if ctx.err != nil {
		s += fmt.Sprintf(", Error: %s", ctx.err.Error())
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
	grctx.fillApp(ctx, cliArgs)
	ctx = context.WithValue(ctx, goradd.PageContext, grctx)

	return r.WithContext(ctx)
}

func (ctx *Context) fillHttp(r *http.Request) (err error) {
	if contentType := r.Header.Get("content-type"); contentType != "" {
		// Per comments in the ResponseWriter, we need to read and process the entire request before attempting to write.
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

	ctx.Request = r
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

// CustomControlValue returns the value of a control that is using the custom control mechanism to report
// its values. You would only call this if your are implementing a control that has custom javascript to
// operate its UI.
func (ctx *Context) CustomControlValue(id string, key string) interface{} {
	if m, ok := ctx.customControlValues[id]; ok {
		if v, ok2 := m[key]; ok2 {
			return v
		}
	}
	return nil
}

// HasCustomControlValue returns true if the given controls has a value for the given key. If you
// are potentially expecting a nil value, you can use this to know that a value is present.
func (ctx *Context) HasCustomControlValue(id string, key string) bool {
	if m, ok := ctx.customControlValues[id]; ok {
		_, ok2 := m[key]
		return ok2
	}
	return false
}

// fillApp fills the app structure with app specific information from the request
// Do not panic here!
func (ctx *Context) fillApp(mainContext context.Context, cliArgs []string) {
	var ok bool
	var v = ""
	//var i interface{}
	var err error

	if ctx.URL != nil {
		if ctx.pageStateId, ok = ctx.FormValue(HtmlVarPagestate); ok {
			v, _ = ctx.FormValue(htmlVarParams)
			if v == "" {
				// javascript is turned off
				// we are in a minimalist environment, where only buttons submit forms

				// If the pagestate is coming from a GET, it is encoded and encrypted
				if _, ok := ctx.Request.PostForm[HtmlVarPagestate]; !ok {
					ctx.pageStateId = crypt.SessionDecryptUrlValue(mainContext, ctx.pageStateId)
				}
				ctx.NoJavaScript = true
				ctx.requestMode = Server
				aId, _ := ctx.FormValue(HtmlVarAction)
				parts := strings.Split(aId, "_")
				ctx.actionControlID = parts[0]
				if len(parts) > 1 {
					ctx.actionValues.Control = []byte(parts[1])
				}
				return
			}
			if h := ctx.Header.Get("X-Requested-With"); strings.ToLower(h) == "xmlhttprequest" {
				ctx.requestMode = Ajax
			} else {
				ctx.requestMode = Server
			}

			type tzParams struct {
				TimezoneOffset int    `json:"o"`
				Timezone       string `json:"z"`
			}

			var params struct {
				ControlValues map[string]map[string]interface{} `json:"controlValues"`
				ControlID     string                            `json:"controlID"`
				EventID       int                               `json:"eventID"`
				Values        action.RawActionValues            `json:"actionValues"`
				RefreshIDs    []string                          `json:"refresh"`
				TimezoneInfo  tzParams                          `json:"tz"`
			}

			dec := json.NewDecoder(strings.NewReader(v))
			dec.UseNumber()
			if err = dec.Decode(&params); err == nil {
				ctx.customControlValues = params.ControlValues
				ctx.actionControlID = params.ControlID
				ctx.refreshIDs = params.RefreshIDs
				if params.EventID != 0 {
					ctx.eventID = event.EventID(params.EventID)
				}
				ctx.actionValues = params.Values
				ctx.clientTimezoneOffset = params.TimezoneInfo.TimezoneOffset
				ctx.clientTimezone = params.TimezoneInfo.Timezone
				ctx.hasTimezoneInfo = true

				// Save in a session for recovery when we have a session but do not have client info
				session.SetInt(mainContext, goradd.SessionTimezoneOffset, params.TimezoneInfo.TimezoneOffset)
				session.SetString(mainContext, goradd.SessionTimezone, params.TimezoneInfo.Timezone)

				if ctx.pageStateId, ok = ctx.FormValue(HtmlVarPagestate); !ok {
					ctx.err = fmt.Errorf("no pagestate found in response")
					return
				}
			} else {
				ctx.err = err
				return
			}

		} else if apistate, ok2 := ctx.FormValue(HtmlVarApistate); ok2 {
			// Allows REST clients to also support the timezone offset in the context
			if offset, err2 := strconv.Atoi(apistate); err2 == nil {
				ctx.clientTimezoneOffset = offset
				ctx.hasTimezoneInfo = true
			}
		} else {
			// Scenarios where we are not posting the form

			if h := ctx.Header.Get("X-Requested-With"); strings.ToLower(h) == "xmlhttprequest" {
				// A custom ajax call
				ctx.requestMode = CustomAjax
			} else {
				// A new call to our web page
				ctx.requestMode = Http

				// Recover client timezone if it was saved earlier
				if session.Has(mainContext, goradd.SessionTimezoneOffset) {
					ctx.hasTimezoneInfo = true
					ctx.clientTimezoneOffset = session.GetInt(mainContext, goradd.SessionTimezoneOffset)
					ctx.clientTimezone = session.GetString(mainContext, goradd.SessionTimezone)
				}
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

// ClientTimezoneOffset returns the number of minutes offset from GMT for the client's timezone.
func (ctx *Context) ClientTimezoneOffset() int {
	return ctx.clientTimezoneOffset
}

// ClientTimezone returns the name of the timezone of the client, if available.
func (ctx *Context) ClientTimezone() string {
	return ctx.clientTimezone
}

// HasTimezoneInfo returns true if timezone info is valid.
func (ctx *Context) HasTimezoneInfo() bool {
	return ctx.hasTimezoneInfo
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
	s := session.NewMock()
	session.SetSessionManager(s)
	r := httptest.NewRequest("", "/", nil)
	ctx = s.With(r.Context())
	r = r.WithContext(ctx)
	r = PutContext(r, nil)
	return r.Context()
}

// OutputLen returns the number of bytes that have been written to the output.
func OutputLen(ctx context.Context) int {
	return http2.OutputLen(ctx)
}

func ResetOutputBuffer(ctx context.Context) []byte {
	return http2.ResetOutputBuffer(ctx)
}
