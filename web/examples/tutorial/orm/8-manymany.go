package orm

import (
	"context"
	"github.com/goradd/goradd/pkg/page"
	. "github.com/goradd/goradd/pkg/page/control"
	"github.com/goradd/goradd/pkg/sys"
	"github.com/goradd/goradd/web/examples/tutorial"
	"path/filepath"
)

type ManyManyPanel struct {
	Panel
}

func NewManyManyPanel(ctx context.Context, parent page.ControlI) page.ControlI {
	p := new(ManyManyPanel)
	p.Init(p, ctx, parent, "")
	return p
}

func (p *ManyManyPanel) Init(self any, ctx context.Context, parent page.ControlI, id string) {
	p.Panel.Init(self, parent, id)
}

func init() {
	page.RegisterControl(&ManyManyPanel{})

	dir := sys.SourceDirectory()
	tutorial.RegisterTutorialPage("orm", 8, "manymany", "Many-to-Many References", NewManyManyPanel,
		[]string{
			sys.SourcePath(),
			filepath.Join(dir, "8-manymany.tpl.got"),
		})
}
