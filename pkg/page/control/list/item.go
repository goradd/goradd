package list

import (
	"fmt"
	"html"
	"strconv"

	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/html5tag"
)

type ValueLabeler interface {
	Value() interface{}
	Label() string
}

// ItemIDer is an interface to a listable object that matches most orm objects
type ItemIDer interface {
	ID() string
	String() string
}

// ItemIntIDer matches orm objects that use an int type for the id
type ItemIntIDer interface {
	ID() int
	String() string
}

type Labeler interface {
	Label() string
}

// An Item is an object that is a member of a list. HTML has a few different kinds of lists, and this can be a member
// of a select list (<select>), or an ordered or unordered list (<ul> or <ol>). It is up to the manager of the list to
// render the item, but this serves as a place to store options about the item. Not all options are pertinent to
// all lists.
//
// A list item generally has a value, and a label. Often, lists will have ids too, that will appear in the html output,
// but the id values are managed by the list manager and generally should not be set by you. In situations where the
// user selects a list item, you would use the id to retrieve the Item selected.
type Item struct {
	List
	value             string
	id                string
	label             string
	attributes        html5tag.Attributes
	shouldEscapeLabel bool
	disabled          bool
	isDivider         bool
	anchorAttributes  html5tag.Attributes
}

// NewListItem creates a new item for a list. Specify an empty string for an item that represents no selection.
func NewListItem(label string, value ...string) *Item {
	l := &Item{label: label}
	if c := len(value); c == 1 {
		l.value = value[0]
	} else if c > 1 {
		panic("Call NewListItem with zero or one value only.")
	} else {
		l.value = label
	}

	l.List = NewList(l)
	return l
}

// NewItemFromValueLabeler creates a new item from any object that has a Value and Label method.
func NewItemFromValueLabeler(i ValueLabeler) *Item {
	var l *Item

	if i.Value() == nil {
		l = &Item{value: "", label: i.Label()}
	} else {
		l = &Item{value: fmt.Sprint(i.Value()), label: i.Label()}
	}
	l.List = NewList(l)
	return l
}

// NewItemFromLabeler creates a new item from any object that has just a Label method.
func NewItemFromLabeler(i Labeler) *Item {
	l := &Item{label: i.Label(), value: i.Label()}
	l.List = NewList(l)
	return l
}

// NewItemFromStringer creates a new item from any object that has just a String method.
// The label and value will be the same.
func NewItemFromStringer(i fmt.Stringer) *Item {
	l := &Item{label: i.String(), value: i.String()}
	l.List = NewList(l)
	return l
}

// NewItemFromItemIDer creates a new item from any object that has an ID and String method.
// Note that the ID() of the ItemIDer will become the value of the select item, and the String()
// will become the label
func NewItemFromItemIDer(i ItemIDer) *Item {
	l := &Item{value: i.ID(), label: i.String()}
	l.List = NewList(l)
	return l
}

func NewItemFromItemIntIDer(i ItemIntIDer) *Item {
	l := &Item{value: strconv.Itoa(i.ID()), label: i.String()}
	l.List = NewList(l)
	return l
}

func (i *Item) SetValue(v string) *Item {
	i.value = v
	return i
}

func (i *Item) Value() string {
	return i.value
}

func (i *Item) IntValue() int {
	v, _ := strconv.Atoi(i.value)
	return v
}

func (i *Item) ID() string {
	return i.id
}

// SetID should not be called by your code typically. It is exported for implementations of item lists. The IDs of an
// item list are completely managed by the list, you cannot have custom ids.
func (i *Item) SetID(id string) {
	i.id = id
	i.Attributes().SetID(id)
	i.List.reindex(0)
}

func (i *Item) HasChildItems() bool {
	return i.List.Len() > 0
}

func (i *Item) Label() string {
	return i.label
}

func (i *Item) SetLabel(l string) {
	i.label = l
}

func (i *Item) SetDisabled(d bool) {
	i.disabled = d
}

