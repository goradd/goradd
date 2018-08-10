package control

import (
	"github.com/spekary/goradd/html"
	"github.com/spekary/goradd/page"
	"goradd-project/override/control_base"
)

type RadioButtonI interface {
	control_base.CheckboxI
}


// RadioButton is a standard html radio button. You can optionally specify a group name for the radiobutton to belong
// to and the browser will make sure only one item in the group is selected.
type RadioButton struct {
	control_base.Checkbox
	group string
}

func NewRadioButton(parent page.ControlI, id string) *RadioButton {
	c := &RadioButton{}
	c.Init(c, parent, id)
	return c
}

func (c *RadioButton) this() RadioButtonI {
	return c.Self.(RadioButtonI)
}

func (c *RadioButton) DrawingAttributes() *html.Attributes {
	a := c.Checkbox.DrawingAttributes()
	a.SetDataAttribute("grctl", "radio")
	a.Set("type", "radio")
	if c.group == "" {
		a.Set("name", c.ID()) // treat it like a checkbox if no group is specified
	} else {
		a.Set("name", c.group)
		a.Set("value", c.ID())
	}
	return a
}

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

func (c *RadioButton) SetChecked(v bool) RadioButtonI {
	if c.group != "" && v {
		if c.Checked() != v {
			c.SetCheckedNoRefresh(v)
			// make sure any other buttons in the group are unchecked
			c.ParentForm().Response().ExecuteJsFunction("goradd.setRadioInGroup", page.PriorityStandard, c.ID())
		}
	} else {
		c.Checkbox.SetChecked(v)
	}
	return c.this()
}
