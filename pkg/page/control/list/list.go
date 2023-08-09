// Package list contains list-type controls. This includes select lists and hierarchical lists.
package list

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/goradd/goradd/pkg/config"
	"github.com/goradd/goradd/pkg/page"
)

// IDer is an object that has an ID.
type IDer interface {
	ID() string
}

// IDSetter is an interface for an item that sets an ID.
type IDSetter interface {
	SetID(id string)
}

// ListI is the interface for all controls that display a list of Items.
type ListI interface {
	Add(label string, value ...string) *Item
	AddAt(index int, label string, value ...string)
	AddItemAt(index int, item *Item)
	AddItems(items ...interface{})
	ItemAt(index int) *Item
	Items() []*Item
	Clear()
	RemoveAt(index int)
	Len() int
	GetItemByID(id string) (foundItem *Item)
	GetItemByValue(value string) (id string, foundItem *Item)
	reindex(start int)
	findItemByValue(value string) (container *List, index int)
}

// List manages a list of *Item items. List is designed to be embedded in another structure, and will
// turn that object into a manager of list items. List will manage the id's of the items in its list, you do not
// have control of that, and it needs to manage those ids in order to efficiently manage the selection process.
//
// Controls that embed this must implement the Serialize and Deserialize methods to call both the ControlBase's
// Serialize and Deserialize methods, and the ones here.
type List struct {
	ownerID string
	items   []*Item
}

// NewList creates a new item list. "owner" is the object that has the list embedded in it, and must be
// an IDer.
func NewList(owner IDer) List {
	return List{ownerID: owner.ID()}
}

// Add adds the given item to the end of the list. The value is optional, but should only be one or zero values.
func (l *List) Add(label string, value ...string) *Item {
	i := NewItem(label, value...)
	l.AddItemAt(len(l.items), i)
	return i
}

// AddAt adds the item at the given index.
// If the index is negative, it counts from the end. -1 would therefore put the item before the last item.
// If the index is bigger or equal to the number of items, it adds it to the end. If the index is zero, or is negative and smaller than
// the negative value of the number of items, it adds to the beginning. This can be an expensive operation in a long
// hierarchical list, so use sparingly.
func (l *List) AddAt(index int, label string, value ...string) {
	l.AddItemAt(index, NewItem(label, value...))
}

// AddItemAt adds the item at the given index. If the index is negative, it counts from the end. If the index is
// -1 or bigger than the number of items, it adds it to the end. If the index is zero, or is negative and smaller than
// the negative value of the number of items, it adds to the beginning. This can be an expensive operation in a long
// hierarchical list, so use sparingly.
func (l *List) AddItemAt(index int, item *Item) {
	if index < 0 {
		index = len(l.items) + index
		if index < 0 {
			index = 0
		}
	} else if index > len(l.items) {
		index = len(l.items)
	}
	l.items = append(l.items, nil)
	copy(l.items[index+1:], l.items[index:])
	l.items[index] = item
	l.reindex(index)
	return
}

// AddItems adds one or more objects to the end of the list. items should be a list of *Item,
// ValueLabeler, ItemIDer, Labeler or Stringer types. This function can accept one or more lists of items, or
// single items
func (l *List) AddItems(items ...interface{}) {
	var start int

	if items == nil {
		return
	}
	for _, item := range items {
		kind := reflect.TypeOf(item).Kind()
		if kind == reflect.Array || kind == reflect.Slice {
			listValue := reflect.ValueOf(item)
			for i := 0; i < listValue.Len(); i++ {
				itemI := listValue.Index(i).Interface()
				l.addItem(itemI)
			}
		} else {
			l.addItem(item)
		}
	}
	l.reindex(start)
	return
}

// Private function to add an interface item to the end of the list. Will need to be reindexed eventually.
func (l *List) addItem(item interface{}) {
	switch v := item.(type) {
	case *Item:
		l.items = append(l.items, v)
	case ValueLabeler:
		item := NewItemFromValueLabeler(v)
		l.items = append(l.items, item)
	case ItemIDer:
		item := NewItemFromItemIDer(v)
		l.items = append(l.items, item)
	case ItemIntIDer:
		item := NewItemFromItemIntIDer(v)
		l.items = append(l.items, item)
	case Labeler:
		item := NewItemFromLabeler(v)
		l.items = append(l.items, item)
	case fmt.Stringer:
		item := NewItemFromStringer(v)
		l.items = append(l.items, item)
	default:
		panic("Unknown object type")
	}
}

// reindex is internal and should get called whenever an item gets added to the list out of order or an id changes.
func (l *List) reindex(start int) {
	if l.ownerID == "" || l.items == nil || len(l.items) == 0 || start >= len(l.items) {
		return
	}
	for i := start; i < len(l.items); i++ {
		id := l.ownerID + "_" + strconv.Itoa(i)
		l.items[i].SetID(id)
	}
}

// ItemAt retrieves an item by index.
func (l *List) ItemAt(index int) *Item {
	if index >= len(l.items) {
		return nil
	}
	return l.items[index]
}

