package page

import (
	"bytes"
	"context"
	"encoding/gob"
	"github.com/goradd/goradd/pkg/html"
)

// Wrapper constants used in the With function
const (
	ErrorWrapper = "page.Error"
	LabelWrapper = "page.Label"
	DivWrapper   = "page.Div"
)

// WrapperI defines the control wrapper interface. Generally you will not call any of these functions.
// The interface is used by the framework to control wrapper drawing. To call wrapper specific functions, cast
// to the specific wrapper type.
type WrapperI interface {
	// ΩWrap is used by the framework to wrap the control and html with the wrapper's html tags.
	ΩWrap(ctx context.Context, ctrl ControlI, html string, buf *bytes.Buffer)
	// ΩNewI returns a new wrapper. It is used by the wrapper registry to return a named wrapper.
	ΩNewI() WrapperI
	// ΩSetValidationMessageChanged notifies the wrapper that the validation message has changed.
	ΩSetValidationMessageChanged()
	// ΩSetValidationStateChanged notifies the wrapper that the validation state has changed.
	ΩSetValidationStateChanged()
	// ΩAjaxRender does an ajax render of the wrapper.
	ΩAjaxRender(ctx context.Context, response *Response, c ControlI)
	// TypeName returns the named type of the wrapper.
	TypeName() string
	// ΩModifyDrawingAttributes is used by the framework to allow the wrapper to modify the attributes of the control and draw time.
	ΩModifyDrawingAttributes(ctrl ControlI, attributes *html.Attributes)
}

var wrapperRegistry = map[string]WrapperI{}

// RegisterControlWrapper registers a wrapper with goradd so that it can be called by name using
// the Control.With() function. It should be called at init time.
func RegisterControlWrapper(name string, w WrapperI) {
	wrapperRegistry[name] = w
	gob.RegisterName(name, w)
}

// NewWrapper returns a newly allocated wrapper from the wrapper registry.
func NewWrapper(name string) WrapperI {
	if w, ok := wrapperRegistry[name]; ok {
		return w.ΩNewI()
	} else {
		panic("Unkown wrapper " + name)
	}
	return nil
}

// An ErrorWrapperType wraps a control with a validation message
type ErrorWrapperType struct {
	ValidationMessageChanged bool
	ValidationStateChanged   bool
	//instructionsChanged bool // do this with a complete redraw. This won't change often.
}

// NewErrorWrapper creates a new ErrorWrapperType.
func NewErrorWrapper() *ErrorWrapperType {
	return &ErrorWrapperType{}
}

// ΩNewI creates a new wrapper. This is used when a new wrapper is created from a named type.
func (w ErrorWrapperType) ΩNewI() WrapperI {
	return NewErrorWrapper()
}

// ΩWrap wraps the given control with an ErrorWrapperTemplate. The ErrorWrapperTemplate adds a validation message
// to a control's html, and also an instructions message.
func (w *ErrorWrapperType) ΩWrap(ctx context.Context, ctrl ControlI, html string, buf *bytes.Buffer) {
	ErrorTmpl(ctx, ctrl, html, buf)
}

// TypeName returns the name of the type of wrapper.
func (w *ErrorWrapperType) TypeName() string {
	return ErrorWrapper
}

// ΩModifyDrawingAttributes should only be called by the framework during a draw.
// It changes attributes of the wrapped control based on the validation state of the control.
func (w *ErrorWrapperType) ΩModifyDrawingAttributes(c ControlI, a *html.Attributes) {
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

	// Reset the validation state since we know the control is being completely redrawn
	// instead of ajax drawn

	switch state {
	case ValidationWaiting:
		fallthrough
	case ValidationValid:
		c.WrapperAttributes().RemoveClass("error")
	case ValidationInvalid:
		c.WrapperAttributes().AddClass("error")
	}

	w.ValidationMessageChanged = false
	w.ValidationStateChanged = false
}

// The following functions enable wrappers to only send changes during the refresh of a control, rather than drawing the
// whole control.

// ΩSetValidationMessageChanged is called by the framework to indicate the validation message changed.
func (w *ErrorWrapperType) ΩSetValidationMessageChanged() {
	w.ValidationMessageChanged = true
}

// ΩSetValidationStateChanged is called by the framework to indicate the validation state changed.
func (w *ErrorWrapperType) ΩSetValidationStateChanged() {
	w.ValidationStateChanged = true
}

