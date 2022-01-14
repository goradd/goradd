package http

import (
	"net/http"
	"strings"
)

// Error represents an error response to an http request.
//
// See http.Status* codes for status code constants
type Error struct {
	message string
	headers map[string]string
	errCode int
}

// SetResponseHeader sets a key-value in the header response.
func (e *Error) SetResponseHeader(key, value string) {
	if e.headers == nil {
		e.headers = map[string]string{key: value}
	} else {
		e.headers[key] = value
	}
}

// SendErrorCode will cause the page to error with the given http error code.
func SendErrorCode(errCode int) {
	e := Error{errCode: errCode}
	panic(e)
}

func SendErrorMessage(message string, errCode int) {
	e := Error{errCode: errCode, message: message}
	panic(e)
}

// Redirect will error such that the server will attempt to access the
// resource at a new location.
//
// This will set the Location header to point to the new location.
//
// Be sure to call page.MakeLocalPath() if the resource is pointing to a
// location on this server
func Redirect(location string, errCode int) {
	e := Error{errCode: errCode}
	e.SetResponseHeader("Location", location)
	panic (e)
}

func SendUnauthorized() {
	e := Error{errCode: http.StatusUnauthorized}
	panic (e)
}

// SendForbidden will tell the user that he/she does not have credentials
// for the given resource.
func SendForbidden() {
	e := Error{errCode: http.StatusForbidden}
	panic (e)
}

// SendMethodNotAllowed will tell the user that the server is not able
// to perform the http method being asked. allowedMethods is a list of the allowed methods.
func SendMethodNotAllowed(allowedMethods ...string) {
	e := Error{errCode: http.StatusMethodNotAllowed}
	e.SetResponseHeader("Allow", strings.Join(allowedMethods, ","))
	panic(e)
}

func SendNotFound() {
	e := Error{errCode: http.StatusNotFound}
	panic (e)
}

func SendNotFoundMessage(message string) {
	e := Error{errCode: http.StatusNotFound, message: message}
	panic (e)
}

func SendBadRequest() {
	e := Error{errCode: http.StatusBadRequest}
	panic (e)
}

func SendBadRequestMessage(message string) {
	e := Error{errCode: http.StatusBadRequest, message: message}
	panic (e)
}

// ErrorHandler wraps the given handler in a default HTTP error handler that
// will respond appropriately to any panics that happen within the given handler.
//
// Panic with an http.Error value to get a specific kind of http error to
// be output.
func ErrorHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				switch v := r.(type) {
				case Error: // A kind of http panic that just returns a response code and headers
					header := w.Header()
					for k,v := range v.headers {
						header.Set(k,v)
					}
					w.WriteHeader(v.errCode)
					if v.message != "" {
						_,_ = w.Write([]byte(v.message))
					}
				case error:
					w.WriteHeader(http.StatusInternalServerError)
					_,_ =w.Write([]byte(v.Error()))
				case string:
					w.WriteHeader(http.StatusInternalServerError)
					_,_ =w.Write([]byte(v))
				case int:
					w.WriteHeader(v)
				default:
					w.WriteHeader(http.StatusInternalServerError)
				}
			}
		}()
		h.ServeHTTP(w, r)
	})
}