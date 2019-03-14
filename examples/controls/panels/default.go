package panels

import (
	"context"
	"github.com/goradd/goradd/pkg/page"
	. "github.com/goradd/goradd/pkg/page/control"
)


type DefaultPanel struct {
	Panel
}

func NewDefaultPanel(ctx context.Context, parent page.ControlI) *DefaultPanel {
	p := &DefaultPanel{}
	p.Panel.Init(p, parent, "defaultPanel")

	return p
}


func init() {
	//browsertest.RegisterTestFunction("Plain Textbox", TestPlain)
}
