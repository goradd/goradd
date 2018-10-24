package page

import (
	"github.com/spekary/goradd/html"
	page2 "github.com/spekary/goradd/page"
)


type BootstrapControl struct {
	page2.Control

	isFormControl bool
}

func (c *BootstrapControl) SetIsFormControl(is bool) {
	c.isFormControl = is
}

func (c *BootstrapControl) DrawingAttributes() *html.Attributes {
	a := c.Control.DrawingAttributes()

	if c.isFormControl {
		c.FormControlAddAttributes(a)
	}

	return a
}

func (c *BootstrapControl) FormControlAddAttributes(attr *html.Attributes) {
	switch c.ValidationState() {
	case page2.Valid:
		attr.AddClass("is-valid")
	case page2.Invalid:
		attr.AddClass("is-invalid")
	}
	attr.AddClass("form-control")
}
