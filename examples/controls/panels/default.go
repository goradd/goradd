package panels

import (
	"github.com/goradd/goradd/pkg/page"
	. "github.com/goradd/goradd/pkg/page/control"
)


type DefaultPanel struct {
	Panel
}

func NewDefaultPanel(parent page.ControlI, id string) *DefaultPanel {
	p := &DefaultPanel{}
	p.Panel.Init(p, parent, id)

	return p
}


func init() {
	//browsertest.RegisterTestFunction("Plain Textbox", TestPlain)
}
