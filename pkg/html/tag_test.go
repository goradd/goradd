package html

import (
	"fmt"
	"testing"
)

func TestTag(t *testing.T) {
	attr := NewAttributes()

	attr.Set("test", "test1")

	s := RenderTag("p", attr, "testInner")
	expected := "<p test=\"test1\">\n  testInner\n</p>\n"
	if s != expected {
		t.Error("Expected " + expected + ", got " + s)
	}
}

func ExampleRenderTag() {
	fmt.Println(RenderTagNoSpace("div", NewAttributesFromMap(map[string]string{"id": "me", "name": "you"}), "Here I am"))
	//Output:<div id="me" name="you">Here I am</div>
}

func ExampleRenderVoidTag() {
	fmt.Println(RenderVoidTag("img", NewAttributesFromMap(map[string]string{"src": "thisFile"})))
	// Output: <img src="thisFile" />
}
