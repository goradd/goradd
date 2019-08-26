package page

import (
	"bytes"
	"github.com/goradd/goradd/pkg/cache"
	"github.com/goradd/goradd/pkg/html"
	"github.com/goradd/goradd/pkg/pool"
)

// PageCacheI is the page cache interface. The PageCache saves and restores pages in between page
// accesses by the user.
type PageCacheI interface {
	Set(pageId string, page *Page)
	Get(pageId string) *Page
	NewPageID() string
	Has(pageId string) bool
}

var pageCache PageCacheI

// PageCacheVersion helps us keep track of when a change to the application changes the pagecache format. It is only needed
// when serializing the pagecache. Some page cache stores may be difficult to invalidate the whole thing, so this lets
// lets us invalidate old pagecaches individually. If you implement your own pagecache, you may want to control this
// independently of goradd, which is why this is exported. Goradd should bump this value whenever the pagecache serialization format
// changes.
var PageCacheVersion int32 = 1

// SetPageCache will set the page cache to the given object.
func SetPageCache(c PageCacheI) {
	pageCache = c
}

// GetPageCache returns the page cache. Used internally by goradd.
func GetPageCache() PageCacheI {
	return pageCache
}

// FastPageCache is an in memory page cache that does no serialization and uses an LRU cache of page objects.
// Objects that are too old are removed, and if the cache is full,
// the oldest item(s) will be removed. When a page is updated, it is moved to the top. Whenever an item is set,
// we could potentially garbage collect. This cache is only appropriate when the pagecache itself is operating on a
// single machine.
type FastPageCache struct {
	cache.LruCache
}

// NewFastPageCache creates a new FastPageCache cache
func NewFastPageCache(maxEntries int, TTL int64) *FastPageCache {
	return &FastPageCache{*cache.NewLruCache(maxEntries, TTL)}
}

// Set puts the page into the page cache and updates its access time, pushing it to the end of the removal queue.
// Page must already be assigned a state ID. Use NewPageId to do that.
func (o *FastPageCache) Set(pageId string, page *Page) {
	o.LruCache.Set(pageId, page)
}

// Get returns the page based on its page id.
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

// NewPageID returns a new page id
func (o *FastPageCache) NewPageID() string {
	s := html.RandomString(40)
	for o.Has(s) { // while it is extremely unlikely that we will get a collision, a collision is such a huge security problem we must make sure
		s = html.RandomString(40)
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
	cache.LruCache
}

func NewSerializedPageCache(maxEntries int, TTL int64) *SerializedPageCache {
	panic("Serialized pages are not ready for prime time yet")
	return &SerializedPageCache{*cache.NewLruCache(maxEntries, TTL)}
}

// Set puts the page into the page cache, and updates its access time, pushing it to the end of the removal queue
// Page must already be assigned a state ID. Use NewPageId to do that.
func (o *SerializedPageCache) Set(pageId string, page *Page) {
	b := pool.GetBuffer()
	defer pool.PutBuffer(b)
	enc := pageEncoder.NewEncoder(b)
	_ = enc.Encode(PageCacheVersion)
	_ = enc.Encode(page.Form().ID())
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
	if ver != PageCacheVersion {
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

// Has returns true if the page with the given pageId is in the cache.
func (o *SerializedPageCache) Has(pageId string) bool {
	return o.LruCache.Has(pageId)
}

// NewPageID returns a new page id
func (o *SerializedPageCache) NewPageID() string {
	s := html.RandomString(40)
	for o.Has(s) { // while it is extremely unlikely that we will get a collision, a collision is such a huge security problem we must make sure
		s = html.RandomString(40)
	}
	return s
}
