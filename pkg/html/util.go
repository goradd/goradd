package html

import (
	html2 "html"
	"math/rand"
	"strings"
	"time"
)

// TextToHtml does a variety of transformations to make standard text presentable as HTML.
// It escapes characters needing to be escaped, turns double-newline characters in paragrpahs, and
// single newlines into breaks.
func TextToHtml(in string) (out string) {
	in = html2.EscapeString(in)
	in = strings.Replace(in, "\n\n", "<p>", -1)
	out = strings.Replace(in, "\n", "<br />", -1)
	return
}

const htmlValueBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789-_()!"

// RandomString generates a pseudo random string of the given length
// Characters are drawn from legal HTML values that do not need encoding.
// The distribution is not perfect, so its not good for crypto, but works for general purposes
// This also works for GET variables
func RandomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = htmlValueBytes[rand.Int63()%int64(len(htmlValueBytes))]
	}
	return string(b)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
