package control_base

import (
	grbase "github.com/spekary/goradd/pkg/page/control/control_base"
	"github.com/spekary/goradd/pkg/page"
	"github.com/spekary/goradd/pkg/html"
	"context"
	"reflect"
)

type CheckboxI interface {
	grbase.CheckboxI
}


// Checkbox is a base class for checkbox-like objects, including bootstrap checkboxes and radio buttons.
type Checkbox struct {
	grbase.Checkbox
}

func (c *Checkbox) Init(self CheckboxI, parent page.ControlI, id string) {
	c.Checkbox.Init(self, parent, id)
	c.LabelMode = html.LabelAfter
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

func (c *Checkbox) Serialize(e page.Encoder) (err error) {
	if err = c.Checkbox.Serialize(e); err != nil {
		return
	}

	return
}

// ΩisSerializer is used by the automated control serializer to determine how far down the control chain the control
// has to go before just calling serialize and deserialize
func (c *Checkbox) ΩisSerializer(i page.ControlI) bool {
	return reflect.TypeOf(c) == reflect.TypeOf(i)
}


func (c *Checkbox) Deserialize(d page.Decoder, p *page.Page) (err error) {
	if err = c.Checkbox.Deserialize(d, p); err != nil {
		return
	}

	return
}
