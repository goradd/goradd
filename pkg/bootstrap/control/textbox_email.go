package control

import (
	"encoding/gob"
	"github.com/spekary/goradd/pkg/html"
	"github.com/spekary/goradd/pkg/page"
	"github.com/spekary/goradd/pkg/page/control"
)

type EmailTextbox struct {
	control.EmailTextbox
}

func NewEmailTextbox(parent page.ControlI, id string) *EmailTextbox {
	t := new (EmailTextbox)
	t.Init(t, parent, id)
	return t
}


func (t *EmailTextbox) DrawingAttributes() *html.Attributes {
	a := t.EmailTextbox.DrawingAttributes()
	a.AddClass("form-control")
	return a
}


func init () {
	gob.RegisterName("bootstrap.emailtextbox", new(EmailTextbox))
}


