package page

import (
	"context"
	"errors"
	"fmt"
	"io"
	"runtime"
	"time"
)

const MaxStackDepth = 50

// ErrorPageFuncType specifies the signature for the error template function. You can replace the built-in error
// function with your own function by setting the config.ErrorPage value. html is the html that was able to be generated before
// the error occurred, which can be helpful in tracking down the source of an error.
type ErrorPageFuncType func(ctx context.Context, html []byte, err *Error, w io.Writer) error

var ErrorPageFunc ErrorPageFuncType

// TODO: move to its own package so it is accessible from the entire framework

// The error structure, specifically designed to manage panics during request handling
type Error struct {
	// the error string
	Err error
	// the time the error occurred
	Time time.Time
	// the copied context when the error occurred
	Ctx *Context
	// unwound Stack info
	Stack []StackFrame
}

// StackFrame holds the file, line and function name in a call chain
type StackFrame struct {
	File string
	Line int
	Func string
}

// DbError represents a database error.
type DbError struct {
	Error
	// DbStatement is the captured database statement if one caused the error, returned by the db adapter
	DbStatement string
}

// NoErr represents no error. A request starts with this.
type NoErr struct {
}

// Known application specific errors

const (
	FrameworkErrNone = iota
	// FrameworkErrNoTemplate indicates a template does not exist. The control will move on to other ways of rendering.
	// No message should be displayed.
	FrameworkErrNoTemplate
	// FrameworkErrRecordNotFound is a rare situation that might come up as a race condition error between viewing a
	// record, and actually editing it. If in the time between clicking on a record to see detail, and viewing the detail,
	// the record was deleted by another user, we would return this error.
	// In a REST environment, this is  404 error
	FrameworkErrRecordNotFound
	// FrameworkErrNotAuthenticated indicates that the user needs to log in, but has not done so.
	// This will result in a 401 http error.
	FrameworkErrNotAuthenticated
	// FrameworkErrNotAuthorized indicates the logged in user does not have permission to access the resource.
	// This will result in a 403 error.
	FrameworkErrNotAuthorized
	// FrameworkErrRedirect indicates the resource is at a different location. This will result in a 303 error
	// telling the browser to go that resource.
	FrameworkErrRedirect
	// FrameworkErrBadRequest is a generic bad request error, but most often it indicates that some FORM or GET
	// parameter was expected, but was either missing or invalid. This is a 400 error.
	FrameworkErrBadRequest
	// FrameworkErrMethodNotAllowed is for REST clients to indicate that the user tried to use a method (i.e. GET or PUT)
	// that is not allowed at the given endpoint. The HTTP spec says the client should also respond with an Allow
	// header with a comma separated list of what methods are allowed. This is a 405 error.
	FrameworkErrMethodNotAllowed
)

// FrameworkError is an expected error that is part of the framework. Usually you would respond to the error
// by displaying a message to the user, but not always.
type FrameworkError struct {
	Err int
	Location string
	Message string	// optional message
}

// NewFrameworkError creates a new FrameworkError
func NewFrameworkError(err int) FrameworkError {
	return FrameworkError{Err: err}
}

func NewRedirectError(location string) FrameworkError {
	return FrameworkError{Err: FrameworkErrRedirect, Location: location}
}

func (e FrameworkError) SetMessage(msg string) FrameworkError {
	e.Message = msg
	return e
}

// Error returns the error string
func (e FrameworkError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	switch e.Err {
	case FrameworkErrNoTemplate:
		return "FormBase or control does not have a template" // just detected, this is not likely to be used
	case FrameworkErrRecordNotFound:
		return "Record does not exist."
	case FrameworkErrNotAuthenticated:
		return "You must log in."
	case FrameworkErrNotAuthorized:
		return "You are not authorized to access this information."
	case FrameworkErrRedirect:
		return "Redirecting to " + e.Location
	case FrameworkErrBadRequest:
		return "Your request did not make sense."
	case FrameworkErrMethodNotAllowed:
		return "The method of your request is not allowed. See the returned Allow header to know what is expected."
	}

	return ""
}

// HttpError returns the corresponding http error
func (e FrameworkError) HttpError() int {
	switch e.Err {
	case FrameworkErrNotAuthenticated: return 401
	case FrameworkErrNotAuthorized: return 403
	case FrameworkErrRecordNotFound: return 404
	case FrameworkErrBadRequest: return 400
	case FrameworkErrMethodNotAllowed: return 405
	}
	return 500
}

func RedirectHtml(loc string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
	<html>
	<head>
	<meta http-equiv="Refresh" content="0; url=%s" />
	</head>
	<body>
	<p>Redirecting to <a href="%[1]s">%[1]s</a>.</p>
	</body>
	</html>`, loc)
}

func (e *NoErr) Error() string {
	return ""
}

func IsError(e error) bool {
	_, ok := e.(*NoErr)
	return !ok
}

// Called by the page manager to record a system error
func newRunError(ctx context.Context, msg interface{}) *Error {
	e := &Error{}

	switch m := msg.(type) {
	case string: // we panic'd
		e.Err = errors.New(m)
		e.fillErr(ctx, 2)

	case error: // system generated error
		e.Err = m
		e.fillErr(ctx, 2)

	default:
		e.Err = fmt.Errorf("Error of type %T: %v", msg, msg)
		e.fillErr(ctx, 1)

	}
	return e
}

// NewError return a generic message error
func NewError(ctx context.Context, msg string) *Error {
	e := &Error{}

	e.Err = errors.New(msg)
	e.fillErr(ctx, 1)
	return e
}

func (e *Error) fillErr(ctx context.Context, skip int) {
	e.Time = time.Now()
	e.Ctx = GetContext(ctx)

	for i := 2 + skip; i < MaxStackDepth; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		name := ""
		if f := runtime.FuncForPC(pc); f != nil {
			name = f.Name()
		}

		frame := StackFrame{file, line, name}
		e.Stack = append(e.Stack, frame)
	}
}

func (e *Error) Error() string {
	return e.Err.Error()
}

// NewDbErr returns a new database error
func NewDbErr(ctx context.Context, msg interface{}, dbStatement string) *DbError {
	e := &DbError{}
	switch m := msg.(type) {
	case string:
		e.Err = errors.New(m)
	case error:
		e.Err = m
	default:
		e.Err = fmt.Errorf("Error of type %T: %v", msg, msg)
	}
	e.fillErr(ctx, 1)
	return e
}
