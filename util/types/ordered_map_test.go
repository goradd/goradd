package types

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"sort"
	"testing"
)

func TestOrderedMap(t *testing.T) {
	var s interface{}
	var ok bool

	m := NewOrderedMap()

	m.Set("B", "This")
	m.Set("A", "That")
	m.Set("C", "Other")

	if m.Values()[1] != "That" {
		t.Errorf("Strings test failed. Expected  (%q) got (%q).", "That", m.Values()[1])
	}

	if m.Keys()[1] != "A" {
		t.Errorf("Keys test failed. Expected  (%q) got (%q).", "A", m.Keys()[1])
	}

	if s = m.GetAt(2); s != "Other" {
		t.Errorf("GetAt test failed. Expected  (%q) got (%q).", "Other", s)
	}

	if s = m.GetAt(3); ok {
		t.Errorf("GetAt test failed. Expected no response, got %q", s)
	}

	// Test that it satisfies the Sort interface
	sort.Sort(m)
	if s = m.GetAt(0); s != "That" {
		t.Error("Sort interface test failed.")
	}

}

func ExampleOrderedMap_Range() {
	m := NewOrderedStringMap()

	m.Set("B", "This")
	m.Set("A", "That")
	m.Set("C", "Other")

	// Iterate by insertion order
	m.Range(func(key string, val string) bool {
		fmt.Printf("%s:%s,", key, val)
		return true // keep iterating to the end
	})
	fmt.Println()

	// Iterate after sorting values
	sort.Sort(m)
	m.Range(func(key string, val string) bool {
		fmt.Printf("%s:%s,", key, val)
		return true // keep iterating to the end
	})
	fmt.Println()

	// Iterate after sorting keys
	sort.Sort(OrderStringMapByKeys(m))
	m.Range(func(key string, val string) bool {
		fmt.Printf("%s:%s,", key, val)
		return true // keep iterating to the end
	})
	fmt.Println()

	// Output: B:This,A:That,C:Other,
	// C:Other,A:That,B:This,
	// A:That,B:This,C:Other,

}

func ExampleOrderedMap_MarshalJSON() {
	m := NewOrderedMap()

	m.Set("A", "That")
	m.Set("B", "This")
	m.Set("C", "Other")

	s, _ := json.Marshal(m)
	os.Stdout.Write(s)

	// Output: {"A":"That","B":"This","C":"Other"}
}

func ExampleOrderedMap_UnmarshalJSON() {
	b := []byte(`{"A":"That","B":"This","C":"Other"}`)
	m := NewOrderedMap()

	json.Unmarshal(b, m)

	for _, i := range m.Values() {
		fmt.Print(i.(string))
	}

	// Output: ThatThisOther
}

func TestOrderedMap_JSON(t *testing.T) {
	b := []byte(`{"A":true,"B":"This","C":{"D":"other","E":"and"}}`)
	m := NewOrderedMap()

	json.Unmarshal(b, &m)

	s, _ := json.Marshal(m)

	if !bytes.Equal(b, s) {
		t.Errorf("JSON test failed.")
	}

	b = []byte(`{"A":true,"B":"This","C":["D","E","F"]}`)
	m = NewOrderedMap()

	json.Unmarshal(b, &m)

	s, _ = json.Marshal(m)

	if !bytes.Equal(b, s) {
		t.Errorf("JSON test failed.")
	}
}

func TestOrderedMap_Nil(t *testing.T) {
	m := NewOrderedMap()

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

func ExampleOrderMap_SetAt() {
	m := NewOrderedMap()

	m.Set("a", 1)
	m.Set("b", 2)

	m.SetAt(1, "c", 3)

	for _, i := range m.Values() {
		fmt.Printf("%d", i.(int))
	}

	// Output: 132
}

func TestOrderedMap_SetAt(t *testing.T) {
	m := NewOrderedMap()

	m.Set("a", 1)
	m.Set("b", 2)

	// Test middle inserts
	m.SetAt(1, "c", 3)
	assert.EqualValues(t, 3, m.GetAt(1))
	m.SetAt(-2, "d", 4)
	assert.EqualValues(t, 4, m.GetAt(2))
	assert.EqualValues(t, 2, m.GetAt(3))

	// Test end inserts
	m.SetAt(-1, "e", 5)
	m.SetAt(1000, "f", 6)
	assert.EqualValues(t, 5, m.GetAt(4))
	assert.EqualValues(t, 6, m.GetAt(5))

	// Test beginning inserts
	m.SetAt(0, "g", 7)
	m.SetAt(-1000, "h", 8)
	assert.EqualValues(t, 8, m.GetAt(0))
	assert.EqualValues(t, 7, m.GetAt(1))
}
