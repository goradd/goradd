package javascript

import (
	"encoding/gob"
	"encoding/json"
)

// Function represents the result of a function call to a global function or function in an object referenced from
// global space. The purpose
// of this is to immediately use the results of the function call, as opposed to a Closure, which stores a pointer
// to a function that is used later.
// context will become the "this" value inside the closure
// args will be passed as values, and strings will be quoted. To pass a variable name, wrap the name with a JsCode call.
func Function(name string, context string, args... interface{}) Ωfunction {
	return Ωfunction{name, context, args}
}

type Ωfunction struct {
	// The function name
	Name string
	// If given, the object in the window object which contains the function and is the context for the function.
	// Use dot '.' notation to traverse the object tree. i.e. "obj1.obj2" refers to window.obj1.obj2 in javascript
	Context string
	// Function arguments. Strings will be quoted. Use a JsCode object to output the name of a javascript variable.
	Args []interface{}
}

func (f Ωfunction) JavaScript() string {
	var args string
	if f.Args != nil {
		args = Arguments(f.Args).JavaScript()
	}

	fName := f.Name
	if f.Context != "" {
		fName = f.Context + "." + fName
	}

	return fName + "(" + args + ")"
}

/**
 * Returns this as a json object to be sent to qcubed.js during ajax drawing.
 * @return mixed
 */
func (f Ωfunction) MarshalJSON() (buf []byte, err error) {
	var obj = map[string]interface{}{}

	obj[JsonObjectType] = "function"
	obj["func"] = f.Name
	if f.Context != "" {
		obj["context"] = f.Context
	}
	if f.Args != nil {
		obj["params"] = f.Args
	}

	buf, err = json.Marshal(obj)
	return
}

func init() {
	// Register objects so they can be serialized
	gob.Register(Ωfunction{})
}