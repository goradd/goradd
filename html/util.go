package html

import (
	"strings"
	html2 "html"
)

// Does a variety of transformations to make standard text presentable as HTML.
func TextToHtml(in string) (out string) {
	in = html2.EscapeString(in)
	in = strings.Replace(in, "\n\n", "<p>", -1)
	out = strings.Replace(in, "\n", "<br />", -1)
	return
}
