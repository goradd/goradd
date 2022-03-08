package control

import (
	"context"
	"github.com/goradd/goradd/pkg/html5tag"
	"github.com/goradd/goradd/pkg/page"
	buf2 "github.com/goradd/goradd/pkg/pool"
	"html"
	"io"
)

type FieldsetI interface {
	PanelI
}

// Fieldset is a Panel that is drawn with a fieldset tag. The panel's text is used as the legend tag.
// Fieldset's cannot have wrappers.
type Fieldset struct {
	Panel
}

// NewFieldset creates a new Fieldset.
func NewFieldset(parent page.ControlI, id string) *Fieldset {
	p := &Fieldset{}
	p.Self = p
	p.Init(parent, id)
	return p
}

// Init is called by subclasses of Fieldset.
func (c *Fieldset) Init(parent page.ControlI, id string) {
	c.Panel.Init(parent, id)
	c.Tag = "fieldset"
}

func (c *Fieldset) this() FieldsetI {
	return c.Self.(FieldsetI)
}

// DrawingAttributes is called by the framework.
func (c *Fieldset) DrawingAttributes(ctx context.Context) html5tag.Attributes {
	a := c.ControlBase.DrawingAttributes(ctx)
	a.SetData("grctl", "fieldset")
	return a
}

// DrawTag is called by the framework.
func (c *Fieldset) DrawTag(ctx context.Context, w io.Writer) {
	var ctrl string

	attributes := c.this().DrawingAttributes(ctx)

	buf := buf2.GetBuffer()
	defer buf2.PutBuffer(buf)

	if l := c.Text(); l != "" {
		ctrl = html5tag.RenderTag("legend", nil, html.EscapeString(l))
	}
	c.this().DrawInnerHtml(ctx, buf)
	ctrl = html5tag.RenderTag(c.Tag, attributes, ctrl+buf.String())
	if _,err := io.WriteString(w, ctrl); err != nil {panic(err)}
}

// FieldsetCreator declares a Fieldset control. Pass it to AddControls or as
// a child of other creators.
type FieldsetCreator struct {
	// ID is the id the tag will have on the page and must be unique on the page
	ID string
	// Legend is the text to use in the legend tag of the fieldset
	Legend string
	// Children are the child creators declaring the controls wrapped by the fieldset
	Children []page.Creator
	page.ControlOptions
}

// Create is called by the framework to create the panel. You do not normally need to call this.
func (c FieldsetCreator) Create(ctx context.Context, parent page.ControlI) page.ControlI {
	ctrl := NewFieldset(parent, c.ID)
	if c.Legend != "" {
		ctrl.SetText(c.Legend)
	}
	ctrl.ApplyOptions(ctx, c.ControlOptions)
	ctrl.AddControls(ctx, c.Children...)
	return ctrl
}

// GetFieldset is a convenience method to return the panel with the given id from the page.
func GetFieldset(c page.ControlI, id string) *Fieldset {
	return c.Page().GetControl(id).(*Fieldset)
}

func init() {
	page.RegisterControl(&Fieldset{})
}