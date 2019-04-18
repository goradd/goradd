package db

type DatabaseIGetter interface {
	Get(key string) (val DatabaseI)
}

type DatabaseILoader interface {
	Load(key string) (val DatabaseI, ok bool)
}

type DatabaseISetter interface {
	Set(string, DatabaseI)
}


// The DatabaseIMapI interface provides a common interface to the many kinds of similar map objects.
//
// Most functions that change the map are omitted so that you can wrap the map in additional functionality that might
// use Set or SetChanged. If you want to use them in an interface setting, you can create your own interface
// that includes them.
type DatabaseIMapI interface {
	Get(key string) (val DatabaseI)
	Has(key string) (exists bool)
	Values() []DatabaseI
	Keys() []string
	Len() int
	// Range will iterate over the keys and values in the map. Pattern is taken from sync.Map
	Range(f func(key string, value DatabaseI) bool)
	Merge(i DatabaseIMapI)
}
