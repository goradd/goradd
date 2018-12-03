package control

import (
	"github.com/spekary/goradd/pkg/html"
	"github.com/spekary/goradd/pkg/page"
)

type SpanI interface {
	PanelI
}

// Span is a Goradd control that is a basic "span" wrapper. Use it to style and listen to events on a span. It
// can also be used as the basis for more advanced javascript controls.
type Span struct {
	Panel
}

func NewSpan(parent page.ControlI, id string) *Span {
	p := &Span{}
	p.Init(p, parent, id)
	return p
}

func (c *Span) Init(self SpanI, parent page.ControlI, id string) {
	c.Panel.Init(self, parent, id)
	c.Tag = "span"
}

func (c *Span) DrawingAttributes() *html.Attributes {
	a := c.Control.DrawingAttributes()
	a.SetDataAttribute("grctl", "span")
	return a
}
