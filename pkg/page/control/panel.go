package control

import (
	"context"
	"github.com/goradd/goradd/pkg/html"
	"github.com/goradd/goradd/pkg/page"
)

// Children is just a helper function for doing declarative control creation for child control creators
func Children(creators ...page.Creator) []page.Creator {
	return creators
}


type PanelI interface {
	page.ControlI
}

// Panel is a Goradd control that is a basic "div" wrapper. Use it to style and listen to events on a div. It
// can also be used as the basis for more advanced javascript controls.
type Panel struct {
	page.Control
}

func NewPanel(parent page.ControlI, id string) *Panel {
	p := &Panel{}
	p.Init(p, parent, id)
	return p
}

func (c *Panel) Init(self PanelI, parent page.ControlI, id string) {
	c.Control.Init(self, parent, id)
	c.Tag = "div"
}

func (c *Panel) this() PanelI {
	return c.Self.(PanelI)
}


func (c *Panel) ΩDrawingAttributes() *html.Attributes {
	a := c.Control.ΩDrawingAttributes()
	a.SetDataAttribute("grctl", "panel")
	return a
}

func (c *Panel) Value() interface{} {
	return c.Text()
}

type PanelCreator struct {
	ID string
	Text string
	TextIsHtml bool
	Children []page.Creator
	page.ControlOptions
}

func (c PanelCreator) Create(ctx context.Context, parent page.ControlI) page.ControlI {
	ctrl := NewPanel(parent, c.ID)
	if c.Text != "" {
		ctrl.SetText(c.Text)
	}
	ctrl.SetTextIsHtml(c.TextIsHtml)
	ctrl.ApplyOptions(c.ControlOptions)
	ctrl.AddControls(ctx, c.Children...)
	return ctrl
}

// GetPanel is a convenience method to return the panel with the given id from the page.
func GetPanel(c page.ControlI, id string) *Panel {
	return c.Page().GetControl(id).(*Panel);
}