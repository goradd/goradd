package page

import (
	"bytes"
	"context"
	"reflect"
)

// WrapperI defines the control wrapper interface. A control wrapper takes the basic html output by a control and wraps
// it in additional html to give it context. See the 2 built-in wrappers: NameWrapper and ErrorWrapper for examples.
// For example, wrappers can be used to add labels that are connected to a control, give additional information, or show error conditions.
type WrapperI interface {
	Init()
	Wrap(ctx context.Context, html string, buf *bytes.Buffer)
	Serialize(buf *bytes.Buffer)
	Unserialize(buf *bytes.Buffer) int
	// Zero means reset. Any other number represents an error condition.
	setControl(ControlI)
	SetModified(bool)
	IsModified() bool
	// Type is the type of the wrapper structure
	TypeName() string
}

var wrapperRegistry = map[string]reflect.Type{}

func RegisterControlWrapper(name string, w WrapperI) {
	t := reflect.TypeOf(w)
	wrapperRegistry[name] = t
}

func GetRegisteredWrapper(name string) WrapperI {
	if w,ok := wrapperRegistry[name]; ok {
		var i interface{}
		i = reflect.New(w)
		return i.(WrapperI)
	}
	return nil
}

type WrapperBase struct {
	isModified bool
	control ControlI
}

func (w *WrapperBase) Init() {
}

func (w *WrapperBase) IsModified() bool  {
	return w.isModified
}

func (w *WrapperBase) SetModified(m bool)  {
	w.isModified = m
}

func (w *WrapperBase) setControl(c ControlI)  {
	w.control = c
}


type ErrorWrapper struct {
	WrapperBase
}

func (w *ErrorWrapper) Wrap(ctx context.Context, html string, buf *bytes.Buffer) {
	ErrorTmpl(ctx, w.control, html, buf)
}

func (w *ErrorWrapper) TypeName() string {
	return "page.Error"
}

func (w *ErrorWrapper) Serialize(buf *bytes.Buffer) {

}

func (w *ErrorWrapper) Unserialize(buf *bytes.Buffer) int {
	return 0
}

type NameWrapper struct {
	ErrorWrapper
}

func (w *NameWrapper) Wrap(ctx context.Context, html string, buf *bytes.Buffer) {
	NameTmpl(ctx, w.control, html, buf)
}

func (w *NameWrapper) TypeName() string {
	return "page.Name"
}

func init() {
	RegisterControlWrapper("page.Error", &ErrorWrapper{})
	RegisterControlWrapper("page.Name", &NameWrapper{})
}