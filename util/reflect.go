package util

import (
	"reflect"
)

func IsSlice(s interface{}) bool {
	return s != nil && reflect.TypeOf(s).Kind() == reflect.Slice
}

// Given an interface that represents a slice of arbitrary values, GetSlice returns a slice of null interfaces to those values
func GetSlice(i interface{}) []interface{} {
	var ret []interface{}
	s := reflect.ValueOf(i)
	for i := 0; i < s.Len(); i++ {
		ret = append(ret, s.Index(i).Interface())
	}
	return ret
}
