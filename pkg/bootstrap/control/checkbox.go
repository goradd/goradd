package control

import (
	"context"
	"encoding/gob"
	"github.com/goradd/goradd/pkg/bootstrap/config"
	"github.com/goradd/goradd/pkg/html"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/control"
	"reflect"
)

type Checkbox struct {
	control.Checkbox
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

func (c *Checkbox) ΩDrawingAttributes() *html.Attributes {
	a := c.Checkbox.ΩDrawingAttributes()
	a.AddClass("form-check-input")
	a.SetDataAttribute("grctl", "bs-checkbox")
	if c.Text() == "" {
		a.AddClass("position-static")
	}
	return a
}

func (c *Checkbox) ΩGetDrawingLabelAttributes() *html.Attributes {
	a := c.Checkbox.ΩGetDrawingLabelAttributes()
	a.AddClass("form-check-label")
	return a
}

func (c *Checkbox) ΩDrawTag(ctx context.Context) (ctrl string) {
	h := c.Checkbox.ΩDrawTag(ctx)
	checkWrapperAttributes := html.NewAttributes().
		AddClass("form-check").
		SetDataAttribute("grel", c.ID()) // make sure the entire control gets removed
	if c.inline {
		checkWrapperAttributes.AddClass("form-check-inline")
	}
	return html.RenderTag("div", checkWrapperAttributes, h) // make sure the entire control gets removed
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

func init() {
	gob.RegisterName("bootstrap.checkbox", new(Checkbox))
}
