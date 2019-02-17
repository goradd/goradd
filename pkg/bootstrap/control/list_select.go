package control

import (
	"encoding/gob"
	"github.com/goradd/goradd/pkg/html"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/control"
)

type SelectList struct {
	control.SelectList
}

func NewSelectList(parent page.ControlI, id string) *SelectList {
	t := new (SelectList)
	t.Init(t, parent, id)
	return t
}


func (t *SelectList) ΩDrawingAttributes() *html.Attributes {
	a := t.SelectList.ΩDrawingAttributes()
	a.AddClass("form-control")
	return a
}


func init () {
	gob.RegisterName("bootstrap.selectlist", new(SelectList))
}


