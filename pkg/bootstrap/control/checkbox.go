package control

import (
	"context"
	"fmt"
	"github.com/goradd/goradd/pkg/bootstrap/config"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/action"
	"github.com/goradd/goradd/pkg/page/control/button"
	"github.com/goradd/goradd/pkg/page/event"
	"github.com/goradd/html5tag"
	"io"
)

type Checkbox struct {
	button.Checkbox
	inline bool
}

func NewCheckbox(parent page.ControlI, id string) *Checkbox {
	c := &Checkbox{}
	c.Init(c, parent, id)
	config.LoadBootstrap(c.ParentForm())
	return c
}

func (c *Checkbox) SetInline(v bool) *Checkbox {
	c.inline = v
	return c
}

func (c *Checkbox) DrawingAttributes(ctx context.Context) html5tag.Attributes {
	a := c.Checkbox.DrawingAttributes(ctx)
	a.AddClass("form-check-input")
	a.SetData("grctl", "bs-checkbox")
	if c.Text() == "" {
		a.AddClass("position-static")
	}
	return a
}

func (c *Checkbox) GetDrawingLabelAttributes() html5tag.Attributes {
	a := c.Checkbox.GetDrawingLabelAttributes()
	a.AddClass("form-check-label")
	return a
}

func (c *Checkbox) DrawTag(ctx context.Context, w io.Writer) {
	checkWrapperAttributes := html5tag.NewAttributes().
		AddClass("form-check").
		SetData("grel", c.ID()) // make sure the entire control gets removed
	if c.inline {
		checkWrapperAttributes.AddClass("form-check-inline")
	}
	if _, err := fmt.Fprint(w, "<div ", checkWrapperAttributes.String(), ">\n"); err != nil {
		panic(err)
	}
	c.Checkbox.DrawTag(ctx, w)
	if _, err := io.WriteString(w, "\n</div>"); err != nil {
		panic(err)
	}
}

func (c *Checkbox) Serialize(e page.Encoder) {
	c.Checkbox.Serialize(e)

	if err := e.Encode(c.inline); err != nil {
		panic(err)
	}
}

func (c *Checkbox) Deserialize(d page.Decoder) {
	c.Checkbox.Deserialize(d)

	if err := d.Decode(&c.inline); err != nil {
		panic(err)
	}
	return
}

type CheckboxCreator struct {
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
	// OnChange is an action to take when the user checks or unchecks the control.
	OnChange action.ActionI
	// Set inline when drawing this checkbox inline or wrapped by an inline FormGroup
	Inline bool
	page.ControlOptions
}

// Create is called by the framework to create a new control from the Creator. You
// do not normally need to call this.
func (c CheckboxCreator) Create(ctx context.Context, parent page.ControlI) page.ControlI {
	ctrl := NewCheckbox(parent, c.ID)
	if c.Text != "" {
		ctrl.SetText(c.Text)
	}
	if c.LabelMode != html5tag.LabelDefault {
		ctrl.LabelMode = c.LabelMode
	}
	if c.LabelAttributes != nil {
		ctrl.LabelAttributes().Merge(c.LabelAttributes)
	}

	ctrl.ApplyOptions(ctx, c.ControlOptions)
	if c.SaveState {
		ctrl.SaveState(ctx, c.SaveState)
	}
	if c.Inline {
		ctrl.SetInline(c.Inline)
	}
	if c.OnChange != nil {
		ctrl.On(event.Change(), c.OnChange)
	}
	return ctrl
}

// GetCheckbox is a convenience method to return the checkbox with the given id from the page.
func GetCheckbox(c page.ControlI, id string) *Checkbox {
	return c.Page().GetControl(id).(*Checkbox)
}

func init() {
	page.RegisterControl(&Checkbox{})
}
