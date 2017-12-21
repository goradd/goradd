package types

import (
	"testing"
	"time"
)


func TestBasicLruCache(t *testing.T) {
	cache := NewLruCache(100, 60 * 60)

	p1:="1"
	p2:="2"

	cache.Set("1", p1)
	cache.Set("2", p2)

	p3 := cache.Get("1")
	if p3 != p1 {
		t.Error("Could not retrieve item")
	}

}

func TestLruCacheExit(t *testing.T) {
	cache := NewLruCache(1, 60 * 60)
	p1:="1"
	p2:="2"

	cache.Set("1", p1)
	cache.Set("2", p2)
	cache.gc()

	p3 := cache.Get("1")

	if p3 != nil {
		t.Error("Item did not fall off end")
	}

	p3 = cache.Get("2")
	if p3 != p2 {
		t.Error("Item was lost")
	}
}

func TestLruCacheTtl(t *testing.T) {
	cache := NewLruCache(10, 1)
	p1:="1"
	p2:="2"

	cache.Set("1", p1)
	cache.Set("2", p2)
	time.Sleep(time.Second)
	cache.gc()

	p3 := cache.Get("1")

	if p3 != nil {
		t.Error("Item did not expire")
	}

	p3 = cache.Get("2")
	if p3 != nil {
		t.Error("Item did not expire")
	}

}
