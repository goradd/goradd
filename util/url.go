package util

import (
	"net/url"
	"fmt"
	"strconv"
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

// Additions to getting values out of a url.Values value. Cast to a UrlValues type, and then call one of the functions.
// Ex: v,err := UrlValues(r.Form).GetInt("name")

type UrlValues url.Values

func (values UrlValues) GetString(name string) (value string, err error) {
	if v,ok := values[name]; !ok {
		err = fmt.Errorf("%s was not found", name)
		return
	} else {
		return v[0],nil
	}
}

func (values UrlValues) GetInt(name string) (value int, err error) {
	if v,ok := values[name]; !ok {
		err = fmt.Errorf("%s was not found", name)
	} else {
		value,err = strconv.Atoi(v[0])
	}
	return
}

func (values UrlValues) GetFloat64(name string) (value float64, err error) {
	if v,ok := values[name]; !ok {
		err = fmt.Errorf("%s was not found", name)
	} else {
		value,err = strconv.ParseFloat(v[0], 64)
	}
	return
}

// GetStrings returns the string array associated with the name
func (values UrlValues) GetStrings(name string) (value []string, err error) {
	if v,ok := values[name]; !ok {
		err = fmt.Errorf("%s was not found", name)
		return
	} else {
		return v,nil
	}
}
