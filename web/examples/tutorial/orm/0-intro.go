package orm

import (
	"context"
	"github.com/goradd/goradd/pkg/page"
	. "github.com/goradd/goradd/pkg/page/control"
	"github.com/goradd/goradd/pkg/sys"
	"github.com/goradd/goradd/web/examples/tutorial"
	"path/filepath"
)

type IntroPanel struct {
	Panel
}

func NewIntroPanel(ctx context.Context, parent page.ControlI) page.ControlI {
	p := &IntroPanel{}
	p.Self = p
	p.Init(ctx, parent, "")
	return p
}

func (p *IntroPanel) Init(ctx context.Context, parent page.ControlI, id string) {
	p.Panel.Init(parent, id)
}


func init() {
	page.RegisterControl(&IntroPanel{})

	dir := sys.SourceDirectory()
	tutorial.RegisterTutorialPage("orm", 0, "intro", "Introduction to the ORM", NewIntroPanel,
		[]string {
			sys.SourcePath(),
			filepath.Join(dir, "0-intro.tpl.got"),
		})
}

