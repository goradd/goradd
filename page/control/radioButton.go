package control

import (
	"github.com/spekary/goradd/page"
	"github.com/spekary/goradd/page/control/control_base"
	"github.com/spekary/goradd/html"
)

// RadioButton is a standard html radio button. You can optionally specify a group name for the radiobutton to belong
// to and the browser will make sure only one item in the group is selected.
type RadioButton struct {
	control_base.Checkbox
	group string
}

func NewRadioButton(parent page.ControlI) *RadioButton {
	c := &RadioButton{}
	c.Init(c, parent)
	return c
}


func (c *RadioButton) DrawingAttributes() *html.Attributes {
	a := c.DrawingAttributes()
	a.SetDataAttribute("grctl", "radio")
	a.Set("type", "radio")
	if c.group == "" {
		a.Set("name", c.Id())	// treat it like a checkbox if no group is specified
	} else {
		a.Set("name", c.group)
		a.Set("value", c.Id())
	}
	return a
}

func (c *RadioButton) UpdateFormValues(ctx *page.Context) {
	id := c.Id()

	if v,ok := ctx.CheckableValue(id); ok {
		c.SetCheckedNoRefresh(v)
	}
}

func (c *RadioButton) SetGroup (g string) page.ControlI {
	c.group = g
	c.Refresh()
	return c.This()
}

func (c *RadioButton) Group() string {
	return c.group
}

func (c *RadioButton) SetChecked(v bool) page.ControlI {
	if c.group != "" && v {
		if c.Checked() != v {
			c.SetCheckedNoRefresh(v)
			// make sure any other buttons in the group are unchecked
			c.Form().Response().ExecuteJsFunction("goradd.setRadioInGroup", page.PriorityStandard, c.Id());
		}
	} else {
		c.Checkbox.SetChecked(v)
	}
	return c.This()
}

