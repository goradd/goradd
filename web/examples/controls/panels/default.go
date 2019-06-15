package panels

import (
	"context"
	"github.com/goradd/goradd/pkg/page"
	. "github.com/goradd/goradd/pkg/page/control"
	"github.com/goradd/goradd/web/examples/controls"
)

type DefaultPanel struct {
	Panel
}

func NewDefaultPanel(ctx context.Context, parent page.ControlI) {
	p := &DefaultPanel{}
	p.Panel.Init(p, parent, "defaultPanel")
}

func init() {
	controls.RegisterPanel("", "Home", NewDefaultPanel, 1)
}
