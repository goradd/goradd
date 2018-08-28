package util

import (
	"net/url"
	"fmt"
)

type UrlBuilder struct {
	url url.URL
	values url.Values
}

func NewUrlBuilder(path string) *UrlBuilder {
	b := &UrlBuilder{values:make(url.Values)}
	b.url.Path = path
	return b
}

func (u *UrlBuilder) AddValue (k string,v interface{}) *UrlBuilder {
	value := fmt.Sprintf("%v", v)
	u.values.Add(k,value)
	return u
}

func (u *UrlBuilder) String () string {
	u.url.RawQuery = u.values.Encode()
	return u.url.String()
}
