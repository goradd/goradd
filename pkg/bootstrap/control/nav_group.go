package control

import (
	"context"
	"github.com/goradd/goradd/pkg/page"
	grctl "github.com/goradd/goradd/pkg/page/control"
	"github.com/goradd/html5tag"
)

type NavGroupI interface {
	grctl.PanelI
}

// NavGroup is a simple div with a "navbar-nav" class that is meant to be a container of NavLink
// items.
type NavGroup struct {
	grctl.Panel
}

// NewNavGroup creates a new NavGroup. NavGroups should only be used as child items of a NavBar.
func NewNavGroup(parent page.ControlI, id string) *NavGroup {
	g := new(NavGroup)
	g.Self = g
	g.Init(parent, id)
	return g
}

// Init initializes the button
func (g *NavGroup) Init(parent page.ControlI, id string) {
	g.Panel.Init(parent, id)
}

func (g *NavGroup) this() NavGroupI {
	return g.Self.(NavGroupI)
}

// DrawingAttributes returns the attributes to add to the tag just before the button is drawn.
func (g *NavGroup) DrawingAttributes(ctx context.Context) html5tag.Attributes {
	a := g.Panel.DrawingAttributes(ctx)
	a.AddClass("navbar-nav")
	return a
}

// Serialize serializes the state of the control for the pagestate
func (g *NavGroup) Serialize(e page.Encoder) {
	g.Panel.Serialize(e)
	return
}

// Deserialize reconstructs the control from the page state.
func (g *NavGroup) Deserialize(d page.Decoder) {
	g.Panel.Deserialize(d)
}

// NavGroupCreator is the initialization structure for the declarative creation of the control.
type NavGroupCreator struct {
	// ID is the control id
	ID string
	// Children is a list of creators to use to create the child controls of the panel.
	// You can wrap your child creators with the Children() function as a helper.
	Children []page.Creator

	page.ControlOptions
}

// Create is called by the framework to create a new control from the Creator. You
// do not normally need to call this.
func (c NavGroupCreator) Create(ctx context.Context, parent page.ControlI) page.ControlI {
	ctrl := NewNavGroup(parent, c.ID)
	c.Init(ctx, ctrl)
	return ctrl
}

// Init is called by implementations of NavGroups to initialize a control with the
// creator. You do not normally need to call this.
func (c NavGroupCreator) Init(ctx context.Context, ctrl NavGroupI) {
	sub := grctl.PanelCreator{
		ID:             c.ID,
		Children:       c.Children,
		ControlOptions: c.ControlOptions,
	}
	sub.Init(ctx, ctrl)
}

// GetNavGroup is a convenience method to return the button with the given id from the page.
func GetNavGroup(c page.ControlI, id string) *NavGroup {
	return c.Page().GetControl(id).(*NavGroup)
}

func init() {
	page.RegisterControl(&NavGroup{})
}
