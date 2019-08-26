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
	p.Panel.Init(p, parent, "")
	return p
}


func init() {
	dir := sys.SourceDirectory()
	tutorial.RegisterTutorialPage("orm", 0, "intro", "Introduction to the ORM", NewIntroPanel,
		[]string {
			sys.SourcePath(),
			filepath.Join(dir, "template_source", "0-intro.tpl.got"),
		})
}

