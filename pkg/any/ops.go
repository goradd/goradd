// Package any has general purpose utility functions for working with interfaces and generic types.
package any

import "reflect"

// If returns the first item if cond is true, or the second item if it is false.
func If[T any](cond bool, i1, i2 T) T {
	if cond {
		return i1
	} else {
		return i2
	}
}

// Zero returns the zero value of a type.
func Zero[T any]() T {
	var v T
	return v
}

// IsNil is a safe test for nil for any kind of variable, and will not panic
// If i points to a nil object, IsNil will return true, as opposed to i==nil which will return false
func IsNil(i any) bool {
	if i == nil {
		return true
	}
	v := reflect.ValueOf(i)
	k := v.Kind()
	switch k {
	case reflect.Chan:
		fallthrough
	case reflect.Func:
		fallthrough
	case reflect.Interface:
		fallthrough
	case reflect.Map:
		fallthrough
	case reflect.Ptr:
		fallthrough
	case reflect.Slice:
		return v.IsNil()
	}
	return false
}
