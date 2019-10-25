package control

import (
	"context"
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

// ΩDrawingAttributes is called by the framework to get the temporary attributes that are specifically set by
// this control just before drawing.
func (c *Canvas) ΩDrawingAttributes(ctx context.Context) html.Attributes {
	a := c.Control.ΩDrawingAttributes(ctx)
	a.SetDataAttribute("grctl", "canvas")
	return a
}

// CanvasCreator is the initialization structure for declarative creation of buttons
type CanvasCreator struct {
	// ID is the control id
	ID string
	page.ControlOptions
}

// Create is called by the framework to create a new control from the Creator. You
// do not normally need to call this.
func (c CanvasCreator) Create(ctx context.Context, parent page.ControlI) page.ControlI {
	ctrl := NewCanvas(parent, c.ID)
	ctrl.ApplyOptions(ctx, c.ControlOptions)
	return ctrl
}

// GetCanvas is a convenience method to return the canvas with the given id from the page.
func GetCanvas(c page.ControlI, id string) *Canvas {
	return c.Page().GetControl(id).(*Canvas)
}

func init() {
	page.RegisterControl(Canvas{})
}