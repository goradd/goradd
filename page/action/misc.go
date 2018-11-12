// The action package defines actions that you can trigger using events.
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
	"github.com/spekary/goradd/javascript"
)

// ActionI is an interface that defines actions that can be triggered by events
type ActionI interface {
	// RenderScript returns the action's javascript
	RenderScript(params ΩrenderParams) string
}

type ΩrenderParams struct {
	TriggeringControlID string
	ControlActionValue  interface{}
	EventID             uint16
	EventActionValue    interface{}
}

type ΩmessageAction struct {
	Message interface{}
}

// Note: Some actions currently depend on a javascript eval if they are introduced to a form during an ajax response.
// One way to fix that would be to register all javascript actions so that they get added to the form at drawing time,
// so that when an event gets attached during an ajax call, the resulting action is already in the browser.
// Another way is to use ajax to change a script tag.

// Message returns an action that will display a standard browser alert message. Specify a string, or one of the
// javascript.* types.
func Message(m interface{}) ΩmessageAction {
	return ΩmessageAction{Message: m}
}

func (a ΩmessageAction) RenderScript(params ΩrenderParams) string {
	return fmt.Sprintf(`goradd.msg(%s)`, javascript.ToJavaScript(a.Message))
}

type ΩconfirmAction struct {
	Message interface{}
}

// Confirm will put up a standard browser confirmation dialog box, and will cancel any following actions if the
// user does not agree.
func Confirm(m interface{}) ΩconfirmAction {
	return ΩconfirmAction{Message: m}
}

func (a ΩconfirmAction) RenderScript(params ΩrenderParams) string {
	return fmt.Sprintf("if (!window.confirm(%s)) return false;\n", javascript.ToJavaScript(a.Message))
}

type ΩblurAction struct {
	ControlID string
}

// Blur will blur the html object specified by the id.
func Blur(controlID string) ΩblurAction {
	return ΩblurAction{ControlID: controlID}
}

func (a ΩblurAction) RenderScript(params ΩrenderParams) string {
	return fmt.Sprintf(`goradd.blur('%s');`, a.ControlID)
}

// Focus will focus the html object specified by the id.
type ΩfocusAction struct {
	ControlID string
}

func Focus(controlID string) ΩfocusAction {
	return ΩfocusAction{ControlID: controlID}
}

func (a ΩfocusAction) RenderScript(params ΩrenderParams) string {
	return fmt.Sprintf(`goradd.focus('%s');`, a.ControlID)
}

type ΩselectAction struct {
	ControlID string
}

// Select will set all of the text in the html object specified by the id. The object should be a text box.
func Select(controlID string) ΩselectAction {
	return ΩselectAction{ControlID: controlID}
}

func (a ΩselectAction) RenderScript(params ΩrenderParams) string {
	return fmt.Sprintf(`goradd.select('%s');`, a.ControlID)
}

type ΩcssPropertyAction struct {
	Property  string
	Value     interface{}
	ControlID string
}

// SetCssProperty will set the css property to the given value on the controlID html object.
func SetCssProperty(id string, property string, value interface{}) ΩcssPropertyAction {
	return ΩcssPropertyAction{ControlID: id}
}

func (a ΩcssPropertyAction) RenderScript(params ΩrenderParams) string {
	return fmt.Sprintf(`goradd.css('%s', '%s', '%s');`, a.ControlID, a.Property, a.Value)
}

type ΩcssAddClassAction struct {
	Classes   string
	ControlID string
}

// AddClass will add the given class, or space separated classes, to the html object specified by the id.
func AddClass(id string, addClasses string) ΩcssAddClassAction {
	return ΩcssAddClassAction{ControlID: id, Classes: addClasses}
}

func (a ΩcssAddClassAction) RenderScript(params ΩrenderParams) string {
	return fmt.Sprintf(`goradd.addClass('%s', '%s');`, a.ControlID, a.Classes)
}

type ΩcssToggleClassAction struct {
	Classes   string
	ControlID string
}

// ToggleClass will turn on or off the given space separated classes in the html object specified by the id.
func ToggleClass(id string, classes string) ΩcssToggleClassAction {
	return ΩcssToggleClassAction{ControlID: id, Classes: classes}
}

func (a ΩcssToggleClassAction) RenderScript(params ΩrenderParams) string {
	return fmt.Sprintf(`goradd.toggleClass('%s', '%s');`, a.ControlID, a.Classes)
}

type ΩredirectAction struct {
	Location string
}

// Redirect will navigate to the given page.
func Redirect(url string) ΩredirectAction {
	return ΩredirectAction{Location: url}
}

func (a ΩredirectAction) RenderScript(params ΩrenderParams) string {
	return fmt.Sprintf(`goradd.redirect("%s");`, a.Location)
}

type ΩtriggerAction struct {
	ControlID string
	Event     string
	Data      interface{}
}

// Trigger will trigger a javascript event on a control
func Trigger(controlID string, event string, data interface{}) ΩtriggerAction {
	return ΩtriggerAction{ControlID: controlID, Event: event, Data: data}
}

func (a ΩtriggerAction) RenderScript(params ΩrenderParams) string {
	return fmt.Sprintf(`$j("#%s").trigger("%s", %s);` + "\n", a.ControlID, a.Event, javascript.ToJavaScript(a.Data))
}

// PrivateAction is used by control implementations to add a private action to a controls action list. Unless you are
// creating a control, you should not use this.
type PrivateAction struct{}

func (a PrivateAction) RenderScript(params ΩrenderParams) string {
	return ""
}

type ΩjavascriptAction struct {
	JavaScript string
}

// Javascript will execute the given javascript
func Javascript(js string) ΩjavascriptAction {
	if js != "" {
		if js[len(js) - 1: len(js)] != ";" {
			js += ";\n"
		}
	}
	return ΩjavascriptAction{JavaScript: js}
}

func (a ΩjavascriptAction) RenderScript(params ΩrenderParams) string {
	return a.JavaScript
}

type ΩsetControlValueAction struct {
	ID    string
	Key   string
	Value interface{}
}

// SetControlValue is primarily used by custom controls to set a value that eventually can get picked
// up by the control in the UpdateFormValues function. It is an aid to tying javascript powered widgets together
// with the go version of the control. Value gets converted to a javascript value, so use the javascript.* helpers
// if you want to interpret a javascript value and pass it on. For example:
//  action.SetControlValue(myControl.ID(), "myKey", javascript.JsCode("event.target.id"))
// will pass the id of the target of an event to the receiver of the action.
func SetControlValue(id string, key string, value interface{}) ΩsetControlValueAction {
	return ΩsetControlValueAction{ID: id, Key:key, Value:value}
}

func (a ΩsetControlValueAction) RenderScript(params ΩrenderParams) string {
	return fmt.Sprintf(`goradd.setControlValue("%s", "%s", %s)`, a.ID, a.Key, javascript.ToJavaScript(a.Value))
}

func init() {
	// Register actions so they can be serialized
	gob.Register(ΩmessageAction{})
	gob.Register(ΩconfirmAction{})
	gob.Register(ΩblurAction{})
	gob.Register(ΩfocusAction{})
	gob.Register(ΩselectAction{})
	gob.Register(ΩcssPropertyAction{})
	gob.Register(ΩcssAddClassAction{})
	gob.Register(ΩcssToggleClassAction{})
	gob.Register(ΩredirectAction{})
	gob.Register(ΩtriggerAction{})
	gob.Register(ΩjavascriptAction{})
	gob.Register(ΩsetControlValueAction{})
}
