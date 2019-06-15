package control

import (
	"encoding/gob"
	"github.com/goradd/goradd/pkg/html"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/control"
)

type IntegerTextbox struct {
	control.IntegerTextbox
}

func NewIntegerTextbox(parent page.ControlI, id string) *IntegerTextbox {
	t := new(IntegerTextbox)
	t.Init(t, parent, id)
	return t
}

func (t *IntegerTextbox) ΩDrawingAttributes() *html.Attributes {
	a := t.IntegerTextbox.ΩDrawingAttributes()
	a.AddClass("form-control")
	return a
}

func init() {
	gob.RegisterName("bootstrap.integertextbox", new(IntegerTextbox))
}
