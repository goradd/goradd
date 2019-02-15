package control

import (
	"github.com/goradd/goradd/pkg/html"
	"github.com/goradd/goradd/pkg/page"
)


type CanvasI interface {
	page.ControlI
}

// Canvas is a Goradd control that is an html canvas control. It currently does not have any primitives
// to draw on the canvas, and is here primarily to create a canvas that you would draw on using JavaScript.
type Canvas struct {
	page.Control
}

// NewCanvas creates a Canvas control
func NewCanvas(parent page.ControlI, id string) *Canvas {
	p := &Canvas{}
	p.Init(p, parent, id)
	return p
}

// Init is called by subcontrols. You do not normally need to call it.
func (c *Canvas) Init(self CanvasI, parent page.ControlI, id string) {
	c.Control.Init(self, parent, id)
	c.Tag = "canvas"
}

// DrawingAttributes is called by the framework to get the temporary attributes that are specifically set by
// this control just before drawing.
func (c *Canvas) DrawingAttributes() *html.Attributes {
	a := c.Control.DrawingAttributes()
	a.SetDataAttribute("grctl", "canvas")
	return a
}
