package control

import (
	"context"
	"github.com/goradd/goradd/pkg/page"
	grctl "github.com/goradd/goradd/pkg/page/control"
	"github.com/goradd/html5tag"
)

type NavLinkI interface {
	grctl.LinkI
	SetFormID(formID string) NavLinkI
	SetValue(eventValue string) NavLinkI
}

// NavLink creates an anchor tag with a nav-link class. It is used to create a link in a NavBar, and
// generally should be a child item of a NavGroup.
type NavLink struct {
	grctl.Link
	// formID is the id that corresponds to the link so that the link can be made active when on that page.
	formID string
	// eventValue is the event value sent when the link is clicked
	eventValue string
}

// NewNavLink creates a new NavLink.
func NewNavLink(parent page.ControlI, id string) *NavLink {
	l := new(NavLink)
	l.Self = l
	l.Init(parent, id)
	return l
}

// Init initializes the button
func (l *NavLink) Init(parent page.ControlI, id string) {
	l.Link.Init(parent, id)
}

func (l *NavLink) this() NavLinkI {
	return l.Self.(NavLinkI)
}

// SetFormID sets the FormID that corresponds to this link for the purpose of making it display as active.
func (l *NavLink) SetFormID(formID string) NavLinkI {
	l.formID = formID
	return l.this()
}

// SetValue sets the EventValue that is sent to the action when the link is clicked.
func (l *NavLink) SetValue(value string) NavLinkI {
	l.eventValue = value
	return l.this()
}

// DrawingAttributes returns the attributes to add to the tag just before the button is drawn.
func (l *NavLink) DrawingAttributes(ctx context.Context) html5tag.Attributes {
	a := l.Link.DrawingAttributes(ctx)
	a.AddClass("nav-link")
	a.SetData("grctl", "navLink")
	if l.eventValue != "" {
		a.SetData("grEv", l.eventValue)
	}
	if l.ParentForm().ID() == l.formID {
		a.AddClass("active")
		a.Set("aria-current", "page")
	}
	return a
}

// Serialize serializes the state of the control for the pagestate
func (l *NavLink) Serialize(e page.Encoder) {
	l.Link.Serialize(e)
	if err := e.Encode(l.formID); err != nil {
		panic(err)
	}
	if err := e.Encode(l.eventValue); err != nil {
		panic(err)
	}
	return
}

// Deserialize reconstructs the control from the page state.
func (l *NavLink) Deserialize(d page.Decoder) {
	l.Link.Deserialize(d)

	if err := d.Decode(&l.formID); err != nil {
		panic(err)
	}
	if err := d.Decode(&l.eventValue); err != nil {
		panic(err)
	}
}

// NavLinkCreator is the initialization structure for the declarative creation of the control.
type NavLinkCreator struct {
	// ID is the control id
	ID string
	// Text is the text displayed in the link
	Text string
	// Location is the content of the href, which is the url where the link will go.
	// Alternatively, add an On handler to the ControlOptions
	Location string
	// FormID is the form ID of the form that corresponds to this link. If this is the current form,
	// the link will get attributes that make it display as active.
	FormID string
	// Value is the event value sent with the event when the link is clicked.
	Value string
	// ControlOptions are additional options to assign to the link.
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
// creator. You do not normally need to call this.
func (c NavLinkCreator) Init(ctx context.Context, ctrl NavLinkI) {
	sub := grctl.LinkCreator{
		ID:             c.ID,
		Text:           c.Text,
		Location:       c.Location,
		ControlOptions: c.ControlOptions,
	}
	sub.Init(ctx, ctrl)
	if c.FormID != "" {
		ctrl.SetFormID(c.FormID)
	}
	if c.Value != "" {
		ctrl.SetValue(c.Value)
	}
}

// GetNavLink is a convenience method to return the button with the given id from the page.
func GetNavLink(c page.ControlI, id string) *NavLink {
	return c.Page().GetControl(id).(*NavLink)
}

func init() {
	page.RegisterControl(&NavLink{})
}
