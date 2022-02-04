package http

import (
	"fmt"
	"github.com/goradd/goradd/pkg/log"
	"io"
	"net/http"
	"runtime"
	"strings"
	"time"
)

var MaxErrorStackDepth = 20

// Error represents an error response to an http request.
//
// See http.Status* codes for status code constants
type Error struct {
	Message string
	Headers map[string]string
	ErrCode int
}

// SetResponseHeader sets a key-value in the header response.
func (e *Error) SetResponseHeader(key, value string) {
	if e.Headers == nil {
		e.Headers = map[string]string{key: value}
	} else {
		e.Headers[key] = value
	}
}

func (e Error) Error() string {
	return e.Message
}

// SendErrorCode will cause the page to error with the given http error code.
func SendErrorCode(errCode int) {
	e := Error{ErrCode: errCode}
	panic(e)
}

func SendErrorMessage(message string, errCode int) {
	e := Error{ErrCode: errCode, Message: message}
	panic(e)
}

// Redirect will error such that the server will attempt to access the
// resource at a new location.
//
// This will set the Location header to point to the new location.
//
// Be sure to call http.MakeLocalPath() if the resource is pointing to a
// location on this server.
//
// errCode should be a 3XX error, like one of the following:
//	StatusMovedPermanently  = 301 // RFC 7231, 6.4.2
//	StatusFound             = 302 // RFC 7231, 6.4.3
//	StatusSeeOther          = 303 // RFC 7231, 6.4.4
//	StatusTemporaryRedirect = 307 // RFC 7231, 6.4.7
//	StatusPermanentRedirect = 308 // RFC 7538, 3
func Redirect(location string, errCode int) {
	e := Error{ErrCode: errCode}
	e.SetResponseHeader("Location", location)
	panic (e)
}

// SendUnauthorized will send an error code indicating that the user is not authenticated (yes,
// even though the title is "authorized", it really means "authenticated", i.e. not logged in.)
// If serving HTML, you likely should redirect to the login page instead.
func SendUnauthorized() {
	e := Error{ErrCode: http.StatusUnauthorized}
	panic (e)
}

// SendForbidden will tell the user that he/she does not have authorization to acceess
// the given resource. The user should be known.
func SendForbidden() {
	e := Error{ErrCode: http.StatusForbidden}
	panic (e)
}

// SendMethodNotAllowed will tell the user that the server is not able
// to perform the http method being asked. allowedMethods is a list of the allowed methods.
func SendMethodNotAllowed(allowedMethods ...string) {
	e := Error{ErrCode: http.StatusMethodNotAllowed}
	e.SetResponseHeader("Allow", strings.Join(allowedMethods, ","))
	panic(e)
}

func SendNotFound() {
	e := Error{ErrCode: http.StatusNotFound}
	panic (e)
}

func SendNotFoundMessage(message string) {
	e := Error{ErrCode: http.StatusNotFound, Message: message}
	panic (e)
}

func SendBadRequest() {
	e := Error{ErrCode: http.StatusBadRequest}
	panic (e)
}

func SendBadRequestMessage(message string) {
	e := Error{ErrCode: http.StatusBadRequest, Message: message}
	panic (e)
}

// ServerError represents an error caused by an unexpected panic
type ServerError struct {
	// the error string
	Err string
	// Mode indicates whether we are serving ajax or not
	Mode string
	// the time the error occurred
	Time time.Time
	Request *http.Request
	// Output will replace what gets written to the output
	Output string
	// How much additional to unwind the stack trace
	StackDepth int
}


// Error returns the string that is sent to the logger
func (s ServerError) Error() string {
	out := s.Err + "\n"
	out += s.Mode + "  " + s.Request.RequestURI + " " + fmt.Sprintf("%v\n", s.Request.PostForm)
	return out
}

func NewServerError(err string, mode string, r *http.Request, skipFrames int, output string) *ServerError {
	e := ServerError{
		Err: err,
		Mode: mode,
		Time: time.Now(),
		Request: r,
		Output: output,
		StackDepth: skipFrames,
	}

	return &e
}


type  ErrorReporter struct {
}

// Use wraps the given handler in a default HTTP error handler that
// will respond appropriately to any panics that happen within the given handler.
//
// Panic with an http.Error value to get a specific kind of http error to
// be output. Otherwise, errors will be sent to the log.Error logger.
func (e ErrorReporter) Use(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				stackDepth := 2
				var newResponse string
				var errMsg string
				switch v := r.(type) {
				case Error: // A kind of http panic that just returns a response code and headers
					header := w.Header()
					for k,v2 := range v.Headers {
						header.Set(k,v2)
					}
					w.WriteHeader(v.ErrCode)
					// Write out the message if it is present as the visible response to the error
					if v.Message != "" {
						_,_ = w.Write([]byte(v.Message))
					}
					return // a normal http error response, so keep going
				case int:
					w.WriteHeader(v)
					return // a normal http error response, so keep going

				case *ServerError:
					newResponse = v.Output
					errMsg = v.Error()
					stackDepth += v.StackDepth
				case error:
					errMsg = v.Error()
				case string:
					errMsg = v
				default:
					errMsg = fmt.Sprintf("%v", v)
				}
				w.WriteHeader(http.StatusInternalServerError)
				buf := ResetOutputBuffer(req.Context())
				errMsg += "\nPartial response written:\n" + string(buf)
				log.Error(errMsg + stackTrace(stackDepth)) // use the application logger to output the error so we know about it
				_,_ = io.WriteString(w, newResponse) // Write the alternate response to client
				return
			}
		}()
		h.ServeHTTP(w, req)

	})
}

func stackTrace(startingDepth int) (out string ){
	for i := 1 + startingDepth; i < MaxErrorStackDepth; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		name := ""
		if f := runtime.FuncForPC(pc); f != nil {
			name = f.Name()
		}

		out += fmt.Sprintf("%s:%d -- %s\n", file, line, name)
	}
	return
}