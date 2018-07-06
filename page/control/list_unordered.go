package control

import (
	"bytes"
	"context"
	"github.com/spekary/goradd/html"
	"github.com/spekary/goradd/page"
	"github.com/spekary/goradd/page/control/data"
)

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
	t := &UnorderedList{}
	t.ItemList = NewItemList(t)
	t.Init(t, parent, id)
	return t
}

func (l *UnorderedList) Init(self page.ControlI, parent page.ControlI, id string) {
	l.Control.Init(self, parent, id)
	l.Tag = "ul"
	l.subItemTag = "li"
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
	h := l.getItemsHtml(l.items)
	buf.WriteString(h)
	return nil
}

func (l *UnorderedList) getItemsHtml(items []ListItemI) string {
	var h = ""

	for _, item := range items {
		if item.HasChildItems() {
			innerhtml := l.getItemsHtml(item.ListItems())
			innerhtml = html.RenderTag(l.Tag, nil, innerhtml)
			h += html.RenderTag(l.subItemTag, item.Attributes(), item.Label()+" "+innerhtml)
		} else {
			h += html.RenderTag(l.subItemTag, item.Attributes(), item.RenderLabel())
		}
	}
	return h
}
