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

// TabSelect is the event generated when a tab is selected.
// The EventValueString of the event is the id of the child item that was selected.
const TabSelect = "gr-bs-tabselect"

// TabHidden is the event generated when a tab is hidden. i.e. another tab is selected when this one was in front.
// The EventValueString of the event is the id of the child item that was hidden.
const TabHidden = "gr-bs-tabhide"

const (
	TabStyleNone      = ""
	TabStyleTabs      = "tab"
	TabStylePills     = "pill"
	TabStyleUnderline = "underline"
)

type TabsI interface {
	control.PanelI
}

// Tabs draws its child controls as a set of tabs. The Text value of the children serve as the tab labels.
//
// This currently draws everything at once, with the current panel visible, but everything else has hidden html.
// By default, it will save its state so that refreshes of the page, or coming back to the same page, will
// cause it to remember what tab it was on last.
//
// Whenever a tab is clicked, the TabSelect event will be fired via an Ajax call. Containers of the Tabs control
// can listen for this event by simply implementing DoAction without needing to create a TabSelect event listener.
//
// If you need to know that a tab is being hidden, you can watch the TabHidden event.
//
// Call SelectedId to return the id of the currently selected tab.
//
// Call SetTabStyle to set the style to be tabs, pills or underline.
//
// The tab structure is surrounded by a Card, and the content is drawn in a div with class card-body.
type Tabs struct {
	control.Panel
	selectedID string // selected child id
	tabStyle   string
}

func NewTabs(ctx context.Context, parent page.ControlI, id string) *Tabs {
	t := new(Tabs)
	t.Init(ctx, t, parent, id)
	return t
}

func (t *Tabs) Init(ctx context.Context, self any, parent page.ControlI, id string) {
	t.Panel.Init(self, parent, id)
	t.tabStyle = TabStyleTabs // default to tabs
	t.SaveState(ctx, true)

	// Trigger tab select
	// The -4 below is to remove the "_tab" suffix at the end of the id
	t.On(event.NewEvent("show.bs.tab").Action(
		action.Trigger(t.ID(), TabSelect, javascript.JsCode(`event.target.id.substring(0, event.target.id.length - 4)`)),
	))

	// Trigger tab hide
	// The -4 below is to remove the "_tab" suffix at the end of the id
	t.On(event.NewEvent("hide.bs.tab").Action(
		action.Trigger(t.ID(), TabHidden, javascript.JsCode(`event.target.id.substring(0, event.target.id.length - 4)`)),
	))

	// Watch tab select so that we can record the value
	t.On(TabSelectEvent().Action(action.Do()))
}

func (t *Tabs) SelectedId() string {
	return t.selectedID
}

// SetTabStyle sets the style of the tabs.
//
// Choose one of:
//   - TabStyleNone
//   - TabStyleTabs
//   - TabStylePills
//   - TabStyleUnderline
func (t *Tabs) SetTabStyle(s string) {
	t.tabStyle = s
}

func (t *Tabs) DrawingAttributes(ctx context.Context) html5tag.Attributes {
	a := t.Panel.DrawingAttributes(ctx)
	a.SetData("grctl", "bs-tabs")
	a.AddClass("card")
	return a
}

func (t *Tabs) DoAction(ctx context.Context, a action.Params) {
	if a.Event == TabSelect {
		t.selectedID = a.EventValueString()
	}
	t.Parent().DoAction(ctx, a)
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

// MarshalState is an internal function to save the state of the control
func (t *Tabs) MarshalState(m page.SavedState) {
	m.Set("sel", t.selectedID)
}

// UnmarshalState is an internal function to restore the state of the control
func (t *Tabs) UnmarshalState(m page.SavedState) {
	if v, ok := m.Load("sel"); ok {
		if s, ok2 := v.(string); ok2 {
			t.selectedID = s
		}
	}
}

// TabsCreator is the declarative definition of a tab pane.
//
// Set the Children to control.Panel objects that will contain the content of each tab shown. The Text value of each
// Panel will be the label of the tab.
//
// TabStyle should be one of:
//   - TabStyleNone
//   - TabStyleTabs
//   - TabStylePills
//   - TabStyleUnderline
type TabsCreator struct {
	// ID is the control id of the html widget and must be unique to the page
	ID string
	page.ControlOptions
	Children []page.Creator
	TabStyle string
}

func (c TabsCreator) Create(ctx context.Context, parent page.ControlI) page.ControlI {
	ctrl := NewTabs(ctx, parent, c.ID)
	ctrl.ApplyOptions(ctx, c.ControlOptions)
	ctrl.AddControls(ctx, c.Children...)
	ctrl.SetTabStyle(c.TabStyle)
	return ctrl
}

// GetTabs is a convenience method to return the control with the given id from the page.
func GetTabs(c page.ControlI, id string) *Tabs {
	return c.Page().GetControl(id).(*Tabs)
}

func init() {
	page.RegisterControl(&Tabs{})
}

func TabSelectEvent() *event.Event {
	return event.NewEvent(TabSelect)
}
