// Package event contains functions that specify various kinds of javascript events that GoRADD controls respond to.
//
// Create an event by using one of the predefined event creation functions like [Click], or
// call [NewEvent] to create an event that responds to any named javascript event. If needed, add
// additional requirements for the event using the builder pattern functions like [Event.Delay], [Event.Selector]
// and [Event.Condition].
//
// For example, the code below will create an event that waits for clicks on the button, but debounces
// the clicks and also prevents all other actions from happening while waiting for the event to fire. This
// would be typically useful in a submit button where you want to prevent multiple submissions of the same button.
//
//	btn := NewButton().On(event.Click().Delay(200).Blocking(), action.Redirect("/mypage"))
package event

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/goradd/goradd/pkg/config"
	"github.com/goradd/goradd/pkg/javascript"
	"github.com/goradd/goradd/pkg/page/action"
	"github.com/goradd/html5tag"
	"strconv"
)

// EventID is used internally by the framework to set a unique id used to specify which event is triggering.
type EventID uint16

// Event represents a javascript event that triggers an action. Create it with a call to NewEvent(), or one of the
// predefined events in the event package, like event.Click()
type Event struct {
	// jsEvent is the JavaScript event that will be triggered by the event.
	jsEvent string
	// condition is a javascript comparison test that if present, must evaluate as true for the event to fire
	condition string
	// delay is the number of milliseconds to delay firing the action after the event triggers
	delay int
	// selector is a css selector that will filter the event bubbling up from a sub-item. The event must originate from
	// an HTML object that the selector specifies. Inside the event handler, event.goradd.match will contain the
	// html item that was selected by the selector. Note that by default, if this object has other objects inside it,
	// events will not bubble up from those objects. Set "bubbles" to true to have sub objects bubble their events.
	selector string
	// blocking specifies that once the event fires, all other events will be blocked until the action associated with
	// the event returns.
	blocking bool
	// actionValue is a static value, or a javascript.* to get a dynamic value when the action returns to us.
	// this value will become the EventValue returned to the action.
	actionValue interface{}
	// action is the action that the event triggers. Multiple actions can be specified using an action group.
	action action.ActionI
	// preventDefault will cause the preventDefault function to be called on the event, which prevents the
	// default action. In particular, this would prevent a submit button from submitting a form.
	preventDefault bool
	// stopPropagation will cause the stopPropagation jQuery function to be called on the event, which prevents the event
	// from bubbling.
	stopPropagation bool
	// validationOverride allows the event to override the control's validation mechanism.
	validationOverride ValidationType
	// validationTargetsOverride allows the event to specify custom targets for validation.
	validationTargetsOverride []string
	// internally assigned event id
	eventID EventID
	// bubbles is used with a selector to determine if the event should bubble up from items contained by the selector,
	// or should only come from the item specified by the selector
	bubbles bool
	// capture indicates an event should fire during the capture phase and not the bubbling phase.
	// This is used in very special situations when you want to not allow a bubbled event to be blocked by a sub-object as it bubbles.
	capture bool
	// private indicates the event is private to the control and cannot be changed or removed. It is responded to in
	// the DoPrivateAction function
	private bool
}

// NewEvent creates an event that triggers on the given javascript event name.
// Use the builder pattern functions from *Event to add delays, conditions, etc.
func NewEvent(name string) *Event {
	return &Event{jsEvent: name}
}

// String returns a debug string listing the contents of the event.
func (e *Event) String() string {
	return fmt.Sprintf("Event: %s, Condition: %s, Delay: %d, Selector: %s, Blocking: %t",
		e.jsEvent, e.condition, e.delay, e.selector, e.blocking)
}

// Condition specifies a javascript condition to check before triggering the event. The given string should be javascript
// code that evaluates to a boolean value.
func (e *Event) Condition(javascript string) *Event {
	e.condition = javascript
	return e
}

// Delay is the time in milliseconds to wait before triggering the actions.
//
// During the delay time, if the event is repeated, the delay timer will restart and
// only one event will eventually be fired. For example, if you have a KeyDown event with a delay,
// and the user enters multiple keys during the delay time, only one keydown event will fire, and it
// will fire delay ms after the last keydown event was received.
func (e *Event) Delay(delay int) *Event {
	e.delay = delay
	return e
}

