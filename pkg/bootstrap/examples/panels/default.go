package panels

import (
	"context"
	"github.com/goradd/goradd/pkg/bootstrap/config"
	"github.com/goradd/goradd/pkg/bootstrap/examples"
	"github.com/goradd/goradd/pkg/page"
	. "github.com/goradd/goradd/pkg/page/control"
)

type DefaultPanel struct {
	Panel
}

func NewDefaultPanel(ctx context.Context, parent page.ControlI) {
	p := &DefaultPanel{}
	p.Panel.Init(p, parent, "defaultPanel")
	config.LoadBootstrap(p.ParentForm())
}

func init() {
	examples.RegisterPanel("", "Home", NewDefaultPanel, 1)

	//browsertest.RegisterTestFunction("Plain Textbox", TestPlain)
}
