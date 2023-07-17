package orm

import (
	"context"
	"github.com/goradd/goradd/pkg/page"
	. "github.com/goradd/goradd/pkg/page/control"
	"github.com/goradd/goradd/pkg/sys"
	"github.com/goradd/goradd/web/examples/tutorial"
	"path/filepath"
)

type RefPanel struct {
	Panel
}

func NewRefPanel(ctx context.Context, parent page.ControlI) page.ControlI {
	p := new(RefPanel)
	p.Init(p, ctx, parent, "")
	return p
}

func (p *RefPanel) Init(self any, ctx context.Context, parent page.ControlI, id string) {
	p.Panel.Init(self, parent, id)
}

func init() {
	page.RegisterControl(&RefPanel{})

	dir := sys.SourceDirectory()
	tutorial.RegisterTutorialPage("orm", 5, "ref", "References", NewRefPanel,
		[]string{
			sys.SourcePath(),
			filepath.Join(dir, "5-ref.tpl.got"),
		})
}
