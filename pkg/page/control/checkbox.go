package control

import (
	"context"
	"github.com/goradd/goradd/pkg/html"
	"github.com/goradd/goradd/pkg/page"
)

// Checkbox is a basic html checkbox input form control.
type Checkbox struct {
	CheckboxBase
}

// NewCheckbox creates a new checkbox control.
func NewCheckbox(parent page.ControlI, id string) *Checkbox {
	c := &Checkbox{}
	c.Init(c, parent, id)
	return c
}

// ΩDrawingAttributes is called by the framework to set the temporary attributes that the control
// needs. Checkboxes set the grctl, name, type and value attributes automatically.
// You do not normally need to call this function.
func (c *Checkbox) ΩDrawingAttributes(ctx context.Context) html.Attributes {
	a := c.CheckboxBase.ΩDrawingAttributes(ctx)
	a.SetDataAttribute("grctl", "checkbox")
	a.Set("name", c.ID()) // needed for posts
	a.Set("type", "checkbox")
	a.Set("value", "1") // required for html validity
	return a
}

// ΩUpdateFormValues is an internal call that lets us reflect the value of the checkbox on the form.
// You do not normally need to call this function.
func (c *Checkbox) ΩUpdateFormValues(ctx *page.Context) {
	c.UpdateCheckboxFormValues(ctx)
}


type CheckboxCreator struct {
	// ID is the id of the control
	ID string
	// Text is the text of the label displayed right next to the checkbox.
	Text string
	// Checked will initialize the checkbox in its checked state.
	Checked bool
	// LabelMode specifies how the label is drawn with the checkbox.
	LabelMode html.LabelDrawingMode
	// LabelAttributes are additional attributes placed on the label tag.
	LabelAttributes html.Attributes
	// SaveState will save the value of the checkbox and restore it when the page is reentered.
	SaveState bool
	page.ControlOptions
}

// Create is called by the framework to create a new control from the Creator. You
// do not normally need to call this.
func (c CheckboxCreator) Create(ctx context.Context, parent page.ControlI) page.ControlI {
	ctrl := NewCheckbox(parent, c.ID)
	if c.Text != "" {
		ctrl.SetText(c.Text)
	}
	if c.LabelMode != html.LabelDefault {
		ctrl.LabelMode = c.LabelMode
	}
	if c.LabelAttributes != nil {
		ctrl.LabelAttributes().Merge(c.LabelAttributes)
	}

	ctrl.ApplyOptions(c.ControlOptions)
	if c.SaveState {
		ctrl.SaveState(ctx, c.SaveState)
	}
	return ctrl
}

// GetCheckbox is a convenience method to return the checkbox with the given id from the page.
func GetCheckbox(c page.ControlI, id string) *Checkbox {
	return c.Page().GetControl(id).(*Checkbox)
}

func init() {
	page.RegisterControl(Checkbox{})
}