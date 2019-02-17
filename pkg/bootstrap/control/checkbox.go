package control

import (
	"encoding/gob"
	"github.com/goradd/goradd/pkg/html"
	"github.com/goradd/goradd/pkg/page"
	"reflect"
)

type Checkbox struct {
	checkboxBase
}

func NewCheckbox(parent page.ControlI, id string) *Checkbox {
	c := &Checkbox{}
	c.Init(c, parent, id)
	return c
}

func (c *Checkbox) ΩDrawingAttributes() *html.Attributes {
	a := c.checkboxBase.ΩDrawingAttributes()
	a.SetDataAttribute("grctl", "bs-checkbox")
	a.Set("name", c.ID()) // needed for posts
	a.Set("type", "checkbox")
	a.Set("value", "1") // required for html validity
	return a
}

// ΩUpdateFormValues is an internal call that lets us reflect the value of the checkbox on the web override
func (c *Checkbox) ΩUpdateFormValues(ctx *page.Context) {
	id := c.ID()

	if v, ok := ctx.CheckableValue(id); ok {
		c.SetCheckedNoRefresh(v)
	}
}

func (c *Checkbox) Serialize(e page.Encoder) (err error) {
	if err = c.checkboxBase.Serialize(e); err != nil {
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
	if err = c.checkboxBase.Deserialize(d, p); err != nil {
		return
	}

	return
}

func init () {
	gob.RegisterName("bootstrap.checkbox", new(Checkbox))
}
