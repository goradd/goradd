package control

import (
	"context"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/html5tag"
)

type ActiveLinkI interface {
	LinkI
	IsActive() bool
	SetIsActive(bool) ActiveLinkI
	ActiveAttributes() html5tag.Attributes
}

// ActiveLink is a link that has an active status.
//
// You can give an ActiveLink a set of attributes to merge into the attributes of the html tag
// to indicate that it is active. This can form the basis of a variety of specialized HTML controls, like
// navbars, tabs, and accordions.
type ActiveLink struct {
	Link
	isActive         bool
	activeAttributes html5tag.Attributes
}

// NewActiveLink creates a new standard html link
func NewActiveLink(parent page.ControlI, id string) *ActiveLink {
	b := new(ActiveLink)
	b.Init(b, parent, id)
	return b
}

// Init is called by subclasses of ActiveLink to initialize the link control structure.
func (l *ActiveLink) Init(self any, parent page.ControlI, id string) {
	l.Link.Init(self, parent, id)
}

func (l *ActiveLink) this() ActiveLinkI {
	return l.Self().(ActiveLinkI)
}

// SetIsActive sets the active state of the control. When active, the ActiveAttributes will be
// merged into the tags other HTML attributes.
func (l *ActiveLink) SetIsActive(isActive bool) ActiveLinkI {
	l.isActive = isActive
	l.Refresh()
	return l.this()
}

// IsActive returns whether the link is active or not.
func (l *ActiveLink) IsActive() bool {
	return l.isActive
}

// ActiveAttributes returns the attributes used to indicate the link is active.
// If you change them, be sure to call Refresh().
func (l *ActiveLink) ActiveAttributes() html5tag.Attributes {
	if l.activeAttributes == nil {
		l.activeAttributes = html5tag.NewAttributes()
	}
	return l.activeAttributes
}

// DrawingAttributes retrieves the tag's attributes at draw time. You should not normally need to call this, and the
// attributes are disposed of after drawing, so they are essentially read-only.
func (l *ActiveLink) DrawingAttributes(ctx context.Context) html5tag.Attributes {
	a := l.ControlBase.DrawingAttributes(ctx)
	a.SetData("grctl", "activelink")
	if l.isActive {
		a.Merge(l.activeAttributes)
	}
	return a
}

// Serialize serializes the state of the control for the pagestate
func (l *ActiveLink) Serialize(e page.Encoder) {
	l.Link.Serialize(e)
	if err := e.Encode(l.isActive); err != nil {
		panic(err)
	}
	if err := e.Encode(l.activeAttributes); err != nil {
		panic(err)
	}
}

// Deserialize reconstructs the control from the page state.
func (l *ActiveLink) Deserialize(d page.Decoder) {
	l.Link.Deserialize(d)

	if err := d.Decode(&l.isActive); err != nil {
		panic(err)
	}
	if err := d.Decode(&l.activeAttributes); err != nil {
		panic(err)
	}
}

// ActiveLinkCreator is the initialization structure for declarative creation of links
type ActiveLinkCreator struct {
	// ID is the control id. Leave blank to have one automatically assigned.
	ID string
	// Text is the text displayed inside the link.
	Text string
	// Location is the destination of the link (the href attribute).
	Location string
	// Download indicates that the "download" attribute should be assigned to the link.
	// Set to true to use the text of the link as the name of the file. Otherwise set to a string
	// to indicate the name of the file.
	Download any
	// ActiveAttributes are the attributes merged in to the current attributes to show that the link is active.
	// Matching attributes will be overridden for the most part, but the class attribute will be merged.
	ActiveAttributes html5tag.Attributes
	// ControlOptions are additional options for the html control.
	page.ControlOptions
}

// Create is called by the framework to create a new control from the Creator. You
// do not normally need to call this.
func (c ActiveLinkCreator) Create(ctx context.Context, parent page.ControlI) page.ControlI {
	ctrl := NewActiveLink(parent, c.ID)
	c.Init(ctx, ctrl)
	return ctrl
}

// Init is called by implementations of ActiveLinks to initialize a control with the creator.
func (c ActiveLinkCreator) Init(ctx context.Context, ctrl ActiveLinkI) {
	sub := LinkCreator{
		ID:             c.ID,
		Text:           c.Text,
		Location:       c.Location,
		Download:       c.Download,
		ControlOptions: c.ControlOptions,
	}

	sub.Init(ctx, ctrl)
	if c.ActiveAttributes != nil {
		ctrl.ActiveAttributes().Merge(c.ActiveAttributes)
	}
}

// GetActiveLink is a convenience method to return the link with the given id from the page.
func GetActiveLink(c page.ControlI, id string) *ActiveLink {
	return c.Page().GetControl(id).(*ActiveLink)
}

func init() {
	page.RegisterControl(&ActiveLink{})
}
