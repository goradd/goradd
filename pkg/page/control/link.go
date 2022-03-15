package control

import (
	"context"
	"github.com/goradd/html5tag"
	"github.com/goradd/goradd/pkg/page"
)

type LinkI interface {
	page.ControlI
	SetLabel(label string) page.ControlI
	SetLocation(href string) LinkI
}

// Link is a standard html link. It corresponds to a <a> tag in html.
type Link struct {
	page.ControlBase
}

// NewLink creates a new standard html link
func NewLink(parent page.ControlI, id string) *Link {
	b := new(Link)
	b.Self = b
	b.Init(parent, id)
	return b
}

// Init is called by subclasses of Link to initialize the link control structure.
func (l *Link) Init(parent page.ControlI, id string) {
	l.ControlBase.Init(parent, id)
	l.Tag = "a"
}

func (l *Link) this() LinkI {
	return l.Self.(LinkI)
}


// SetLabel sets the text that appears between the a tags.
func (l *Link) SetLabel(label string) page.ControlI {
	l.SetText(label)
	return l.this()
}

// SetLocation sets the href attribute of the link
func (l *Link) SetLocation(url string) LinkI {
	l.SetAttribute("href", url)
	return l.this()
}

// SetDownload sets the download attribute of the link.
//
// When a user clicks on the link, the browser will cause the
// destination to be downloaded. Pass a value to name the download
// file, or pass the empty string to cause the browser to use the
// link name as the file name.
func (l *Link) SetDownload(v string) {
	if v == "" {
		l.SetAttribute("download", true)
	} else {
		l.SetAttribute("download", v)
	}
}

// DrawingAttributes retrieves the tag's attributes at draw time. You should not normally need to call this, and the
// attributes are disposed of after drawing, so they are essentially read-only.
func (l *Link) DrawingAttributes(ctx context.Context) html5tag.Attributes {
	a := l.ControlBase.DrawingAttributes(ctx)
	a.SetData("grctl", "link")
	return a
}

// LinkCreator is the initialization structure for declarative creation of links
type LinkCreator struct {
	// ID is the control id
	ID string
	// Text is the text displayed inside the link
	Text string
	// Location is the destination of the link (the href attribute)
	Location string
	page.ControlOptions
}

// Create is called by the framework to create a new control from the Creator. You
// do not normally need to call this.
func (c LinkCreator) Create(ctx context.Context, parent page.ControlI) page.ControlI {
	ctrl := NewLink(parent, c.ID)
	c.Init(ctx, ctrl)
	return ctrl
}

// Init is called by implementations of Links to initialize a control with the
// creator. You do not normally need to call this.
func (c LinkCreator) Init(ctx context.Context, ctrl LinkI) {
	ctrl.SetLabel(c.Text)
	if c.Location != "" {
		ctrl.SetLocation(c.Location)
	}
	ctrl.ApplyOptions(ctx, c.ControlOptions)
}

// GetLink is a convenience method to return the link with the given id from the page.
func GetLink(c page.ControlI, id string) *Link {
	return c.Page().GetControl(id).(*Link)
}

func init() {
	page.RegisterControl(&Link{})
}