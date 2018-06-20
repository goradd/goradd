package control

import (
	"bytes"
	"context"
	"github.com/spekary/goradd/html"
	"github.com/spekary/goradd/page"
)

// DivWrapper is a wrapper similar to a form group, but simply without the FormGroup class added. Use this for
// wrapping inline elements and other special situations listed in the Bootstrap doc under the Forms component.
// https://getbootstrap.com/docs/4.1/components/forms/ as of this writing
type DivWrapper struct {
	page.LabelWrapper
	innerDivAttributes *html.Attributes
	useTooltips        bool // uses tooltips for the error class
}

func NewDivWrapper() *DivWrapper {
	return &DivWrapper{}
}

func (w *DivWrapper) Wrap(ctx context.Context, ctrl page.ControlI, html string, buf *bytes.Buffer) {
	FormGroupTmpl(ctx, w, ctrl, html, buf)
}

func (w DivWrapper) TypeName() string {
	return "bootstrap.Div"
}

// InnerDivAttributes returns attributes for the innerDiv. Changes will be remembered. If you set these, the control
// itself will be wrapped with a div with these attributes. This is useful for layouts that have the label next to
// the control.
func (w *DivWrapper) InnerDivAttributes() *html.Attributes {
	if w.innerDivAttributes == nil {
		w.innerDivAttributes = html.NewAttributes()
	}
	return w.innerDivAttributes
}

func (w DivWrapper) HasInnerDivAttributes() bool {
	if w.innerDivAttributes == nil || w.innerDivAttributes.Len() == 0 {
		return false
	}
	return true
}

func (w *DivWrapper) SetUseTooltips(t bool) *DivWrapper {
	w.useTooltips = t
	return w
}

type FormGroupWrapper struct {
	DivWrapper
}

func NewFormGroupWrapper() *FormGroupWrapper {
	return &FormGroupWrapper{}
}

func (w *FormGroupWrapper) Wrap(ctx context.Context, ctrl page.ControlI, html string, buf *bytes.Buffer) {
	ctrl.WrapperAttributes().AddClass("form-group")
	FormGroupTmpl(ctx, &w.DivWrapper, ctrl, html, buf)
}

func (w FormGroupWrapper) TypeName() string {
	return "bootstrap.FormGroup"
}


type FieldsetWrapper struct {
	page.LabelWrapper
	useTooltips        bool // uses tooltips for the error class
}

// https://getbootstrap.com/docs/4.1/components/forms/#horizontal-form
func NewFieldsetWrapper() *FieldsetWrapper {
	return &FieldsetWrapper{}
}


func (w *FieldsetWrapper) Wrap(ctx context.Context, ctrl page.ControlI, html string, buf *bytes.Buffer) {
	FieldsetTmpl(ctx, w, ctrl, html, buf)
}

func (w *FieldsetWrapper) SetUseTooltips(t bool) *FieldsetWrapper {
	w.useTooltips = t
	return w
}

func (w FieldsetWrapper) TypeName() string {
	return "bootstrap.Fieldset"
}


func init() {
	page.RegisterControlWrapper("bootstrap.FormGroup", &DivWrapper{})
}

// TODO: will need to serialize this when we are ready to serialize formstate
