package orm

import (
	"context"
	"github.com/goradd/goradd/pkg/page"
	. "github.com/goradd/goradd/pkg/page/control"
	"github.com/goradd/goradd/pkg/sys"
	"github.com/goradd/goradd/web/examples/tutorial"
	"path/filepath"
)

type OneOnePanel struct {
	Panel
}

func NewOneOnePanel(ctx context.Context, parent page.ControlI) page.ControlI {
	p := &OneOnePanel{}
	p.Self = p
	p.Init(ctx, parent, "")
	return p
}

func (p *OneOnePanel) Init(ctx context.Context, parent page.ControlI, id string) {
	p.Panel.Init(parent, id)
}


func init() {
	page.RegisterControl(&OneOnePanel{})

	dir := sys.SourceDirectory()
	tutorial.RegisterTutorialPage("orm", 7, "onetoone", "One-to-One References", NewOneOnePanel,
		[]string {
			sys.SourcePath(),
			filepath.Join(dir, "7-onetoone.tpl.got"),
		})
}

