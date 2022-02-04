package page

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"time"
)

const MaxStackDepth = 50

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
	}

	return ""
}

// HttpError returns the corresponding http error
func (e FrameworkError) HttpError() int {
	switch e.Err {
	case FrameworkErrRecordNotFound: return 404
	}
	return 500
}

func (e *NoErr) Error() string {
	return ""
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
