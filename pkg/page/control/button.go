package control

import (
	"goradd-project/override/control_base"
	"github.com/spekary/goradd/pkg/page"
	"reflect"
)

type ButtonI interface {
	control_base.ButtonI
}

// Button is a standard html button. It derives from the button in the override class, allowing you to customize the
// behavior of all buttons in your application.
type Button struct {
	control_base.Button
}

// Creates a new standard html button
func NewButton(parent page.ControlI, id string) *Button {
	b := &Button{}
	b.Init(b, parent, id)
	return b
}

func (b *Button) Serialize(e page.Encoder) (err error) {
	if err = b.Control.Serialize(e); err != nil {
		return
	}

	return
}

// ΩisSerializer is used by the automated control serializer to determine how far down the control chain the control
// has to go before just calling serialize and deserialize
func (b *Button) ΩisSerializer(i page.ControlI) bool {
	return reflect.TypeOf(b) == reflect.TypeOf(i)
}


func (b *Button) Deserialize(d page.Decoder, p *page.Page) (err error) {
	if err = b.Control.Deserialize(d, p); err != nil {
		return
	}

	return
}
