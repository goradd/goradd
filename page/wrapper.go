package page

import (
	"bytes"
	"context"
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
	CopyI() WrapperI
}

var wrapperRegistry = map[string]WrapperI{}

func RegisterControlWrapper(name string, w WrapperI) {
	wrapperRegistry[name] = w
}

// NewWrapper returns a newly allocated wrapper from the wrapper registry.
func NewWrapper(name string) WrapperI {
	if w, ok := wrapperRegistry[name]; ok {
		return w.CopyI()
	} else {
		panic ("Unkown wrapper " + name)
	}
	return nil
}

type ErrorWrapperType struct {
}

func NewErrorWrapper() ErrorWrapperType {
	return ErrorWrapperType{}
}

// Copy copies itself and returns it. This is used when a new wrapper is created from a named type.
func (w ErrorWrapperType) CopyI() WrapperI {
	return w // Since we are not a pointer type, a copy was sent in
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

func (w LabelWrapperType) Copy() LabelWrapperType {
	w.labelAttributes = w.labelAttributes.Copy()
	return w
}

func (w LabelWrapperType) CopyI() WrapperI {
	w.Copy()
	return w
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

func (w DivWrapperType) CopyI() WrapperI {
	return w
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
