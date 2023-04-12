package reflect

import (
	"fmt"
	"reflect"
	"strings"
)

// InterfaceSlice converts a slice of any object to a slice of interfaces of those internal objects
// if in is not an addressable item, it will panic
func InterfaceSlice(in any) (o []any) {
	if in == nil {
		return
	}
	v := reflect.ValueOf(in)
	for i := 0; i < v.Len(); i++ {
		o = append(o, v.Index(i).Interface())
	}
	return o
}

func IsSlice(in any) bool {
	return reflect.TypeOf(in).Kind() == reflect.Slice
}

// JoinStringers will join a slice of any kind of objects that are Stringers.
// If the passed in object is not addressable, it will panic.
// Any objects that are not stringers will result in a missing item in the result.
func JoinStringers(in any, sep string) (out string) {
	if in == nil {
		return
	}
	v := reflect.ValueOf(in)
	var o []string
	for i := 0; i < v.Len(); i++ {
		if s, ok := v.Index(i).Interface().(fmt.Stringer); ok {
			o = append(o, s.String())
		}
	}
	return strings.Join(o, sep)
}
