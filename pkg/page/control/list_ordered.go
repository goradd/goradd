package control

import (
	"bytes"
	"context"
	"github.com/goradd/goradd/pkg/html"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/control/data"
	html2 "html"
	"strconv"
)

type OrderedListI interface {
	UnorderedListI
	SetNumberType(t string) OrderedListI
	SetStart(start int) OrderedListI
}

// OrderedList is a dynamically generated html ordered list (ol). Such lists are often used as the basis for
// javascript and css widgets. If you use a data provider to set the data, you should call AddItems to the list
// in your LoadData function.
type OrderedList struct {
	UnorderedList
}

const (
	OrderedListNumberTypeNumber      = "1" // default
	OrderedListNumberTypeUpperLetter = "A"
	OrderedListNumberTypeLowerLetter = "a"
	OrderedListNumberTypeUpperRoman  = "I"
	OrderedListNumberTypeLowerRoman  = "i"
)

// NewOrderedList creates a new ordered list (ol tag).
func NewOrderedList(parent page.ControlI, id string) *OrderedList {
	t := &OrderedList{}
	t.Init(t, parent, id)
	return t
}

func (l *OrderedList) Init(self page.ControlI, parent page.ControlI, id string) {
	l.UnorderedList.Init(self, parent, id)
	l.Tag = "ol"
}

// this() supports object oriented features by giving easy access to the virtual function interface.
func (l *OrderedList) this() OrderedListI {
	return l.Self.(OrderedListI)
}


// SetNumberType sets the top level number style for the list. Choose from the OrderedListNumberType* constants.
// To set a number type for a sublevel, set the "type" attribute on the list item that is the parent of the sub list.
func (l *OrderedList) SetNumberType(t string) OrderedListI {
	l.SetAttribute("type", l)
	return l.this()
}

// SetStart sets the starting number for the numbers in the top level list. To set the start of a sub-list, set
// the "start" attribute on the list item that is the parent of the sub-list.
func (l *OrderedList) SetStart(start int) OrderedListI {
	l.SetAttribute("start", strconv.Itoa(start))
	return l.this()
}

// NumberType returns the string used for the type attribute.
func (l *OrderedList) NumberType() string {
	if a := l.Attribute("type"); a == "" {
		return OrderedListNumberTypeNumber
	} else {
		return a
	}
}

func (l *OrderedList) ΩDrawInnerHtml(ctx context.Context, buf *bytes.Buffer) (err error) {
	h := l.getItemsHtml(l.items)
	buf.WriteString(h)
	return nil
}

func (l *OrderedList) getItemsHtml(items []ListItemI) string {
	var h = ""

	for _, item := range items {
		if item.HasChildItems() {
			innerhtml := l.getItemsHtml(item.ListItems())
			a := item.Attributes().Copy()

			// Certain attributes apply to the sub list and not the list item, so we split them here
			a2 := html.NewAttributes()
			if a.Has("type") {
				a2.Set("type", a.Get("type"))
				a.RemoveAttribute("type")
			}

			if a.Has("start") {
				a2.Set("start", a.Get("start"))
				a.RemoveAttribute("start")
			}

			innerhtml = html.RenderTag(l.Tag, a2, innerhtml)
			h += html.RenderTag(l.itemTag, a, item.Label()+" "+innerhtml)
		} else {
			h += html.RenderTag(l.itemTag, item.Attributes(), html2.EscapeString(item.Label()))
		}
	}
	return h
}

type OrderedListCreator struct {
	ID string
	// Items is a static list of labels and values that will be in the list. Or, use a DataProvider to dynamically generate the items.
	Items []ListValue
	// DataProvider is the control that will dynamically provide the data for the list and that implements the DataBinder interface.
	DataProvider data.DataBinder
	// DataProviderID is the id of a control that will dynamically provide the data for the list and that implements the DataBinder interface.
	DataProviderID string
	// NumberType is the type attribute and defaults to OrderedListNumberTypeNumber.
	NumberType string
	// StartAt sets the number to start counting from. The default is 1.
	StartAt int
	page.ControlOptions
}

// Create is called by the framework to create a new control from the Creator. You
// do not normally need to call this.
func (c OrderedListCreator) Create(ctx context.Context, parent page.ControlI) page.ControlI {
	ctrl := NewOrderedList(parent, c.ID)
	c.Init(ctx, ctrl)
	return ctrl
}

func (c OrderedListCreator) Init(ctx context.Context, ctrl OrderedListI) {
	if c.Items != nil {
		ctrl.AddListItems(c.Items)
	}

	if c.DataProvider != nil {
		ctrl.SetDataProvider(c.DataProvider)
	} else if c.DataProviderID != "" {
		provider := ctrl.Page().GetControl(c.DataProviderID).(data.DataBinder)
		ctrl.SetDataProvider(provider)
	}

	if c.NumberType != "" {
		ctrl.SetNumberType(c.NumberType)
	}
	if c.StartAt != 0 {
		ctrl.SetStart(c.StartAt)
	}
	ctrl.ApplyOptions(c.ControlOptions)
}

// GetOrderedList is a convenience method to return the control with the given id from the page.
func GetOrderedList(c page.ControlI, id string) *OrderedList {
	return c.Page().GetControl(id).(*OrderedList)
}
