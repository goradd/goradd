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

type ListGroupI interface {
	list.UnorderedListI
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
}

func NewListGroup(parent page.ControlI, id string) *ListGroup {
	l := &ListGroup{}
	l.Self = l
	l.Init(parent, id)
	return l
}

func (l *ListGroup) Init(parent page.ControlI, id string) {
	l.UnorderedList.Init(parent, id)
	l.Tag = "div"
	l.SetItemTag("a") // default to anchor tags. Change it to something else if needed.
	l.AddClass("list-group")
}

func (l *ListGroup) GetItemsHtml(items []*list.Item) string {
	// make sure the list items have the correct classes before drawing them
	for _, item := range items {
		item.Attributes().AddClass("list-group-item list-group-item-action")
	}
	return l.UnorderedList.GetItemsHtml(items)
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
}

// GetListGroup is a convenience method to return the control with the given id from the page.
func GetListGroup(c page.ControlI, id string) *ListGroup {
	return c.Page().GetControl(id).(*ListGroup)
}

func init() {
	page.RegisterControl(&ListGroup{})
}
