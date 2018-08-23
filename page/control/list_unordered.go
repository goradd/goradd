package control

import (
	"bytes"
	"context"
	"github.com/spekary/goradd/html"
	"github.com/spekary/goradd/page"
	"github.com/spekary/goradd/page/control/data"
	page2 "goradd-project/override/page"
)


type UnorderedListI interface {
	page2.ControlI
	GetItemsHtml(items []ListItemI) string

}


// UnorderedList is a dynamically generated html unordered list (ul). Such lists are often used as the basis for
// javascript and css widgets. If you use a data provider to set the data, you should call AddItems to the list
// in your GetData function. After drawing, the items will be removed.
type UnorderedList struct {
	page.Control
	ItemList
	subItemTag string
	data.DataManager
}

const (
	UnorderedListStyleDisc   = "disc" // default
	UnorderedListStyleCircle = "circle"
	UnorderedListStyleSquare = "square"
	UnorderedListStyleNone   = "none"
)

func NewUnorderedList(parent page.ControlI, id string) *UnorderedList {
	l := &UnorderedList{}
	l.Init(l, parent, id)
	return l
}

func (l *UnorderedList) Init(self page.ControlI, parent page.ControlI, id string) {
	l.Control.Init(self, parent, id)
	l.ItemList = NewItemList(l)
	l.Tag = "ul"
	l.subItemTag = "li"
}

// this() supports object oriented features by giving easy access to the virtual function interface
// Subclasses should provide a duplicate. Calls that implement chaining should return the result of this function.
func (c *UnorderedList) this() UnorderedListI {
	return c.Self.(UnorderedListI)
}

func (l *UnorderedList) SetSubTag(s string) {
	l.subItemTag = s
}

// SetBulletType sets the bullet type. Choose from the UnorderedListStyle* constants.
func (l *UnorderedList) SetBulletStyle(s string) {
	l.Control.SetStyle("list-style-type", s)
}

func (l *UnorderedList) DrawTag(ctx context.Context) string {
	l.GetData(ctx, l)
	defer l.Clear()
	return l.Control.DrawTag(ctx)
}

// DrawingAttributes retrieves the tag's attributes at draw time. You should not normally need to call this, and the
// attributes are disposed of after drawing, so they are essentially read-only.
func (l *UnorderedList) DrawingAttributes() *html.Attributes {
	a := l.Control.DrawingAttributes()
	a.SetDataAttribute("grctl", "hlist")
	return a
}

func (l *UnorderedList) DrawInnerHtml(ctx context.Context, buf *bytes.Buffer) (err error) {
	h := l.this().GetItemsHtml(l.items)
	buf.WriteString(h)
	return nil
}

func (l *UnorderedList) GetItemsHtml(items []ListItemI) string {
	var h = ""

	for _, item := range items {
		if item.HasChildItems() {
			innerhtml := l.this().GetItemsHtml(item.ListItems())
			innerhtml = html.RenderTag(l.Tag, nil, innerhtml)
			h += html.RenderTag(l.subItemTag, item.Attributes(), item.Label()+" "+innerhtml)
		} else {
			h += html.RenderTag(l.subItemTag, item.Attributes(), item.RenderLabel())
		}
	}
	return h
}
