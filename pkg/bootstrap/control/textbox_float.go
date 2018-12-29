package control

import (
	"encoding/gob"
	"github.com/goradd/goradd/pkg/html"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/control"
)

type FloatTextbox struct {
	control.FloatTextbox
}

func NewFloatTextbox(parent page.ControlI, id string) *FloatTextbox {
	t := new (FloatTextbox)
	t.Init(t, parent, id)
	return t
}


func (t *FloatTextbox) DrawingAttributes() *html.Attributes {
	a := t.FloatTextbox.DrawingAttributes()
	a.AddClass("form-control")
	return a
}


func init () {
	gob.RegisterName("bootstrap.floattextbox", new(FloatTextbox))
}


