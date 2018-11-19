package control

import (
	"github.com/spekary/goradd/pkg/html"
	"github.com/spekary/goradd/pkg/page"
	"goradd-project/override/control_base"
)

type Checkbox struct {
	control_base.Checkbox
}

func NewCheckbox(parent page.ControlI, id string) *Checkbox {
	c := &Checkbox{}
	c.Init(c, parent, id)
	return c
}

func (c *Checkbox) DrawingAttributes() *html.Attributes {
	a := c.Checkbox.DrawingAttributes()
	a.SetDataAttribute("grctl", "checkbox")
	a.Set("name", c.ID()) // needed for posts
	a.Set("type", "checkbox")
	a.Set("value", "1") // required for html validity
	return a
}

// UpdateFormValues is an internal call that lets us reflect the value of the checkbox on the web override
func (c *Checkbox) UpdateFormValues(ctx *page.Context) {
	id := c.ID()

	if v, ok := ctx.CheckableValue(id); ok {
		c.SetCheckedNoRefresh(v)
	}
}