// Items returns a slice of the *Item items, in the order they were added or arranged.
func (l *List) Items() []*Item {
	return l.items
}

// Clear removes all the items from the list.
func (l *List) Clear() {
	l.items = nil
}

// RemoveAt removes an item at the given index.
func (l *List) RemoveAt(index int) {
	if index < 0 || index >= len(l.items) {
		panic("Index out of range.")
	}
	l.items = append(l.items[:index], l.items[index+1:]...)
}

// Len returns the length of the item list at the current level. In other words, it does not count items in sublists.
func (l *List) Len() int {
	if l.items == nil {
		return 0
	}
	return len(l.items)
}

// GetItemByID recursively searches for and returns the item corresponding to the given id. Since we are managing the
// id numbers, we can efficiently find the item. Note that if you add items to the list, the ids may change.
func (l *List) GetItemByID(id string) (foundItem *Item) {
	if l.items == nil {
		return nil
	}

	parts := strings.SplitN(id, "_", 3) // first item is our own id, 2nd is id from the list, 3rd is a level beyond the list

	var countParts int
	if countParts = len(parts); countParts <= 1 {
		return nil
	}

	l1Id, err := strconv.Atoi(parts[1])
	if err != nil || l1Id < 0 || l1Id >= len(l.items) {
		panic("Bad id")
	}

	item := l.items[l1Id]

	if countParts == 2 {
		return item
	}

	return item.GetItemByID(parts[1] + "_" + parts[2])
}

// GetItemByValue recursively searches the list to find the item with the given value.
// It starts with the current list, and if not found, will search in sublists.
func (l *List) GetItemByValue(value string) (id string, foundItem *Item) {
	container, index := l.findItemByValue(value)

	if container != nil {
		foundItem = container.items[index]
		id = foundItem.ID()
		return
	}
	return "", nil
}

// findItemByValue searches for the item by value, and returns the index of the found item,
// and the List that the item was found in. The returned List could be the current
// item list, or a sublist.
func (l *List) findItemByValue(value string) (container *List, index int) {
	if len(l.items) == 0 {
		return nil, -1 // no sub items, so its not here
	}
	var item *Item

	for index, item = range l.items {
		v := item.Value()
		if v == value {
			container = l
			return
		}
	}

	for index, item = range l.items {
		container, index = item.findItemByValue(value)
		if container != nil {
			return
		}
	}

	return nil, -1 // not found
}

func (l *List) Serialize(e page.Encoder) {
	if err := e.Encode(l.ownerID); err != nil {
		panic(err)
	}
	var count int = len(l.items)
	if err := e.Encode(count); err != nil {
		panic(err)
	}

	// Opt for our own serialization method, rather than using gob
	for _, i := range l.items {
		i.Serialize(e)
	}
}

func (l *List) Deserialize(dec page.Decoder) {
	if err := dec.Decode(&l.ownerID); err != nil {
		panic(err)
	}

	var count int
	if err := dec.Decode(&count); err != nil {
		panic(err)
	}

	for i := 0; i < count; i++ {
		item := Item{}
		item.Deserialize(dec)
		l.items = append(l.items, &item)
	}

	return
}

// SortIds sorts a list of auto-generated ids in numerical and hierarchical order.
// This is normally just called by the framework.
func SortIds(ids []string) {
	if len(ids) > 1 {
		sort.Sort(IdSlice(ids))
	}
}

// IdSlice is a slice of string ids, and is used to sort a list of ids
// that the item list uses.
type IdSlice []string

func (p IdSlice) Len() int { return len(p) }
func (p IdSlice) Less(i, j int) bool {
	// First ones are always the main control id, and should be equal
	vals1 := strings.SplitN(p[i], "_", 2)
	vals2 := strings.SplitN(p[j], "_", 2)
	if vals1[0] != vals2[0] {
		panic("The first part of an id should be equal when sorting.")
	}

	for {
		vals1 = strings.SplitN(vals1[1], "_", 2)
		vals2 = strings.SplitN(vals2[1], "_", 2)
		i1, _ := strconv.Atoi(vals1[0])
		i2, _ := strconv.Atoi(vals2[0])
		if i1 < i2 {
			return true
		} else if i1 > i2 {
			return false
		} else if len(vals1) < len(vals2) {
			return true
		} else if len(vals1) > len(vals2) || len(vals2) <= 1 {
			return false
		}
	}
}

func (p IdSlice) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

// NoSelectionItemList returns a default item list to start a selection list that allows no selection
func NoSelectionItemList() []interface{} {
	return []interface{}{NewItem(config.NoSelectionString, "")}
}

// SelectOneItemList returns a default item list to start a selection list that asks the user to select an item
func SelectOneItemList() []interface{} {
	return []interface{}{NewItem(config.SelectOneString, "")}
}

// IDerStringListCompare is a utility function that will compare a list of IDers with a list
// of strings to see if their values are equal.
func IDerStringListCompare[T IDer](ids []T, values []string) bool {
	if len(ids) != len(values) {
		return false
	}
	for i := range ids {
		if ids[1].ID() != values[i] {
			return false
		}
	}
	return true
}
