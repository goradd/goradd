package control

import (
	"encoding/gob"
	"github.com/spekary/goradd/pkg/html"
	"github.com/spekary/goradd/pkg/page"
	"github.com/spekary/goradd/pkg/page/control"
)

type DateTextbox struct {
	control.DateTextbox
}

func NewDateTextbox(parent page.ControlI, id string) *DateTextbox {
	t := new (DateTextbox)
	t.Init(t, parent, id)
	return t
}


func (t *DateTextbox) DrawingAttributes() *html.Attributes {
	a := t.DateTextbox.DrawingAttributes()
	a.AddClass("form-control")
	return a
}


func init () {
	gob.RegisterName("bootstrap.datetextbox", new(DateTextbox))
}


