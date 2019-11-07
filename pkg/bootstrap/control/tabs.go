package control

import (
	"context"
	"github.com/goradd/goradd/pkg/html"
	"github.com/goradd/goradd/pkg/javascript"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/action"
	"github.com/goradd/goradd/pkg/page/control"
	"github.com/goradd/goradd/pkg/page/event"
)

type TabsI interface {
	control.PanelI
}

// A Tabs pane draws its child controls as a set of tabs. The labels of the children serve as the tab labels.
// This currently draws everything at once, with the current panel visible, but everything else has hidden html.
type Tabs struct {
	control.Panel
	selectedID string // selected child id
}

// TODO: Modify this so that you can optionally show each panel through ajax

func NewTabs(parent page.ControlI, id string) *Tabs {
	t := &Tabs{}
	t.Self = t
	t.Init(parent, id)
	return t
}

func (l *Tabs) Init(parent page.ControlI, id string) {
	l.Panel.Init(parent, id)
	l.On(event.Event("show.bs.tab"), action.SetControlValue(l.ID(), "selectedId", javascript.JsCode("event.target.id")))
}

func (c *Tabs) DrawingAttributes(ctx context.Context) html.Attributes {
	a := c.Panel.DrawingAttributes(ctx)
	a.SetDataAttribute("grctl", "bs-tabs")
	return a
}


type TabsCreator struct {
	// ID is the control id of the html widget and must be unique to the page
	ID string
	page.ControlOptions
	Children []page.Creator
}

func (c TabsCreator) Create(ctx context.Context, parent page.ControlI) page.ControlI {
	ctrl := NewTabs(parent, c.ID)
	ctrl.ApplyOptions(ctx, c.ControlOptions)
	ctrl.AddControls(ctx, c.Children...)
	return ctrl
}


// GetTabs is a convenience method to return the control with the given id from the page.
func GetTabs(c page.ControlI, id string) *Tabs {
	return c.Page().GetControl(id).(*Tabs)
}

func init() {
	page.RegisterControl(&Tabs{})
}
