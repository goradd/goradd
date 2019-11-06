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
	p.Panel.Init(p, parent, "")
	return p
}


func init() {
	page.RegisterControl(RefPanel{})

	dir := sys.SourceDirectory()
	tutorial.RegisterTutorialPage("orm", 3, "ref", "References", NewRefPanel,
		[]string {
			sys.SourcePath(),
			filepath.Join(dir, "template_source", "5-ref.tpl.got"),
		})
}

