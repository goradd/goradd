package controls

import (
	"context"
	"github.com/goradd/goradd/pkg/javascript"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/action"
	. "github.com/goradd/goradd/pkg/page/control"
	. "github.com/goradd/goradd/pkg/page/control/button"
	. "github.com/goradd/goradd/pkg/page/control/textbox"
	"github.com/goradd/goradd/pkg/page/event"
	"github.com/goradd/goradd/pkg/sys"
	"github.com/goradd/goradd/web/examples/tutorial"
	"path/filepath"
	"time"
)

type ActionsPanel struct {
	Panel
}

func NewActionsPanel(ctx context.Context, parent page.ControlI) page.ControlI {
	p := &ActionsPanel{}
	p.Init(p, ctx, parent, "")
	return p
}

func (p *ActionsPanel) Init(self any, ctx context.Context, parent page.ControlI, id string) {
	p.Panel.Init(self, parent, id)

	textbox1 := NewTextbox(p, "textbox1")
	textbox1.On(event.Input(), action.Javascript("event.target.value = event.target.value.toUpperCase()"))

	btn1 := NewButton(p, "serverTimeButton").SetText("Get Server Time")
	btn1.On(event.Click().Action(action.Do().ID(1000)))

	btn2 := NewButton(p, "clientTimeButton").SetText("Get Client Time")
	btn2.On(event.Click().Action(action.Do().ID(1001).ActionValue(javascript.NewClosureCall(
		`var today = new DateTextbox(); return today.getHours() + ':' + today.getMinutes();`, "",
	))))

	span1 := NewSpan(p, "timeSpan")
	span1.SetText("Unknown - click the button")
}

func (p *ActionsPanel) DoAction(ctx context.Context, a action.Params) {
	span1 := GetSpan(p, "timeSpan")
	switch a.ID {
	case 1000:
		t := time.Now()
		span1.SetText("Server time = " + t.String())
	case 1001:
		span1.SetText("Client time = " + a.ActionValueString())
	}
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
