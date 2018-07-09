package page

import (
	"bytes"
	"context"
	"reflect"
	"github.com/spekary/goradd/html"
)

// Wrapper constants used in the With function
const (
	ErrorWrapper = "page.Error"
	LabelWrapper = "page.Label"
	DivWrapper = "page.Div"
)

// WrapperI defines the control wrapper interface. A control wrapper takes the basic html output by a control and wraps
// it in additional html to give it context.
// For example, wrappers can be used to add labels that are connected to a control, give additional information,
// or show error conditions. See the built-in wrappers LabelWrapperType and ErrorWrapperType for examples.
type WrapperI interface {
	Wrap(ctx context.Context, ctrl ControlI, html string, buf *bytes.Buffer)
	ModifyDrawingAttributes(ctrl ControlI, attributes *html.Attributes)
}

var wrapperRegistry = map[string]reflect.Type{}

func RegisterControlWrapper(name string, w WrapperI) {
	t := reflect.TypeOf(w)
	wrapperRegistry[name] = t
}

// NewRegisteredWrapper returns a newly allocated named wrapper.
func NewRegisteredWrapper(name string) WrapperI {
	if w, ok := wrapperRegistry[name]; ok {
		var i interface{}
		i = reflect.New(w)
		return i.(WrapperI)
	}
	return nil
}

type ErrorWrapperType struct {
}

func NewErrorWrapper() ErrorWrapperType {
	return ErrorWrapperType{}
}

func (w ErrorWrapperType) Wrap(ctx context.Context, ctrl ControlI, html string, buf *bytes.Buffer) {
	ErrorTmpl(ctx, ctrl, html, buf)
}

func (w ErrorWrapperType) TypeName() string {
	return ErrorWrapper
}

func (w ErrorWrapperType) ModifyDrawingAttributes(c ControlI, a *html.Attributes) {
	state := c.control().validationState
	if state != NotValidated {
		a.Set("aria-describedby", c.ID() + "_err")
		if state == Valid {
			a.Set("aria-invalid", "false")
		} else {
			a.Set("aria-invalid", "true")
		}
	} else if c.control().instructions != "" {
		a.Set("aria-describedby", c.ID() + "_inst")
	}
}


type LabelWrapperType struct {
	ErrorWrapperType
	labelAttributes *html.Attributes
}

func NewLabelWrapper() LabelWrapperType {
	return LabelWrapperType{}
}

func (w LabelWrapperType) Wrap(ctx context.Context, ctrl ControlI, html string, buf *bytes.Buffer) {
	LabelTmpl(ctx, w, ctrl, html, buf)
}

// LabelAttributes returns attributes that will apply to the label. Changes will be remembered.
func (w *LabelWrapperType) LabelAttributes() *html.Attributes {
	if w.labelAttributes == nil {
		w.labelAttributes = html.NewAttributes()
	}
	return w.labelAttributes
}

func (w LabelWrapperType) HasLabelAttributes() bool {
	if w.labelAttributes == nil || w.labelAttributes.Len() == 0 {
		return false
	}
	return true
}


func (w LabelWrapperType) TypeName() string {
	return LabelWrapper
}

func (w LabelWrapperType) ModifyDrawingAttributes(c ControlI, a *html.Attributes) {
	state := c.control().validationState
	if state != NotValidated {
		a.Set("aria-describedby", c.ID() + "_err")
		if state == Valid {
			a.Set("aria-invalid", "false")
		} else {
			a.Set("aria-invalid", "true")
		}
	} else if c.control().instructions != "" {
		a.Set("aria-describedby", c.ID() + "_inst")
	}
	if c.control().label != "" && !c.control().hasFor { // if it has a for, then screen readers already know about the label
		a.Set("aria-labeledby", c.ID() + "_lbl")
	}
}


type DivWrapperType struct {
}

func NewDivWrapper() DivWrapperType {
	return DivWrapperType{}
}

func (w DivWrapperType) Wrap(ctx context.Context, ctrl ControlI, html string, buf *bytes.Buffer) {
	DivTmpl(ctx, ctrl, html, buf)
}

func (w DivWrapperType) TypeName() string {
	return DivWrapper
}

func (w DivWrapperType) ModifyDrawingAttributes(ctrl ControlI, a *html.Attributes) {
}



func init() {
	RegisterControlWrapper(ErrorWrapper, &ErrorWrapperType{})
	RegisterControlWrapper(LabelWrapper, &LabelWrapperType{})
	RegisterControlWrapper(DivWrapper, &DivWrapperType{})
}
