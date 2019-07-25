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
func (c *Checkbox) ΩDrawingAttributes() *html.Attributes {
	a := c.CheckboxBase.ΩDrawingAttributes()
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
	ID string
	Text string
	Checked bool
	LabelMode html.LabelDrawingMode
	LabelAttributes html.AttributeCreator
	SaveState bool
	page.ControlOptions
}

func (c CheckboxCreator) Create(ctx context.Context, parent page.ControlI) page.ControlI {
	ctrl := NewCheckbox(parent, c.ID)
	if c.Text != "" {
		ctrl.SetText(c.Text)
	}
	if c.LabelMode != html.LabelDefault {
		ctrl.LabelMode = c.LabelMode
	}
	if c.LabelAttributes != nil {
		ctrl.LabelAttributes().MergeMap(c.LabelAttributes)
	}

	ctrl.ApplyOptions(c.ControlOptions)
	if c.SaveState {
		ctrl.SaveState(ctx, c.SaveState)
	}
	return ctrl
}

// GetCheckbox is a convenience method to return the checkbox with the given id from the page.
func GetCheckbox(c page.ControlI, id string) *Checkbox {
	return c.Page().GetControl(id).(*Checkbox);
}