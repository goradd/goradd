package control

import (
	"github.com/spekary/goradd/html"
	"github.com/spekary/goradd/page"
	"github.com/spekary/goradd/bootstrap/control/control_base"
)

type RadioButtonI interface {
	control_base.CheckboxI
}

type RadioButton struct {
	control_base.Checkbox
	group string

}

func NewRadioButton(parent page.ControlI) *RadioButton {
	c := &RadioButton{}
	c.Init(c, parent)
	return c
}

func (c *RadioButton) this() RadioButtonI {
	return c.Self.(RadioButtonI)
}

func (c *RadioButton) DrawingAttributes() *html.Attributes {
	a := c.Checkbox.DrawingAttributes()
	a.SetDataAttribute("grctl", "bs-radio")
	a.Set("type", "radio")
	if c.group == "" {
		a.Set("name", c.ID()) // treat it like a checkbox if no group is specified
	} else {
		a.Set("name", c.group)
		a.Set("value", c.ID())
	}
	return a
}

// UpdateFormValues is an internal call that lets us reflect the value of the checkbox on the web page
func (c *RadioButton) UpdateFormValues(ctx *page.Context) {
	id := c.ID()

	if v, ok := ctx.CheckableValue(id); ok {
		c.SetCheckedNoRefresh(v)
	}
}

func (c *RadioButton) SetGroup(g string) RadioButtonI {
	c.group = g
	c.Refresh()
	return c.this()
}

func (c *RadioButton) Group() string {
	return c.group
}
