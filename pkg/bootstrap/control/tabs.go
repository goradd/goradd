package control

import (
	"context"
	"github.com/goradd/goradd/pkg/javascript"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/action"
	"github.com/goradd/goradd/pkg/page/control"
	"github.com/goradd/goradd/pkg/page/event"
	"github.com/goradd/html5tag"
)

type TabsI interface {
	control.PanelI
}

// Tabs draws its child controls as a set of tabs. The labels of the children serve as the tab labels.
// This currently draws everything at once, with the current panel visible, but everything else has hidden html.
type Tabs struct {
	control.Panel
	selectedID string // selected child id
}

// TODO: Modify this so that you can optionally show each panel through ajax

func NewTabs(parent page.ControlI, id string) *Tabs {
	t := new(Tabs)
	t.Init(t, parent, id)
	return t
}

func (t *Tabs) Init(self any, parent page.ControlI, id string) {
	t.Panel.Init(self, parent, id)
	t.On(event.NewEvent("show.bs.tab"), action.SetControlValue(t.ID(), "selectedId", javascript.JsCode("event.target.id")))
}

func (t *Tabs) DrawingAttributes(ctx context.Context) html5tag.Attributes {
	a := t.Panel.DrawingAttributes(ctx)
	a.SetData("grctl", "bs-tabs")
	return a
}

func (t *Tabs) Serialize(e page.Encoder) {
	t.Panel.Serialize(e)

	if err := e.Encode(t.selectedID); err != nil {
		panic(err)
	}
}

func (t *Tabs) Deserialize(d page.Decoder) {
	t.Panel.Deserialize(d)

	if err := d.Decode(&t.selectedID); err != nil {
		panic(err)
	}
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
