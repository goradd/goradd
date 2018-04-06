package page

import (
	"bytes"
	"context"
	"reflect"
)

// WrapperI defines the control wrapper interface. A control wrapper takes the basic html output by a control and wraps
// it in additional html to give it context. See the 2 built-in wrappers: LabelWrapper and ErrorWrapper for examples.
// For example, wrappers can be used to add labels that are connected to a control, give additional information, or show error conditions.
type WrapperI interface {
	Wrap(ctx context.Context, ctrl ControlI, html string, buf *bytes.Buffer)
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

type ErrorWrapper struct {
}

func NewErrorWrapper() ErrorWrapper {
	return ErrorWrapper{}
}

func (w ErrorWrapper) Wrap(ctx context.Context, ctrl ControlI, html string, buf *bytes.Buffer) {
	ErrorTmpl(ctx, ctrl, html, buf)
}

func (w ErrorWrapper) TypeName() string {
	return "page.Error"
}

type LabelWrapper struct {
	ErrorWrapper
}

func NewLabelWrapper() LabelWrapper {
	return LabelWrapper{}
}

func (w LabelWrapper) Wrap(ctx context.Context, ctrl ControlI, html string, buf *bytes.Buffer) {
	NameTmpl(ctx, ctrl, html, buf)
}

func (w LabelWrapper) TypeName() string {
	return "page.Label"
}

func init() {
	RegisterControlWrapper("page.Error", &ErrorWrapper{})
	RegisterControlWrapper("page.Label", &LabelWrapper{})
}