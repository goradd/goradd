package control

import (
	"encoding/gob"
	"github.com/spekary/goradd/pkg/html"
	"github.com/spekary/goradd/pkg/page"
	"github.com/spekary/goradd/pkg/page/control"
	"reflect"
)

type Textbox struct {
	control.Textbox
	BootstrapControl
}

func NewTextbox(parent page.ControlI, id string) *Textbox {
	t := new (Textbox)
	t.Init(t, parent, id)
	return t
}


func (t *Textbox) DrawingAttributes() *html.Attributes {
	a := t.Textbox.DrawingAttributes()
	t.BootstrapControl.AddBootstrapAttributes(t, a)
	return a
}

func (t *Textbox) Serialize(e page.Encoder) (err error) {
	if err = t.Textbox.Serialize(e); err != nil {
		return
	}
	if err = t.BootstrapControl.Serialize(e); err != nil {
		return
	}
	return
}

// ΩisSerializer is used by the automated control serializer to determine how far down the control chain the control
// has to go before just calling serialize and deserialize
func (t *Textbox) ΩisSerializer(i page.ControlI) bool {
	return reflect.TypeOf(t) == reflect.TypeOf(i)
}


func (t *Textbox) Deserialize(d page.Decoder, p *page.Page) (err error) {
	if err = t.Textbox.Deserialize(d, p); err != nil {
		return
	}
	if err = t.BootstrapControl.Deserialize(d, p); err != nil {
		return
	}
	return
}


func init () {
	gob.RegisterName("bootstrap.textbox", new(Textbox))
}


// TODO: The same for the other kinds of textboxes