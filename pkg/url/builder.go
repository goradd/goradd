// package url contains url utilities beyond what is available in the net/url package
package url

import (
	"fmt"
	"net/url"
)

// Builder uses a builder pattern to create a URL.
// Call NewBuilder to start off a URL, then add values to the URL using AddValue.
type Builder struct {
	url    url.URL
	values url.Values
}

// NewBuilder starts a URL builder
func NewBuilder(path string) *Builder {
	b := &Builder{values: make(url.Values)}
	b.url.Path = path
	return b
}

// AddValue adds a GET value to the URL
func (u *Builder) AddValue(k string, v interface{}) *Builder {
	value := fmt.Sprintf("%v", v)
	u.values.Add(k, value)
	return u
}

// String returns the encoded URL.
func (u *Builder) String() string {
	u.url.RawQuery = u.values.Encode()
	return u.url.String()
}