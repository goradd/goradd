package control

import (
	"bytes"
	"context"
	"github.com/goradd/goradd/pkg/html"
	"github.com/goradd/goradd/pkg/page"
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
	ΩinnerDivAttr *html.Attributes
	UseTooltips   bool // uses tooltips for the error class
}

func NewDivWrapper() *DivWrapperType {
	return &DivWrapperType{}
}

/*
func (w *DivWrapperType) Copy()  *DivWrapperType {
	wNew := &DivWrapperType{}
	wNew.LabelWrapperType = *w.LabelWrapperType.Copy()
	wNew.ΩinnerDivAttr = w.ΩinnerDivAttr.Copy()
	wNew.UseTooltips = w.UseTooltips
	return wNew
}*/

func (w *DivWrapperType) ΩNewI() page.WrapperI {
	return NewDivWrapper()
}

func (w *DivWrapperType) ΩWrap(ctx context.Context, ctrl page.ControlI, html string, buf *bytes.Buffer) {
	FormGroupTmpl(ctx, w, ctrl, html, buf)
}

func (w DivWrapperType) TypeName() string {
	return DivWrapper
}

// InnerDivAttributes returns attributes for the innerDiv. Changes will be remembered, but only drawn if you redraw the
// control. If you set these, the control
// itself will be wrapped with a div with these attributes. This is useful for layouts that have the label next to
// the control.
func (w *DivWrapperType) InnerDivAttributes() *html.Attributes {
	if w.ΩinnerDivAttr == nil {
		w.ΩinnerDivAttr = html.NewAttributes()
	}
	return w.ΩinnerDivAttr
}

// HasInnerDivAttributes returns true if the wrapper has attributes on the inner div.
func (w *DivWrapperType) HasInnerDivAttributes() bool {
	if w.ΩinnerDivAttr == nil || w.ΩinnerDivAttr.Len() == 0 {
		return false
	}
	return true
}

// SetUseTooltips sets whether to use tooltips to display validation messages.
func (w *DivWrapperType) SetUseTooltips(t bool) *DivWrapperType {
	w.UseTooltips = t
	return w
}

// ΩAjaxRender is called by the framework to draw any changes to the wrapper that we have recorded.
// This has to work closely with the wrapper template so that it would create the same effect as if that
// entire control had been redrawn
func (w *DivWrapperType) ΩAjaxRender(ctx context.Context, response *page.Response, c page.ControlI) {
	var class string
	if w.ValidationStateChanged {
		switch c.ValidationState() {
		case page.ValidationWaiting:
			response.RemoveClass(c.ID(), "is-valid")
			response.RemoveClass(c.ID(), "is-invalid")
			if w.UseTooltips {
				class = "valid-tooltip"
			} else {
				class = "valid-feedback"
			}

		case page.ValidationValid:
			response.AddClass(c.ID(), "is-valid")
			response.RemoveClass(c.ID(), "is-invalid")
			if w.UseTooltips {
				class = "valid-tooltip"
			} else {
				class = "valid-feedback"
			}

		case page.ValidationInvalid:
			response.RemoveClass(c.ID(), "is-valid")
			response.AddClass(c.ID(), "is-invalid")
			if w.UseTooltips {
				class = "invalid-tooltip"
			} else {
				class = "invalid-feedback"
			}
		}
		response.SetClass(c.ID() + "_err", class)
	}
	w.LabelWrapperType.ΩAjaxRender(ctx, response, c)
}

func (w *DivWrapperType) ΩModifyDrawingAttributes(c page.ControlI, attr *html.Attributes) {
	switch c.ValidationState() {
	case page.ValidationValid:
		attr.AddClass("is-valid")
	case page.ValidationInvalid:
		attr.AddClass("is-invalid")
	}
}


type FormGroupWrapperType struct {
	DivWrapperType
}

func NewFormGroupWrapper() *FormGroupWrapperType {
	w := new(FormGroupWrapperType)
	return w
}

func (w *FormGroupWrapperType)ΩNewI() page.WrapperI {
	return NewFormGroupWrapper()
}

func (w *FormGroupWrapperType) ΩWrap(ctx context.Context, ctrl page.ControlI, html string, buf *bytes.Buffer) {
	ctrl.WrapperAttributes().AddClass("form-group")
	FormGroupTmpl(ctx, &w.DivWrapperType, ctrl, html, buf)
}

func (w FormGroupWrapperType) TypeName() string {
	return FormGroupWrapper
}


type FieldsetWrapperType struct {
	page.LabelWrapperType
	UseTooltips bool // uses tooltips for the error class
}

// https://getbootstrap.com/docs/4.1/components/forms/#horizontal-form
func NewFieldsetWrapper() *FieldsetWrapperType {
	return new(FieldsetWrapperType)
}

func (w *FieldsetWrapperType) ΩNewI() page.WrapperI {
	return NewFieldsetWrapper()
}

func (w *FieldsetWrapperType) ΩWrap(ctx context.Context, ctrl page.ControlI, html string, buf *bytes.Buffer) {
	FieldsetTmpl(ctx, w, ctrl, html, buf)
}

func (w *FieldsetWrapperType) SetUseTooltips(t bool) *FieldsetWrapperType {
	w.UseTooltips = t
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

