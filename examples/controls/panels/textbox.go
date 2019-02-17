package panels

import (
	"github.com/goradd/goradd/pkg/page"
	. "github.com/goradd/goradd/pkg/page/control"
)


type TextboxPanel struct {
	Panel
	PlainText   *Textbox
	IntegerText *IntegerTextbox
	FloatText   *FloatTextbox
	Submit      *Button
}

func NewTextboxPanel(parent page.ControlI, id string) *TextboxPanel {
	p := &TextboxPanel{}
	p.Panel.Init(p, parent, id)

	p.PlainText = NewTextbox(p, "plain")
	return p
}


func init() {
	//browsertest.RegisterTestFunction("Plain Textbox", TestPlain)
}
