// Package action defines actions that you can trigger using events.
// Normally you would do this with the .On() function that all goradd controls have.
//
// Defining Your Own Actions
// You can define your own actions by creating a class that implements the ActionI interface, AND that is
// encodable by gob.Serialize, meaning it either implements the gob.Encoder interface or exports its structure, AND
// registers itself with gob.Register so that the gob.Decoder knows how to deserialize it into an interface.
// We have chosen to export the structures that represent an action here, but we prefix the name of the structures with a
// greek capital Omega (Ω). We do this to call out that these exported structures and variables are not for general use.
package action

import (
	"encoding/gob"
	"fmt"
	"github.com/goradd/goradd/pkg/javascript"
)

// ActionI is an interface that defines actions that can be triggered by events
type ActionI interface {
	// ΩRenderScript is called by the framework to return the action's javascript
	ΩRenderScript(params ΩrenderParams) string
}

type ΩrenderParams struct {
	TriggeringControlID string
	ControlActionValue  interface{}
	EventID             uint16
	EventActionValue    interface{}
}

type ΩwidgetAction struct {
	ControlID string
	Op string
	Args []interface{}
}

// WidgetFunction calls a goradd widget function in javascript on an html control with the given id.
// Available functions are defined by the widget object in goradd.js
func WidgetFunction (controlID string, operation string, arguments ...interface{}) ΩwidgetAction {
	return ΩwidgetAction{controlID, operation, arguments}
}

func (a ΩwidgetAction) ΩRenderScript(params ΩrenderParams) string {
	return fmt.Sprintf(`g$('%s').%s(%s);`, a.ControlID, a.Op, javascript.Arguments(a.Args).JavaScript())
}

type ΩgoraddAction struct {
	Op string
	Args []interface{}
}

// GoraddFunction calls a goradd function with the given parameters. This is a function defined on the goradd
// object in goradd.js.
func GoraddFunction (operation string, arguments ...interface{}) ΩgoraddAction {
	return ΩgoraddAction{operation, arguments}
}

func (a ΩgoraddAction) ΩRenderScript(params ΩrenderParams) string {
	return fmt.Sprintf(`goradd.%s(%s);`, a.Op, javascript.Arguments(a.Args).JavaScript())
}


// Message returns an action that will display a standard browser alert message. Specify a string, or one of the
// javascript.* types.
func Message(m interface{}) ΩgoraddAction {
	return GoraddFunction("msg", m)
}

type ΩconfirmAction struct {
	Message interface{}
}

// Confirm will put up a standard browser confirmation dialog box, and will cancel any following actions if the
// user does not agree.
func Confirm(m interface{}) ΩconfirmAction {
	return ΩconfirmAction{Message: m}
}

func (a ΩconfirmAction) ΩRenderScript(params ΩrenderParams) string {
	return fmt.Sprintf("if (!window.confirm(%s)) return false;\n", javascript.ToJavaScript(a.Message))
}

// Blur will blur the html object specified by the id.
func Blur(controlID string) ΩwidgetAction {
	return WidgetFunction(controlID, "blur")
}

// Focus will focus the html object specified by the id.
func Focus(controlID string) ΩwidgetAction {
	return WidgetFunction(controlID, "focus")
}

// Select will set all of the text in the html object specified by the id. The object should be a text box.
func Select(controlID string) ΩwidgetAction {
	return WidgetFunction(controlID, "selectAll")
}

// Css will set the css property to the given value on the given html object.
func Css(controlID string, property string, value interface{}) ΩwidgetAction {
	return WidgetFunction(controlID, "css", property, value)
}

// AddClass will add the given class, or space separated classes, to the html object specified by the id.
func AddClass(controlID string, classes string) ΩwidgetAction {
	return WidgetFunction(controlID, "class", "+" + classes)
}

// ToggleClass will turn on or off the given space separated classes in the html object specified by the id.
func ToggleClass(controlID string, classes string) ΩwidgetAction {
	return WidgetFunction(controlID, "toggleClass", classes)
}

// RemoveClass will turn off the given space separated classes in the html object specified by the id.
func RemoveClass(controlID string, classes string) ΩwidgetAction {
	return WidgetFunction(controlID, "class", "-" + classes)
}

type ΩredirectAction struct {
	Location string
}

// Redirect will navigate to the given page.
// TODO: If javascript is turned off, this should still work. We would need to detect the presence of javascript,
// and then emit a server action instead
func Redirect(url string) ΩredirectAction {
	return ΩredirectAction{Location: url}
}

func (a ΩredirectAction) ΩRenderScript(params ΩrenderParams) string {
	return fmt.Sprintf(`goradd.redirect("%s");`, a.Location)
}

// Trigger will trigger a javascript event on a control
func Trigger(controlID string, event string, data interface{}) ΩwidgetAction {
	return WidgetFunction(controlID, "trigger", event, data)
}

// PrivateAction is used by control implementations to add a private action to a controls action list. Unless you are
// creating a control, you should not use this.
type PrivateAction struct{}

func (a PrivateAction) ΩRenderScript(params ΩrenderParams) string {
	return ""
}

type ΩjavascriptAction struct {
	JavaScript string
}

// Javascript will execute the given javascript
func Javascript(js string) ΩjavascriptAction {
	if js != "" {
		if js[len(js)-1:len(js)] != ";" {
			js += ";\n"
		}
	}
	return ΩjavascriptAction{JavaScript: js}
}

func (a ΩjavascriptAction) ΩRenderScript(params ΩrenderParams) string {
	return a.JavaScript
}

// SetControlValue is primarily used by custom controls to set a value that eventually can get picked
// up by the control in the ΩUpdateFormValues function. It is an aid to tying javascript powered widgets together
// with the go version of the control. Value gets converted to a javascript value, so use the javascript.* helpers
// if you want to interpret a javascript value and pass it on. For example:
//   action.SetControlValue(myControl.ID(), "myKey", javascript.JsCode("event.target.id"))
// will pass the id of the target of an event to the receiver of the action.
func SetControlValue(id string, key string, value interface{}) ΩgoraddAction {
	return GoraddFunction("setControlValue", id, key, value)
}

func init() {
	// Register actions so they can be serialized
	gob.Register(ΩgoraddAction{})
	gob.Register(ΩwidgetAction{})
	gob.Register(ΩconfirmAction{})
	gob.Register(ΩredirectAction{})
	gob.Register(ΩjavascriptAction{})
}
