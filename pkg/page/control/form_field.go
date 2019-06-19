package control

import (
	"github.com/goradd/goradd/pkg/html"
	"github.com/goradd/goradd/pkg/page"
)

type FormFieldI interface {
	page.ControlI
}

// FormField is a Goradd control that wraps other controls, and provides common companion
// functionality like a form label, validation state display, and help text.
type FormField struct {
	page.Control

	// textLabelMode describes how to draw the internal label
	textLabelMode html.LabelDrawingMode
	// isInline is true to use a span for the wrapper, false for a div
	isInline bool
	// hasFor tells us if we should draw a for attribute in the label tag. This is helpful for screen readers and navigation on certain kinds of tags.
	hasFor bool
	// instructions is text associated with the control for extra explanation. You could also try adding a tooltip to the wrapper.
	instructions string
}

func NewFormField(parent page.ControlI, id string) *FormField {
	p := &FormField{}
	p.Init(p, parent, id)
	return p
}

func (c *FormField) Init(self PanelI, parent page.ControlI, id string) {
	c.Control.Init(self, parent, id)
	c.Tag = "div"
}

func (c *FormField) ΩDrawingAttributes() *html.Attributes {
	a := c.Control.ΩDrawingAttributes()
	a.SetDataAttribute("grctl", "panel")
	return a
}

