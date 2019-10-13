package cache

import (
	"testing"
	"time"
)

func TestBasicLruCache(t *testing.T) {
	c := NewLruCache(100, 60*60)

	p1 := "1"
	p2 := "2"

	c.Set("1", p1)
	c.Set("2", p2)

	p3 := c.Get("1")
	if p3 != p1 {
		t.Error("Could not retrieve item")
	}

}

func TestLruReplace(t *testing.T) {
	c := NewLruCache(100, 60*60)

	p1 := "1"
	p2 := "2"

	c.Set("1", p1)
	c.Set("2", p2)

	p3 := c.Get("1")
	if p3 != p1 {
		t.Error("Could not retrieve item")
	}

	p4 := "4"
	c.Set("1", p4)

	if c.Get("1") != "4" {
		t.Error("Could not replace item")
	}
}

func TestLruCacheExit(t *testing.T) {
	c := NewLruCache(1, 60*60)
	p1 := "1"
	p2 := "2"

	c.Set("1", p1)
	c.Set("2", p2)
	c.gc()

	p3 := c.Get("1")

	if p3 != nil {
		t.Error("Item did not fall off end")
	}

	p3 = c.Get("2")
	if p3 != p2 {
		t.Error("Item was lost")
	}
}

func TestLruCacheTtl(t *testing.T) {
	c := NewLruCache(10, 1)
	p1 := "1"
	p2 := "2"

	c.Set("1", p1)
	c.Set("2", p2)
	time.Sleep(time.Second)
	c.gc()

	p3 := c.Get("1")

	if p3 != nil {
		t.Error("Item did not expire")
	}

	p3 = c.Get("2")
	if p3 != nil {
		t.Error("Item did not expire")
	}

}

func TestLruReset(t *testing.T) {
	c := NewLruCache(2, 60*60)
	p1 := "1"
	p2 := "2"
	p3 := "3"

	c.Set("1", p1)
	c.Set("2", p2)
	c.Set("1", p1)
	c.Set("3", p3)
	c.Set("1", p1)
	c.gc()

	p4 := c.Get("2")

	if p4 != nil {
		t.Error("Item did not fall off end")
	}

	p4 = c.Get("1").(string)
	if p4 != p1 {
		t.Error("Item was lost")
	}
}
