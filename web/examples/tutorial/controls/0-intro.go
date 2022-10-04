package controls

import (
	"context"
	"github.com/goradd/goradd/pkg/javascript"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/action"
	. "github.com/goradd/goradd/pkg/page/control"
	"github.com/goradd/goradd/pkg/page/event"
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

	textBox := NewTextbox(p, "textField")
	textBox.On(event.Input().Delay(1000), action.Message(javascript.JsCode("event.target.value")))
}

func init() {
	page.RegisterControl(&IntroPanel{})

	dir := sys.SourceDirectory()
	tutorial.RegisterTutorialPage("controls", 0, "intro", "Introduction to GoRADD Controls", NewIntroPanel,
		[]string{
			sys.SourcePath(),
			filepath.Join(dir, "0-intro.tpl.got"),
		})
}
