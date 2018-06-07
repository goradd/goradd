package control

import (
	"bytes"
	"context"
	"github.com/spekary/goradd/html"
	"github.com/spekary/goradd/page"
	html2 "html"
	"strconv"
)

// UnorderedList is a dynamically generated html unordered list (ul). Such lists are often used as the basis for
// javascript and css widgets. If you use a data provider to set the data, you should call AddItems to the list
// in your GetData function. After drawing, the items will be removed.
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

func NewOrderedList(parent page.ControlI) *OrderedList {
	t := &OrderedList{}
	t.ItemList = NewItemList(t)
	t.Init(t, parent)
	return t
}

func (l *OrderedList) Init(self page.ControlI, parent page.ControlI) {
	l.UnorderedList.Init(self, parent)
	l.Tag = "ol"
}

// SetNumberType sets the top level number style for the list. Choose from the OrderedListNumberType* constants.
// To set a number type for a sublevel, set the "type" attribute on the list item that is the parent of the sub list.
func (l *OrderedList) SetNumberType(t string) *OrderedList {
	l.SetAttribute("type", l)
	return l
}

// SetStart sets the starting number for the numbers in the top level list. To set the start of a sub-list, set
// the "start" attribute on the list item that is the parent of the sub-list.
func (l *OrderedList) SetStart(start int) *OrderedList {
	l.SetAttribute("start", strconv.Itoa(start))
	return l
}

func (l *OrderedList) NumberType() string {
	if a := l.Attribute("type"); a == "" {
		return OrderedListNumberTypeNumber
	} else {
		return a
	}
}

func (l *OrderedList) DrawInnerHtml(ctx context.Context, buf *bytes.Buffer) (err error) {
	h := l.getItemsHtml(l.items)
	buf.WriteString(h)
	return nil
}

func (l *OrderedList) getItemsHtml(items []ListItemI) string {
	var h = ""

	for _, item := range items {
		if item.HasChildItems() {
			innerhtml := l.getItemsHtml(item.ListItems())
			a := item.Attributes().Clone()

			// Certain attributes apply to the sub list and not the list item, so we split them here
			a2 := html.NewAttributes()
			if a.Has("type") {
				a2.Set("type", a.Get("type"))
				a.Remove("type")
			}

			if a.Has("start") {
				a2.Set("start", a.Get("start"))
				a.Remove("start")
			}

			innerhtml = html.RenderTag(l.Tag, a2, innerhtml)
			h += html.RenderTag(l.subItemTag, a, item.Label()+" "+innerhtml)
		} else {
			h += html.RenderTag(l.subItemTag, item.Attributes(), html2.EscapeString(item.Label()))
		}
	}
	return h
}
