package control

import (
	"context"
	"github.com/goradd/html5tag"
	"github.com/goradd/goradd/pkg/page"
)

type CanvasI interface {
	page.ControlI
}

// Canvas is a Goradd control that is an html canvas control. It currently does not have any primitives
// to draw on the canvas, and is here primarily to create a canvas that you would draw on using JavaScript.
type Canvas struct {
	page.ControlBase
}

// NewCanvas creates a Canvas control
func NewCanvas(parent page.ControlI, id string) *Canvas {
	p := &Canvas{}
	p.Self = p
	p.Init(parent, id)
	return p
}

// Init is called by subcontrols. You do not normally need to call it.
func (c *Canvas) Init(parent page.ControlI, id string) {
	c.ControlBase.Init(parent, id)
	c.Tag = "canvas"
}

// DrawingAttributes is called by the framework to get the temporary attributes that are specifically set by
// this control just before drawing.
func (c *Canvas) DrawingAttributes(ctx context.Context) html5tag.Attributes {
	a := c.ControlBase.DrawingAttributes(ctx)
	a.SetData("grctl", "canvas")
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
	page.RegisterControl(&Canvas{})
}