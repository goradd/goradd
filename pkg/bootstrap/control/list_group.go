package control

import (
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/control"
)

// A ListGroup implements the Bootstrap ListGroup control.
// Since just a static list isn't all that interesting, this is a dynamic list whose
// individual items are considered clickable. To conform with the
// html standard and accessibility rules, items should appear as anchors if they link to another page, but as buttons
// if they cause a different action, like popping up a dialog.
//
// Use the data provider to AddItems to the list, assigning attributes as needed to produce the items you want.
// You can also use a proxy control to create the attributes.
type ListGroup struct {
	control.UnorderedList
}

func NewListGroup(parent page.ControlI, id string) *ListGroup {
	l := &ListGroup{}
	l.Init(l, parent, id)
	return l
}

func (l *ListGroup) Init(self page.ControlI, parent page.ControlI, id string) {
	l.UnorderedList.Init(self, parent, id)
	l.Tag = "div"
	l.SetItemTag("a") // default to anchor tags. Change it to something else if needed.
	l.AddClass("list-group")
}

func (l *ListGroup) GetItemsHtml(items []control.ListItemI) string {
	// make sure the list items have the correct classes before drawing them

	for _, item := range items {
		item.Attributes().AddClass("list-group-item list-group-item-action")
	}
	return l.UnorderedList.GetItemsHtml(items)
}
