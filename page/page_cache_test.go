package page

import (
	"testing"
)


func TestBasicPageCache(t *testing.T) {
	cache := NewPageCache()

	p1:=NewPage("1")
	p2:=NewPage("2")

	cache.Set("1", p1)
	cache.Set("2", p2)

	p3 := cache.Get("1")

	if p3 != p1 {
		t.Error("Could not retrieve page")
	}
}
