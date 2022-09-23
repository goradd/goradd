// package url contains url utilities beyond what is available in the net/url package
package url

import (
	"fmt"
	url2 "net/url"
)

// Builder uses a builder pattern to create a URL.
type Builder struct {
	url    *url2.URL
	values url2.Values
}

// NewBuilder starts a URL builder from a basic path.
func NewBuilder(path string) *Builder {
	b := &Builder{}
	var err error
	b.url, err = url2.Parse(path)
	if err != nil {
		panic(err)
	}
	b.values = b.url.Query()
	return b
}

// NewBuilderFromUrl creates a new builder from a URL.
func NewBuilderFromUrl(u *url2.URL) *Builder {
	b := &Builder{url: u}
	b.values, _ = url2.ParseQuery(u.RawQuery)
	return b
}

// SetValue sets the GET value in the URL
func (u *Builder) SetValue(k string, v interface{}) *Builder {
	value := fmt.Sprint(v)
	u.values.Set(k, value)
	return u
}

// RemoveValue removes the value from the URL
func (u *Builder) RemoveValue(k string) *Builder {
	u.values.Del(k)
	return u
}

// SetFragment sets the fragment portion of the URL
func (u *Builder) SetFragment(f string) *Builder {
	u.url.Fragment = f
	return u
}

// ClearFragment removes the fragment from the URL.
func (u *Builder) ClearFragment() *Builder {
	u.url.Fragment = ""
	return u
}

// String returns the encoded URL.
func (u *Builder) String() string {
	u.url.RawQuery = u.values.Encode()
	return u.url.String()
}
