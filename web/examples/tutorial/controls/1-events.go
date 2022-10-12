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

	textBox1 := NewTextbox(p, "textBox1")
	textBox1.On(event.Input().Delay(1000), action.Message(javascript.JsCode("event.target.value")))

	btn1 := NewButton(p, "btn1")
	btn1.SetText("Click Me")
	btn1.OnClick(action.Message("btn1 clicked"))

	textBox2 := NewTextbox(p, "textBox2")
	textBox2.On(event.NewEvent("cut"), action.Message("textbox2 cut"))
	textBox2.SetText("Cut Me")

}
func init() {
	page.RegisterControl(&EventsPanel{})

	dir := sys.SourceDirectory()
	tutorial.RegisterTutorialPage("controls", 1, "events", "Events", NewEventsPanel,
		[]string{
			sys.SourcePath(),
			filepath.Join(dir, "1-events.tpl.got"),
		})
}
