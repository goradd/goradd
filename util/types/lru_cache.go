package types

import (
	"sync"
	"sort"
	"time"
)

// LruCache is a kind of LRU cache. Objects that are too old are removed, and if the cache is full,
// the oldest item(s) will be removed. When an item is set more than once, it is pushed to the end so its last to be removed.
// Limits are approximate, as garbage collecting will randomly happen. Also, in order to prevent memory thrashing, strict order is not
// preserved, but items will fall out more or less in LRU order.
// Someday we may offer a backing store version to extend the size of the cache to disk or some other kind of storage
type LruCache struct {
	sync.RWMutex
	maxItemCount int
	ttl int64
	items map[string] lruItem
	order []string
}

type lruItem struct {
	timestamp int64
	v interface{}
}

// Create and return a new cache.
// maxItemCount is the maximum number of items the cache can hold
// ttl is the age in seconds past when items will be removed
func NewLruCache(maxItemCount int, ttl int64) *LruCache {
	return &LruCache{
		maxItemCount:maxItemCount,
		ttl: ttl * (1000 * 1000 * 1000),	// we compare against nanos
		items: make(map[string] lruItem, maxItemCount),
		order: make([]string, 0, maxItemCount),
	}
}

// Puts the item into the cache, and updates its access time, pushing it to the end of the removal queue
func (o *LruCache) Set(key string, v interface{})  {
	o.Lock()
	defer o.Unlock()

	if v == nil {
		panic ("Cannot put a nil pointer into the lru cache")
	}
	if key == "" {
		panic ("Cannot use a blank key in the lru cache")
	}

	t := time.Now().UnixNano()

	if item, ok :=  o.items[key]; ok {
		// already exists in the cache
		i := sort.Search(len(o.order), func(i int) bool { return o.items[o.order[i]].timestamp >= item.timestamp })
		if o.items[o.order[i]].timestamp != item.timestamp {
			panic ("Lru cache is not in sync")	// this would be a bug in the cache if this happens
		}
		// To prevent some memory thrashing with an active cache, we only update the order if the item is slightly stale
		if i < o.maxItemCount / 2 || item.timestamp < t - o.ttl / 8 { // if 1/8th of the ttl has passed, or the item is getting close to getting pushed off the end, bring it to front
			o.order = append(o.order[:i], o.order[i + 1:]...)
			o.order = append(o.order, key)
			item.timestamp = t
			o.items[key] = item
		}
	} else {
		// new item
		o.items[key] = lruItem{t, v}
		o.order = append(o.order, key)
	}

	// garbage collect
	if t % ((int64(o.maxItemCount) / 8) + 1 ) == 1 {
		go o.gc()
	}
}

// Garbage collect
// However we do this, we must MAKE SURE that any recent Set does not get garbage collected here
func (o *LruCache) gc() {
	o.Lock()
	defer o.Unlock()

	// remove based on TTL
	for len(o.order) > 0 && o.items[o.order[0]].timestamp < time.Now().UnixNano() - o.ttl {
		delete(o.items, o.order[0])
		o.order = o.order[1:]
	}

	// remove based on size
	for len(o.order) > o.maxItemCount {
		// TODO: If this is happening, we are throwing out items before they expire. We should log this, and it means
		// either allocating more memory, reducing TTL, or implementing a backing store to keep items on disk
		// a backing store will require serialization of the v objects
		delete(o.items, o.order[0])
		o.order = o.order[1:]
	}
}


// Get returns the item based on its id.
// If not found, it will return nil.
func (o *LruCache) Get(key string)  interface{} {
	o.Lock()
	defer o.Unlock()
	var i lruItem
	var ok bool
	if i, ok =  o.items[key]; !ok {
		return nil
	}
	return i.v
}


// Has tests for the existence of the key
func (o *LruCache) Has(key string) (exists bool) {
	o.Lock()
	defer o.Unlock()
	_, exists =  o.items[key]
	return
}

