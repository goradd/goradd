package html

import (
	"fmt"
	"testing"
)

func ExampleStyle_SetTo() {
	s := NewStyle()
	s.SetTo("height: 9em; width: 100%; position:absolute")
	fmt.Print(s)
	//Output: height:9em;position:absolute;width:100%
}

func ExampleStyle_Set_a() {
	s := NewStyle()
	s.Set("height", "9")
	fmt.Print(s)
	//Output: height:9px
}

func ExampleStyle_Get() {
	s := NewStyle()
	s.SetTo("height: 9em; width: 100%; position:absolute")
	fmt.Print(s.Get("width"))
	//Output: 100%
}

func ExampleStyle_Delete() {
	s := NewStyle()
	s.SetTo("height: 9em; width: 100%; position:absolute")
	s.Delete("position")
	fmt.Print(s)
	//Output: height:9em;width:100%
}

func ExampleStyle_RemoveAll() {
	s := NewStyle()
	s.SetTo("height: 9em; width: 100%; position:absolute")
	s.RemoveAll()
	fmt.Print(s)
	//Output:
}

func ExampleStyle_Has() {
	s := NewStyle()
	s.SetTo("height: 9em; width: 100%; position:absolute")
	found := s.Has("display")
	fmt.Print(found)
	//Output:false
}

func ExampleStyle_Set_b() {
	s := NewStyle()
	s.SetTo("height:9px")
	s.Set("height", "+ 10")
	fmt.Print(s)
	//Output: height:19px
}

func TestStyleSet(t *testing.T) {
	s := NewStyle()

	changed, err := s.SetChanged("height", "4")

	if !changed {
		t.Error("Expected a change")
	}
	if err != nil {
		t.Error(err)
	}

	s.RemoveAll()
	if s.Has("height") {
		t.Error("Expected no height")
	}

	s.Set("height", "4")

	changed, err = s.SetTo("height: 3; width: 5")

	if !changed {
		t.Error("Expected a change")
	}
	if err != nil {
		t.Error(err)
	}
	v := s.Get("width")
	if v != "5px" {
		t.Error("Expect a width of 5px, got " + v)
	}
	if s.Get("height") != "3px" {
		t.Error("Expect a height of 3px")
	}

}

func TestStyleLengths(t *testing.T) {
	s := NewStyle()

	s.Set("height", "4px")
	changed, err := s.SetChanged("height", "4")

	if changed {
		t.Error("Expected no change")
	}
	if err != nil {
		t.Error(err)
	}

	changed, err = s.SetChanged("height", "4em")
	if !changed {
		t.Error("Expected change")
	}
	if err != nil {
		t.Error(err)
	}

	changed, err = s.SetChanged("width", "0")
	if !changed {
		t.Error("Expected change")
	}
	if err != nil {
		t.Error(err)
	}

	if w := s.Get("width"); w != "0" {
		t.Error("Expected a 0")
	}

	changed, err = s.SetChanged("width", "1")
	if w := s.Get("width"); w != "1px" {
		t.Error("Expected a 1px")
	}

	// test a non-length numeric
	changed, err = s.SetChanged("volume", "4")
	if w := s.Get("volume"); w != "4" {
		t.Error("Expected a 4")
	}
}

func TestStyleMath(t *testing.T) {
	s := NewStyle()

	s.Set("height", "4em")
	s.Set("height", "* 2")
	if h := s.Get("height"); h != "8em" {
		t.Error("Expected 8em, got " + h)
	}

	s.Set("height", "2em 9px")
	s.Set("height", "/ 2")
	if h := s.Get("height"); h != "1em 4.5px" {
		t.Error("Expected 1em 4.5px, got " + h)
	}

	s.Set("width", "7.6in")
	s.Set("width", "+ 2")
	if h := s.Get("width"); h != "9.6in" {
		t.Error("Expected 9.6in, got " + h)
	}

	s.Set("width", "1.6in")
	s.Set("width", "- 2")
	if h := s.Get("width"); h != "-0.4in" { // this test in particular can produce a rounding error if not handled carefully
		t.Error("Expected -0.4in, got " + h)
	}
}

func TestStyle(t *testing.T) {
	s := NewStyle()

	s.Set("position", "9")

	if a := s.Get("position"); a != "9px" {
		t.Error("Style test failed: " + a)
	}
}
