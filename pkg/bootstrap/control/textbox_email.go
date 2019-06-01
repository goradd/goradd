package control

import (
	"encoding/gob"
	"github.com/goradd/goradd/pkg/html"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/control"
)

type EmailTextbox struct {
	control.EmailTextbox
}

func NewEmailTextbox(parent page.ControlI, id string) *EmailTextbox {
	t := new(EmailTextbox)
	t.Init(t, parent, id)
	return t
}

func (t *EmailTextbox) ΩDrawingAttributes() *html.Attributes {
	a := t.EmailTextbox.ΩDrawingAttributes()
	a.AddClass("form-control")
	return a
}

func init() {
	gob.RegisterName("bootstrap.emailtextbox", new(EmailTextbox))
}
