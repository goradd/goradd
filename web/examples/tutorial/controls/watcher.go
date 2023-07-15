package controls

import (
	"context"
	"github.com/goradd/goradd/pkg/page"
	. "github.com/goradd/goradd/pkg/page/control"
	"github.com/goradd/goradd/pkg/sys"
	"github.com/goradd/goradd/web/examples/tutorial"
	"path/filepath"
)

type WatcherPanel struct {
	Panel
}

func NewWatcherPanel(ctx context.Context, parent page.ControlI) page.ControlI {
	p := new(WatcherPanel)
	p.Panel.Init(p, parent, "")
	return p
}

func init() {
	dir := sys.SourceDirectory()
	tutorial.RegisterTutorialPage("controls", 0, "watcher", "Watching Database Changes", NewWatcherPanel,
		[]string{
			sys.SourcePath(),
			filepath.Join(dir, "watcher.tpl.got"),
		})
}

func init() {
	page.RegisterControl(&WatcherPanel{})
}
