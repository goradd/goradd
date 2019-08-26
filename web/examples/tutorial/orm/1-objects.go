package orm

import (
	"context"
	"github.com/goradd/goradd/pkg/page"
	. "github.com/goradd/goradd/pkg/page/control"
	"github.com/goradd/goradd/pkg/sys"
	"github.com/goradd/goradd/web/examples/tutorial"
	"path/filepath"
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
	dir := sys.SourceDirectory()
	tutorial.RegisterTutorialPage("orm", 1, "objects", "Code-generated objects", NewObjectsPanel,
		[]string {
			sys.SourcePath(),
			filepath.Join(dir, "template_source", "1-objects.tpl.got"),
		})
}

