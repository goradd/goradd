package javascript

import "strings"

// Arguments represents a list of javascript function arguments. We can output this as javascript, or as JSON, which
// gets sent to the goradd javascript during Ajax calls and unpacked there.
// primitive types get expressed as constant values in javascript. If you want to represent the name of variable,
// us a VarName object. Function can be represented using the Function object or the Closure object, depending on whether
// you want the output of the function now, or later.

type Arguments []interface{}

// Implements the JavaScripter interface
func (a Arguments) JavaScript() string {
	var values []string
	for _, v := range a {
		values = append(values, ToJavaScript(v))
	}

	return strings.Join(values, ", ")
}
