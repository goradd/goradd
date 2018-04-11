package control

import (
	"goradd/page/control_base"
	"github.com/spekary/goradd/page"
)

const (
	TEXTBOX_TYPE_DEFAULT   = "text"
	TEXTBOX_TYPE_PASSWORD    = "password"
	TEXTBOX_TYPE_SEARCH    = "search"
	TEXTBOX_TYPE_NUMBER = "number" // Puts little arrows in box, will need to widen it.
	TEXTBOX_TYPE_EMAIL = "email" // see TextEmail. Prevents submission of RFC5322 email addresses (Gogh Fir <gf@example.com>)
	TEXTBOX_TYPE_TEL = "tel" // not well supported
	TEXTBOX_TYPE_URL = "url"
)

// Text is a basic text entry form item.
type Textbox struct {
	control_base.Textbox
}

func NewTextbox(parent page.ControlI) *Textbox {
	t := &Textbox{}
	t.Init(t, parent)
	return t
}
