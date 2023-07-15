package orm

import (
	"context"
	"github.com/goradd/goradd/pkg/page"
	. "github.com/goradd/goradd/pkg/page/control"
	"github.com/goradd/goradd/pkg/sys"
	"github.com/goradd/goradd/web/examples/tutorial"
	"path/filepath"
)

type LoadPanel struct {
	Panel
}

func NewLoadPanel(ctx context.Context, parent page.ControlI) page.ControlI {
	p := new(LoadPanel)
	p.Init(p, ctx, parent, "")
	return p
}

func (p *LoadPanel) Init(self any, ctx context.Context, parent page.ControlI, id string) {
	p.Panel.Init(self, parent, id)
}

func init() {
	page.RegisterControl(&LoadPanel{})

	dir := sys.SourceDirectory()
	tutorial.RegisterTutorialPage("orm", 2, "load", "Loading Individual Records", NewLoadPanel,
		[]string{
			sys.SourcePath(),
			filepath.Join(dir, "2-load.tpl.got"),
		})
}
