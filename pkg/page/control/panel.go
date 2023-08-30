package control

import (
	"context"
	"fmt"

	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/html5tag"
)

// Children is a helper function for doing declarative control creation for child control creators.
// It returns the creators passed to it as a slice.
func Children(creators ...page.Creator) []page.Creator {
	return creators
}

type PanelI interface {
	page.ControlI
}

// Panel is a GoRADD control that is a basic "div" wrapper.
//
// Panel can be used for any kind of HTML tag by simply changing the Tag attribute. For example,
//
//	panel.Tag = "nav"
//
// Turns a panel into a "nav" tag.
//
// Customize how the tag is drawn by calling functions inherited from ControlBase. With these, you can
// set the class, data attributes, or any attribute. Call SetText() to set a string that will be
// drawn inside the div tag. Call SetTextIsHtml() to tell the control to treat the text as HTML and
// not escape it.
//
// One typical use for a Panel is as a container for custom HTML and child controls. Child controls assigned
// to the Panel will automatically be drawn by default in the order they were assigned.
//
// To customize this behavior, embed a Panel into your own custom struct, and then define the
// DrawInnerHtml() function on your struct. The framework will automatically call that function.
// Or, use a template to create the DrawTemplate function on your struct and the framework will use that
// instead. Examples can be found in the tutorial, in the bootstrap package, and in the code-generated panels.
type Panel struct {
	page.ControlBase
}

func NewPanel(parent page.ControlI, id string) *Panel {
	p := &Panel{}
	p.Init(p, parent, id)
	return p
}

func (c *Panel) Init(self any, parent page.ControlI, id string) {
	c.ControlBase.Init(self, parent, id)
	c.Tag = "div"
}

func (c *Panel) this() PanelI {
	return c.Self().(PanelI)
}

func (c *Panel) DrawingAttributes(ctx context.Context) html5tag.Attributes {
	a := c.ControlBase.DrawingAttributes(ctx)
	a.SetData("grctl", "panel")
	return a
}

// Value satisfies the Valuer interface and returns the text of the panel.
func (c *Panel) Value() interface{} {
	return c.Text()
}

// SetValue satisfies the Valuer interface and sets the text of the panel.
func (c *Panel) SetValue(v interface{}) page.ControlI {
	return c.SetText(fmt.Sprint(v))
}

// PanelCreator creates a div control with child controls.
// Pass it to AddControls or as a child of a parent control.
type PanelCreator struct {
	// ID is the id the tag will have on the page and must be unique on the page
	ID string
	// Tag replaces the tag of the div object with the given tag.
	Tag string
	// Text is text that will become the innerhtml part of the tag.
	Text string
	// If you set TextIsHtml, the Text will not be escaped prior to drawing.
	TextIsHtml bool
	// Children is a list of creators to use to create the child controls of the panel.
	// You can wrap your child creators with the Children() function as a helper. For example:
	//   Children: Children(
	//     TextboxCreator{...},
	//     ButtonCreator{...},
	//   )
	Children []page.Creator
	page.ControlOptions
}

// Create is called by the framework to create the panel.
func (c PanelCreator) Create(ctx context.Context, parent page.ControlI) page.ControlI {
	ctrl := NewPanel(parent, c.ID)
	c.Init(ctx, ctrl)
	return ctrl
}

// Init is called by implementations to initialize a control with the creator.
func (c PanelCreator) Init(ctx context.Context, ctrl PanelI) {
	if c.Text != "" {
		ctrl.SetText(c.Text)
	}
	if c.Tag != "" {
		ctrl.SetTag(c.Tag)
	}
	ctrl.SetTextIsHtml(c.TextIsHtml)
	ctrl.ApplyOptions(ctx, c.ControlOptions)
	ctrl.AddControls(ctx, c.Children...)
}

// GetPanel is a convenience method to return the panel with the given id from the page.
func GetPanel(c page.ControlI, id string) *Panel {
	return c.Page().GetControl(id).(*Panel)
}

func init() {
	page.RegisterControl(&Panel{})
}
