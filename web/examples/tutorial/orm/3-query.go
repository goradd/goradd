package orm

import (
	"context"
	"github.com/goradd/goradd/pkg/page"
	. "github.com/goradd/goradd/pkg/page/control"
	"github.com/goradd/goradd/pkg/sys"
	"github.com/goradd/goradd/web/examples/tutorial"
	"path/filepath"
)

type QueryPanel struct {
	Panel
}

func NewQueryPanel(ctx context.Context, parent page.ControlI) page.ControlI {
	p := &QueryPanel{}
	p.Self = p
	p.Init(ctx, parent, "")
	return p
}

func (p *QueryPanel) Init(ctx context.Context, parent page.ControlI, id string) {
	p.Panel.Init(parent, id)
}


func init() {
	page.RegisterControl(&QueryPanel{})

	dir := sys.SourceDirectory()
	tutorial.RegisterTutorialPage("orm", 3, "query", "Querying the Database Using a QueryBuilder", NewQueryPanel,
		[]string {
			sys.SourcePath(),
			filepath.Join(dir, "3-query.tpl.got"),
		})
}

