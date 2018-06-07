package control

import (
	"bytes"
	"context"
	"github.com/spekary/goradd/html"
	"github.com/spekary/goradd/page"
)

type FormGroupWrapper struct {
	innerDivAttributes *html.Attributes
	labelAttributes    *html.Attributes
	useTooltips        bool // uses tooltips for the error class
}

func NewFormGroupWrapper() *FormGroupWrapper {
	return &FormGroupWrapper{}
}

func (w *FormGroupWrapper) Wrap(ctx context.Context, ctrl page.ControlI, html string, buf *bytes.Buffer) {
	FormGroupTmpl(ctx, w, ctrl, html, buf)
}

func (w FormGroupWrapper) TypeName() string {
	return "bootstrap.FormGroup"
}

// InnerDivAttributes returns attributes for the innerDiv. Changes will be remembered. If you set these, the control
// itself will be wrapped with a div with these attributes. This is useful for layouts that have the label next to
// the control.
func (w *FormGroupWrapper) InnerDivAttributes() *html.Attributes {
	if w.innerDivAttributes == nil {
		w.innerDivAttributes = html.NewAttributes()
	}
	return w.innerDivAttributes
}

func (w FormGroupWrapper) HasInnerDivAttributes() bool {
	if w.innerDivAttributes == nil || w.innerDivAttributes.Len() == 0 {
		return false
	}
	return true
}

// LabelAttributes returns attributes that will apply to the label. Changes will be remembered.
func (w *FormGroupWrapper) LabelAttributes() *html.Attributes {
	if w.labelAttributes == nil {
		w.labelAttributes = html.NewAttributes()
	}
	return w.labelAttributes
}

func (w FormGroupWrapper) HasLabelAttributes() bool {
	if w.labelAttributes == nil || w.labelAttributes.Len() == 0 {
		return false
	}
	return true
}

func (w *FormGroupWrapper) SetUseTooltips(t bool) *FormGroupWrapper {
	w.useTooltips = t
	return w
}

func init() {
	page.RegisterControlWrapper("bootstrap.FormGroup", &FormGroupWrapper{})
}

// TODO: will need to serialize this when we are ready to serialize formstate
