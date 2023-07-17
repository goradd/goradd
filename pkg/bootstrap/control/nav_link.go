package control

import (
	"context"
	"github.com/goradd/goradd/pkg/page"
	grctl "github.com/goradd/goradd/pkg/page/control"
	"github.com/goradd/html5tag"
)

type NavLinkI interface {
	grctl.ActiveLinkI
}

// NavLink creates an anchor tag with a nav-link class. It is used to create a link in a NavBar, and
// generally should be a child item of a NavGroup.
type NavLink struct {
	grctl.ActiveLink
}

// NewNavLink creates a new NavLink.
func NewNavLink(parent page.ControlI, id string) *NavLink {
	l := new(NavLink)
	l.Init(l, parent, id)
	return l
}

// Init initializes the button
func (l *NavLink) Init(self any, parent page.ControlI, id string) {
	l.ActiveLink.Init(self, parent, id)
	l.ActiveAttributes().
		AddClass("active").
		Set("aria-current", "page")
}

func (l *NavLink) this() NavLinkI {
	return l.Self().(NavLinkI)
}

// DrawingAttributes returns the attributes to add to the tag just before the button is drawn.
func (l *NavLink) DrawingAttributes(ctx context.Context) html5tag.Attributes {
	a := l.ActiveLink.DrawingAttributes(ctx)
	a.AddClass("nav-link")
	a.SetData("grctl", "bs-navlink")
	return a
}

// NavLinkCreator is the initialization structure for the declarative creation of the control.
type NavLinkCreator struct {
	// ID is the control id. This is also the eventValue sent by the enclosing Navbar.
	ID string
	// Text is the text displayed in the link
	Text string
	// Location is the content of the href, which is the url where the link will go.
	// Alternatively, add an On handler to the ControlOptions
	Location string
	// ActiveAttributes are additional attributes to add to the active link. By default, the ActiveAttributes
	// will add the "active" class and aria-current=page.
	ActiveAttributes html5tag.Attributes
	page.ControlOptions
}

// Create is called by the framework to create a new control from the Creator. You
// do not normally need to call this.
func (c NavLinkCreator) Create(ctx context.Context, parent page.ControlI) page.ControlI {
	ctrl := NewNavLink(parent, c.ID)
	c.Init(ctx, ctrl)
	return ctrl
}

// Init is called by implementations of NavLinks to initialize a control with the
// creator.
func (c NavLinkCreator) Init(ctx context.Context, ctrl NavLinkI) {
	sub := grctl.ActiveLinkCreator{
		ID:               c.ID,
		Text:             c.Text,
		Location:         c.Location,
		ActiveAttributes: c.ActiveAttributes,
		ControlOptions:   c.ControlOptions,
	}
	sub.Init(ctx, ctrl)
}

// GetNavLink is a convenience method to return the button with the given id from the page.
func GetNavLink(c page.ControlI, id string) *NavLink {
	return c.Page().GetControl(id).(*NavLink)
}

func init() {
	page.RegisterControl(&NavLink{})
}
