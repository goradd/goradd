// Package javascript converts go objects and types to javascript code, suitable for embedding in html or sending to the browser via a specialized ajax call.
package javascript

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Interface that allows an object to be converted to javascript (not JSON!)
type JavaScripter interface {
	JavaScript() string
}

const JsonObjectType = "goraddObject"

// ToJavaScript will convert the given value to javascript such that it can be embedded in a browser. If it can, it will
// use the JavaScripter interface to do the conversion. Otherwise it generally follows json encoding rules. Strings are
// escaped. Nil pointers become null objects. String maps become javascript objects.
func ToJavaScript(v interface{}) string {
	switch s := v.(type) {
	case JavaScripter:
		return s.JavaScript()
	case VarName:
		return string(s)	// we are going to assume variable names are not going to containing anything needing to be escaped
	case json.Marshaler:
		b, _ := json.Marshal(s)
		s1 := strings.Replace(string(b), "/", `\/`, -1) // Replace forward slashes to avoid potential confusion in browser from closing html tags
		return fmt.Sprintf("%s", s1)
	case string:
		b, _ := json.Marshal(s)                         // This does a good job of handling most escape sequences we might need
		s1 := strings.Replace(string(b), "/", `\/`, -1) // Replace forward slashes to avoid potential confusion in browser from closing html tags
		return fmt.Sprintf("%v", s1) // will surround with quotes
	case []interface{}:
		a := Arguments(s)
		return "[" + a.JavaScript() + "]"
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

// VarName represents a global variable name that can be used as a value in a javascript argument list. Normally,
// string values would be quoted. This outputs a string without a quote, which in javascript will be treated as
// a variable or function name.
type VarName string


// Undefined explicitly outputs as "undefined" in javascript. Generally, nil pointers become "null" in javascript, so
// use this if you would rather have an undefined value.
type Undefined struct {
}

func (n Undefined) JavaScript() string {
	return "undefined"
}