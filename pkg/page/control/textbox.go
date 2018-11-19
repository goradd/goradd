package control

import (
	"encoding/gob"
	"github.com/spekary/goradd/pkg/page"
	"goradd-project/override/control_base"
)

const (
	TextboxTypeDefault  = "text"
	TextboxTypePassword = "password"
	TextboxTypeSearch   = "search"
	TextboxTypeNumber   = "number" // Puts little arrows in box, will need to widen it.
	TextboxTypeEmail    = "email"  // see TextEmail. Prevents submission of RFC5322 email addresses (Gogh Fir <gf@example.com>)
	TextboxTypeTel    = "tel"    // not well supported
	TextboxTypeUrl    = "url"
)

// Text is a basic text entry form item.
type Textbox struct {
	control_base.Textbox
}

type TextboxI interface {
	control_base.TextboxI
}

func NewTextbox(parent page.ControlI, id string) *Textbox {
	t := &Textbox{}
	t.Init(t, parent, id)
	return t
}

func init () {
	gob.Register(&Textbox{})
}


