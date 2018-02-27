package page

import (
	"goradd/config"
	"github.com/spekary/goradd/util/types"
	"github.com/spekary/goradd/util"
)

// The page cache is an LRU cache of page objects. Objects that are too old are removed, and if the cache is full,
// the oldest item(s) will be removed. Pages that are set multiple times will be pushed to the top. Whenever an item
// is retrieved from the cache, it is removed. Whenever an item is set, we garbage collect.
type pageCache struct {
	types.LruCache
}

func NewPageCache() *pageCache {
	return &pageCache{*types.NewLruCache(config.PAGE_CACHE_MAX_SIZE, config.PAGE_CACHE_TTL)}
}

// Puts the page into the page cache, and updates its access time, pushing it to the end of the removal queue
// Page must already be assigned a state ID. Use NewPageId to do that.
func (o *pageCache) Set(pageId string, page PageI)  {
	o.LruCache.Set(pageId, page)
}


// Get returns the page based on its page id.
// If not found, will return null.
func (o *pageCache) Get(pageId string) (PageI) {
	var p PageI
	p = o.LruCache.Get(pageId).(PageI)

	if p != nil && p.GetPageBase().stateId != pageId {
		panic("pageId does not match")
	}
	return p
}

// Returns a new page id
func (o *pageCache) NewPageId() string {
	s := util.RandomHtmlValueString(40)
	for o.Has(s) {	// while it is extremely unlikely that we will get a collision, a collision is such a huge security problem we must make sure
		s = util.RandomHtmlValueString(40)
	}
	return s
}