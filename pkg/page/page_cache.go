package page

import (
	"github.com/goradd/goradd/pkg/cache"
	"github.com/goradd/goradd/pkg/config"
	"github.com/goradd/goradd/pkg/html5tag"
	"github.com/goradd/goradd/pkg/log"
)

// PagestateCacheI is the page cache interface. The PageCache saves and restores pages in between page
// accesses by the user.
type PagestateCacheI interface {
	Set(pageId string, page *Page)
	Get(pageId string) *Page
	NewPageID() string
	Has(pageId string) bool
}

var pageCache PagestateCacheI


// SetPagestateCache will set the page cache to the given object.
func SetPagestateCache(c PagestateCacheI) {
	pageCache = c
}

// GetPagestateCache returns the page cache. Used internally by goradd.
func GetPagestateCache() PagestateCacheI {
	return pageCache
}

func HasPage(pageStateId string) bool {
	return pageCache.Has(pageStateId)
}

// FastPagestateCache is an in memory page cache that does no serialization and uses an LRU cache of page objects.
// Objects that are too old are removed, and if the cache is full,
// the oldest item(s) will be removed. When a page is updated, it is moved to the top. Whenever an item is set,
// we could potentially garbage collect. This cache can be used in a production environment if the
// application is guaranteed to only work on a single machine. If you want scalability, use a serializing
// page cache that serializes directly to a database that is accessible from all instances of the app.
type FastPagestateCache struct {
	cache.LruCache
}

// NewFastPageCache creates a new FastPagestateCache cache
func NewFastPageCache(maxEntries int, TTL int64) *FastPagestateCache {
	return &FastPagestateCache{*cache.NewLruCache(maxEntries, TTL)}
}

// Set puts the page into the page cache and updates its access time, pushing it to the end of the removal queue.
// Page must already be assigned a state ID. Use NewPageId to do that.
func (o *FastPagestateCache) Set(pageId string, page *Page) {
	o.LruCache.Set(pageId, page)
}

// Get returns the page based on its page id.
// If not found, will return null.
func (o *FastPagestateCache) Get(pageId string) *Page {
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
func (o *FastPagestateCache) Has(pageId string) bool {
	return o.LruCache.Has(pageId)
}

// NewPageID returns a new page id
func (o *FastPagestateCache) NewPageID() string {
	s := html5tag.RandomString(40)
	for o.Has(s) { // while it is extremely unlikely that we will get a collision, a collision is such a huge security problem we must make sure
		s = html5tag.RandomString(40)
	}
	return s
}

// SerializedPagestateCache is an in memory page cache that does serialization and uses an LRU cache of page objects.
// Use the serialized page cache during development to ensure that you can eventually move your page cache to a database
// or a separate machine so that your application is scalable.
//
// This also uses an in memory, non-serialized map to keep the pages in memory so that the testing harness
// can perform browser-based tests. Essentially this cache should only be used for development purposes
// and not production.
type SerializedPagestateCache struct {
	cache.LruCache
	testPageID string // Used for testing serialization using automated testing. Not for production.
	testPage   *Page

}

// This special interface is used by our test harness to prevent the serialization of the
// test form.
type TestFormI interface {
	NoSerialize() bool
}

func NewSerializedPageCache(maxEntries int, TTL int64) *SerializedPagestateCache {
	return &SerializedPagestateCache{
		LruCache:*cache.NewLruCache(maxEntries, TTL),
	}
}

// Set puts the page into the page cache, and updates its access time, pushing it to the end of the removal queue
// Page must already be assigned a state ID. Use NewPageId to do that.
func (o *SerializedPagestateCache) Set(pageId string, page *Page) {
	if _,ok := page.Form().(TestFormI); ok {
		o.testPageID = pageId
		o.testPage = page
		return
	}

	if b, err := page.MarshalBinary(); err == nil {
		o.LruCache.Set(pageId, b)
		log.FrameworkDebug("Write page to cache: ", pageId)
	}
}

// Get returns the page based on its page id.
// If not found, will return null.
func (o *SerializedPagestateCache) Get(pageId string) *Page {
	if pageId == o.testPageID {
		return o.testPage
	}

	b := o.LruCache.Get(pageId)
	log.FrameworkDebug("Get page from cache: ", pageId)

	if b == nil {
		log.FrameworkDebug("Page not found: ", pageId)
		return nil
	}

	var p Page

	// write over the top of the previous page to reuse the memory
	if err := p.UnmarshalBinary(b.([]byte)); err != nil {
		if config.Debug {
			panic("Page unmarshal error: " + err.Error())
		} else {
			log.FrameworkDebug("Page unmarshal error: ", err.Error())
		}
		return nil
	}

	if p.stateId != pageId {
		panic("pageId does not match") // or return nil?
	}
	p.Restore()
	return &p
}

// Has returns true if the page with the given pageId is in the cache.
func (o *SerializedPagestateCache) Has(pageId string) bool {
	if pageId == o.testPageID {
		return true
	}
	return o.LruCache.Has(pageId)
}

// NewPageID returns a new page id
func (o *SerializedPagestateCache) NewPageID() string {
	s := html5tag.RandomString(40)
	for o.Has(s) { // while it is extremely unlikely that we will get a collision, a collision is such a huge security problem we must make sure
		s = html5tag.RandomString(40)
	}
	return s
}
