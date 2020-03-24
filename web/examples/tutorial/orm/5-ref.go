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
	p := &RefPanel{}
	p.Self = p
	p.Init(ctx, parent, "")
	return p
}

func (p *RefPanel) Init(ctx context.Context, parent page.ControlI, id string) {
	p.Panel.Init(parent, id)
}


func init() {
	page.RegisterControl(&RefPanel{})

	dir := sys.SourceDirectory()
	tutorial.RegisterTutorialPage("orm", 5, "ref", "References", NewRefPanel,
		[]string {
			sys.SourcePath(),
			filepath.Join(dir, "5-ref.tpl.got"),
		})
}

