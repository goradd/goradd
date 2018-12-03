package control

import (
	"github.com/spekary/goradd/pkg/html"
	"github.com/spekary/goradd/pkg/page"
)

// BootstrapControl is a mixin to help a control become a bootstrap control
type BootstrapControl struct {
	IsFormControl bool
}

func (b *BootstrapControl) SetIsFormControl(is bool) {
	b.IsFormControl = is
}


// Adds the special bootstrap attributes to controls that should be drawn with the form-control attribute
// Should be called during DrawingAttributes
func (b *BootstrapControl) AddBootstrapAttributes(c page.ControlI, attr *html.Attributes) {
	if b.IsFormControl {
		switch c.ValidationState() {
		case page.ValidationValid:
			attr.AddClass("is-valid")
		case page.ValidationInvalid:
			attr.AddClass("is-invalid")
		}
		attr.AddClass("form-control")
	}
}

func (c *BootstrapControl) Serialize(e page.Encoder) (err error) {
	return e.Encode(c.IsFormControl)
}

func (c *BootstrapControl) Deserialize(d page.Decoder, p *page.Page) (err error) {
	if err = d.Decode(&c.IsFormControl); err != nil {
		return
	}

	return
}
