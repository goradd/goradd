package control

import (
	"github.com/goradd/goradd/pkg/html"
	"github.com/goradd/goradd/pkg/page"
)

type RadioButtonI interface {
	CheckboxI
}


// RadioButton is a standard html radio button. You can optionally specify a group name for the radiobutton to belong
// to and the browser will make sure only one item in the group is selected.
type RadioButton struct {
	CheckboxBase
	group string
}

// NewRadioButton creates a new radio button
func NewRadioButton(parent page.ControlI, id string) *RadioButton {
	c := &RadioButton{}
	c.Init(c, parent, id)
	return c
}

func (c *RadioButton) this() RadioButtonI {
	return c.Self.(RadioButtonI)
}

// SetGroup sets the name of the group that the control will belong to. Set all the radio buttons
// that represent a selection from a group to this same group name.
func (c *RadioButton) SetGroup(g string) RadioButtonI {
	c.group = g
	c.Refresh()
	return c.this()
}

// Group returns the name of the group that the control belongs to.
func (c *RadioButton) Group() string {
	return c.group
}

// SetChecked will set the checked status of this radio button to the given value.
func (c *RadioButton) SetChecked(v bool) RadioButtonI {
	if c.group != "" && v {
		if c.Checked() != v {
			c.SetCheckedNoRefresh(v)
			// make sure any other buttons in the group are unchecked
			// TODO: This requires a round trip from the client, so doesn't work that great. In other words, eventually
			// we will get this response, but not right away. Since its more common to make a RadioList rather than
			// separate radio buttons in a group, we are not going to worry about it for now. It if becomes an issue,
			// the code would need to change to look through the forms control list for other buttons in the group, and
			// update those buttons in the go code here.
			c.ParentForm().Response().ExecuteJsFunction("goradd.setRadioInGroup", page.PriorityStandard, c.ID())
		}
	} else {
		c.CheckboxBase.SetChecked(v)
	}
	return c.this()
}

// ΩDrawingAttributes is called by the framework to create temporary attributes for the input tag.
func (c *RadioButton) ΩDrawingAttributes() *html.Attributes {
	a := c.CheckboxBase.ΩDrawingAttributes()
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

// ΩUpdateFormValues is called by the framework to update the value of the control based on
// values sent by the browser.
func (c *RadioButton) ΩUpdateFormValues(ctx *page.Context) {
	c.UpdateRadioFormValues(ctx, c.Group())
}

