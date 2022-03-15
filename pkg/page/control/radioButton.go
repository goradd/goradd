package control

import (
	"context"
	"github.com/goradd/html5tag"
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
	c.Self = c
	c.Init(parent, id)
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

// DrawingAttributes is called by the framework to create temporary attributes for the input tag.
func (c *RadioButton) DrawingAttributes(ctx context.Context) html5tag.Attributes {
	a := c.CheckboxBase.DrawingAttributes(ctx)
	a.SetData("grctl", "radio")
	a.Set("type", "radio")
	if c.group == "" {
		a.Set("name", c.ID()) // treat it like a checkbox if no group is specified
	} else {
		a.Set("name", c.group)
		a.Set("value", c.ID())
	}
	return a
}

// UpdateFormValues is used by the framework to cause the control to retrieve its values from the form
func (c *RadioButton) UpdateFormValues(ctx context.Context) {
	c.UpdateRadioFormValues(ctx, c.Group())
}

func (l *RadioButton) Serialize(e page.Encoder) {
	l.CheckboxBase.Serialize(e)
	if err := e.Encode(l.group); err != nil {
		panic(err)
	}
}

func (l *RadioButton) Deserialize(dec page.Decoder) {
	l.CheckboxBase.Deserialize(dec)
	if err := dec.Decode(&l.group); err != nil {
		panic(err)
	}
}

type RadioButtonCreator struct {
	// ID is the id of the control
	ID string
	// Text is the text of the label displayed right next to the checkbox.
	Text string
	// Checked will initialize the checkbox in its checked state.
	Checked bool
	// LabelMode specifies how the label is drawn with the checkbox.
	LabelMode html5tag.LabelDrawingMode
	// LabelAttributes are additional attributes placed on the label tag.
	LabelAttributes html5tag.Attributes
	// SaveState will save the value of the checkbox and restore it when the page is reentered.
	SaveState bool
	// Group is the name of the group that the button belongs to
	Group string
	page.ControlOptions
}

// Create is called by the framework to create a new control from the Creator. You
// do not normally need to call this.
func (c RadioButtonCreator) Create(ctx context.Context, parent page.ControlI) page.ControlI {
	ctrl := NewRadioButton(parent, c.ID)
	if c.Text != "" {
		ctrl.SetText(c.Text)
	}
	if c.LabelMode != html5tag.LabelDefault {
		ctrl.LabelMode = c.LabelMode
	}
	if c.LabelAttributes != nil {
		ctrl.LabelAttributes().Merge(c.LabelAttributes)
	}
	if c.Group != "" {
		ctrl.group = c.Group
	}

	ctrl.ApplyOptions(ctx, c.ControlOptions)
	if c.SaveState {
		ctrl.SaveState(ctx, c.SaveState)
	}
	return ctrl
}

// GetRadioButton is a convenience method to return the radio button with the given id from the page.
func GetRadioButton(c page.ControlI, id string) *RadioButton {
	return c.Page().GetControl(id).(*RadioButton)
}

func init() {
	page.RegisterControl(&RadioButton{})
}