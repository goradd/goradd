package control

import (
	"context"
	"github.com/goradd/goradd/pkg/html"
	"github.com/goradd/goradd/pkg/page"
)

type FieldsetI interface {
	PanelI
}
// Fieldset is a Panel that is drawn with a fieldset tag. The panel's label is used as the legend tag.
// Fieldset's cannot have wrappers.
type Fieldset struct {
	Panel
}

// NewFieldset creates a new Fieldset.
func NewFieldset(parent page.ControlI, id string) *Fieldset {
	p := &Fieldset{}
	p.Init(p, parent, id)
	return p
}

// Init is called by subclasses of Fieldset.
func (c *Fieldset) Init (self FieldsetI, parent page.ControlI, id string) {
	c.Panel.Init(self, parent, id)
	c.Tag = "fieldset"
}

func (c *Fieldset) this() FieldsetI {
	return c.Self.(FieldsetI)
}

// DrawingAttributes is called by the framework.
func (c *Fieldset) DrawingAttributes() *html.Attributes {
	a := c.Control.DrawingAttributes()
	a.SetDataAttribute("grctl", "fieldset")
	return a
}

// DrawTag is called by the framework.
func (c *Fieldset) DrawTag(ctx context.Context) string {
	var ctrl string

	attributes := c.this().DrawingAttributes()
	if c.HasWrapper() {
		panic("Fieldsets cannot have wrappers.")
	}

	buf := page.GetBuffer()
	defer page.PutBuffer(buf)

	if l := c.Label(); l != "" {
		ctrl = html.RenderTag("legend", nil, l)
	}
	if err := c.this().DrawInnerHtml(ctx, buf); err != nil {
		panic(err)
	}
	ctrl = html.RenderTag(c.Tag, attributes, ctrl+buf.String())
	return ctrl
}
