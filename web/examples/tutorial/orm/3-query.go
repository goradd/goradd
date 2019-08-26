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
	p.Panel.Init(p, parent, "")
	return p
}


func init() {
	dir := sys.SourceDirectory()
	tutorial.RegisterTutorialPage("orm", 3, "query", "Querying the Database Using a QueryBuilder", NewQueryPanel,
		[]string {
			sys.SourcePath(),
			filepath.Join(dir, "template_source", "3-query.tpl.got"),
		})
}

