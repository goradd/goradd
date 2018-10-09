package old

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSafeMap(t *testing.T) {
	m := NewSafeMap()

	m.Set("B", "This")
	m.Set("A", "That")
	m.Set("C", "Other")

	if m.Get("A").(string) != "That" {
		t.Errorf("Strings test failed. Expected  (%q) got (%q).", "That", m.Get("A").(string))
	}
}

func TestSafeMap_Nil(t *testing.T) {
	m := NewSafeMap()

	m.Set("a", nil)

	b := m.Get("a")

	assert.Nil(t, b)
	assert.True(t, b == nil)

	var c *int

	m.Set("c", c)

	d := m.Get("c")

	assert.Nil(t, d)
	//assert.True(t, d==nil)

	e := d.(*int)

	assert.Nil(t, e)
	assert.True(t, e == nil)
}

func TestSafeMap_MarshalBinary(t *testing.T) {
	m := NewSafeMap()

	m.Set("B", "This")
	m.Set("A", "That")
	m.Set("C", "Other")

	m3 := NewSafeMap()

	m3.Set("D", 1)
	m3.Set("E", 2)
	m3.Set("F", 3)

	m.Set("map", m3)

	data, _ := m.MarshalBinary()

	m2 := NewSafeMap()
	m2.UnmarshalBinary(data)

	assert.Equal(t, 4, m2.Len())
	assert.Equal(t, "That", m2.Get("A").(string))
}
