package page

import (
	"bytes"
	"github.com/spekary/goradd/pkg/html"
	"github.com/spekary/goradd/ideas/types"
	"goradd-project/config"
)

type PageCacheI interface {
	Set(pageId string, page *Page)
	Get(pageId string) *Page
	NewPageID() string
	Has(pageId string) bool
}

var pageCache PageCacheI

func SetPageCache(c PageCacheI) {
	pageCache = c
}

// GetPageCache returns the page cache. Used internally by goradd.
func GetPageCache() PageCacheI {
	return pageCache
}

// FastPageCache is an in memory override cache that does no serialization and uses an LRU cache of override objects.
// Objects that are too old are removed, and if the cache is full,
// the oldest item(s) will be removed. Pages that are set multiple times will be pushed to the top. Whenever an item is set,
// we could potentially garbage collect. This cache is only appropriate when the pagecache itself is operating on a
// single machine.
type FastPageCache struct {
	types.LruCache
}

func NewFastPageCache() *FastPageCache {
	return &FastPageCache{*types.NewLruCache(config.PageCacheMaxSize, config.PageCacheTTL)}
}

// Puts the override into the override cache, and updates its access time, pushing it to the end of the removal queue
// Page must already be assigned a state ID. Use NewPageId to do that.
func (o *FastPageCache) Set(pageId string, page *Page) {
	o.LruCache.Set(pageId, page)
}

// Get returns the override based on its override id.
// If not found, will return null.
func (o *FastPageCache) Get(pageId string) *Page {
	var p *Page

	if i := o.LruCache.Get(pageId); i != nil {
		p = i.(*Page)
	}

	if p != nil && p.stateId != pageId {
		panic("pageId does not match")
	}
	return p
}

// Has tests to see if the given page id is in the page cache, without actually loading the page
func (o *FastPageCache) Has(pageId string) bool {
	return o.LruCache.Has(pageId)
}


// Returns a new override id
func (o *FastPageCache) NewPageID() string {
	s := html.RandomString(40)
	for o.Has(s) { // while it is extremely unlikely that we will get a collision, a collision is such a huge security problem we must make sure
		s = html.RandomString(40)
	}
	return s
}

// SerializedPageCache is an in memory override cache that does serialization and uses an LRU cache of override objects.
// Use the serialized override cache during development to ensure that you can eventually move your override cache to a database
// or a separate machine so that your application is scalable.
// Objects that are too old are removed, and if the cache is full,
// the oldest item(s) will be removed. Pages that are set multiple times will be pushed to the top. Whenever an item is set,
// we could potentially garbage collect.
type SerializedPageCache struct {
	types.LruCache
}

func NewSerializedPageCache() *SerializedPageCache {
	panic("Serialized pages are not ready for prime time yet")
	return &SerializedPageCache{*types.NewLruCache(config.PageCacheMaxSize, config.PageCacheTTL)}
}

// Puts the override into the override cache, and updates its access time, pushing it to the end of the removal queue
// Page must already be assigned a state ID. Use NewPageId to do that.
func (o *SerializedPageCache) Set(pageId string, page *Page) {
	b := GetBuffer()
	defer PutBuffer(b)
	enc := pageEncoder.NewEncoder(b)
	enc.Encode(config.PageCacheVersion)
	enc.Encode(page.Form().ID())
	err := page.Encode(enc)
	if err != nil {
		o.LruCache.Set(pageId, b.Bytes())
	}
}

// Get returns the page based on its page id.
// If not found, will return null.
func (o *SerializedPageCache) Get(pageId string) *Page {
	b := o.LruCache.Get(pageId).([]byte)
	dec := pageEncoder.NewDecoder(bytes.NewBuffer(b))
	var ver int32
	if err := dec.Decode(&ver); err != nil {
		panic(err)
	}
	if ver != config.PageCacheVersion {
		return nil
	}

	var formId string
	var p Page
	if err := dec.Decode(&formId); err != nil {
		panic(err)
	}
	if _, ok := pageManager.formIdRegistry[formId]; !ok {
		panic("Form id not found")
	}

	if err := dec.Decode(&p); err != nil {
		panic(err)
	}

	if p.stateId != pageId {
		panic("pageId does not match")
	}
	//p.Restore()
	return &p
}

func (o *SerializedPageCache) Has(pageId string) bool {
	return o.LruCache.Has(pageId)
}

// Returns a new override id
func (o *SerializedPageCache) NewPageID() string {
	s := html.RandomString(40)
	for o.Has(s) { // while it is extremely unlikely that we will get a collision, a collision is such a huge security problem we must make sure
		s = html.RandomString(40)
	}
	return s
}
