package page

import (
	"github.com/spekary/goradd/pkg/html"
	page2 "github.com/spekary/goradd/pkg/page"
)


type BootstrapControl struct {
	page2.Control

	IsFormControl bool
}

func (c *BootstrapControl) SetIsFormControl(is bool) *BootstrapControl {
	c.IsFormControl = is
	return c
}

func (c *BootstrapControl) DrawingAttributes() *html.Attributes {
	a := c.Control.DrawingAttributes()

	if c.IsFormControl {
		c.FormControlAddAttributes(a)
	}

	return a
}

func (c *BootstrapControl) FormControlAddAttributes(attr *html.Attributes) {
	switch c.ValidationState() {
	case page2.ValidationValid:
		attr.AddClass("is-valid")
	case page2.ValidationInvalid:
		attr.AddClass("is-invalid")
	}
	attr.AddClass("form-control")
}

/*
func (c *BootstrapControl) Serialize(e page2.Encoder) (err error) {
	if err = c.Control.Serialize(e); err != nil {
		return
	}

	err = e.Serialize(c.IsFormControl)
	return
}

func (c *BootstrapControl) Deserialize(d page2.Decoder, p *page2.Page) (err error) {
	if err = c.Control.Deserialize(d, p); err != nil {
		return
	}

	err = d.Deserialize(&c.IsFormControl)
	return
}
*/