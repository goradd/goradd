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
	p.Self = p
	p.Init(ctx, parent, "")
	return p
}

func (p *CrudPanel) Init(ctx context.Context, parent page.ControlI, id string) {
	p.Panel.Init(parent, id)
}


func init() {
	page.RegisterControl(&CrudPanel{})

	dir := sys.SourceDirectory()
	tutorial.RegisterTutorialPage("orm", 3, "crud", "Creating, Modifying and Deleting Database Objects", NewCrudPanel,
		[]string {
			sys.SourcePath(),
			filepath.Join(dir, "template_source", "4-crud.tpl.got"),
		})
}

