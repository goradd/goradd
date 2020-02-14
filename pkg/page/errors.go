package page

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"runtime"
	"time"
)

const MaxStackDepth = 50

// The ErrorPageTemplate type specifies the signature for the error template function. You can replace the built-in error
// function with your own function by setting the config.ErrorPage value. html is the html that was able to be generated before
// the error occurred, which can be helpful in tracking down the source of an error.
type ErrorPageFuncType func(ctx context.Context, html string, err *Error, buf *bytes.Buffer)

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
	FrameworkErrRecordNotFound
	// A standard situation when someone tries to go to a page they are not authorized to view.
	FrameworkErrNotAuthorized
	FrameworkErrRedirect
)

// FrameworkError is an expected error that is part of the framework. Usually you would respond to the error
// by displaying a message to the user, but not always.
type FrameworkError struct {
	Err int
	Location string
}

// NewFrameworkError creates a new FrameworkError
func NewFrameworkError(err int) FrameworkError {
	return FrameworkError{Err: err}
}

// Error returns the error string
func (e FrameworkError) Error() string {
	switch e.Err {
	case FrameworkErrNoTemplate:
		return "FormBase or control does not have a template" // just detected, this is not likely to be used
	case FrameworkErrRecordNotFound:
		return "Record does not exist. Perhaps it has been deleted by someone else?"
	case FrameworkErrNotAuthorized:
		return "You are not authorized to view this information."
	case FrameworkErrRedirect:
		return "Redirecting to " + e.Location
	}

	return ""
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
