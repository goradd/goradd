package generator

import (
	"fmt"
	"github.com/goradd/goradd/pkg/strings"
	"reflect"
)

// Exports the given creator so that it can be embedded in a go file.
// Empty items are not exported
// Not all creators can be exported. This function is mainly a helper for code generation of controls.
// Specifically, events and actions do not export cleanly currently. But that should not be a problem for code generation.
func ExportCreator(creator interface{}) string {
	v := reflect.ValueOf(creator)
	t := reflect.TypeOf(creator)

	if t.Kind() != reflect.Struct {
		panic("a creator must be a struct type")
	}
	var s string
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		val := v.Field(i)
		if isZero(val) {
			continue
		}

		switch field.Type.Kind() {
		case reflect.Chan: panic("creator cannot have a channel, since channels are not serializable")
		case reflect.Func: panic("creator cannot have a function, since functions are not serializable. Try an interface instead.")
		case reflect.Struct:
			if f := val.MethodByName("Export"); f.IsValid() {
				result := f.Call(nil)
				s += field.Name + ":" + result[0].String() + ",\n"
			} else {
				s += field.Name + ":" + ExportCreator(val.Interface()) + ",\n"
			}
		case reflect.Interface:
			if f := reflect.ValueOf(val.Interface()).MethodByName("Export"); f.IsValid() {
				result := f.Call(nil)
				s += field.Name + ":" + result[0].String() + ",\n"
			}
		case reflect.Array:fallthrough
		case reflect.Slice:
			s += field.Name + ":" + exportSlice(val) + ",\n"
		default:
			s += field.Name + ":" + fmt.Sprintf("%#v", val.Interface()) + ",\n"
		}
	}
	s = strings.Indent(s)
	s = t.String() + "{\n" + s + "}"
	return s
}

func exportSlice(slice reflect.Value) string {
	var s = slice.Type().String() + "{\n"
	for i := 0; i < slice.Len(); i++ {
		val := slice.Index(i)

		switch val.Type().Kind() {
		case reflect.Chan: panic("creator cannot have a channel, since channels are not serializable")
		case reflect.Func: panic("creator cannot have a function, since functions are not serializable. Try an interface instead.")
		case reflect.Struct:
			s += ExportCreator(val.Interface()) + ",\n"
		case reflect.Array:fallthrough
		case reflect.Slice:
			s += exportSlice(val)

		default:
			s += fmt.Sprintf("%#v", val.Interface()) + ",\n"
		}
	}
	s += "}\n"
	return s
}

func isZero(v reflect.Value) bool {
	return !v.IsValid() || reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface())
}

