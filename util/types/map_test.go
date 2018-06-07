package types

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMap(t *testing.T) {
	m := NewMap()

	m.Set("B", "This")
	m.Set("A", "That")
	m.Set("C", "Other")

	if m.Get("A").(string) != "That" {
		t.Errorf("Strings test failed. Expected  (%q) got (%q).", "That", m.Get("A").(string))
	}
}

func TestMap_Nil(t *testing.T) {
	m := NewMap()

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

func TestMap_MarshalBinary(t *testing.T) {
	m := NewMap()
	m3 := NewMap()

	m.Set("B", "This")
	m.Set("A", "That")
	m.Set("C", "Other")

	m3.Set("D", 1)
	m3.Set("E", 2)
	m3.Set("F", 3)

	m.Set("map", m3)

	data, err := m.MarshalBinary()
	assert.NoError(t, err)

	m2 := NewMap()
	m2.UnmarshalBinary(data)

	assert.Equal(t, 4, m2.Len())
	assert.Equal(t, "That", m2.Get("A").(string))
	assert.Equal(t, 2, m2.Get("map").(*Map).Get("E").(int))
}
