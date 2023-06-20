package any

import (
	"reflect"
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
