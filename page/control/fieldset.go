package control

import (
	"github.com/spekary/goradd/page/control/control_base"
	"github.com/spekary/goradd/page"
	"github.com/spekary/goradd/html"
	"context"
)


// Fieldset is a Panel that is drawn with a fieldset tag. The panel's label is used as the legend tag.
// Fieldset's cannot have wrappers.
type Fieldset struct {
	control_base.Panel
}

func NewFieldset(parent page.ControlI) *Fieldset {
	p := &Fieldset{}
	p.Tag = "fieldset"
	p.Init(p, parent)
	return p
}

func (c *Fieldset) DrawingAttributes() *html.Attributes {
	a := c.Control.DrawingAttributes()
	a.SetDataAttribute("grctl", "fieldset")
	return a
}


func (c *Fieldset) DrawTag(ctx context.Context) string {
	var ctrl string

	attributes := c.This().DrawingAttributes()
	if c.HasWrapper()  {
		panic ("Fieldsets cannot have wrappers.")
	}

	buf := page.GetBuffer()
	defer page.PutBuffer(buf)

	if l := c.Label(); l != "" {
		ctrl = html.RenderTag("legend", nil, l)
	}
	if err := c.This().DrawInnerHtml(ctx, buf); err != nil {
		panic (err)
	}
	ctrl = html.RenderTag(c.Tag, attributes, ctrl + buf.String())
	return ctrl
}
