package orm

import (
	"context"
	"github.com/goradd/goradd/pkg/page"
	. "github.com/goradd/goradd/pkg/page/control"
	"github.com/goradd/goradd/web/examples/tutorial"
)

type ObjectsPanel struct {
	Panel
}

func NewObjectsPanel(ctx context.Context, parent page.ControlI) page.ControlI {
	p := &ObjectsPanel{}
	p.Panel.Init(p, parent, "")
	return p
}


func init() {
	tutorial.RegisterTutorialPage("orm", 1, "Code-generated objects", NewObjectsPanel)
}

