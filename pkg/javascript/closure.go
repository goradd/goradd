package javascript

import (
	"encoding/gob"
	"encoding/json"
	"strings"
)

// NewClosure creates a new Closure object.
func NewClosure(body string, args ...string) Closure {
	return Closure{body, args}
}

// Closure represents a javascript function pointer that can be called by javascript at a later time.
type Closure struct {
	// Body is the body javascript of the Closure
	Body string
	// Args are the names of the arguments in the argument list of the Closure
	Args []string
}

// JavaScript implements the JavaScripter interface and returns the Closure as javascript code.
func (c Closure) JavaScript() string {
	var args string

	if c.Args != nil {
		args = strings.Join(c.Args, ", ")
	}

	return "function(" + args + ") {" + c.Body + "}"
}

// MarshalJSON implements the json.Marshaller interface.
// The output of this is designed to be unpacked by the goradd.js javascript file during Ajax calls.
func (c Closure) MarshalJSON() (buf []byte, err error) {
	var obj = map[string]interface{}{}

	obj[JsonObjectType] = "closure"
	obj["func"] = c.Body
	obj["params"] = c.Args

	buf, err = json.Marshal(obj)
	return
}

// NewClosureCall creates a new ClosureCall.
func NewClosureCall(body string, context string, args ...string) ClosureCall {
	return ClosureCall{body, context, args}
}

// ClosureCall represents the result of a javascript Closure that is called immediately.
// context will become the "this" variable inside the Closure when called.
type ClosureCall struct {
	// Body is the body javascript of the Closure
	Body string
	// Context is what will become the "this" var inside the Closure when called. Specifying "this" will bring the "this" from the outer context in to the Closure.
	Context string
	// Args are the names of the arguments in the argument list of the Closure
	Args []string
}

// JavaScript implements the JavaScripter interface and returns the Closure as javascript code.
func (c ClosureCall) JavaScript() string {
	var args string

	if c.Args != nil {
		args = strings.Join(c.Args, ", ")
	}

	return "(function(" + args + ") {" + c.Body + "}).call(" + c.Context + ")"
}

// MarshalJSON implements the json.Marshaller interface.
// The output of this is designed to be unpacked by the goradd.js javascript file.
func (c ClosureCall) MarshalJSON() (buf []byte, err error) {
	var obj = map[string]interface{}{}

	obj[JsonObjectType] = "closure"
	obj["func"] = c.Body
	obj["params"] = c.Args
	obj["call"] = c.Context

	buf, err = json.Marshal(obj)
	return
}

func init() {
	// Register objects so they can be serialized
	gob.Register(Closure{})
	gob.Register(ClosureCall{})
}
