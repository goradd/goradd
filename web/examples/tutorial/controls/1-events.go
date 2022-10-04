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

type EventsPanel struct {
	Panel
}

func NewEventsPanel(ctx context.Context, parent page.ControlI) page.ControlI {
	p := &EventsPanel{}
	p.Self = p
	p.Init(ctx, parent, "")
	return p
}

func (p *EventsPanel) Init(ctx context.Context, parent page.ControlI, id string) {
	p.Panel.Init(parent, id)

	textBox := NewTextbox(p, "textField")
	textBox.On(event.Input().Delay(1000), action.Message(javascript.JsCode("event.target.value")))

	btn := NewButton(p, "btn1")
	btn.OnClick(action.Ajax(p.ID(), 1))
}

func init() {
	page.RegisterControl(&EventsPanel{})

	dir := sys.SourceDirectory()
	tutorial.RegisterTutorialPage("controls", 0, "intro", "Introduction to GoRADD Controls", NewEventsPanel,
		[]string{
			sys.SourcePath(),
			filepath.Join(dir, "0-intro.tpl.got"),
		})
}
