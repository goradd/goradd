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

type ActionsPanel struct {
	Panel
}

func NewActionsPanel(ctx context.Context, parent page.ControlI) page.ControlI {
	p := &ActionsPanel{}
	p.Self = p
	p.Init(ctx, parent, "")
	return p
}

func (p *ActionsPanel) Init(ctx context.Context, parent page.ControlI, id string) {
	p.Panel.Init(parent, id)

	textbox1 := NewTextbox(p, "textbox1")
	textbox1.On(event.Input(), action.Javascript("event.target.value = event.target.value.toUpperCase()"))

	btn1 := NewButton(p, "okButton").SetText("OK")
	btn1.On(event.Click(), action.Message(javascript.JsCode("event.target.value")+" was clicked"))
}

func init() {
	page.RegisterControl(&ActionsPanel{})

	dir := sys.SourceDirectory()
	tutorial.RegisterTutorialPage("controls", 2, "actions", "Actions", NewActionsPanel,
		[]string{
			sys.SourcePath(),
			filepath.Join(dir, "2-actions.tpl.got"),
		})
}
