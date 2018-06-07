package html

import (
	"fmt"
	"github.com/spekary/goradd/util/types"
	"testing"
)

func TestTag(t *testing.T) {
	attr := NewAttributes()

	attr.Set("test", "test1")

	s := RenderTag("p", attr, "testInner")

	if s != "<p test=\"test1\">\ntestInner\n</p>\n" {
		t.Error("Expected <p test=\"test1\">\ntestInner\n</p>, got " + s)
	}
}

func ExampleRenderTag() {
	fmt.Println(RenderTagNoSpace("div", NewAttributesFrom(types.StringMap{"id": "me", "name": "you"}), "Here I am"))
	//Output:<div id="me" name="you">Here I am</div>
}

func ExampleRenderVoidTag() {
	fmt.Println(RenderVoidTag("img", NewAttributesFrom(types.StringMap{"src": "thisFile"})))
	// Output: <img src="thisFile" />
}
