package control

import (
	"bytes"
	"context"
	"github.com/spekary/goradd/html"
	"github.com/spekary/goradd/page"
)

const (
	DivWrapper = "bootstrap.Div"
	FormGroupWrapper = "bootstrap.FormGroup"
	FieldsetWrapper = "bootstrap.Fieldset"
)

// DivWrapperType is a wrapper similar to a form group, but simply without the FormGroup class added. Use this for
// wrapping inline elements and other special situations listed in the Bootstrap doc under the Forms component.
// https://getbootstrap.com/docs/4.1/components/forms/ as of this writing
type DivWrapperType struct {
	page.LabelWrapperType
	innerDivAttributes *html.Attributes
	useTooltips        bool // uses tooltips for the error class
}

func NewDivWrapper() *DivWrapperType {
	return &DivWrapperType{}
}

func (w *DivWrapperType) Copy()  *DivWrapperType {
	wNew := &DivWrapperType{}
	wNew.LabelWrapperType = w.LabelWrapperType.Copy()
	wNew.innerDivAttributes = w.innerDivAttributes.Copy()
	wNew.useTooltips = w.useTooltips
	return wNew
}

func (w *DivWrapperType) CopyI() page.WrapperI {
	return w.Copy()
}

func (w *DivWrapperType) Wrap(ctx context.Context, ctrl page.ControlI, html string, buf *bytes.Buffer) {
	FormGroupTmpl(ctx, w, ctrl, html, buf)
}

func (w DivWrapperType) TypeName() string {
	return DivWrapper
}

// InnerDivAttributes returns attributes for the innerDiv. Changes will be remembered. If you set these, the control
// itself will be wrapped with a div with these attributes. This is useful for layouts that have the label next to
// the control.
func (w *DivWrapperType) InnerDivAttributes() *html.Attributes {
	if w.innerDivAttributes == nil {
		w.innerDivAttributes = html.NewAttributes()
	}
	return w.innerDivAttributes
}

func (w *DivWrapperType) HasInnerDivAttributes() bool {
	if w.innerDivAttributes == nil || w.innerDivAttributes.Len() == 0 {
		return false
	}
	return true
}

func (w *DivWrapperType) SetUseTooltips(t bool) *DivWrapperType {
	w.useTooltips = t
	return w
}

type FormGroupWrapperType struct {
	DivWrapperType
}

func NewFormGroupWrapper() *FormGroupWrapperType {
	return new(FormGroupWrapperType)
}

func (w *FormGroupWrapperType)CopyI() page.WrapperI {
	wNew := new(FormGroupWrapperType)
	wNew.DivWrapperType = *w.Copy()
	return wNew
}

func (w *FormGroupWrapperType) Wrap(ctx context.Context, ctrl page.ControlI, html string, buf *bytes.Buffer) {
	ctrl.WrapperAttributes().AddClass("form-group")
	FormGroupTmpl(ctx, &w.DivWrapperType, ctrl, html, buf)
}

func (w FormGroupWrapperType) TypeName() string {
	return FormGroupWrapper
}


type FieldsetWrapperType struct {
	page.LabelWrapperType
	useTooltips        bool // uses tooltips for the error class
}

// https://getbootstrap.com/docs/4.1/components/forms/#horizontal-form
func NewFieldsetWrapper() *FieldsetWrapperType {
	return new(FieldsetWrapperType)
}

func (w *FieldsetWrapperType) CopyI() page.WrapperI {
	wNew := NewFieldsetWrapper()
	wNew.LabelWrapperType = w.LabelWrapperType.Copy()
	wNew.useTooltips = w.useTooltips
	return w
}

func (w *FieldsetWrapperType) Wrap(ctx context.Context, ctrl page.ControlI, html string, buf *bytes.Buffer) {
	FieldsetTmpl(ctx, w, ctrl, html, buf)
}

func (w *FieldsetWrapperType) SetUseTooltips(t bool) *FieldsetWrapperType {
	w.useTooltips = t
	return w
}

func (w *FieldsetWrapperType) TypeName() string {
	return FieldsetWrapper
}


func init() {
	page.RegisterControlWrapper(DivWrapper, &DivWrapperType{})
	page.RegisterControlWrapper(FormGroupWrapper, &FormGroupWrapperType{})
	page.RegisterControlWrapper(FieldsetWrapper, &FieldsetWrapperType{})
}

