package page

import (
	"goradd/config"
	"github.com/spekary/goradd/util/types"
	"github.com/spekary/goradd/util"
	"bytes"
)

type PageCacheI interface {
	Set(pageId string, page PageI)
	Get(pageId string) PageI
	NewPageId() string
}

var pageCache PageCacheI


func SetPageCache(c PageCacheI) {
	if pageCache != nil {
		panic("Only set the page cache when the application is initialized, and only once.")
	}
	pageCache = c
}

// FastPageCache is an in memory page cache that does no serialization and uses an LRU cache of page objects.
// Objects that are too old are removed, and if the cache is full,
// the oldest item(s) will be removed. Pages that are set multiple times will be pushed to the top. Whenever an item is set,
// we could potentially garbage collect. This cache is only appropriate when the pagecache itself is operating on a
// single machine.
type FastPageCache struct {
	types.LruCache
}

func NewFastPageCache() *FastPageCache {
	return &FastPageCache{*types.NewLruCache(config.PAGE_CACHE_MAX_SIZE, config.PAGE_CACHE_TTL)}
}

// Puts the page into the page cache, and updates its access time, pushing it to the end of the removal queue
// Page must already be assigned a state ID. Use NewPageId to do that.
func (o *FastPageCache) Set(pageId string, page PageI)  {
	o.LruCache.Set(pageId, page)
}


// Get returns the page based on its page id.
// If not found, will return null.
func (o *FastPageCache) Get(pageId string) (PageI) {
	var p PageI

	if i:= o.LruCache.Get(pageId); i != nil {
		p = i.(PageI)
	}

	if p != nil && p.GetPageBase().stateId != pageId {
		panic("pageId does not match")
	}
	return p
}

// Returns a new page id
func (o *FastPageCache) NewPageId() string {
	s := util.RandomHtmlValueString(40)
	for o.Has(s) {	// while it is extremely unlikely that we will get a collision, a collision is such a huge security problem we must make sure
		s = util.RandomHtmlValueString(40)
	}
	return s
}

// SerializedPageCache is an in memory page cache that does serialization and uses an LRU cache of page objects.
// Use the serialized page cache during development to ensure that you can eventually move your page cache to a database
// or a separate machine so that your application is scalable.
// Objects that are too old are removed, and if the cache is full,
// the oldest item(s) will be removed. Pages that are set multiple times will be pushed to the top. Whenever an item is set,
// we could potentially garbage collect.
type SerializedPageCache struct {
	types.LruCache
}

func NewSerializedPageCache() *SerializedPageCache {
	panic ("Serialized pages are not ready for prime time yet")
	return &SerializedPageCache{*types.NewLruCache(config.PAGE_CACHE_MAX_SIZE, config.PAGE_CACHE_TTL)}
}

// Puts the page into the page cache, and updates its access time, pushing it to the end of the removal queue
// Page must already be assigned a state ID. Use NewPageId to do that.
func (o *SerializedPageCache) Set(pageId string, page PageI)  {
	b := GetBuffer()
	defer PutBuffer(b)
	enc := pageEncoder.NewEncoder(b)
	enc.Encode(config.PageCacheVersion)
	enc.Encode(page.Type())
	err := page.Encode(enc)
	if err != nil {
		o.LruCache.Set(pageId, b.Bytes())
	}
}


// Get returns the page based on its page id.
// If not found, will return null.
func (o *SerializedPageCache) Get(pageId string) (PageI) {
	b := o.LruCache.Get(pageId).([]byte)
	dec := pageEncoder.NewDecoder(bytes.NewBuffer(b))
	var ver int32
	if err := dec.Decode(&ver); err != nil {
		panic(err)
	}
	if ver != config.PageCacheVersion {
		return nil
	}

	var pageType string
	var p PageI
	if err := dec.Decode(&pageType); err != nil {
		panic(err)
	}
	if newPageFunc, ok := pageManager.typeRegistry[pageType]; !ok {
		panic("Page type not found")
	} else {
		p = newPageFunc()
	}

	if err := p.Decode(dec); err != nil {
		panic(err)
	}

	if p != nil && p.GetPageBase().stateId != pageId {
		panic("pageId does not match")
	}
	p.Restore()
	return p
}

// Returns a new page id
func (o *SerializedPageCache) NewPageId() string {
	s := util.RandomHtmlValueString(40)
	for o.Has(s) {	// while it is extremely unlikely that we will get a collision, a collision is such a huge security problem we must make sure
		s = util.RandomHtmlValueString(40)
	}
	return s
}
