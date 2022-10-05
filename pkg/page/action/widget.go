package action

import (
	"encoding/gob"
	"fmt"
	"github.com/goradd/goradd/pkg/javascript"
)

type widgetAction struct {
	ControlID string
	Op        string
	Args      []interface{}
}

// WidgetFunction calls a goradd widget function in javascript on an HTML control with the given id.
// Available functions are defined by the widget object in goradd.js
func WidgetFunction(controlID string, operation string, arguments ...interface{}) ActionI {
	return widgetAction{controlID, operation, arguments}
}

// RenderScript is called by the framework to output the action as JavaScript.
func (a widgetAction) RenderScript(_ RenderParams) string {
	return fmt.Sprintf(`g$('%s').%s(%s);`, a.ControlID, a.Op, javascript.Arguments(a.Args).JavaScript())
}

// Blur will blur the html object specified by the id.
func Blur(controlID string) ActionI {
	return WidgetFunction(controlID, "blur")
}

// Focus will focus the html object specified by the id.
func Focus(controlID string) ActionI {
	return WidgetFunction(controlID, "focus")
}

// Select will select all the text in the html object specified by the id. The object should be a text box.
func Select(controlID string) ActionI {
	return WidgetFunction(controlID, "selectAll")
}

// Css will set the css property to the given value on the given html object.
func Css(controlID string, property string, value interface{}) ActionI {
	return WidgetFunction(controlID, "css", property, value)
}

// AddClass will add the given class, or space separated classes, to the html object specified by the id.
func AddClass(controlID string, classes string) ActionI {
	return WidgetFunction(controlID, "class", "+"+classes)
}

// ToggleClass will turn on or off the given space separated classes in the html object specified by the id.
func ToggleClass(controlID string, classes string) ActionI {
	return WidgetFunction(controlID, "toggleClass", classes)
}

// RemoveClass will turn off the given space separated classes in the html object specified by the id.
func RemoveClass(controlID string, classes string) ActionI {
	return WidgetFunction(controlID, "class", "-"+classes)
}

// Show will show the given control if it is hidden
func Show(controlID string) ActionI {
	return WidgetFunction(controlID, "show")
}

// Hide will hide the given control
func Hide(controlID string) ActionI {
	return WidgetFunction(controlID, "hide")
}

// Trigger will trigger a javascript event on a control
func Trigger(controlID string, event string, data interface{}) ActionI {
	return WidgetFunction(controlID, "trigger", event, data)
}

func init() {
	gob.Register(widgetAction{})
}
