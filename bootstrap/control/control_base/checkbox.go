package control_base

import (
	grbase "github.com/spekary/goradd/page/control/control_base"
	"github.com/spekary/goradd/page"
	"github.com/spekary/goradd/html"
	"context"
)

type CheckboxI interface {
	grbase.CheckboxI
}


// Checkbox is a base class for checkbox-like objects, including bootstrap checkboxes and radio buttons.
type Checkbox struct {
	grbase.Checkbox
}

func (c *Checkbox) Init(self CheckboxI, parent page.ControlI) {
	c.Checkbox.Init(self, parent)
	c.LabelMode = html.LABEL_AFTER
}

func (c *Checkbox) this() CheckboxI {
	return c.Self.(CheckboxI)
}

func (c *Checkbox) DrawingAttributes() *html.Attributes {
	a := c.Checkbox.DrawingAttributes()
	a.AddClass("form-check-input")

	if c.Text() == "" {
		a.AddClass("position-static")
	}
	return a
}

func (c *Checkbox) GetDrawingInputLabelAttributes() *html.Attributes {
	a := c.Checkbox.GetDrawingInputLabelAttributes()
	a.AddClass("form-check-label")
	return a
}

func (c *Checkbox) DrawTag(ctx context.Context) (ctrl string) {
	h := c.Checkbox.DrawTag(ctx)
	return html.RenderTag("div", html.NewAttributes().
		AddClass("form-check").
		SetDataAttribute("grel", c.ID()), h)	// make sure the entire control gets removed
}