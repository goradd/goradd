package page

import (
	"bytes"
	"context"
	"reflect"
	"github.com/spekary/goradd/html"
)

// WrapperI defines the control wrapper interface. A control wrapper takes the basic html output by a control and wraps
// it in additional html to give it context. See the 2 built-in wrappers: LabelWrapper and ErrorWrapper for examples.
// For example, wrappers can be used to add labels that are connected to a control, give additional information, or show error conditions.
type WrapperI interface {
	Wrap(ctx context.Context, ctrl ControlI, html string, buf *bytes.Buffer)
	ModifyDrawingAttributes(ctrl ControlI, attributes *html.Attributes)
}

var wrapperRegistry = map[string]reflect.Type{}

func RegisterControlWrapper(name string, w WrapperI) {
	t := reflect.TypeOf(w)
	wrapperRegistry[name] = t
}

func GetRegisteredWrapper(name string) WrapperI {
	if w, ok := wrapperRegistry[name]; ok {
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

func (w ErrorWrapper) ModifyDrawingAttributes(c ControlI, a *html.Attributes) {
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


type LabelWrapper struct {
	ErrorWrapper
	labelAttributes *html.Attributes
}

func NewLabelWrapper() LabelWrapper {
	return LabelWrapper{}
}

func (w LabelWrapper) Wrap(ctx context.Context, ctrl ControlI, html string, buf *bytes.Buffer) {
	LabelTmpl(ctx, w, ctrl, html, buf)
}

// LabelAttributes returns attributes that will apply to the label. Changes will be remembered.
func (w *LabelWrapper) LabelAttributes() *html.Attributes {
	if w.labelAttributes == nil {
		w.labelAttributes = html.NewAttributes()
	}
	return w.labelAttributes
}

func (w LabelWrapper) HasLabelAttributes() bool {
	if w.labelAttributes == nil || w.labelAttributes.Len() == 0 {
		return false
	}
	return true
}


func (w LabelWrapper) TypeName() string {
	return "page.Label"
}

func (w LabelWrapper) ModifyDrawingAttributes(c ControlI, a *html.Attributes) {
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


type DivWrapper struct {
}

func NewDivWrapper() DivWrapper {
	return DivWrapper{}
}

func (w DivWrapper) Wrap(ctx context.Context, ctrl ControlI, html string, buf *bytes.Buffer) {
	DivTmpl(ctx, ctrl, html, buf)
}

func (w DivWrapper) TypeName() string {
	return "page.Div"
}

func (w DivWrapper) ModifyDrawingAttributes(ctrl ControlI, a *html.Attributes) {
}



func init() {
	RegisterControlWrapper("page.Error", &ErrorWrapper{})
	RegisterControlWrapper("page.Label", &LabelWrapper{})
	RegisterControlWrapper("page.Div", &DivWrapper{})
}