func (i *Item) Disabled() bool {
	return i.disabled
}

func (i *Item) SetIsDivider(d bool) {
	i.isDivider = d
}

func (i *Item) IsDivider() bool {
	return i.isDivider
}

// SetAnchor assigns the value of the href attribute. This will cause the item to output with an
// anchor tag (a).
func (i *Item) SetAnchor(a string) *Item {
	i.AnchorAttributes().Set("href", a)
	return i
}

func (i *Item) HasAnchor() bool {
	return i.AnchorAttributes().Has("href")
}

func (i *Item) Anchor() string {
	if i.anchorAttributes == nil || !i.anchorAttributes.Has("href") {
		return ""
	}
	return i.anchorAttributes.Get("href")
}

// AnchorAttributes returns the attributes that will be used for each of the anchor tags.
//
// You can directly set attributes on the returned value.
func (i *Item) AnchorAttributes() html5tag.Attributes {
	if i.anchorAttributes == nil {
		i.anchorAttributes = html5tag.NewAttributes()
	}
	return i.anchorAttributes
}

func (i *Item) SetShouldEscapeLabel(e bool) *Item {
	i.shouldEscapeLabel = e
	return i
}

// RenderLabel is called by list implementations to render the item.
func (i *Item) RenderLabel() (h string) {
	if i.shouldEscapeLabel {
		h = html.EscapeString(i.label)
	} else {
		h = i.label
	}
	if i.Anchor() != "" && !i.disabled {
		h = html5tag.RenderTag("a", i.anchorAttributes, h)
	}
	return
}

// Attributes returns a pointer to the attributes of the item for customization. You can directly set the attributes
// on the returned object.
func (i *Item) Attributes() html5tag.Attributes {
	if i.attributes == nil {
		i.attributes = html5tag.NewAttributes()
	}
	return i.attributes
}

func (i *Item) SetAttribute(key, value string) *Item {
	i.Attributes().Set(key, value)
	return i
}

func (i *Item) AddClass(class string) *Item {
	i.Attributes().AddClass(class)
	return i
}

// IsEmptyValue returns true if the value is empty, meaning it does not satisfy a selection being made
// if the list has IsRequired turned on.
func (i *Item) IsEmptyValue() bool {
	return i.value == ""
}

func (i *Item) Serialize(e page.Encoder) {
	if err := e.Encode(&i.value); err != nil {
		panic(err)
	}
	if err := e.Encode(i.id); err != nil {
		panic(err)
	}
	if err := e.Encode(i.label); err != nil {
		panic(err)
	}
	if err := e.Encode(i.attributes); err != nil {
		panic(err)
	}
	if err := e.Encode(i.shouldEscapeLabel); err != nil {
		panic(err)
	}
	if err := e.Encode(i.disabled); err != nil {
		panic(err)
	}
	if err := e.Encode(i.isDivider); err != nil {
		panic(err)
	}
	if err := e.Encode(i.anchorAttributes); err != nil {
		panic(err)
	}
	i.List.Serialize(e)
}

func (i *Item) Deserialize(dec page.Decoder) {
	if err := dec.Decode(&i.value); err != nil {
		panic(err)
	}
	if err := dec.Decode(&i.id); err != nil {
		panic(err)
	}
	if err := dec.Decode(&i.label); err != nil {
		panic(err)
	}
	if err := dec.Decode(&i.attributes); err != nil {
		panic(err)
	}
	if err := dec.Decode(&i.shouldEscapeLabel); err != nil {
		panic(err)
	}
	if err := dec.Decode(&i.disabled); err != nil {
		panic(err)
	}
	if err := dec.Decode(&i.isDivider); err != nil {
		panic(err)
	}
	if err := dec.Decode(&i.anchorAttributes); err != nil {
		panic(err)
	}

	i.List.Deserialize(dec)
}

// ListValue is a helper for initializing a control based on List.
// It satisfies the ValueLabeler interface. To use it, create a slice of ListValue's and
// pass the list to AddItems or SetData.
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
