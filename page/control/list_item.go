package control

import (
	"github.com/spekary/goradd/html"
	"fmt"
)

type ListItemI interface {
	ItemListI
	Value() interface{}
	ID() string
	SetID(string)
	Label() string
	SetLabel(string)
	IntValue() int
	StringValue() string
	HasChildItems() bool
	Attributes() *html.Attributes
}

type ItemLister interface {
	Value() interface{}
	Label() string
}

type Labeler interface {
	Label() string
}


type ListItem struct {
	value interface{}
	id string
	ItemList
	label string
	attributes *html.Attributes
}

// NewListItem creates a new item for a list. Specify an empty value for an item that represents no selection.
func NewListItem(label string, value ...interface{}) *ListItem {
	l := &ListItem{attributes:html.NewAttributes(), label: label}
	if c := len(value); c == 1 {
		l.value = value[0]
	} else if c > 1 {
		panic ("Call NewListItem with zero or one value only.")
	}

	l.ItemList = NewItemList(l)
	return l
}

// NewItemFromItemLister creates a new item from any object that has a Value and Label method.
func NewItemFromItemLister(i ItemLister) *ListItem {
	l := &ListItem{attributes:html.NewAttributes(), value: i.Value(), label: i.Label()}
	l.ItemList = NewItemList(l)
	return l
}

func NewItemFromLabeler(i Labeler) *ListItem {
	l := &ListItem{attributes:html.NewAttributes(), label: i.Label()}
	l.ItemList = NewItemList(l)
	return l
}


func (i *ListItem) SetValue(v interface{}) {
	i.value = v
}

func (i *ListItem) Value() interface{} {
	return i.value
}

func (i *ListItem) IntValue() int {
	return i.value.(int)
}

func (i *ListItem) StringValue() string {
	if s,ok := i.value.(fmt.Stringer); ok {
		return s.String()
	} else {
		return i.value.(string)
	}
}

func (i *ListItem) ID() string {
	return i.id
}

func (i *ListItem) SetID(id string) {
	i.id = id
	i.attributes.SetID(id)
	i.ItemList.reindex(0)
}

func (i *ListItem) HasChildItems() bool {
	return i.ItemList.Len() > 0
}

func (i *ListItem) Label() string {
	return i.label
}

func (i *ListItem) SetLabel(l string) {
	i.label = l
}

// Attributes returns a pointer to the attributes of the item for customization. You can directly set the attributes
// on the returned object.
func (i *ListItem) Attributes() *html.Attributes {
	return i.attributes
}