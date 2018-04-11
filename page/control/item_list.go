package control

import (
	"strings"
	"strconv"
)

// Ider is an object that can embed an ItemList
type Ider interface {
	Id() string
	SetId(id string)
}

type ItemListI interface {
	AddItem(value interface{}, label string) ListItemI
	AddItemAt(index int, value interface{}, label string)
	AddListItem(item ListItemI)
	AddListItemAt(index int, item ListItemI)
	AddListItems(items... ListItemI)
	AddItemListers(items... ItemLister)
	GetItemAt(index int) ListItemI
	ListItems() []ListItemI
	Clear()
	RemoveAt(index int)
	Len() int
	FindById(id string) (foundItem ListItemI)
	FindByValue(value interface{}) (id string, foundItem ListItemI)
	reindex(start int)
}

// ItemList manages a list of ListItemI list items. ItemList is designed to be embedded in another structure, and will
// turn that object into a manager of list items. ItemList will manage the id's of the items in its list, you do not
// have control of that, and it needs to manage those ids in order to efficiently manage the selection process.
type ItemList struct {
	owner Ider
	items []ListItemI
}

// NewItemList creates a new item list. "owner" is the object that has the ItemList embedded in it, and must be
// an Ider.
func NewItemList(owner Ider) ItemList {
	return ItemList{owner: owner, items: []ListItemI{}}
}

// AddItem adds the given item to the end of the list.
func (l *ItemList) AddItem(value interface{}, label string) ListItemI {
	i := NewListItem(value, label)
	l.AddListItem(i)
	return i
}

// AddItemAt adds the item at the given index. If the index is negative, it counts from the end. If the index is
// -1 or bigger than the number of items, it adds it to the end. If the index is zero, or is negative and smaller than
// the negative value of the number of items, it adds to the beginning. This can be an expensive operation in a long
// hierarchical list, so use sparingly.
func (l *ItemList) AddItemAt(index int, value interface{}, label string) {
	l.AddListItemAt(index, NewListItem(value, label))
}


// AddItem adds the given item to the end of the list.
func (l *ItemList) AddListItem(item ListItemI) {
	id := l.owner.Id() + "_" + strconv.Itoa(len(l.items))
	item.SetId(id)
	l.items = append(l.items, item)
	item.reindex(0)
}

// AddItemAt adds the item at the given index. If the index is negative, it counts from the end. If the index is
// -1 or bigger than the number of items, it adds it to the end. If the index is zero, or is negative and smaller than
// the negative value of the number of items, it adds to the beginning. This can be an expensive operation in a long
// hierarchical list, so use sparingly.
func (l *ItemList) AddListItemAt(index int, item ListItemI) {
	l.items = append(l.items, nil)
	copy(l.items[index+1:], l.items[index:])
	l.items[index] = item
	l.reindex(index)
}

// AddItems adds a slice of ListItemI objects to the end of the list.
func (l *ItemList) AddListItems(items... ListItemI) {
	for _,item := range items {
		id := l.owner.Id() + "_" + strconv.Itoa(len(l.items))
		item.SetId(id)
	}
	l.items = append(l.items, items...)
}

// AddItems adds to the list a slice of objects that contain a Value and Label method.
func (l *ItemList) AddItemListers(items... ItemLister) {
	iList := []ListItemI{}
	for _,i := range items {
		item := NewItemFromItemLister(i)
		iList = append(iList, item)
	}
	l.AddListItems(iList...)
}


// reindex is internal and should get called whenever an item gets added to the list out of order or an id changes.
func (l *ItemList) reindex(start int) {
	if l.owner.Id() == "" || len(l.items) == 0 || start >= len (l.items) {
		return
	}
	for i,item := range l.items[start:] {
		id := l.owner.Id() + "_" + strconv.Itoa(i)
		item.SetId(id)
	}
}

// GetItemAt retrieves an item by index.
func (l *ItemList) GetItemAt(index int) ListItemI {
	if index >= len(l.items) {
		return nil
	}
	return l.items[index]
}

// Items returns a slice of the ListItemI items in the ItemList, in the order they were added or arranged.
func (l *ItemList) ListItems() []ListItemI {
	return l.items
}

// Clear removes all the items from the list.
func (l *ItemList) Clear() {
	l.items = []ListItemI{}
}

// RemoveAt removes an item at the given index.
func (l *ItemList) RemoveAt(index int) {
	if index < 0 || index >= len(l.items) {
		panic ("Index out of range.")
	}
	l.items = append(l.items[:index], l.items[index+1:]...)
}

// Len returns the length of the item list at the current level. In other words, it does not count items in sublists.
func (l *ItemList) Len() int {
	return len(l.items)
}

// FindById recursively searches for and returns the item corresponding to the given id. Since we are managing the
// id numbers, we can efficiently find the item. Note that if you add items to the list, the ids may change.
func (l *ItemList) FindById(id string) (foundItem ListItemI) {
	parts := strings.SplitN(id, "_", 3) // first item is our own id, 2nd is id from the list, 3rd is a level beyond the list

	l1Id,err := strconv.Atoi(parts[1])
	if err != nil || l1Id < 0 {
		panic("Bad id")
	}

	var countParts int
	if countParts = len(parts); countParts <= 1 || l1Id >= len(l.items) {
		return nil
	}

	item := l.items[l1Id]

	if countParts == 2 {
		return item
	}

	return item.FindById(parts[1] + "_" + parts[2])
}

// FindByValue recursively searches the list to find the item with the given value.
// It starts with the current list, and if not found, will search in sublists.
func (l *ItemList) FindByValue(value interface{}) (id string, foundItem ListItemI) {
	if len(l.items) == 0 {
		return "", nil
	}

	for _,foundItem = range l.items {
		if foundItem.Value() == value {
			id = foundItem.Id()
			return
		}
	}

	for _,item := range l.items {
		id, foundItem = item.FindByValue(value)
		if foundItem != nil {
			return
		}
	}

	return "", nil
}