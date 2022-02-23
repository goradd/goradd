package javascript

import (
	"encoding/gob"
	"encoding/json"
)

// NewFunctionCall creates a new FunctionCall object.
//
// context will become the "this" value inside the Closure.
// args will be passed as values, and strings will be quoted. To pass a variable name, wrap the name with a NewJsCode call.
func NewFunctionCall(name string, context string, args ...interface{}) FunctionCall {
	return FunctionCall{name, context, args}
}

// FunctionCall represents the result of a function call to a global function or function in an object referenced from
// global space.
//
// The purpose of this is to immediately use the results of the function call, as opposed to a Closure,
// which stores a pointer to a function that is used later.
type FunctionCall struct {
	// Name is the function name
	Name string
	// Context, if given, is the object in the window object which contains the function and is the context for the function.
	// Use dot '.' notation to traverse the object tree. i.e. "obj1.obj2" refers to window.obj1.obj2 in javascript
	Context string
	// Args is the list of arguments of the function call.
	// Strings will be quoted. Use a NewJsCode object to output the name of a javascript variable.
	Args []interface{}
}

// JavaScript implements the JavaScripter interface and outputs the function call as embedded JavaScript.
func (f FunctionCall) JavaScript() string {
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

// MarshalJSON implements the json.Marshaller interface.
func (f FunctionCall) MarshalJSON() (buf []byte, err error) {
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
	gob.Register(FunctionCall{})
}
