package control

import (
	"github.com/spekary/goradd/html"
	"github.com/spekary/goradd/page"
	"github.com/spekary/goradd/page/control/control_base"
)

// Panel is a Goradd control that is a basic "div" wrapper. Use it to style and listen to events on a div. It
// can also be used as the basis for more advanced javascript controls.
type Span struct {
	control_base.Panel
}

func NewSpan(parent page.ControlI) *Span {
	p := &Span{}
	p.Tag = "span"
	p.Init(p, parent)
	return p
}

func (c *Span) DrawingAttributes() *html.Attributes {
	a := c.Control.DrawingAttributes()
	a.SetDataAttribute("grctl", "span")
	return a
}
