// Package javascript converts go objects and types to javascript code, suitable for embedding in html or sending to the browser via a specialized ajax call.
package javascript

import (
	"encoding/json"
	"fmt"
	"github.com/spekary/gengen/maps"
	"strconv"
	"strings"
)

// Interface that allows an object to be converted to javascript (not JSON!)
type JavaScripter interface {
	JavaScript() string
}

const JsonObjectType = "goraddObject"

// ToJavaScript will convert the given value to javascript such that it can be embedded in a browser. If it can, it will
// use the JavaScripter interface to do the conversion. Otherwise it generally follows json encoding rules. Strings are
// escaped. Nil pointers become null objects. String maps become javascript objects. To convert a fairly complex object,
// like a map or slice of objects, convert the inner objects to interfaces
// TODO: Add some introspection to handle any kind of complex object that is a Javascripter at the inner level
func ToJavaScript(v interface{}) string {
	switch s := v.(type) {
	case JavaScripter:
		return s.JavaScript()
	case string:
		b, _ := json.Marshal(s)                         // This does a good job of handling most escape sequences we might need
		s1 := strings.Replace(string(b), "/", `\/`, -1) // Replace forward slashes to avoid potential confusion in browser from closing html tags
		return fmt.Sprintf("%v", s1)                    // will surround with quotes
	case []string:
		var values []string
		for _, item := range s {
			values = append(values, `"` + item + `"`)
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
		for k, v := range s {
			if v2, ok := v.(NoQuoteKey); ok {
				out += k + ":" + ToJavaScript(v2.Value) + ","
			} else {
				out += ToJavaScript(k) + ":" + ToJavaScript(v) + ","
			}
		}
		if len(out) == 0 {
			return "{}"
		} else {
			return "{" + out[:len(out)-1] + "}" // remove final comma and wrap in a javascript object
		}
	case map[int]interface{}:
		var out string
		for k, item := range s {
			out += fmt.Sprintf("%d:%s,", k, ToJavaScript(item))
		}
		if len(out) == 0 {
			return "{}"
		} else {
			return "{" + out[:len(out)-1] + "}" // remove final comma and wrap in a javascript object
		}

	case maps.MapI:
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
	case maps.StringMapI:
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
	case nil:
		return "null"
	default:
		return fmt.Sprint(s)
	}
}

// A value wrapper to specify a value in a map whose key should not be quoted when converting to javascript
// In some situations, a quoted key has a different meaning from a non-quoted key.
// For example, when making a list of parameters to pass when calling the jQuery $() command,
// (i.e. $j(selector, params)), quoted words are turned into parameters, and non-quoted words
// are turned into functions. For example, "size" will set the size attribute of the object, and
// size (no quotes), will call the size() function on the object.
type NoQuoteKey struct {
	Value interface{}
}

// Prevent using this as a general value.
func (n NoQuoteKey) JavaScript() string {
	panic("NoQuoteKey should only be used as a value in a string map.")
}

type jsCode struct {
	code  string
	isInt bool
}

// JsCode represents straight javascript code that should not be escaped or quoted.
// Normally, string values would be quoted. This outputs a string without quoting or escaping.
func JsCode(s string) jsCode {
	return jsCode{code: s}
}

func (c jsCode) JavaScript() string {
	return c.code
}

// Undefined explicitly outputs as "undefined" in javascript. Generally, nil pointers become "null" in javascript, so
// use this if you would rather have an undefined value.
type Undefined struct {
}

func (n Undefined) JavaScript() string {
	return "undefined"
}

func (n Undefined) MarshalJSON() ([]byte, error) {
	return []byte("undefined"), nil
}

// NumberInt is a helper function to convert an expected integer that is returned from a json Unmarshal as a Number,
// into an actual integer without returning any errors. If there is an error, it just returns 0. Use this when you absolutely
// know you are expecting an integer. Can convert strings too.
func NumberInt(i interface{}) int {
	switch n := i.(type) {
	case json.Number:
		v, _ := n.Int64()
		return int(v)
	case string:
		v, _ := strconv.Atoi(n)
		return v
	}
	return 0
}

// NumberFloat is a helper function to convert an expected float that is returned from a json Unmarshal as a Number,
// into an actual float64 without returning any errors. If there is an error, it just returns 0. Use this when you absolutely
// know you are expecting a float. Can convert strings too.
func NumberFloat(i interface{}) float64 {
	switch n := i.(type) {
	case json.Number:
		v, _ := n.Float64()
		return v
	case string:
		v, _ := strconv.ParseFloat(n, 64)
		return v
	}
	return 0
}

// NumberString is a helper function to convert a value that might get cast as a Json Number into a string.
// If there is an error, it just returns 0. Use this when you absolutely
// know you are expecting a string.
func NumberString(i interface{}) string {
	switch n := i.(type) {
	case json.Number:
		v := n.String()
		return v
	case string:
		return n
	}
	panic("Unknown type for NumberString")
	return ""
}

