package orm

import (
	"context"
	"github.com/goradd/goradd/pkg/page"
	. "github.com/goradd/goradd/pkg/page/control"
	"github.com/goradd/goradd/pkg/sys"
	"github.com/goradd/goradd/web/examples/tutorial"
	"path/filepath"
)

type TypesPanel struct {
	Panel
}

func NewTypesPanel(ctx context.Context, parent page.ControlI) page.ControlI {
	p := &TypesPanel{}
	p.Self = p
	p.Init(ctx, parent, "")
	return p
}

func (p *TypesPanel) Init(ctx context.Context, parent page.ControlI, id string) {
	p.Panel.Init(parent, id)
}


func init() {
	page.RegisterControl(&TypesPanel{})

	dir := sys.SourceDirectory()
	tutorial.RegisterTutorialPage("orm", 9, "types", "Type Tables", NewTypesPanel,
		[]string {
			sys.SourcePath(),
			filepath.Join(dir, "9-types.tpl.got"),
		})
}

