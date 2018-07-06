package control

import (
	"goradd/override/control_base"
	"github.com/spekary/goradd/page"
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
