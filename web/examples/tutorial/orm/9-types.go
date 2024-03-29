package orm

import (
	"context"
	"github.com/goradd/goradd/pkg/page"
	. "github.com/goradd/goradd/pkg/page/control"
	"github.com/goradd/goradd/pkg/sys"
	"github.com/goradd/goradd/web/examples/tutorial"
	"path/filepath"
)

type EnumsPanel struct {
	Panel
}

func NewEnumsPanel(ctx context.Context, parent page.ControlI) page.ControlI {
	p := new(EnumsPanel)
	p.Init(p, ctx, parent, "")
	return p
}

func (p *EnumsPanel) Init(self any, ctx context.Context, parent page.ControlI, id string) {
	p.Panel.Init(self, parent, id)
}

func init() {
	page.RegisterControl(&EnumsPanel{})

	dir := sys.SourceDirectory()
	tutorial.RegisterTutorialPage("orm", 9, "types", "Type Tables", NewEnumsPanel,
		[]string{
			sys.SourcePath(),
			filepath.Join(dir, "9-types.tpl.got"),
		})
}
