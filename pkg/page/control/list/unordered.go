package list

import (
	"context"
	control2 "github.com/goradd/goradd/pkg/page/control"
	"io"
	"reflect"

	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/html5tag"
)

type UnorderedListI interface {
	page.ControlI
	ListI
	control2.DataManagerI
	GetItemsHtml(items []*Item) string
	SetBulletStyle(s string) UnorderedListI
	SetItemTag(s string) UnorderedListI
	ItemTag() string
}

// UnorderedList is a dynamically generated html unordered list (ul). Such lists are often used as the basis for
// javascript and css widgets. If you use a data provider to set the data, you should call AddItems to the list
// in your LoadData function.
type UnorderedList struct {
	page.ControlBase
	List
	control2.DataManager
	itemTag string
}

const (
	// UnorderedListStyleDisc is the default list style for main items and is a bullet
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

func (l *UnorderedList) Init(self any, parent page.ControlI, id string) {
	l.ControlBase.Init(self, parent, id)
	l.List = NewList(l)
	l.Tag = "ul"
	l.itemTag = "li"
}

// this() supports object oriented features by giving easy access to the virtual function interface.
func (l *UnorderedList) this() UnorderedListI {
	return l.Self().(UnorderedListI)
}

// SetItemTag sets the tag that will be used for items in the list. By default this is "li".
func (l *UnorderedList) SetItemTag(s string) UnorderedListI {
	l.itemTag = s
	return l.this()
}

// ItemTag returns the HTML tag for an item in the list.
func (l *UnorderedList) ItemTag() string {
	return l.itemTag
}

// SetBulletStyle sets the list-style-type attribute of the list. Choose from the UnorderedListStyle* constants.
func (l *UnorderedList) SetBulletStyle(s string) UnorderedListI {
	l.ControlBase.SetStyle("list-style-type", s)
	return l.this()
}

func (l *UnorderedList) DrawTag(ctx context.Context, w io.Writer) {
	if l.HasDataProvider() {
		l.this().LoadData(ctx, l.this())
		defer l.ResetData()
	}
	l.ControlBase.DrawTag(ctx, w)
}

// DrawingAttributes retrieves the tag's attributes at draw time. You should not normally need to call this, and the
// attributes are disposed of after drawing, so they are essentially read-only.
func (l *UnorderedList) DrawingAttributes(ctx context.Context) html5tag.Attributes {
	a := l.ControlBase.DrawingAttributes(ctx)
	a.SetData("grctl", "hlist")
	return a
}

// DrawInnerHtml is called by the framework to draw the content of the tag.
func (l *UnorderedList) DrawInnerHtml(_ context.Context, w io.Writer) {
	h := l.this().GetItemsHtml(l.items)
	page.WriteString(w, h)
	return
}

// GetItemsHtml returns the HTML for the items. It is exported so that
// it can be overridden by other implementations of an UnorderedList.
func (l *UnorderedList) GetItemsHtml(items []*Item) string {
	var h = ""

	for _, item := range items {
		if item.HasChildItems() {
			innerhtml := l.this().GetItemsHtml(item.Items())
			innerhtml = html5tag.RenderTag(l.Tag, nil, innerhtml)
			h += html5tag.RenderTag(l.itemTag, item.Attributes(), item.Label()+" "+innerhtml)
		} else {
			h += html5tag.RenderTag(l.itemTag, item.Attributes(), item.RenderLabel())
		}
	}
	return h
}

// SetData replaces the current list with the given data.
// ValueLabeler, ItemIDer, Labeler or Stringer types are accepted.
// This function can accept one or more lists of items, or
// single items. They will all get added to the top level of the list. To add sub items, get a list item
// and add items to it.
func (l *UnorderedList) SetData(data interface{}) {
	kind := reflect.TypeOf(data).Kind()
	if !(kind == reflect.Slice || kind == reflect.Array) {
		panic("you must call SetData with a slice or array")
	}

	l.List.Clear()
	l.AddItems(data)
}

func (l *UnorderedList) Serialize(e page.Encoder) {
	l.ControlBase.Serialize(e)
	l.List.Serialize(e)
	l.DataManager.Serialize(e)
	if err := e.Encode(l.itemTag); err != nil {
		panic(err)
	}
}

func (l *UnorderedList) Deserialize(dec page.Decoder) {
	l.ControlBase.Deserialize(dec)
	l.List.Deserialize(dec)
	l.DataManager.Deserialize(dec)
	if err := dec.Decode(&l.itemTag); err != nil {
		panic(err)
	}
}

// UnorderedListCreator is a declarative helper to create an UnorderedListControl.
// Pass it to control.AddControls() to add the declared list to the control.
type UnorderedListCreator struct {
	ID string
	// Items is a static list of labels and values that will be in the list. Or, use a DataProvider to dynamically generate the items.
	Items []ListValue
	// DataProvider is the control that will dynamically provide the data for the list and that implements the DataBinder interface.
	DataProvider control2.DataBinder
	// DataProviderID is the id of a control that will dynamically provide the data for the list and that implements the DataBinder interface.
	DataProviderID string
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
		ctrl.AddItems(c.Items)
	}
	if c.DataProvider != nil {
		ctrl.SetDataProvider(c.DataProvider)
	} else if c.DataProviderID != "" {
		provider := ctrl.Page().GetControl(c.DataProviderID).(control2.DataBinder)
		ctrl.SetDataProvider(provider)
	}
	if c.BulletStyle != "" {
		ctrl.SetBulletStyle(c.BulletStyle)
	}
	ctrl.ApplyOptions(ctx, c.ControlOptions)
}

// GetUnorderedList is a convenience method to return the control with the given id from the page.
func GetUnorderedList(c page.ControlI, id string) *UnorderedList {
	return c.Page().GetControl(id).(*UnorderedList)
}

func init() {
	page.RegisterControl(&UnorderedList{})
}
