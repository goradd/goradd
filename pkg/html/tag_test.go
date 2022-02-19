package html

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func ExampleVoidTag_Render() {
	v := VoidTag{"br", Attributes{"id":"hi"}}
	fmt.Println(v.Render())
	//Output: <br id="hi" />
}


func ExampleRenderTag() {
	fmt.Println(RenderTagNoSpace("div", NewAttributesFrom(map[string]string{"id": "me", "name": "you"}), "Here I am"))
	// Output: <div id="me" name="you">Here I am</div>
}

func ExampleRenderVoidTag() {
	fmt.Println(RenderVoidTag("img", NewAttributesFrom(map[string]string{"src": "thisFile"})))
	// Output: <img src="thisFile" />
}

func ExampleRenderLabel() {
	s1 := RenderLabel(nil, "Title", "<input />", LabelBefore)
	s2 := RenderLabel(nil, "Title", "<input />", LabelAfter)
	s3 := RenderLabel(nil, "Title", "<input />", LabelWrapBefore)
	s4 := RenderLabel(nil, "Title", "<input />", LabelWrapAfter)
	fmt.Println(s1)
	fmt.Println(s2)
	fmt.Println(s3)
	fmt.Println(s4)
	// Output: <label>Title</label> <input />
	// <input /> <label>Title</label>
	// <label>
	//   Title <input />
	// </label>
	//
	// <label>
	//   <input /> Title
	// </label>
}

func TestTag(t *testing.T) {
	attr := NewAttributes()

	attr.Set("test", "test1")

	s := RenderTag("p", attr, "testInner")
	expected := "<p test=\"test1\">\n  testInner\n</p>\n"
	if s != expected {
		t.Error("Expected " + expected + ", got " + s)
	}
}

func TestRenderTagNoSpace(t *testing.T) {
	type args struct {
		tag       string
		attr      Attributes
		innerHtml string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"Test empty", args{"p", nil, ""}, `<p></p>`},
		{"Test empty with attributes", args{"p", Attributes{"height":"10"}, ""}, `<p height="10"></p>`},
		{"Test text", args{"p", nil, "I am here"}, `<p>I am here</p>`},
		{"Test html", args{"p", nil, "<p>I am here</p>"}, `<p>
  <p>I am here</p>
</p>`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RenderTagNoSpace(tt.args.tag, tt.args.attr, tt.args.innerHtml); got != tt.want {
				t.Errorf("RenderTagNoSpace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIndent(t *testing.T) {
	s :=
`a
  b
    c`
	s1:=Indent(s)
	assert.Equal(t,
`  a
    b
      c`, s1)

	s =
`<textarea height="10">a
  b
    c</textarea>`
	assert.Equal(t, s, Indent(s))

	// check for html error
	s = `<textarea height="10">a
  b
    c`
	assert.Equal(t, s, Indent(s))

}

func ExampleRenderImage() {
	s := RenderImage("http://abc.com/img.jpg", "my image", Attributes{"height":"10", "width":"20"})
	fmt.Print(s)
	//Output: <img src="http://abc.com/img.jpg" alt="my image" width="20" height="10" />
}

func ExampleComment() {
	s := Comment("This is a test")
	fmt.Print(s)
	//Output: <!-- This is a test -->
}