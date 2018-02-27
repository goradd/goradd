package types

import (
	"testing"
	"fmt"
	"sort"
	"encoding/json"
	"bytes"
	"encoding/gob"
	"os"
)


func TestOrderedStringMap(t *testing.T) {
	var s string
	var ok bool

	m := NewOrderedStringMap()

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


	s = m.Join("+")

	if s != "This+That+Other" {
		t.Error("Failed Join.")
	}

	m.Remove("A")

	s = m.Join("-")

	if s != "This-Other" {
		t.Error("Remove Failed.")
	}

	if m.Len() != 2 {
		t.Error("Len Failed.")
	}

	if m.Has ("NOT THERE") {
		t.Error("Getting non-existant value did not return false")
	}

	val := m.Get("B")
	if val != "This" {
		t.Error("Get failed")
	}

	// Test that it satisfies the StringMapI interface
	var i StringMapI = m
	if s = i.Get("B"); s != "This" {
		t.Error("StringMapI interface test failed.")
	}

	// Test that it satisfies the Sort interface
	sort.Sort(m)
	if s = m.GetAt(0); s != "Other" {
		t.Error("Sort interface test failed.")
	}

	if changed, _ := m.SetChanged("F", "9");  !changed {
		t.Error("Add non-string value failed.")
	}
	if m.Get("F") != "9" {
		t.Error("Add non-string value failed.")
	}
}


func TestOrderedStringMapChange(t *testing.T) {
	m := NewOrderedStringMap()

	m.Set("B", "This")
	m.Set("A","That")
	m.Set("C", "Other")

	if changed, _ := m.SetChanged("D", "And another"); !changed {
		t.Error("Set did not produce a change flag")
	}

	if changed, _ := m.SetChanged("D", "And another"); changed {
		t.Error("Set again erroneously produced a change flag")
	}
}


func ExampleOrderedStringMap_Range() {
	m := NewOrderedStringMap()

	m.Set("B", "This")
	m.Set("A","That")
	m.Set("C", "Other")

	// Iterate by insertion order
	m.Range(func (key string, val string) bool {
		fmt.Printf("%s:%s,", key, val)
		return true	// keep iterating to the end
	})
	fmt.Println()

	// Iterate after sorting values
	sort.Sort(m)
	m.Range(func (key string, val string) bool {
		fmt.Printf("%s:%s,", key, val)
		return true	// keep iterating to the end
	})
	fmt.Println()

	// Iterate after sorting keys
	sort.Sort(OrderStringMapByKeys(m))
	m.Range(func (key string, val string) bool {
		fmt.Printf("%s:%s,", key, val)
		return true	// keep iterating to the end
	})
	fmt.Println()


	// Output: B:This,A:That,C:Other,
	// C:Other,A:That,B:This,
	// A:That,B:This,C:Other,
}


func ExampleOrderedStringMap_MarshalBinary() {
	// You would rarely call MarshallBinary directly, but rather would use an encoder, like GOB for binary encoding

	m := NewOrderedStringMap()
	var m2 OrderedStringMap

	m.Set("B", "This")
	m.Set("A","That")
	m.Set("C", "Other")

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf) // Will write
	dec := gob.NewDecoder(&buf) // Will read

	enc.Encode(m)
	dec.Decode(&m2)
	s := m2.Get("A")
	fmt.Println(s)
	s = m2.GetAt(2)
	fmt.Println(s)
	// Output: That
	// Other
}

func ExampleOrderedStringMap_MarshalJSON() {
	// You don't normally call MarshallJSON directly, but rather use the Marshall and Unmarshall json commands
	m := NewOrderedStringMap()

	m.Set("B", "This")
	m.Set("A","That")
	m.Set("C", "Other")

	s, _ := json.Marshal(m)
	os.Stdout.Write(s)

	// Note: The below output is what is produced, but isn't guaranteed. go seems to currently be sorting keys
	// Output: {"A":"That","B":"This","C":"Other"}
}

func ExampleOrderedStringMap_UnmarshalJSON() {
	b := []byte(`{"A":"That","B":"This","C":"Other"}`)
	var m OrderedStringMap

	json.Unmarshal(b, &m)
	sort.Sort(OrderStringMapByKeys(&m))

	fmt.Println(&m)

	// Output: {"A":"That","B":"This","C":"Other"}
}

func ExampleOrderedStringMap_Merge() {
	m := NewOrderedStringMap()

	m.Set("B", "This")
	m.Set("A","That")
	m.Set("C", "Other")

	m.Merge(StringMap{"D":"Last"})

	fmt.Println(m.GetAt(3))
	//Output: Last
}

func ExampleOrderedStringMap_Values() {
	m := NewOrderedStringMap()
	m.Set("B", "This")
	m.Set("A","That")
	m.Set("C", "Other")

	values := m.Values()
	fmt.Println(values)
	//Output: [This That Other]
}

func ExampleOrderedStringMap_Keys() {
	m := NewOrderedStringMap()
	m.Set("B", "This")
	m.Set("A","That")
	m.Set("C", "Other")

	values := m.Keys()
	fmt.Println(values)
	//Output: [B A C]
}

func ExampleNewOrderedStringMapFrom() {
	m := NewOrderedStringMapFrom(StringMap{"a":"this","b":"that"})
	fmt.Println(m.Get("b"))
	//Output: that
}


func ExampleOrderedStringMap_Equals() {
	m := NewOrderedStringMapFrom(StringMap{"A":"This","B":"That"})
	n := StringMap{"B":"That", "A":"This"}
	if m.Equals(n) {
		fmt.Print("Equal")
	} else {
		fmt.Print("Not Equal")
	}
	//Output: Equal
}