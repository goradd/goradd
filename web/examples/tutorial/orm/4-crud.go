package orm

import (
	"context"
	"github.com/goradd/goradd/pkg/page"
	. "github.com/goradd/goradd/pkg/page/control"
	"github.com/goradd/goradd/pkg/sys"
	"github.com/goradd/goradd/web/examples/tutorial"
	"path/filepath"
)

type CrudPanel struct {
	Panel
}

func NewCrudPanel(ctx context.Context, parent page.ControlI) page.ControlI {
	p := &CrudPanel{}
	p.Panel.Init(p, parent, "")
	return p
}


func init() {
	page.RegisterControl(CrudPanel{})

	dir := sys.SourceDirectory()
	tutorial.RegisterTutorialPage("orm", 3, "crud", "Creating, Modifying and Deleting Database Objects", NewCrudPanel,
		[]string {
			sys.SourcePath(),
			filepath.Join(dir, "template_source", "4-crud.tpl.got"),
		})
}

