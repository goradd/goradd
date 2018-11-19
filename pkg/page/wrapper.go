package page

import (
	"bytes"
	"context"
	"encoding/gob"
	"github.com/spekary/goradd/pkg/html"
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
	SetValidationMessageChanged()
	SetValidationStateChanged()
	AjaxRender(ctx context.Context, response *Response, c ControlI)
}

var wrapperRegistry = map[string]WrapperI{}

func RegisterControlWrapper(name string, w WrapperI) {
	wrapperRegistry[name] = w
	gob.RegisterName(name, w)
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
	ValidationMessageChanged bool
	ValidationStateChanged   bool
	//instructionsChanged bool // do this with a complete redraw. This won't change often.
}

func NewErrorWrapper() *ErrorWrapperType {
	return &ErrorWrapperType{}
}

// Copy copies itself and returns it. This is used when a new wrapper is created from a named type.
func (w *ErrorWrapperType) CopyI() WrapperI {
	return w // Since we are not a pointer type, a copy was sent in
}

func (w *ErrorWrapperType) Wrap(ctx context.Context, ctrl ControlI, html string, buf *bytes.Buffer) {
	ErrorTmpl(ctx, ctrl, html, buf)
}

func (w *ErrorWrapperType) TypeName() string {
	return ErrorWrapper
}

// ModifyDrawingAttributes should only be called by the framework during a draw.
// It changes attributes of the wrapped control based on the validation state of the control.
func (w *ErrorWrapperType) ModifyDrawingAttributes(c ControlI, a *html.Attributes) {
	var describedBy string
	state := c.control().validationState
	if state != ValidationNever {
		describedBy = c.ID() + "_err"
	}
	if c.control().instructions != "" {
		if describedBy != "" {
			describedBy += " "
		}
		describedBy += c.ID() + "_inst"
	}
	if describedBy != "" {
		a.Set("aria-describedby", describedBy)
	}

	// has the side effect of resetting the validation state since we know the control is being completely redrawn
	// instead of ajax drawn
	w.ValidationMessageChanged = false
	w.ValidationStateChanged = false

	if w.ValidationStateChanged {
		switch c.control().validationState {
		case ValidationWaiting:fallthrough
		case ValidationValid:
			c.WrapperAttributes().RemoveClass("error")
		case ValidationInvalid:
			c.WrapperAttributes().AddClass("error")
		}
	}
}

// The following functions enable wrappers to only send changes during the refresh of a control, rather than drawing the
// whole control.

func (w *ErrorWrapperType) SetValidationMessageChanged() {
	w.ValidationMessageChanged = true
}

func (w *ErrorWrapperType) SetValidationStateChanged() {
	w.ValidationStateChanged = true
}


// Called by the framework to draw any changes to the wrapper that we have recorded.
// This has to work closely with the wrapper template so that it would create the same effect as if that
// entire control had been redrawn
func (w *ErrorWrapperType) AjaxRender(ctx context.Context, response *Response, c ControlI) {
	if w.ValidationMessageChanged {
		response.ExecuteControlCommand(c.ID() + "_err", "text", c.ValidationMessage())
		w.ValidationMessageChanged = false
	}

	if w.ValidationStateChanged {
		switch c.control().validationState {
		case ValidationWaiting:fallthrough
		case ValidationValid:
			response.ExecuteControlCommand(c.ID() + "_ctl", "removeClass", "error")
		case ValidationInvalid:
			response.ExecuteControlCommand(c.ID() + "_ctl", "addClass", "error")
		}
		w.ValidationStateChanged = false
	}
}


type LabelWrapperType struct {
	ErrorWrapperType
	ΩlabelAttr *html.Attributes
}

func NewLabelWrapper() *LabelWrapperType {
	return &LabelWrapperType{}
}

func (w *LabelWrapperType) Copy() *LabelWrapperType {
	w.ΩlabelAttr = w.ΩlabelAttr.Copy()
	return w
}

func (w *LabelWrapperType) CopyI() WrapperI {
	w.Copy()
	return w
}


func (w *LabelWrapperType) Wrap(ctx context.Context, ctrl ControlI, html string, buf *bytes.Buffer) {
	LabelTmpl(ctx, w, ctrl, html, buf)
}

// LabelAttributes returns attributes that will apply to the label. Changes will be remembered, but will not
// be applied unless you redraw the control.
func (w *LabelWrapperType) LabelAttributes() *html.Attributes {
	if w.ΩlabelAttr == nil {
		w.ΩlabelAttr = html.NewAttributes()
	}
	return w.ΩlabelAttr
}

func (w *LabelWrapperType) HasLabelAttributes() bool {
	if w.ΩlabelAttr == nil || w.ΩlabelAttr.Len() == 0 {
		return false
	}
	return true
}


func (w *LabelWrapperType) TypeName() string {
	return LabelWrapper
}

func (w *LabelWrapperType) ModifyDrawingAttributes(c ControlI, a *html.Attributes) {
	w.ErrorWrapperType.ModifyDrawingAttributes(c, a)
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

func (w DivWrapperType) SetValidationMessageChanged() {
}

func (w DivWrapperType) SetValidationStateChanged() {
}

func (w DivWrapperType) AjaxRender(ctx context.Context, response *Response, c ControlI) {
}

func init() {
	RegisterControlWrapper(ErrorWrapper, &ErrorWrapperType{})
	RegisterControlWrapper(LabelWrapper, &LabelWrapperType{})
	RegisterControlWrapper(DivWrapper, &DivWrapperType{})
}
