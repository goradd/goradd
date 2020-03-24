package javascript

import (
	"encoding/gob"
	"strings"
)

// Arguments represents a list of javascript function arguments. We can output this as javascript, or as JSON, which
// gets sent to the goradd javascript during Ajax calls and unpacked there.
// Primitive types get expressed as constant values in javascript. If you want to represent the name of variable,
// us a JsCode object. Function can be represented using the Function object or the Closure object, depending on whether
// you want the output of the function now, or later.
type Arguments []interface{}

// JavaScript returns the arguments as a comma separated list of values suitable to put in JavaScript function arguments.
func (a Arguments) JavaScript() string {
	var values []string
	for _, v := range a {
		values = append(values, ToJavaScript(v))
	}

	return strings.Join(values, ",")
}

func init() {
	// Register objects so they can be serialized
	gob.Register(Arguments{})
}
