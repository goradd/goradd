// Package javascript converts go objects and types to javascript code, suitable for embedding in html or sending to the
// browser via a specialized ajax call.
package javascript

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	json2 "github.com/goradd/goradd/pkg/json"
	"sort"
	"strings"
	"time"

	"github.com/goradd/maps"
)

// JavaScripter specifies that an object can be converted to javascript (not JSON!). These objects should also be
// gob encodable and registered with gob, since they might be embedded in a control and need to be serialized.
type JavaScripter interface {
	JavaScript() string
}

// JsonObjectType is used by the ajax processor in goradd.js to indicate that we are sending a special kind of
// object to the browser. These are things like dates, closures, etc. that are not easily represented by JSON.
const JsonObjectType = "goraddObject"

// ToJavaScript will convert the given value to javascript such that it can be embedded in a browser. If it can, it will
// use the JavaScripter interface to do the conversion. Otherwise, it generally follows json encoding rules. Strings are
// escaped. Nil pointers become null objects. String maps become javascript objects. To convert a fairly complex object,
// like a map or slice of objects, convert the inner objects to interfaces
func ToJavaScript(v interface{}) string {
	// TODO: Add some introspection to handle any kind of complex object that is a JavaScripter at the inner level
	// TODO: This is a good place for Go 1.18 templates to expand on arrays of items.
	switch s := v.(type) {
	case JavaScripter:
		return s.JavaScript()
	case string:
		// Note that we cannot use template literals here (backticks) because not all browsers support them
		b, _ := json.Marshal(s)      // This does a good job of handling most escape sequences we might need
		return fmt.Sprint(string(b)) // will surround with quotes
	case []string:
		var values []string
		for _, item := range s {
			values = append(values, ToJavaScript(item))
		}
		return "[" + strings.Join(values, ",") + "]"
	case []interface{}:
		var values []string
		for _, item := range s {
			values = append(values, ToJavaScript(item))
		}
		return "[" + strings.Join(values, ",") + "]"
	case map[string]interface{}:
		var out string
		// For testing and consistency, we always return maps in order sorted by key if possible
		keys := make([]string, 0, len(s))
		for k := range s {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			v2 := s[k]
			if v3, ok := v2.(NoQuoteKey); ok {
				out += k + ":" + ToJavaScript(v3.Value) + ","
			} else {
				out += ToJavaScript(k) + ":" + ToJavaScript(v2) + ","
			}
		}
		if len(out) == 0 {
			return "{}"
		} else {
			return "{" + out[:len(out)-1] + "}" // remove final comma and wrap in a javascript object
		}
	case map[int]interface{}:
		var out string
		// For testing and consistency, we always return maps in order sorted by key
		keys := make([]int, 0, len(s))
		for k := range s {
			keys = append(keys, k)
		}
		sort.Ints(keys)
		for _, k := range keys {
			out += fmt.Sprintf("%d:%s,", k, ToJavaScript(s[k]))
		}
		if len(out) == 0 {
			return "{}"
		} else {
			return "{" + out[:len(out)-1] + "}" // remove final comma and wrap in a javascript object
		}

	case maps.MapI[string, any]:
		var out string
		s.Range(func(k string, v interface{}) bool {
			if v2, ok := v.(NoQuoteKey); ok {
				out += k + ":" + ToJavaScript(v2.Value) + ","
			} else {
				out += ToJavaScript(k) + ":" + ToJavaScript(v) + ","
			}
			return true
		})
		if len(out) == 0 {
			return "{}"
		} else {
			return "{" + out[:len(out)-1] + "}" // remove final comma and wrap in a javascript object
		}
	case maps.MapI[string, string]:
		var out string
		s.Range(func(k string, v string) bool {
			out += ToJavaScript(k) + ":" + ToJavaScript(v) + ","
			return true
		})
		if len(out) == 0 {
			return "{}"
		} else {
			return "{" + out[:len(out)-1] + "}" // remove final comma and wrap in a javascript object
		}
	case time.Time:
		return fmt.Sprintf(`new Date(%d)`, s.UnixMilli())
	case nil:
		return "null"
	default:
		/** TBD
		rv := reflect.ValueOf(v)
		if rv.Kind() == reflect.StdMap {

		}*/
		return fmt.Sprint(s)
	}
}

// NoQuoteKey is a value wrapper to specify a value in a map whose key should not be quoted when converting to javascript.
// In some situations, a quoted key has a different meaning from a non-quoted key.
// For example, when making a list of parameters to pass when calling a javascript command,
// quoted words are turned into parameters, and non-quoted words
// are turned into functions. For example, "size" will set the size attribute of the object, and
// size (no quotes), will call the size() function on the object.
// i.e. map[string]string {"size":4, "size":NoQuoteKey{JsCode{"obj"}}}
type NoQuoteKey struct {
	Value interface{}
}

// JsCode represents straight JavaScript code that should not be escaped or quoted.
// Normally, string values would be quoted. This outputs a string without quoting or escaping.
type JsCode string

// JavaScript implements the JavaScripter interface.
func (c JsCode) JavaScript() string {
	return string(c)
}

// Undefined explicitly outputs as "undefined" in javascript. Generally, nil pointers become "null" in javascript, so
// use this if you would rather have an undefined value.
type Undefined struct {
}

// JavaScript implements the JavaScripter interface
func (n Undefined) JavaScript() string {
	return "undefined"
}

/*
If this is needed, it will need to be coordinated with goradd.js and unpacked there. JSON will not accept
an undefined value.
func (n Undefined) MarshalJSON() ([]byte, error) {
	return []byte(`"undefined"`), nil
}
*/

// NumberInt is a helper function to convert an expected integer that is returned from a json Unmarshal as a Number,
// into an actual integer without returning any errors. If there is an error, it just returns 0. Use this when you absolutely
// know you are expecting an integer. Can convert strings too.
// Deprecated: Use json.NumberInt
func NumberInt(i interface{}) int {
	return json2.NumberInt(i)
}

// NumberFloat is a helper function to convert an expected float that is returned from a json Unmarshal as a Number,
// into an actual float64 without returning any errors. If there is an error, it just returns 0. Use this when you absolutely
// know you are expecting a float. Can convert strings too.
// Deprecated: Use json.NumberFLoat
func NumberFloat(i interface{}) float64 {
	return json2.NumberFloat(i)
}

// NumberString is a helper function to convert a value that might get cast as a Json Number into a string.
// If there is an error, it just returns 0. Use this when you absolutely
// know you are expecting a string.
//Deprecated: Use json.NumberString
func NumberString(i interface{}) string {
	return json2.NumberString(i)
}

func init() {
	// Register objects so they can be serialized
	gob.Register(NoQuoteKey{})
	gob.Register(JsCode(""))
	gob.Register(Undefined{})
}
