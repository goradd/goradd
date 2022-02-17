package sql

type JoinTreeItemGetter interface {
	Get(key string) (val *JoinTreeItem)
}

type JoinTreeItemLoader interface {
	Load(key string) (val *JoinTreeItem, ok bool)
}

type JoinTreeItemSetter interface {
	Set(string, *JoinTreeItem)
}


// The JoinTreeItemMapI interface provides a common interface to the many kinds of similar map objects.
//
// Most functions that change the map are omitted so that you can wrap the map in additional functionality that might
// use Set or SetChanged. If you want to use them in an interface setting, you can create your own interface
// that includes them.
type JoinTreeItemMapI interface {
	Get(key string) (val *JoinTreeItem)
	Has(key string) (exists bool)
	Values() []*JoinTreeItem
	Keys() []string
	Len() int
	// Range will iterate over the keys and values in the map. Pattern is taken from sync.Map
	Range(f func(key string, value *JoinTreeItem) bool)
	Merge(i JoinTreeItemMapI)
	String() string
}
