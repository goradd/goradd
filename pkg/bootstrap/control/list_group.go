package control

import (
	"context"
	"github.com/goradd/goradd/pkg/javascript"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/action"
	"github.com/goradd/goradd/pkg/page/control"
	"github.com/goradd/goradd/pkg/page/control/list"
	"github.com/goradd/goradd/pkg/page/event"
)

const (
	refreshAction = iota + 100
)

type ListGroupI interface {
	list.UnorderedListI
	SetIsSelectable(bool)
}

// A ListGroup implements the Bootstrap ListGroup control.
// Since just a static list isn't all that interesting, this is a dynamic list whose
// individual items are considered clickable. To conform with the
// html standard and accessibility rules, items should appear as anchors if they link to another page, but as buttons
// if they cause a different action, like popping up a dialog.
//
// Use the data provider to AddItems to the list, assigning attributes as needed to produce the items you want.
// You can also use a proxy control to create the attributes.
type ListGroup struct {
	list.UnorderedList
	isSelectable bool
	selectedID   string
}

func NewListGroup(parent page.ControlI, id string) *ListGroup {
	l := new(ListGroup)
	l.Init(l, parent, id)
	return l
}

func (l *ListGroup) Init(self any, parent page.ControlI, id string) {
	l.UnorderedList.Init(self, parent, id)
	l.Tag = "div"
	l.SetItemTag("a") // default to anchor tags. Change it to something else if needed.
	l.AddClass("list-group")

	// Set the Control action value to the item clicked on
	l.SetActionValue(javascript.JsCode("event.target.id"))
}

// SetIsSelectable sets whether the list group will remember and show the
// most recently selected item as selected.
// Do this AFTER you set the item tag.
func (l *ListGroup) SetIsSelectable(canSelect bool) {
	l.isSelectable = canSelect
	l.Refresh()
	l.PrivateOff()
	if canSelect {
		l.On(event.Click().Selector(l.ItemTag()).Private(), action.Ajax(l.ID(), refreshAction))
	}
}

// SelectedID returns the id of the currently selected item.
func (l *ListGroup) SelectedID() string {
	return l.selectedID
}

func (l *ListGroup) GetItemsHtml(items []*list.Item) string {
	// make sure the list items have the correct classes before drawing them
	for _, item := range items {
		item.Attributes().AddClass("list-group-item list-group-item-action")
		if l.isSelectable && l.selectedID == item.ID() {
			item.Attributes().AddClass(" active")
			item.Attributes().AddValues("aria-current", "true")
		}
	}
	return l.UnorderedList.GetItemsHtml(items)
}

func (l *ListGroup) DoPrivateAction(_ context.Context, a action.Params) {
	switch a.ID {
	case refreshAction:
		l.selectedID = a.ControlValueString()
		l.Refresh()
	}
}

func (l *ListGroup) Serialize(e page.Encoder) {
	l.UnorderedList.Serialize(e)
	if err := e.Encode(l.isSelectable); err != nil {
		panic(err)
	}
	if err := e.Encode(l.selectedID); err != nil {
		panic(err)
	}

}

func (l *ListGroup) Deserialize(dec page.Decoder) {
	l.UnorderedList.Deserialize(dec)
	if err := dec.Decode(&l.isSelectable); err != nil {
		panic(err)
	}
	if err := dec.Decode(&l.selectedID); err != nil {
		panic(err)
	}
}

type ListGroupCreator struct {
	ID string
	// Items is a static list of labels and values that will be in the list. Or, use a DataProvider to dynamically generate the items.
	Items []list.ListValue
	// DataProvider is the control that will dynamically provide the data for the list and that implements the DataBinder interface.
	DataProvider control.DataBinder
	// DataProviderID is the id of a control that will dynamically provide the data for the list and that implements the DataBinder interface.
	DataProviderID string
	page.ControlOptions
	// ItemTag is the tag of the items. It defaults to "a".
	ItemTag string
	// IsSelectable determines whether a clicked item will be shown as selected.
	IsSelectable bool
	// OnClick is the action to take when an item is clicked.
	// The id of the item will appear as the action's ControlValue.
	OnClick action.ActionI
}

// Create is called by the framework to create a new control from the Creator. You
// do not normally need to call this.
func (c ListGroupCreator) Create(ctx context.Context, parent page.ControlI) page.ControlI {
	ctrl := NewListGroup(parent, c.ID)
	c.Init(ctx, ctrl)
	return ctrl
}

func (c ListGroupCreator) Init(ctx context.Context, ctrl ListGroupI) {
	sub := list.UnorderedListCreator{
		Items:          c.Items,
		DataProvider:   c.DataProvider,
		ControlOptions: c.ControlOptions,
	}

	sub.Init(ctx, ctrl)

	if c.ItemTag != "" {
		ctrl.SetItemTag(c.ItemTag)
	}

	if c.OnClick != nil {
		ctrl.On(event.Click().Selector(ctrl.ItemTag()), c.OnClick)
		// Set the Control action value to the item clicked on
		ctrl.SetActionValue(javascript.JsCode("event.target.id"))
	}
	ctrl.SetIsSelectable(c.IsSelectable)
}

// GetListGroup is a convenience method to return the control with the given id from the page.
func GetListGroup(c page.ControlI, id string) *ListGroup {
	return c.Page().GetControl(id).(*ListGroup)
}

func init() {
	page.RegisterControl(&ListGroup{})
}
