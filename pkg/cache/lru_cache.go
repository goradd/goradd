package cache

import (
	"sort"
	"sync"
	"time"
)

// LruCache is a kind of LRU cache. Objects that are too old are removed, and if the cache is full,
// the oldest item(s) will be removed. When an item is set more than once, it is pushed to the end so its last to be removed.
// Limits are approximate, as garbage collecting will randomly happen. Also, in order to prevent memory thrashing, strict order is not
// preserved, but items will fall out more or less in LRU order.
// Someday we may offer a backing store version to extend the size of the cache to disk or some other kind of storage
// If the item has a "Removed" function, that function will be called when the item falls out of the cache.
// If the item has a "Cleanup" function, that function will be called when the item is removed from memory. If
// the cache has a backing store, it may be removed from memory, but still in the disk-based cache.
type LruCache struct {
	sync.RWMutex
	maxItemCount int
	ttl          int64
	items        map[string]lruItem
	order        []string
}

type lruItem struct {
	timestamp int64
	v         interface{}
}

type Remover interface {
	Removed()
}

type Cleanuper interface {
	Cleanup()
}


// Create and return a new cache.
// maxItemCount is the maximum number of items the cache can hold
// ttl is the age in seconds past when items will be removed
func NewLruCache(maxItemCount int, ttl int64) *LruCache {
	return &LruCache{
		maxItemCount: maxItemCount,
		ttl:          ttl * (1000 * 1000 * 1000), // we compare against nanos
		items:        make(map[string]lruItem, maxItemCount),
		order:        make([]string, 0, maxItemCount),
	}
}

// Puts the item into the cache, and updates its access time, pushing it to the end of the removal queue
func (o *LruCache) Set(key string, v interface{}) {
	o.Lock()

	if v == nil {
		panic("Cannot put a nil pointer into the lru cache")
	}
	if key == "" {
		panic("Cannot use a blank key in the lru cache")
	}

	t := time.Now().UnixNano()

	if item, ok := o.items[key]; ok {
		// already exists in the cache
		i := sort.Search(len(o.order), func(i int) bool { return o.items[o.order[i]].timestamp >= item.timestamp })
		if o.items[o.order[i]].timestamp != item.timestamp {
			panic("Lru cache is not in sync") // this would be a bug in the cache if this happens
		}
		// To prevent some memory thrashing with an active cache, we only update the order if the item is slightly stale
		if i < o.maxItemCount/2 || item.timestamp < t-o.ttl/8 { // if 1/8th of the ttl has passed, or the item is getting close to getting pushed off the end, bring it to front
			o.order = append(o.order[:i], o.order[i+1:]...)
			o.order = append(o.order, key)
			item.timestamp = t
			o.items[key] = item
		}
	} else {
		// new item
		o.items[key] = lruItem{t, v}
		o.order = append(o.order, key)
	}
	o.Unlock()

	// garbage collect
	if t%((int64(o.maxItemCount)/8)+1) == 1 {
		go o.gc()
	}
}

// Garbage collect
func (o *LruCache) gc() {
	// However we do this, we must MAKE SURE that any recent Set does not get garbage collected here
	o.Lock()


	// remove based on TTL
	for len(o.order) > 0 && o.items[o.order[0]].timestamp < time.Now().UnixNano()-o.ttl {
		v := o.items[o.order[0]].v
		delete(o.items, o.order[0])
		o.order = o.order[1:]
		if r,ok := v.(Remover); ok {
			r.Removed()
		}
		if r,ok := v.(Cleanuper); ok {
			r.Cleanup()
		}
	}

	// remove based on size
	for len(o.order) > o.maxItemCount {
		// TODO: If this is happening, we are throwing out items before they expire. We should log this, and it means
		// either allocating more memory, reducing TTL, or implementing a backing store to keep items on disk
		// a backing store will require serialization of the v objects
		v := o.items[o.order[0]].v
		delete(o.items, o.order[0])
		o.order = o.order[1:]
		if r,ok := v.(Remover); ok {
			r.Removed()
		}
		if r,ok := v.(Cleanuper); ok {
			r.Cleanup()
		}
	}
	o.Unlock()
}

// Get returns the item based on its id.
// If not found, it will return nil.
func (o *LruCache) Get(key string) interface{} {
	var i lruItem
	var ok bool

	o.RLock()
	i, ok = o.items[key]
	o.RUnlock()
	if !ok {
		return nil
	}
	return i.v
}

// Has tests for the existence of the key
func (o *LruCache) Has(key string) (exists bool) {
	o.RLock()
	_, exists = o.items[key]
	o.RUnlock()
	return
}
