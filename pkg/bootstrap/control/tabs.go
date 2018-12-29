package control

import (
	"github.com/goradd/goradd/pkg/html"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/control"
	"github.com/goradd/goradd/pkg/page/event"
	"github.com/goradd/goradd/pkg/page/action"
	"github.com/goradd/goradd/pkg/javascript"
)

// A Tabs pane draws its child controls as a set of tabs. The labels of the children serve as the tab labels.
// This currently draws everything at once, with the current panel visible, but everything else has hidden html.
type Tabs struct {
	control.Panel
	selectedID string // selected child id
}

// TODO: Modify this so that you can optionally show each panel through ajax

func NewTabs(parent page.ControlI, id string) *Tabs {
	t := &Tabs{}
	t.Init(t, parent, id)
	return t
}

func (l *Tabs) Init(self page.ControlI, parent page.ControlI, id string) {
	l.Panel.Init(self, parent, id)
	l.On(event.Event("show.bs.tab"), action.SetControlValue(l.ID(), "selectedId", javascript.JsCode("event.target.id")))
}

func (c *Tabs) DrawingAttributes() *html.Attributes {
	a := c.Panel.DrawingAttributes()
	a.SetDataAttribute("grctl", "bs-tabs")
	return a
}


