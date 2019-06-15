package control

import (
	"context"
	"github.com/goradd/goradd/pkg/bootstrap/config"
	"github.com/goradd/goradd/pkg/html"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/control"
)

type RadioButtonI interface {
	control.RadioButtonI
}

type RadioButton struct {
	control.RadioButton
	inline bool
}

func NewRadioButton(parent page.ControlI, id string) *RadioButton {
	c := &RadioButton{}
	c.Init(c, parent, id)
	config.LoadBootstrap(c.ParentForm())
	return c
}

func (c *RadioButton) this() RadioButtonI {
	return c.Self.(RadioButtonI)
}

func (c *RadioButton) SetInline(v bool) *RadioButton {
	c.inline = v
	return c
}

func (c *RadioButton) ΩDrawingAttributes() *html.Attributes {
	a := c.RadioButton.ΩDrawingAttributes()
	a.SetDataAttribute("grctl", "bs-radio")
	a.AddClass("form-check-input")

	if c.Text() == "" {
		a.AddClass("position-static")
	}
	return a
}

func (c *RadioButton) ΩGetDrawingLabelAttributes() *html.Attributes {
	a := c.RadioButton.ΩGetDrawingLabelAttributes()
	a.AddClass("form-check-label")
	return a
}

func (c *RadioButton) ΩDrawTag(ctx context.Context) (ctrl string) {
	h := c.RadioButton.ΩDrawTag(ctx)
	checkWrapperAttributes := html.NewAttributes().
		AddClass("form-check").
		SetDataAttribute("grel", c.ID()) // make sure the entire control gets removed
	if c.inline {
		checkWrapperAttributes.AddClass("form-check-inline")
	}
	return html.RenderTag("div", checkWrapperAttributes, h)
}

// TODO: Serialize
