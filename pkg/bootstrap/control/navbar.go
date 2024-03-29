package control

import (
	"context"
	"github.com/goradd/goradd/pkg/bootstrap/config"
	"github.com/goradd/goradd/pkg/javascript"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/action"
	"github.com/goradd/goradd/pkg/page/event"
	"github.com/goradd/html5tag"
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

const NavbarSelect = "gr-bs-navbarselect"

type NavbarI interface {
	page.ControlI
	SetNavbarStyle(style NavbarStyle) NavbarI
	SetBrand(label string, anchor string, p NavbarCollapsedBrandPlacement) NavbarI
	SetBackgroundClass(c BackgroundColorClass) NavbarI
	SetExpand(e NavbarExpandClass) NavbarI
	SetContainerClass(class ContainerClass) NavbarI
	OnClick(action action.ActionI) NavbarI
}

// Navbar is a bootstrap navbar object. Use SetText() to set the logo text of the navbar, and
// SetTextIsHtml() to true to turn off encoding if needed. Add child controls to populate it.
// When adding NavLink objects, they should be grouped together using the NavGroup object.
type Navbar struct {
	page.ControlBase
	brandAnchor   string
	brandLocation NavbarCollapsedBrandPlacement

	style          NavbarStyle
	background     BackgroundColorClass
	expand         NavbarExpandClass
	containerClass ContainerClass
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

func (b *Navbar) Init(self any, parent page.ControlI, id string) {
	b.ControlBase.Init(self, parent, id)
	b.Tag = "nav"
	b.style = NavbarDark // default
	b.background = BackgroundColorDark
	b.expand = NavbarExpandLarge
	b.containerClass = Container
	config.LoadBootstrap(b.ParentForm())

	b.On(event.Click().Selector(`a[href="#"][data-grctl="bs-navlink"]`).Capture(),
		action.Trigger(b.ID(), NavbarSelect, javascript.JsCode(`g$(event.target).id`)))

}

func (b *Navbar) this() NavbarI {
	return b.Self().(NavbarI)
}

func (b *Navbar) SetNavbarStyle(style NavbarStyle) NavbarI {
	b.style = style
	return b.this()
}

func (b *Navbar) SetBackgroundClass(c BackgroundColorClass) NavbarI {
	b.background = c
	return b.this()
}

func (b *Navbar) SetBrand(label string, anchor string, p NavbarCollapsedBrandPlacement) NavbarI {
	b.SetText(label)
	b.brandAnchor = anchor
	b.brandLocation = p
	return b.this()
}

func (b *Navbar) SetExpand(e NavbarExpandClass) NavbarI {
	b.expand = e
	return b.this()
}

func (b *Navbar) SetContainerClass(class ContainerClass) NavbarI {
	b.containerClass = class
	return b.this()
}

// OnClick sets the action to take when a link in the Navbar is selected.
// It will only respond to links whose href is "#", which indicates its an empty link.
// The EventValue will be the id of the link clicked.
func (b *Navbar) OnClick(a action.ActionI) NavbarI {
	b.On(NavbarSelectEvent(), a)
	return b.this()
}

func (b *Navbar) DrawingAttributes(ctx context.Context) html5tag.Attributes {
	a := b.ControlBase.DrawingAttributes(ctx)
	a.AddClass("navbar")
	a.AddClass(string(b.style))
	a.AddClass(string(b.expand))
	a.AddClass(string(b.background))
	a.SetData("grctl", "bs-navbar")
	return a
}

func (b *Navbar) Serialize(e page.Encoder) {
	b.ControlBase.Serialize(e)

	if err := e.Encode(b.brandAnchor); err != nil {
		panic(err)
	}
	if err := e.Encode(b.brandLocation); err != nil {
		panic(err)
	}
	if err := e.Encode(b.style); err != nil {
		panic(err)
	}
	if err := e.Encode(b.background); err != nil {
		panic(err)
	}
	if err := e.Encode(b.expand); err != nil {
		panic(err)
	}

	return
}

func (b *Navbar) Deserialize(d page.Decoder) {
	b.ControlBase.Deserialize(d)

	if err := d.Decode(&b.brandAnchor); err != nil {
		panic(err)
	}
	if err := d.Decode(&b.brandLocation); err != nil {
		panic(err)
	}
	if err := d.Decode(&b.style); err != nil {
		panic(err)
	}
	if err := d.Decode(&b.background); err != nil {
		panic(err)
	}
	if err := d.Decode(&b.expand); err != nil {
		panic(err)
	}
}

type NavbarCreator struct {
	ID string
	// Brand is the string to use for the brand
	Brand string
	// BrandIsHtml can be set to true to specify that the Brand string should not be escaped
	BrandIsHtml bool
	// BrandAnchor is the url to go to when the main logo in the navbar is clicked
	BrandAnchor string
	// Style is either NavbarDark or NavbarLight
	Style NavbarStyle
	// BackgroundColorClass is one of the background colors that you can assign the navbar
	BackgroundColorClass BackgroundColorClass
	// Expand determines at what screen width the navbar will be expanded.
	Expand NavbarExpandClass
	// BrandLocation controls the placement of the brand item
	BrandLocation NavbarCollapsedBrandPlacement
	// ContainerClass is the class of the container that will wrap the contents of the navbar
	ContainerClass ContainerClass

	// OnClick is the action to take when a link is clicked. It will only respond
	// to nav-link items that have an href of "#". The EventValue will be the id of the item clicked.
	OnClick action.ActionI

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
	if c.Brand != "" {
		ctrl.SetBrand(c.Brand, c.BrandAnchor, c.BrandLocation)
	}
	if c.BrandIsHtml {
		ctrl.SetTextIsHtml(true)
	}
	if c.Style != "" {
		ctrl.SetNavbarStyle(c.Style)
	}
	if c.BackgroundColorClass != "" {
		ctrl.SetBackgroundClass(c.BackgroundColorClass)
	}
	if c.Expand != "" {
		ctrl.SetExpand(c.Expand)
	}
	if c.OnClick != nil {
		ctrl.OnClick(c.OnClick)
	}
	if c.ContainerClass != "" {
		ctrl.SetContainerClass(c.ContainerClass)
	}
	ctrl.AddControls(ctx, c.Children...)
}

// GetNavbar is a convenience method to return the control with the given id from the page.
func GetNavbar(c page.ControlI, id string) *Navbar {
	return c.Page().GetControl(id).(*Navbar)
}

func init() {
	page.RegisterControl(&Navbar{})
}

func NavbarSelectEvent() *event.Event {
	return event.NewEvent(NavbarSelect)
}
