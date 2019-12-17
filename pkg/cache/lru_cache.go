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
	sync.Mutex
	maxItemCount int
	ttl          int64
	items        map[string]lruItem
}

type lruItem struct {
	timestamp int64
	value interface{}
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
		items: make(map[string]lruItem),
	}
}

// Puts the item into the cache, and updates its access time
func (o *LruCache) Set(key string, v interface{}) {
	o.Lock()
	if v == nil {
		panic("Cannot put a nil pointer into the lru cache")
	}
	if key == "" {
		panic("Cannot use a blank key in the lru cache")
	}

	t := time.Now().UnixNano()
	i := lruItem{t, v}
	o.items[key] = i
	o.Unlock()
	// garbage collect

	if t%((int64(o.maxItemCount)/8)+1) == 1 {
		go o.gc()
	}
}

// Garbage collect. Garbage collection requires significant time, so it is done in a go routine.
func (o *LruCache) gc() {
	o.Lock()
	var keys []string
	for k,_ := range o.items {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i,j int) bool {
		return o.items[keys[i]].timestamp < o.items[keys[j]].timestamp
	})
	now := time.Now().UnixNano()
	var itemNum int
	var k string
	for itemNum,k = range keys {
		if o.items[k].timestamp + o.ttl < now {
			delete(o.items, k)
			if c,ok := o.items[k].value.(Remover); ok {
				c.Removed()
			}
		} else {
			break
		}
	}

	if len(o.items) > o.maxItemCount {
		// TODO: log that the cache is filling up
		for _,k := range keys[itemNum:] {
			delete(o.items, k)
			if c,ok := o.items[k].value.(Remover); ok {
				c.Removed()
			}
			if len(o.items) <= o.maxItemCount {
				break
			}
		}
	}
	o.Unlock()
}

// Get returns the item based on its id, and updates its access time.
// If not found, it will return nil.
func (o *LruCache) Get(key string) interface{} {
	o.Lock()
	i, ok := o.items[key]
	if !ok {
		o.Unlock()
		return nil
	}
	i.timestamp = time.Now().UnixNano()
	o.items[key] = i
	o.Unlock()
	return i.value
}

// Has tests for the existence of the key. It does not update the access time though.
func (o *LruCache) Has(key string) (exists bool) {
	o.Lock()
	_, ok := o.items[key]
	o.Unlock()
	return ok
}
