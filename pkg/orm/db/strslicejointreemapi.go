package db

type joinTreeItemGetter interface {
	Get(key string) (val *joinTreeItem)
}

type joinTreeItemLoader interface {
	Load(key string) (val *joinTreeItem, ok bool)
}

type joinTreeItemSetter interface {
	Set(string, *joinTreeItem)
}

// The joinTreeItemMapI interface provides a common interface to the many kinds of similar map objects.
//
// Most functions that change the map are omitted so that you can wrap the map in additional functionality that might
// use Set or SetChanged. If you want to use them in an interface setting, you can create your own interface
// that includes them.
type joinTreeItemMapI interface {
	Get(key string) (val *joinTreeItem)
	Has(key string) (exists bool)
	Values() []*joinTreeItem
	Keys() []string
	Len() int
	// Range will iterate over the keys and values in the map. Pattern is taken from sync.Map
	Range(f func(key string, value *joinTreeItem) bool)
	Merge(i joinTreeItemMapI)
}
