package types

import (
	"testing"
	"fmt"
	"bytes"
	"encoding/gob"
	"encoding/json"
	"os"
	"sort"
)


func TestSafeStringMap(t *testing.T) {
	var s string

	m := NewSafeStringMap()

	m.Set("B", "This")
	m.Set("A","That")
	m.Set("C", "Other")

	if s = m.Get("B"); s != "This" {
		t.Errorf("Strings test failed. Expected  (%q) got (%q).", "This", s)
	}

	if s = m.Get("C"); s != "Other" {
		t.Errorf("Strings test failed. Expected  (%q) got (%q).", "Other", s)
	}

	m.Remove("A")

	if m.Len() != 2 {
		t.Error("Len Failed.")
	}


	if  m.Has ("NOT THERE") {
		t.Error("Getting non-existant value did not return false")
	}

	s = m.Get("B")
	if s != "This" {
		t.Error("Get failed")
	}

	if !m.Has("B") {
		t.Error ("Existance test failed.")
	}


	// Verify it satisfies the StringMapI interface
	var i StringMapI = m
	if s := i.Get("B"); s != "This" {
		t.Error("StringMapI interface test failed.")
	}
}

func TestSafeStringMapChange(t *testing.T) {
	m := NewSafeStringMap()

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

func ExampleOrderedSafeStringMap_Set() {
	m := NewSafeStringMap()
	m.Set("a", "Here")
	fmt.Println(m.Get("a"))
	// Output Here
}


func ExampleSafeStringMap_Range() {
	m := NewSafeStringMap()
	a :=[]string{}

	m.Set("B", "This")
	m.Set("A","That")
	m.Set("C", "Other")

	m.Range(func (key string, val string) bool {
		a = append(a, val)
		return true	// keep iterating to the end
	})
	fmt.Println()

	sort.Sort(sort.StringSlice(a))
	fmt.Println(a)
	//Output: [Other That This]
}

func ExampleSafeStringMap_MarshalBinary() {
	// You would rarely call MarshallBinary directly, but rather would use an encoder, like GOB for binary encoding

	m := NewSafeStringMap()
	var m2 SafeStringMap

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
	//Output: That
}

func ExampleSafeStringMap_MarshalJSON() {
	// You don't normally call MarshallJSON directly, but rather use the Marshall and Unmarshall json commands
	m := NewSafeStringMap()

	m.Set("B", "This")
	m.Set("A","That")
	m.Set("C", "Other")

	s, _ := json.Marshal(m)
	os.Stdout.Write(s)

	// Note: The below output is what is produced, but isn't guaranteed. go seems to currently be sorting keys
	// Output: {"A":"That","B":"This","C":"Other"}
}

func ExampleSafeStringMap_UnmarshalJSON() {
	b := []byte(`{"A":"That","B":"This","C":"Other"}`)
	var m SafeStringMap

	json.Unmarshal(b, &m)

	fmt.Println(m.items)

	// Note: The below output is what is produced, but isn't guaranteed. go seems to currently be sorting keys
	// Output: {"A":"That","B":"This","C":"Other"}
}

func ExampleSafeStringMap_Merge() {
	m := NewSafeStringMap()

	m.Set("B", "This")
	m.Set("A","That")
	m.Set("C", "Other")

	m.Merge(StringMap{"D":"Last"})

	fmt.Println(m.Get("D"))
	//Output: Last
}

func ExampleSafeStringMap_Values() {
	m := NewSafeStringMap()
	m.Set("B", "This")
	m.Set("A","That")
	m.Set("C", "Other")

	values := m.Values();
	sort.Sort(sort.StringSlice(values))
	fmt.Println(values)
	//Output: [Other That This]
}

func ExampleSafeStringMap_Keys() {
	m := NewSafeStringMap()
	m.Set("B", "This")
	m.Set("A","That")
	m.Set("C", "Other")

	values := m.Keys();
	sort.Sort(sort.StringSlice(values))
	fmt.Println(values)
	//Output: [A B C]
}

func ExampleNewSafeStringMapFrom() {
	m := NewSafeStringMapFrom(StringMap{"a":"this","b":"that"})
	fmt.Println(m.Get("b"))
	//Output: that
}

func ExampleSafeStringMap_Equals() {
	m := NewSafeStringMapFrom(StringMap{"A":"This","B":"That"})
	n := StringMap{"B":"That", "A":"This"}
	if m.Equals(n) {
		fmt.Print("Equal")
	} else {
		fmt.Print("Not Equal")
	}
	//Output: Equal
}