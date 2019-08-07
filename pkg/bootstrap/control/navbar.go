package control

import (
	"context"
	"github.com/goradd/goradd/pkg/bootstrap/config"
	"github.com/goradd/goradd/pkg/html"
	"github.com/goradd/goradd/pkg/page"
)

const (
	NavTabs      = "nav-tabs"
	NavPills     = "nav-pills"
	NavJustified = "nav-justified"

	NavbarHeader   = "navbar-header"
	NavbarCollapse = "navbar-collapse"
	NavbarBrand    = "navbar-brand"
	NavbarToggle   = "navbar-toggle"
	NavbarNav      = "navbar-nav"
	NavbarLeft     = "navbar-left"
	NavbarRight    = "navbar-right"
	NavbarForm     = "navbar-form"
)

type NavbarExpandClass string

const (
	NavbarExpandExtraLarge NavbarExpandClass = "navbar-expand-xl"
	NavbarExpandLarge                        = "navbar-expand-lg"
	NavbarExpandMedium                       = "navbar-expand-md"
	NavbarExpandSmall                        = "navbar-expand-sm"
	// NavbarExpandNone will always show the navbar as collapsed at any size
	NavbarExpandNone = ""
)

// NavbarCollapsedBrandPlacement controls the location of the brand when the navbar is collapsed
type NavbarCollapsedBrandPlacement int

const (
	// NavbarCollapsedBrandLeft will place the brand on the left and the toggle button on the right when collapsed
	NavbarCollapsedBrandLeft NavbarCollapsedBrandPlacement = iota
	// NavbarCollapsedBrandRight will place the brand on the right and the toggle button on the left when collapsed
	NavbarCollapsedBrandRight
	// NavbarCollapsedBrandHidden means the brand will be hidden when collapsed, and shown when expanded
	NavbarCollapsedBrandHidden
)

type NavbarI interface {
	page.ControlI
	SetHeaderAnchor(a string) NavbarI
	SetNavbarStyle(style NavbarStyle) NavbarI
	SetBrandPlacement(p NavbarCollapsedBrandPlacement) NavbarI
	SetBackgroundClass(c BackgroundColorClass) NavbarI
	SetExpand(e NavbarExpandClass) NavbarI
}

// Navbar is a bootstrap navbar object. Use SetText() to set the logo text of the navbar, and
// SetTextIsHtml() to true to turn off encoding if needed. Add child controls to populate it.
type Navbar struct {
	page.Control
	headerAnchor string

	style NavbarStyle
	//container ContainerClass ??
	background    BackgroundColorClass
	expand        NavbarExpandClass
	brandLocation NavbarCollapsedBrandPlacement
}

type NavbarStyle string

const (
	NavbarDark  NavbarStyle = "navbar-dark" // black on white
	NavbarLight             = "navbar-light"
)

// NewNavbar returns a newly created Bootstrap Navbar object
func NewNavbar(parent page.ControlI, id string) *Navbar {
	b := &Navbar{}
	b.Init(b, parent, id)
	return b
}

func (b *Navbar) Init(self page.ControlI, parent page.ControlI, id string) {
	b.Control.Init(self, parent, id)
	b.Tag = "nav"
	b.style = NavbarDark // default
	b.background = BackgroundColorDark
	b.expand = NavbarExpandLarge
	config.LoadBootstrap(b.ParentForm())
}

func (b *Navbar) this() NavbarI {
	return b.Self.(NavbarI)
}

func (b *Navbar) SetNavbarStyle(style NavbarStyle) NavbarI {
	b.style = style
	return b.this()
}

func (b *Navbar) SetBackgroundClass(c BackgroundColorClass) NavbarI {
	b.background = c
	return b.this()
}

func (b *Navbar) SetHeaderAnchor(a string) NavbarI {
	b.headerAnchor = a
	return b.this()
}

func (b *Navbar) SetExpand(e NavbarExpandClass) NavbarI {
	b.expand = e
	return b.this()
}


// SetBrandPlacement places the brand left, right, or hidden (meaning inside the collapse area).
// The expand button location will be affected by the placement
func (b *Navbar) SetBrandPlacement(p NavbarCollapsedBrandPlacement) NavbarI {
	b.brandLocation = p
	return b.this()
}

func (b *Navbar) ΩDrawingAttributes() html.Attributes {
	a := b.Control.ΩDrawingAttributes()
	a.AddClass("navbar")
	a.AddClass(string(b.style))
	a.AddClass(string(b.expand))
	a.AddClass(string(b.background))
	a.SetDataAttribute("grctl", "bs-navbar")
	return a
}

type NavbarCreator struct {
	ID string
	HeaderAnchor string
	Style NavbarStyle
	BackgroundColorClass    BackgroundColorClass
	Expand        NavbarExpandClass
	BrandLocation NavbarCollapsedBrandPlacement

	page.ControlOptions
	Children []page.Creator
}

// Create is called by the framework to create a new control from the Creator. You
// do not normally need to call this.
func (c NavbarCreator) Create(ctx context.Context, parent page.ControlI) page.ControlI {
	ctrl := NewNavbar(parent, c.ID)
	c.Init(ctx, ctrl)
	return ctrl
}

func (c NavbarCreator) Init(ctx context.Context, ctrl NavbarI) {
	if c.HeaderAnchor != "" {
		ctrl.SetHeaderAnchor(c.HeaderAnchor)
	}
	if c.Style != "" {
		ctrl.SetNavbarStyle(c.Style)
	}
	if c.BackgroundColorClass != "" {
		ctrl.SetBackgroundClass(c.BackgroundColorClass)
	}
	if c.Expand != "" {
		ctrl.SetExpand (c.Expand)
	}
	ctrl.SetBrandPlacement(c.BrandLocation)
	ctrl.AddControls(ctx, c.Children...)
}

// GetNavbar is a convenience method to return the control with the given id from the page.
func GetNavbar(c page.ControlI, id string) *Navbar {
	return c.Page().GetControl(id).(*Navbar)
}
