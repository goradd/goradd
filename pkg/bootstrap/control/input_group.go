package control

import (
	"context"
	"github.com/goradd/goradd/pkg/bootstrap/config"
	"github.com/goradd/goradd/pkg/page"
	grctl "github.com/goradd/goradd/pkg/page/control"
	"github.com/goradd/html5tag"
)

type InputGroupI interface {
	grctl.PanelI
}

type InputGroup struct {
	grctl.Panel
}

// NewInputGroup creates a new input group
func NewInputGroup(parent page.ControlI, id string) *InputGroup {
	b := new(InputGroup)
	b.Init(b, parent, id)
	return b
}

func (g *InputGroup) Init(self any, parent page.ControlI, id string) {
	g.Panel.Init(self, parent, id)
	config.LoadBootstrap(g.ParentForm())
}

func (g *InputGroup) this() InputGroupI {
	return g.Self().(InputGroupI)
}

func (g *InputGroup) DrawingAttributes(ctx context.Context) html5tag.Attributes {
	a := g.Panel.DrawingAttributes(ctx)
	a.AddClass("input-group")
	return a
}

type InputGroupCreator struct {
	// ID is the control id
	ID       string
	Children []page.Creator
	page.ControlOptions
}

// Create is called by the framework to create a new control from the Creator. You
// do not normally need to call this.
func (c InputGroupCreator) Create(ctx context.Context, parent page.ControlI) page.ControlI {
	ctrl := NewInputGroup(parent, c.ID)
	c.Init(ctx, ctrl)
	return ctrl
}

// Init is called by implementations to initialize a control with the creator.
func (c InputGroupCreator) Init(ctx context.Context, ctrl InputGroupI) {
	ctrl.ApplyOptions(ctx, c.ControlOptions)
	ctrl.AddControls(ctx, c.Children...)
}

func init() {
	page.RegisterControl(new(InputGroup))
}
