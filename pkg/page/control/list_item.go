package control

import (
	"fmt"
	"github.com/goradd/goradd/pkg/html"
	"github.com/goradd/goradd/pkg/page"
	html2 "html"
	"strconv"
)


type ItemLister interface {
	Value() interface{}
	Label() string
}

// ItemIDer is an interface to a listable object that matches orm objects
type ItemIDer interface {
	ID() string
	String() string
}

type Labeler interface {
	Label() string
}

// A ListItem is an object that is a member of a list. HTML has a few different kinds of lists, and this can be a member
// of a select list (<select>), or an ordered or unordered list (<ul> or <ol>). It is up to the manager of the list to
// render the item, but this serves as a place to store options about the item. Not all options are pertinent to
// all lists.
//
// A list item generally has a value, and a label. Often, lists will have ids too, that will appear in the html output,
// but the id values are managed by the list manager and generally should not be set by you. In situations where the
// user selects a list item, you would use the id to retrieve the ListItem selected.
type ListItem struct {
	ItemList
	value string // html list item values eventually get expressed as a string in the html
	id    string
	label             string
	attributes        html.Attributes
	shouldEscapeLabel bool
	disabled          bool
	isDivider         bool
	anchorAttributes  html.Attributes
}

// NewListItem creates a new item for a list. Specify an empty string for an item that represents no selection.
func NewListItem(label string, value ...string) *ListItem {
	l := &ListItem{label: label}
	if c := len(value); c == 1 {
		l.value = value[0]
	} else if c > 1 {
		panic("Call NewListItem with zero or one value only.")
	} else {
		l.value = label
	}

	l.ItemList = NewItemList(l)
	return l
}

// NewItemFromItemLister creates a new item from any object that has a Value and Label method.
func NewItemFromItemLister(i ItemLister) *ListItem {
	var l *ListItem

	if i.Value() == nil {
		l = &ListItem{value: "", label: i.Label()}
	} else {
		l = &ListItem{value: fmt.Sprintf("%v", i.Value()), label: i.Label()}
	}
	l.ItemList = NewItemList(l)
	return l
}

// NewItemFromLabeler creates a new item from any object that has just a Label method.
func NewItemFromLabeler(i Labeler) *ListItem {
	l := &ListItem{label: i.Label(), value: i.Label()}
	l.ItemList = NewItemList(l)
	return l
}

// NewItemFromStringer creates a new item from any object that has just a String method.
// The label and value will be the same.
func NewItemFromStringer(i fmt.Stringer) *ListItem {
	l := &ListItem{label: i.String(), value: i.String()}
	l.ItemList = NewItemList(l)
	return l
}

// NewItemFromItemIDer creates a new item from any object that has an ID and String method.
// Note that the ID() of the ItemIDer will become the value of the select item, and the String()
// will become the label
func NewItemFromItemIDer(i ItemIDer) *ListItem {
	l := &ListItem{value: i.ID(), label: i.String()}
	l.ItemList = NewItemList(l)
	return l
}

func (i *ListItem) SetValue(v string) *ListItem {
	i.value = v
	return i
}

func (i *ListItem) Value() string {
	return i.value
}

func (i *ListItem) IntValue() int {
	v,_ := strconv.Atoi(i.value)
	return v
}


func (i *ListItem) ID() string {
	return i.id
}

// SetID should not be called by your code typically. It is exported for implementations of item lists. The IDs of an
// item list are completely managed by the list, you cannot have custom ids.
func (i *ListItem) SetID(id string) {
	i.id = id
	i.Attributes().SetID(id)
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

func (i *ListItem) SetDisabled(d bool) {
	i.disabled = d
}

func (i *ListItem) Disabled() bool {
	return i.disabled
}

func (i *ListItem) SetIsDivider(d bool) {
	i.isDivider = d
}

func (i *ListItem) IsDivider() bool {
	return i.isDivider
}

func (i *ListItem) SetAnchor(a string) {
	i.AnchorAttributes().Set("href", a)
}

func (i *ListItem) Anchor() string {
	if i.anchorAttributes == nil || !i.anchorAttributes.Has("href") {
		return ""
	}
	return i.anchorAttributes.Get("href")
}

func (i *ListItem) AnchorAttributes() html.Attributes {
	if i.anchorAttributes == nil {
		i.anchorAttributes = html.NewAttributes()
	}
	return i.anchorAttributes
}

func (i *ListItem) SetShouldEscapeLabel(e bool) *ListItem {
	i.shouldEscapeLabel = e
	return i
}

func (i *ListItem) RenderLabel() (h string) {
	if i.shouldEscapeLabel {
		h = html2.EscapeString(i.label)
	} else {
		h = i.label
	}
	if i.Anchor() != "" && !i.disabled {
		h = html.RenderTag("a", i.anchorAttributes, h)
	}
	return
}

// Attributes returns a pointer to the attributes of the item for customization. You can directly set the attributes
// on the returned object.
func (i *ListItem) Attributes() html.Attributes {
	if i.attributes == nil {
		i.attributes = html.NewAttributes()
	}
	return i.attributes
}

// IsEmptyValue returns true if the value is empty, meaning it does not satisfy a selection being made
// if the list has IsRequired turned on.
func (i *ListItem) IsEmptyValue() bool {
	return i.value == ""
}

func (i *ListItem) Serialize(e page.Encoder) (err error) {
	if err = e.Encode(&i.value); err != nil {
		return
	}
	if err = e.Encode(i.id); err != nil {
		return
	}
	if err = e.Encode(i.label); err != nil {
		return
	}
	if err = e.Encode(i.attributes); err != nil {
		return
	}
	if err = e.Encode(i.shouldEscapeLabel); err != nil {
		return
	}
	if err = e.Encode(i.disabled); err != nil {
		return
	}
	if err = e.Encode(i.isDivider); err != nil {
		return
	}
	if err = e.Encode(i.anchorAttributes); err != nil {
		return
	}
	return i.ItemList.Serialize(e)
}

func (i *ListItem) Deserialize(dec page.Decoder) (err error) {
	if err = dec.Decode(&i.value); err != nil {
		panic(err)
	}
	if err = dec.Decode(&i.id); err != nil {
		panic(err)
	}
	if err = dec.Decode(&i.label); err != nil {
		panic(err)
	}
	if err = dec.Decode(&i.attributes); err != nil {
		panic(err)
	}
	if err = dec.Decode(&i.shouldEscapeLabel); err != nil {
		panic(err)
	}
	if err = dec.Decode(&i.disabled); err != nil {
		panic(err)
	}
	if err = dec.Decode(&i.isDivider); err != nil {
		panic(err)
	}
	if err = dec.Decode(&i.anchorAttributes); err != nil {
		panic(err)
	}

	return i.ItemList.Deserialize(dec)
}

// ListValue is a helper for initializing a control based on ItemList.
// It satisfies the ItemLister interface. To use it, create a slice of ListValue's and
// pass the list to AddListItems or SetData.
type ListValue struct {
	L string
	// V is the value
	V interface{}
}

func (l ListValue) Value() interface{} {
	return l.V
}

func (l ListValue) Label() string {
	return l.L
}
