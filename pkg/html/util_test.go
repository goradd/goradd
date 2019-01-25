package html

import (
	"fmt"
	"testing"
)

func TestRandomHtmlValueString(t *testing.T) {
	//rand.Seed(1) // reset random seed

	s := RandomString(40)
	fmt.Printf(s + " ")

	if len(s) != 40 {
		t.Error("Wrong size")
	}
}

func ExampleTextToHtml() {
	s := TextToHtml("This is a & test.\n\nA paragraph\nwith a forced break.")
	fmt.Println(s)
	// Output: This is a &amp; test.<p>A paragraph<br />with a forced break.
}