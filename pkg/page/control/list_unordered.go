package control

import (
	"bytes"
	"context"
	"github.com/goradd/goradd/pkg/html"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/control/data"
	"reflect"
)

type UnorderedListI interface {
	page.ControlI
	ItemListI
	GetItemsHtml(items []ListItemI) string
	data.DataManagerEmbedder
	SetBulletStyle(s string) UnorderedListI
	SetItemTag(s string) UnorderedListI
}

// UnorderedList is a dynamically generated html unordered list (ul). Such lists are often used as the basis for
// javascript and css widgets. If you use a data provider to set the data, you should call AddItems to the list
// in your LoadData function.
type UnorderedList struct {
	page.Control
	ItemList
	itemTag string
	data.DataManager
}

const (
	// UnoderedListStyleDisc is the default list style for main items and is a bullet
	UnorderedListStyleDisc = "disc" // default
	// UnorderedListStyleCircle is the default list style for 2nd level items and is an open circle
	UnorderedListStyleCircle = "circle"
	// UnorderedListStyleSquare sets a square as the bullet
	UnorderedListStyleSquare = "square"
	// UnorderedListStyleNone removes the bullet from the list
	UnorderedListStyleNone = "none"
)

// NewUnorderedList creates a new ul type list.
func NewUnorderedList(parent page.ControlI, id string) *UnorderedList {
	l := &UnorderedList{}
	l.Init(l, parent, id)
	return l
}

func (l *UnorderedList) Init(self page.ControlI, parent page.ControlI, id string) {
	l.Control.Init(self, parent, id)
	l.ItemList = NewItemList(l)
	l.Tag = "ul"
	l.itemTag = "li"
}

// this() supports object oriented features by giving easy access to the virtual function interface.
func (l *UnorderedList) this() UnorderedListI {
	return l.Self.(UnorderedListI)
}

// SetItemTag sets the tag that will be used for items in the list. By default this is "li".
func (l *UnorderedList) SetItemTag(s string) UnorderedListI {
	l.itemTag = s
	return l.this()
}

// SetBulletType sets the list-style-type attribute of the list. Choose from the UnorderedListStyle* constants.
func (l *UnorderedList) SetBulletStyle(s string) UnorderedListI {
	l.Control.SetStyle("list-style-type", s)
	return l.this()
}

func (l *UnorderedList) ΩDrawTag(ctx context.Context) string {
	l.LoadData(ctx, l)
	defer l.Clear()
	return l.Control.ΩDrawTag(ctx)
}

// ΩDrawingAttributes retrieves the tag's attributes at draw time. You should not normally need to call this, and the
// attributes are disposed of after drawing, so they are essentially read-only.
func (l *UnorderedList) ΩDrawingAttributes() *html.Attributes {
	a := l.Control.ΩDrawingAttributes()
	a.SetDataAttribute("grctl", "hlist")
	return a
}

func (l *UnorderedList) ΩDrawInnerHtml(ctx context.Context, buf *bytes.Buffer) (err error) {
	h := l.this().GetItemsHtml(l.items)
	buf.WriteString(h)
	return nil
}

// GetItemsHtml is used by the framework to get the items for the html. It is exported so that
// it can be overridden by other implementations of an UnorderedList.
func (l *UnorderedList) GetItemsHtml(items []ListItemI) string {
	var h = ""

	for _, item := range items {
		if item.HasChildItems() {
			innerhtml := l.this().GetItemsHtml(item.ListItems())
			innerhtml = html.RenderTag(l.Tag, nil, innerhtml)
			h += html.RenderTag(l.itemTag, item.Attributes(), item.Label()+" "+innerhtml)
		} else {
			h += html.RenderTag(l.itemTag, item.Attributes(), item.RenderLabel())
		}
	}
	return h
}

// SetData replaces the current list with the given data.
// The result is kept in memory currently.
// ItemLister, ItemIDer, Labeler or Stringer types. This function can accept one or more lists of items, or
// single items. They will all get added to the top level of the list. To add sub items, get a list item
// and add items to it.
func (l *UnorderedList) SetData(data interface{}) {
	kind := reflect.TypeOf(data).Kind()
	if !(kind == reflect.Slice || kind == reflect.Array) {
		panic("you must call SetData with a slice or array")
	}

	l.ItemList.Clear()
	l.AddListItems(data)
}


type UnorderedListCreator struct {
	ID string
	// Items is a static list of labels and values that will be in the list. Or, use a DataProvider to dynamically generate the items.
	Items []ListValue
	// DataProvider is the control that will dynamically provide the data for the list and that implements the DataBinder interface.
	// This can be either an id of a control, or the control itself.
	DataProvider interface{}
	// BulletStyle is the list-style-type property.
	BulletStyle string
	page.ControlOptions
}

// Create is called by the framework to create a new control from the Creator. You
// do not normally need to call this.
func (c UnorderedListCreator) Create(ctx context.Context, parent page.ControlI) page.ControlI {
	ctrl := NewUnorderedList(parent, c.ID)
	c.Init(ctx, ctrl)
	return ctrl
}

func (c UnorderedListCreator) Init(ctx context.Context, ctrl UnorderedListI) {
	if c.Items != nil {
		ctrl.AddListItems(c.Items)
	}
	if c.DataProvider != nil {
		// If this fails, then perhaps you are giving a data provider id for a control that is not yet created. Create the control first.
		var provider data.DataBinder
		if s,ok := c.DataProvider.(string); ok {
			provider = ctrl.Page().GetControl(s).(data.DataBinder)
		} else {
			provider = c.DataProvider.(data.DataBinder)
		}
		ctrl.SetDataProvider(provider)
	}
	if c.BulletStyle != "" {
		ctrl.SetBulletStyle(c.BulletStyle)
	}
	ctrl.ApplyOptions(c.ControlOptions)
}

// GetUnorderedList is a convenience method to return the control with the given id from the page.
func GetUnorderedList(c page.ControlI, id string) *UnorderedList {
	return c.Page().GetControl(id).(*UnorderedList)
}
