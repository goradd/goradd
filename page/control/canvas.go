package control

import (
	"github.com/spekary/goradd/html"
	"github.com/spekary/goradd/page"
)


type CanvasI interface {
	page.ControlI
}

// Canvas is a Goradd control that is a basic "div" wrapper. Use it to style and listen to events on a div. It
// can also be used as the basis for more advanced javascript controls.
type Canvas struct {
	page.Control
}

func NewCanvas(parent page.ControlI) *Canvas {
	p := &Canvas{}
	p.Init(p, parent)
	return p
}

func (c *Canvas) Init(self CanvasI, parent page.ControlI) {
	c.Control.Init(self, parent)
	c.Tag = "canvas"
}


func (c *Canvas) DrawingAttributes() *html.Attributes {
	a := c.Control.DrawingAttributes()
	a.SetDataAttribute("grctl", "canvas")
	return a
}
