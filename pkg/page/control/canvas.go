package control

import (
	"github.com/goradd/goradd/pkg/html"
	"github.com/goradd/goradd/pkg/page"
)


type CanvasI interface {
	page.ControlI
}

// Canvas is a Goradd control that is a basic "div" wrapper. Use it to style and listen to events on a div. It
// can also be used as the basis for more advanced javascript controls.
type Canvas struct {
	page.Control
}

func NewCanvas(parent page.ControlI, id string) *Canvas {
	p := &Canvas{}
	p.Init(p, parent, id)
	return p
}

func (c *Canvas) Init(self CanvasI, parent page.ControlI, id string) {
	c.Control.Init(self, parent, id)
	c.Tag = "canvas"
}


func (c *Canvas) DrawingAttributes() *html.Attributes {
	a := c.Control.DrawingAttributes()
	a.SetDataAttribute("grctl", "canvas")
	return a
}