// Selector specifies a CSS filter that is used to check for bubbled events. This allows the event to be fired from
// child controls. By default, the event will not come from sub-controls of the
// specified child controls. Use Bubbles or Capture to change that.
func (e *Event) Selector(s string) *Event {
	e.selector = s
	return e
}

// Bubbles works with a Selector to allow events to come from a sub-control
// of the selected control. The event could be blocked by the sub-control if
// the sub-control issues a preventPropagation on the event.
func (e *Event) Bubbles() *Event {
	e.bubbles = true
	return e
}

// Capture works with a Selector to allow events to come from a sub-control
// of the selected control. The event never actually reaches the sub-control
// for processing, and instead is captured and handled by the selected control.
// This is generally used in special situations where you do not want to allow
// sub-controls to prevent bubbling.
func (e *Event) Capture() *Event {
	e.capture = true
	return e
}

// Blocking prevents other events from firing after this fires, but before it processes.
// If another event fires between the time when this event fires and when a response is received, it will be lost.
func (e *Event) Blocking() *Event {
	e.blocking = true
	return e
}

// Terminating prevents the event from bubbling or doing the default action. It is essentially a combination
// of calling PreventDefault and StopPropagation.
func (e *Event) Terminating() *Event {
	e.preventDefault = true
	e.stopPropagation = true
	return e
}

// PreventDefault causes the event not to do the default action.
func (e *Event) PreventDefault() *Event {
	e.preventDefault = true
	return e
}

// PreventBubbling causes the event to not bubble to enclosing objects.
func (e *Event) PreventBubbling() *Event {
	e.stopPropagation = true
	return e
}

// ActionValue sets the event value in the actions triggered by this event. Specify a static
// value, or javascript objects that will gather data at the time the event fires. The event will appear in the
// Params as the EventValue. By default, this will be the value passed in to the javascript event as event data.
//
// See on: and trigger: in goradd.js.
//
// For example:
//
//	ActionValue(javascript.JsCode{"event.target.id"})
//
// will cause the EventValue for the action to be the HTML id of the target object of the event.
func (e *Event) ActionValue(r interface{}) *Event {
	e.actionValue = r
	return e
}

// EventValueTargetID will set the event value of the resulting action to the HTML id of the target of the event.
func (e *Event) EventValueTargetID() *Event {
	e.actionValue = javascript.JsCode("event.target.id")
	return e
}

// GetActionValue returns the event value associated with the actions resulting from the event.
func GetActionValue(e *Event) interface{} {
	return e.actionValue
}

// Validate overrides the controls validation setting just for this event.
func (e *Event) Validate(v ValidationType) *Event {
	e.validationOverride = v
	return e
}

// ValidationTargets overrides the control's validation targets just for this event.
func (e *Event) ValidationTargets(targets ...string) *Event {
	e.validationTargetsOverride = targets
	return e
}

// Private makes the event private to the control and not removable. This should generally only be used by
// control implementations to add events that are required by the control and that should not be removed by Off()
func (e *Event) Private() *Event {
	e.private = true
	return e
}

// Name returns the name of the javascript event being triggered.
func Name(e *Event) string {
	return e.jsEvent
}

// ID returns the event id.
func ID(e *Event) EventID {
	return e.eventID
}

type getCallbackActioner interface {
	GetCallbackAction() action.CallbackActionI
}

// GetCallbackAction will return the action associated with the event if it is a callback action.
// Otherwise, it will return nil.
func GetCallbackAction(e *Event) action.CallbackActionI {
	switch a := e.action.(type) {
	case action.CallbackActionI:
		return a
	case getCallbackActioner:
		return a.GetCallbackAction()
	default:
		return nil
	}
}

// IsPrivate returns whether this is an event private to the control or one created from outside the control.
func IsPrivate(e *Event) bool {
	return e.private
}

// GetValidationOverride returns the validation type of the event that will override the control's validation type.
func GetValidationOverride(e *Event) ValidationType {
	return e.validationOverride
}

// eventEncoded contains exported types. We use this to serialize an event for the page serializer.
type eventEncoded struct {
	JsEvent                   string
	Condition                 string
	Delay                     int
	Selector                  string
	Blocking                  bool
	Bubbles                   bool
	Capture                   bool
	Private                   bool
	ActionValue               interface{} // A static value, or js to get a dynamic value when the action returns to us.
	Action                    action.ActionI
	PreventDefault            bool
	StopPropagation           bool
	ValidationOverride        ValidationType
	ValidationTargetsOverride []string
}

