package html

import (
	"fmt"
	"github.com/goradd/gengen/maps"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestBasicAttributes(t *testing.T) {
	cases := []struct {
		attr, val    string
		changed, err bool
	}{
		{"a", "t", true, false},
		{"b", "t", true, false},
		{"id", "9", true, false},
		{"a", "t", false, false},
		{"id", "t y", true, true},
		{"class", "t y", true, false},
		{"bad name", "t y", true, true},
	}

	a := NewAttributes()
	for _, c := range cases {
		changed, err := a.SetChanged(c.attr, c.val)
		if err != nil {
			if c.err { // expected an error
				continue
			} else {
				t.Errorf("Unexpected error on (%q, %q): %v", c.attr, c.val, err)
				continue
			}
		}

		if c.err && err == nil { // expected an error, but didn't get one
			t.Errorf("Expected error on (%q, %q)", c.attr, c.val)
			continue // no sense in checking other things, since we were expecting an error
		}

		got := a.Get(c.attr)
		if got != c.val {
			t.Errorf("Basic set/get (%q, %q)", c.attr, c.val)
		}

		if changed != c.changed {
			t.Errorf("Basic changed test (%q, %q)", c.attr, c.val)
		}
	}
}

func TestBasicStyles(t *testing.T) {
	cases := []struct {
		attr         string
		val          string
		changed, err bool
	}{
		{"color", "blue", true, false},
		{"color", "red", true, false},
		{"color", "red", false, false},
		{"height", "4px", true, false},
	}

	var changed bool
	var err error

	a := NewAttributes()
	for _, c := range cases {
		changed, err = a.SetStyleChanged(c.attr, c.val)
		if err != nil {
			if c.err { // expected an error
				continue
			} else {
				t.Errorf("Unexpected error on (%q, %q): %v", c.attr, c.val, err)
				continue
			}
		}

		if c.err && err == nil { // expected an error, but didn't get one
			t.Errorf("Expected error on (%q, %q)", c.attr, c.val)
			continue // no sense in checking other things, since we were expecting an error
		}

		got := a.GetStyle(c.attr)
		assert.Equal(t, got, c.val, "BasicStyle set/get (%q, %q), got %q", c.attr, c.val, got)

		if changed != c.changed {
			t.Errorf("BasicStyle changed test (%q, %q)", c.attr, c.val)
		}
	}

	a = NewAttributes()
	if changed, err = a.SetChanged("style", "height:4px;width:6px"); !changed || err != nil {
		t.Error("Problem setting style")
	}

	a.SetStyle("width", "+ 2")

	if changed, err = a.SetChanged("style", "width:8px; height:4px"); changed || err != nil {
		t.Errorf("Problem setting same style again: (%v, %v)", changed, err)
	}

	if a.GetStyle("width") != "8px" {
		t.Error("Problem with setting style")
	}

	a = NewAttributes()
	b := []bool{}
	changed, _ = a.SetStyleChanged("height", "9")
	b = append(b, changed)
	changed, _ = a.SetStyleChanged("height", "10")
	b = append(b, changed)
	changed, _ = a.SetStyleChanged("height", "10")
	b = append(b, changed)
	changed, _ = a.SetStyleChanged("width", "4")
	b = append(b, changed)
	out := fmt.Sprint(b)
	if out != "[true true false true]" {
		t.Error("Changing styles failed. Got: " + out)
	}

}

func TestClass(t *testing.T) {
	cases := []struct {
		val     string
		got     string
		changed bool
	}{
		{"c1", "c1", true},
		{"c2", "c2", true},
		{"c1 c2", "c1 c2", true},
		{"+ c3", "c1 c2 c3", true},
		{"+ c3", "c1 c2 c3", false},
		{"- c1", "c2 c3", true},
	}

	a := NewAttributes()
	for _, c := range cases {
		changed := a.SetClassChanged(c.val)
		/*
			if err != nil {
				if c.err {	// expected an error
					continue
				} else {
					t.Errorf("Unexpected error on (%q): %v", c.val, err)
					continue
				}
			}

			if (c.err && err == nil) { // expected an error, but didn't get one
				t.Errorf("Expected error on (%q)", c.val)
				continue // no sense in checking other things, since we were expecting an error
			}
		*/
		got := a.Class()
		if got != c.got {
			t.Errorf("Class set (%q), expected (%q), got (%q)", c.val, c.got, got)
		}

		if changed != c.changed {
			t.Errorf("Class changed test (%q)", c.val)
		}
	}

}

func TestDataAttributes(t *testing.T) {
	cases := []struct {
		attr, val    string
		changed, err bool
	}{
		{"data-a", "t", true, false},
		{"data-b", "t", true, false},
		{"data-b", "t", false, false},
		{"data-id", "9", true, false},
		{"data-$a", "t", false, true},
		{"data-bad name", "t y", true, true},
	}

	a := NewAttributes()
	for _, c := range cases {
		changed, err := a.SetChanged(c.attr, c.val)
		if err != nil {
			if c.err { // expected an error
				continue
			} else {
				t.Errorf("Unexpected error on (%q, %q): %v", c.attr, c.val, err)
				continue
			}
		}

		if c.err && err == nil { // expected an error, but didn't get one
			t.Errorf("Expected error on (%q, %q)", c.attr, c.val)
			continue // no sense in checking other things, since we were expecting an error
		}

		got := a.DataAttribute(c.attr[5:])
		if got != c.val {
			t.Errorf("Data Attribute set/get (%q, %q)", c.attr, c.val)
		}

		if changed != c.changed {
			t.Errorf("Data Attribute changed test (%q, %q)", c.attr, c.val)
		}

		if !a.HasDataAttribute(c.attr[5:]) {
			t.Errorf("Has data attribute (%q)", c.attr)
		}
	}

	a.RemoveDataAttribute("data-id")

	if a.HasDataAttribute("data-id") {
		t.Error("Removed data attribute (data-id)")
	}

}

func TestOutput(t *testing.T) {
	var s string
	a := NewAttributes()
	a.Set("class", "a")
	a.Set("id", "b")

	s = a.String()

	if !(s == `class="a" id="b"` || s == `id="b" class="a"`) {
		t.Error("No style test")
	}

	a.RemoveAttribute("class")
	a.RemoveAttribute("id")
	a.SetStyle("height", "4")

	s = a.String()

	if !(s == `style="height:4px"`) {
		t.Error("With style test: " + s)
	}

	// Test some escapes here
	a = NewAttributes()
	a.Set("placeholder", "<& I Am Here >")
	expected := "placeholder=\"&lt;&amp; I Am Here &gt;\""
	if s = a.String(); s != expected {
		t.Errorf("Not escaping. Expected (%q) got (%q)", expected, s)
	}
}

func TestOverride(t *testing.T) {
	a := NewAttributes()
	a.Set("class", "a")
	a.Set("id", "b")
	a.Set("style", "height:4px; width:3px")

	m := a.Override(maps.NewStringMapFromMap(map[string]string{"id": "c", "style": "height:7px"}))

	if m.Get("id") != "c" {
		t.Error("Error overriding id")
	}

	if m.GetStyle("height") != "7px" {
		t.Error("Error overriding height style")
	}

	if m.GetStyle("width") != "3px" {
		t.Error("Error leaving old style untouched")
	}

}

// Examples

func ExampleNewAttributesFrom() {
	a := NewAttributesFromMap(map[string]string{"id": "1", "name": "test"})
	fmt.Println(a.Get("id"))
	//Output: 1
}

func ExampleAttributes_SetAttribute() {
	a := NewAttributes()
	a.Set("class", "a")
	a.Set("id", "b")
	fmt.Println(a)
}

func ExampleAttributes_SetClass() {
	a := NewAttributes()
	a.SetClass("this")
	a.SetClass("+ that")
	s := a.Class()
	fmt.Println(s)
	// Output: this that
}

func ExampleAtributes_SetStyle() {
	a := NewAttributes()
	a.SetStyle("height", "4em")
	a.SetStyle("width", "8")
	a.SetStyle("width", "- 2")
	fmt.Println(a.GetStyle("height"))
	fmt.Println(a.GetStyle("width"))
	// Output:
	// 4em
	// 6px
}

func ExampleAttributes_Override() {
	a := NewAttributes()
	a.SetClass("this")
	a.SetStyle("height", "4em")

	b := NewAttributes()
	b.Set("class", "that")
	b.SetStyle("width", strconv.Itoa(6))

	a = a.Override(b)
	fmt.Println(a)
}

func ExampleAttributes_AddClass() {
	a := NewAttributes()
	a.AddClass("this")
	a.AddClass("that")
	fmt.Println(a)
	//Output: class="this that"
}

func ExampleAttributes_HasClass() {
	a := NewAttributes()
	if !a.HasClass("that") {
		fmt.Println("Not found")
	}
	a.SetClass("this that other")
	if a.HasClass("that") {
		fmt.Println("found")
	}
	// Output: Not found
	// found
}

func ExampleAttributes_HasStyle() {
	a := NewAttributes()
	b := []bool{}
	var found bool
	found = a.HasStyle("height")
	b = append(b, found)
	a.SetStyle("height", strconv.Itoa(10))
	found = a.HasStyle("height")
	b = append(b, found)
	fmt.Println(b)
	// Output: [false true]
}

func ExampleAttributes_RemoveStyle() {
	a := NewAttributes()
	a.SetStyle("height", "10")
	a.SetStyle("width", strconv.Itoa(5))
	a.RemoveStyle("height")
	fmt.Println(a)
	// Output: style="width:5px"
}