// ΩAjaxRender is called by the framework to draw any changes to the wrapper that we have recorded.
// This has to work closely with the wrapper template so that it would create the same effect as if that
// entire control had been redrawn
func (w *ErrorWrapperType) ΩAjaxRender(ctx context.Context, response *Response, c ControlI) {
	if w.ValidationMessageChanged {
		response.ExecuteControlCommand(c.ID()+"_err", "text", c.ValidationMessage())
		w.ValidationMessageChanged = false
	}

	if w.ValidationStateChanged {
		switch c.control().validationState {
		case ValidationWaiting:
			fallthrough
		case ValidationValid:
			response.RemoveClass(c.ID() + "_ctl", "error")
		case ValidationInvalid:
			response.AddClass(c.ID() + "_ctl", "error")
		}
		w.ValidationStateChanged = false
	}
}

// LabelWrapperType is a wrapper that associates a label tag with the control, as well
// as instuctions and an error message.
type LabelWrapperType struct {
	ErrorWrapperType
	labelAttr *html.Attributes
}

// NewLabelWrapper returns a new LabelWrapperType.
func NewLabelWrapper() *LabelWrapperType {
	return &LabelWrapperType{}
}

// Copy returns a copy of the wrapper
func (w *LabelWrapperType) Copy() *LabelWrapperType {
	return &LabelWrapperType{w.ErrorWrapperType, w.labelAttr.Copy()}
}

// ΩNewI is used when a new wrapper is created from a named type.
func (w *LabelWrapperType) ΩNewI() WrapperI {
	return NewLabelWrapper()
}

// ΩWrap wraps the given control and html output from that control with html from the LabelWrapperTemplate.
// The LabelWrapperTemplate will associate a label tag with the control, and also add tags to display validation
// errors and information or instructions associated with the control. See the template_source/wrapper_label.tpl.got
// source file for details on how this is done.
func (w *LabelWrapperType) ΩWrap(ctx context.Context, ctrl ControlI, html string, buf *bytes.Buffer) {
	LabelTmpl(ctx, w, ctrl, html, buf)
}

// LabelAttributes returns attributes that will apply to the label. Changes will be remembered, but will not
// be applied unless you redraw the control.
func (w *LabelWrapperType) LabelAttributes() *html.Attributes {
	if w.labelAttr == nil {
		w.labelAttr = html.NewAttributes()
	}
	return w.labelAttr
}

// HasLabelAttributes returns true if attributes are defined on the wrapper itself
func (w *LabelWrapperType) HasLabelAttributes() bool {
	if w.labelAttr == nil || w.labelAttr.Len() == 0 {
		return false
	}
	return true
}

// TypeName returns the name used to identify a LabelWrapper.
func (w *LabelWrapperType) TypeName() string {
	return LabelWrapper
}

// ΩModifyDrawingAttributes is a framework function that allows wrappers to modify a control's attributes at draw time.
// Label wrappers will set the aria-labeledby attribute in the control if needed.
func (w *LabelWrapperType) ΩModifyDrawingAttributes(c ControlI, a *html.Attributes) {
	w.ErrorWrapperType.ΩModifyDrawingAttributes(c, a)
	if c.control().label != "" {
		a.AddAttributeValue("aria-labelledby", c.ID()+"_lbl")
	}
}

// DivWrapperType is a wrapper that only wraps the control in a div tag. No instructions, error message
// or label is included.
type DivWrapperType struct {
}

func NewDivWrapper() *DivWrapperType {
	return &DivWrapperType{}
}

func (w DivWrapperType) ΩNewI() WrapperI {
	return NewDivWrapper()
}

func (w DivWrapperType) ΩWrap(ctx context.Context, ctrl ControlI, html string, buf *bytes.Buffer) {
	DivTmpl(ctx, ctrl, html, buf)
}

func (w DivWrapperType) TypeName() string {
	return DivWrapper
}

func (w DivWrapperType) ΩModifyDrawingAttributes(ctrl ControlI, a *html.Attributes) {
}

func (w DivWrapperType) ΩSetValidationMessageChanged() {
}

func (w DivWrapperType) ΩSetValidationStateChanged() {
}

func (w DivWrapperType) ΩAjaxRender(ctx context.Context, response *Response, c ControlI) {
}

func init() {
	RegisterControlWrapper(ErrorWrapper, &ErrorWrapperType{})
	RegisterControlWrapper(LabelWrapper, &LabelWrapperType{})
	RegisterControlWrapper(DivWrapper, &DivWrapperType{})
}
