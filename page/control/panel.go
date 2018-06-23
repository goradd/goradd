package control

import (
	"github.com/spekary/goradd/html"
	"github.com/spekary/goradd/page"
	"github.com/spekary/goradd/page/control/control_base"
)


type PanelI interface {
	control_base.PanelI
}

// Panel is a Goradd control that is a basic "div" wrapper. Use it to style and listen to events on a div. It
// can also be used as the basis for more advanced javascript controls.
type Panel struct {
	control_base.Panel
}

func NewPanel(parent page.ControlI, id string) *Panel {
	p := &Panel{}
	p.Init(p, parent, id)
	return p
}

func (c *Panel) Init(self PanelI, parent page.ControlI, id string) {
	c.Panel.Init(self, parent, id)
	c.Tag = "div"
}


func (c *Panel) DrawingAttributes() *html.Attributes {
	a := c.Control.DrawingAttributes()
	a.SetDataAttribute("grctl", "panel")
	return a
}
