package page

import (
	"github.com/goradd/goradd/pkg/cache"
	"github.com/goradd/goradd/pkg/html"
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


// SetPageCache will set the page cache to the given object.
func SetPageCache(c PageCacheI) {
	pageCache = c
}

// GetPageCache returns the page cache. Used internally by goradd.
func GetPageCache() PageCacheI {
	return pageCache
}

func HasPage(pageStateId string) bool {
	return pageCache.Has(pageStateId)
}

// FastPageCache is an in memory page cache that does no serialization and uses an LRU cache of page objects.
// Objects that are too old are removed, and if the cache is full,
// the oldest item(s) will be removed. When a page is updated, it is moved to the top. Whenever an item is set,
// we could potentially garbage collect. This cache can be used in a production environment if the
// application is guaranteed to only work on a single machine. If you want scalability, use a serializing
// page cache that serializes directly to a database that is accessible from all instances of the app.
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
//
// This also uses an in memory, non-serialized map to keep the pages in memory so that the testing harness
// can perform browser-based tests. Essentially this cache should only be used for development purposes
// and not production.
type SerializedPageCache struct {
	cache.LruCache
	testPageID string // Used for testing serialization using automated testing. Not for production.
	testPage   *Page

}

// This special interface is used by our test harness to prevent the serialization of the
// test form.
type TestFormI interface {
	NoSerialize() bool
}

func NewSerializedPageCache(maxEntries int, TTL int64) *SerializedPageCache {
	return &SerializedPageCache{
		LruCache:*cache.NewLruCache(maxEntries, TTL),
	}
}

// Set puts the page into the page cache, and updates its access time, pushing it to the end of the removal queue
// Page must already be assigned a state ID. Use NewPageId to do that.
func (o *SerializedPageCache) Set(pageId string, page *Page) {
	if _,ok := page.Form().(TestFormI); ok {
		o.testPageID = pageId
		o.testPage = page
		return
	}

	if b, err := page.MarshalBinary(); err == nil {
		o.LruCache.Set(pageId, b)
	}
}

// Get returns the page based on its page id.
// If not found, will return null.
func (o *SerializedPageCache) Get(pageId string) *Page {
	if pageId == o.testPageID {
		return o.testPage
	}

	b := o.LruCache.Get(pageId)
	if b == nil {
		return nil
	}

	var p Page

	// write over the top of the previous page to reuse the memory
	if err := p.UnmarshalBinary(b.([]byte)); err != nil {
		return nil
	}

	if p.stateId != pageId {
		panic("pageId does not match") // or return nil?
	}
	p.Restore()
	return &p
}

// Has returns true if the page with the given pageId is in the cache.
func (o *SerializedPageCache) Has(pageId string) bool {
	if pageId == o.testPageID {
		return true
	}
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
