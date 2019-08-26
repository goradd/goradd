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
	p := &LoadPanel{}
	p.Panel.Init(p, parent, "")
	return p
}


func init() {
	dir := sys.SourceDirectory()
	tutorial.RegisterTutorialPage("orm", 2, "load", "Loading Individual Records", NewLoadPanel,
		[]string {
			sys.SourcePath(),
			filepath.Join(dir, "template_source", "2-load.tpl.got"),
		})
}

