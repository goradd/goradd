package types

import (
	"testing"
	"fmt"
	"sort"
	"encoding/json"
	"os"
	"bytes"
)


func TestOrderedMap(t *testing.T) {
	var s interface{}
	var ok bool

	m := NewOrderedMap()

	m.Set("B", "This")
	m.Set("A","That")
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


func ExampleOrderedMap_Iter() {
	m := NewOrderedMap()

	m.Set("B", "This")
	m.Set("A","That")
	m.Set("C", "Other")

	// Iterate by insertion order
	for s := range m.Iter() {
		fmt.Print(s)
	}
	fmt.Println()

	// Iterate after sorting
	sort.Sort(m)
	for s := range m.Iter() {
		fmt.Print(s)
	}
	fmt.Println()

	// Output: ThisThatOther
	// ThatThisOther
}

func ExampleOrderedMap_IterKeys() {
	m := NewOrderedMap()

	m.Set("B", "This")
	m.Set("A","That")
	m.Set("C", "Other")

	// Iterate on keys in order added
	for s := range m.IterKeys() {
		fmt.Print(s)
	}
	fmt.Println()

	// Iterate after sorting keys
	sort.Sort(OrderMapByKeys(m))
	for s := range m.IterKeys() {
		fmt.Print(s)
	}
	fmt.Println()


	// Output: BAC
	// ABC
}

func ExampleOrderedMap_MarshalJSON() {
	m := NewOrderedMap()

	m.Set("A","That")
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

	for _,i := range m.Values() {
		fmt.Print(i.(string))
	}

	// Output: ThatThisOther
}

func TestOrderedMap_JSON(t *testing.T) {
	b := []byte(`{"A":true,"B":"This","C":{"D":"other","E":"and"}}`)
	m := NewOrderedMap()

	json.Unmarshal(b, &m)

	s, _ := json.Marshal(m)

	if !bytes.Equal(b,s) {
		t.Errorf("JSON test failed.")
	}

	b = []byte(`{"A":true,"B":"This","C":["D","E","F"]}`)
	m = NewOrderedMap()

	json.Unmarshal(b, &m)

	s, _ = json.Marshal(m)

	if !bytes.Equal(b,s) {
		t.Errorf("JSON test failed.")
	}
}
