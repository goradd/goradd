package action

import (
	"encoding/gob"
	"fmt"
	"github.com/goradd/goradd/pkg/javascript"
)

type goraddAction struct {
	Op   string
	Args []interface{}
}

// GoraddFunction calls a function defined on the global goradd object that is in goradd.js.
func GoraddFunction(operation string, arguments ...interface{}) ActionI {
	return goraddAction{operation, arguments}
}

// RenderScript outputs the action as JavaScript.
func (a goraddAction) RenderScript(params RenderParams) string {
	return fmt.Sprintf(`goradd.%s(%s);`, a.Op, javascript.Arguments(a.Args).JavaScript())
}

// Message returns an action that will display a standard browser alert message. Specify a string, or one of the
// javascript.* types.
func Message(m interface{}) ActionI {
	return GoraddFunction("msg", m)
}

// Refresh will cause the given control to redraw
func Refresh(id string) ActionI {
	return GoraddFunction("refresh", id)
}

// SetControlValue is primarily used by custom controls to set a value that eventually can get picked
// up by the control in the UpdateFormValues function. It is an aid to tying javascript powered widgets together
// with the go version of the control. Value gets converted to a javascript value, so use the javascript.* helpers
// if you want to interpret a javascript value and pass it on. For example:
//
//	action.SetControlValue(myControl.ID(), "myKey", javascript.NewJsCode("event.target.id"))
//
// will pass the id of the target of an event to the receiver of the action.
//
// Note that if you need to save and restore state by calling MarshalState and UnmarshalState, you will need to
// have your control listen for the event that sets this value and call action.Do() so that UpdateFormValues will
// eventually get called.
func SetControlValue(id string, key string, value interface{}) ActionI {
	return GoraddFunction("setControlValue", id, key, value)
}

func init() {
	gob.Register(goraddAction{})
}
