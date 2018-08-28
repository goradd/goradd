package control

import (
	"github.com/spekary/goradd/html"
	"github.com/spekary/goradd/page"
	localPage "goradd-project/override/page"
	"github.com/spekary/goradd/bootstrap/config"
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
	localPage.ControlI
}

// Navbar is a bootstrap navbar object. Use SetText() to set the logo text of the navbar, and
// SetEscapeText() to false to turn off encoding if needed. Add child controls to populate it.
type Navbar struct {
	localPage.Control
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

// SetBrandPlacement places the brand left, right, or hidden (meaning inside the collapse area).
// The expand button location will be affected by the placement
func (b *Navbar) SetBrandPlacement(p NavbarCollapsedBrandPlacement) NavbarI {
	b.brandLocation = p
	return b.this()
}

func (b *Navbar) DrawingAttributes() *html.Attributes {
	a := b.Control.DrawingAttributes()
	a.AddClass("navbar")
	a.AddClass(string(b.style))
	a.AddClass(string(b.expand))
	a.AddClass(string(b.background))
	a.SetDataAttribute("grctl", "bs-navbar")
	return a
}