// GobEncode is called by the framework to binary encode the event.
func (e *Event) GobEncode() (data []byte, err error) {
	s := eventEncoded{
		JsEvent:                   e.jsEvent,
		Condition:                 e.condition,
		Delay:                     e.delay,
		Selector:                  e.selector,
		Blocking:                  e.blocking,
		Bubbles:                   e.bubbles,
		Capture:                   e.capture,
		Private:                   e.private,
		ActionValue:               e.actionValue,
		Action:                    e.action,
		PreventDefault:            e.preventDefault,
		StopPropagation:           e.stopPropagation,
		ValidationOverride:        e.validationOverride,
		ValidationTargetsOverride: e.validationTargetsOverride,
	}
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)

	err = encoder.Encode(s)
	data = buf.Bytes()
	return
}

// GobDecode is called by the framework to binary decode the event.
func (e *Event) GobDecode(data []byte) (err error) {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	s := eventEncoded{}
	if err = dec.Decode(&s); err != nil {
		return
	}
	e.jsEvent = s.JsEvent
	e.condition = s.Condition
	e.delay = s.Delay
	e.selector = s.Selector
	e.blocking = s.Blocking
	e.bubbles = s.Bubbles
	e.capture = s.Capture
	e.private = s.Private
	e.actionValue = s.ActionValue
	e.action = s.Action
	e.preventDefault = s.PreventDefault
	e.stopPropagation = s.StopPropagation
	e.validationOverride = s.ValidationOverride
	e.validationTargetsOverride = s.ValidationTargetsOverride

	return nil
}

// SetEventItems is used internally by the framework to set up an event.
// You should not normally need to call this function.
func SetEventItems(e *Event, action action.ActionI, eventId EventID) {
	e.action = action
	e.eventID = eventId
}

// renderer is used to help avoid import loops by defining the functions that impact rendering
// without importing the page or control packages
type renderer interface {
	ID() string
	ActionValue() any
	WrapEvent(eventName string, selector string, eventJs string, options map[string]interface{}) string
}

// RenderActions is used internally by the framework to render the javascript
// associated with the event and connected actions.
// You should not normally need to call this function.
func RenderActions(e *Event, control renderer, eventID EventID) string {
	if e.action == nil {
		return ""
	}

	var js string
	var options map[string]interface{}

	if e.capture || e.bubbles {
		options = make(map[string]interface{})
		if e.bubbles {
			options["bubbles"] = true
		}
		if e.capture {
			options["capture"] = true
		}
	}

	if e.preventDefault {
		js += "event.preventDefault();\n"
	}
	if e.stopPropagation {
		js += "event.stopPropagation();\n"
	}

	var params = action.RenderParams{TriggeringControlID: control.ID(), ControlActionValue: control.ActionValue(), EventID: uint16(eventID), EventActionValue: e.actionValue}

	var actionJs = e.action.RenderScript(params)

	if e.blocking {
		actionJs += "goradd.blockEvents = true;\n"
	}

	actionJs += "if (event.goradd && event.goradd.postFunc) {event.goradd.postFunc();}\n"

	if !config.Minify {
		actionJs = html5tag.Indent(actionJs)
	}
	if e.delay == 0 {
		actionJs = fmt.Sprintf(
			"goradd.queueAction({f: (function(){\n%s\n}).bind(this.element), name: '%s'});\n",
			actionJs, e.jsEvent,
		)
	} else {
		actionJs = fmt.Sprintf(
			"goradd.queueAction({f: (function(){\n%s\n}).bind(this.element), name: '%s', d: %d, k: '%s'});\n",
			actionJs, e.jsEvent, e.delay,
			// An event key unique to events used on the page, for preventing the same event from repeating during the delay time
			control.ID()+"-"+strconv.Itoa(int(e.eventID)))
	}

	if e.condition != "" {
		js = fmt.Sprintf("if (%s) {%s%s\n};", e.condition, js, actionJs)
	} else {
		js = js + actionJs
	}

	js = control.WrapEvent(e.jsEvent, e.selector, js, options)

	if !config.Minify {
		// Render a comment
		js = fmt.Sprintf("/*** Event: %s  ControlBase Type: %T, ControlBase Id: %s  ***/\n%s\n", e.jsEvent, control, control.ID(), js)
	}

	return js
}

func init() {
	gob.Register(map[string]interface{}{}) // This is so that action values that use this can be serialized
}
