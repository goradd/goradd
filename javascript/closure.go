package javascript

import (
	"strings"
	"encoding/json"
)

// Closure represents a javascript function pointer that can be called by javascript at a later time.
type Closure struct {
	// Body is the body javascript of the closure
	Body string
	// Args are the names of the arguments in the argument list of the closure
	Args []string
}

// JavaScript implements the JavsScripter interface and returns the closure as javascript code.
func (c Closure) JavaScript() string {
	var args string

	if c.Args != nil {
		args = strings.Join(c.Args, ", ")
	}

	return "function(" + args + ") {" + c.Body + "}"
}

// MarshalJSON implements the json.Marshaller interface. The output of this is designed to be unpacked by the goradd javascript file.
func (c Closure) MarshalJSON() (buf []byte, err error) {
	var obj = map[string]interface{}{}

	obj[JsonObjectType] = "closure"
	obj["func"] = c.Body
	obj["params"] = c.Args

	buf,err = json.Marshal(obj)
	return
}

// ClosureCall represents the result of a javascript closure that is called immediately.
type ClosureCall struct {
	// Body is the body javascript of the closure
	Body string
	// Args are the names of the arguments in the argument list of the closure
	Args []string
	// Context is what will become the "this" argument inside of the closure when called. Specifying "this" will bring the "this" from the outer context in to the closure.
	Context string
}

// JavaScript implements the JavsScripter interface and returns the closure as javascript code.
func (c ClosureCall) JavaScript() string {
	var args string

	if c.Args != nil {
		args = strings.Join(c.Args, ", ")
	}

	return "(function(" + args + ") {" + c.Body + "}).call(" + c.Context + ")"
}

// Implements the json.Marshaller interface. The output of this is designed to be unpacked by the goradd javascript file.
func (c ClosureCall) MarshalJSON() (buf []byte, err error) {
	var obj = map[string]interface{}{}

	obj[JsonObjectType] = "closure"
	obj["func"] = c.Body
	obj["params"] = c.Args
	obj["call"] = c.Context

	buf,err = json.Marshal(obj)
	return
}
