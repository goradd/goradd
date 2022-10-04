package page

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/goradd/goradd/pkg/config"
	action2 "github.com/goradd/goradd/pkg/page/action"
	"github.com/goradd/html5tag"
	"strconv"
)

// EventID is a unique id used to specify which event is triggering.
type EventID uint16
type EventMap map[EventID]*Event

// Event represents a javascript event that triggers an action. Create it with a call to NewEvent(), or one of the
// predefined events in the event package, like event.Click()
type Event struct {
	// JsEvent is the JavaScript event that will be triggered by the event.
	JsEvent string
	// condition is a javascript comparison test that if present, must evaluate as true for the event to fire
	condition string
	// delay is the number of milliseconds to delay firing the action after the event triggers
	delay int
	// selector is a css selector that will filter the event bubbling up from a sub-item. The event must originate from
	// an html object that the selector specifies. Inside the event handler, event.goradd.match will contain the
	// html item that was selected by the selector. Note that by default, if this object has other objects inside of it,
	// events will not bubble up from those objects. Set "bubbles" to true to have sub objects bubble their events.
	selector string
	// blocking specifies that once the event fires, all other events will be blocked until the action associated with
	// the event returns.
	blocking bool
	// actionValue is a static value, or a javascript.* to get a dynamic value when the action returns to us.
	// this value will become the EventValue returned to the action.
	actionValue interface{}
	// action is the action that the event triggers. Multiple actions can be specified using an action group.
	action action2.ActionI
	// preventDefault will cause the preventDefault function to be called on the event, which prevents the
	// default action. In particular, this would prevent a submit button from submitting a form.
	preventDefault bool
	// stopPropogation will cause the stopPropogation jQuery function to be called on the event, which prevents the event
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
	// the PrivateAction function
	private bool
}

// NewEvent creates an event that triggers on the given javascript event name.
// Use the builder pattern functions from *Event to add delays, conditions, etc.
func NewEvent(name string) *Event {
	return &Event{JsEvent: name}
}

// String returns a debug string listing the contents of the event.
func (e *Event) String() string {
	return fmt.Sprintf("Event: %s, Condition: %s, Delay: %d, Selector: %s, Blocking: %t",
		e.JsEvent, e.condition, e.delay, e.selector, e.blocking)
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

// Bubbles works with a Selector to allow events to come from sub-control
// of the selected control. The event could be blocked by the sub-control if
// the subcontrol issues a preventPropagation on the event.
func (e *Event) Bubbles() *Event {
	e.bubbles = true
	return e
}

// Capture works with a Selector to allow events to come from a sub-control
// of the selected control. The event never actually reaches the sub-control
// for processing, and instead is captured and handled by the selected control.
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

// Call PreventingDefault to cause the event not to do the default action.
func (e *Event) PreventingDefault() *Event {
	e.preventDefault = true
	return e
}

// Call PreventBubbling to cause the event to not bubble to enclosing objects.
func (e *Event) PreventBubbling() *Event {
	e.stopPropagation = true
	return e
}

// ActionValue is a value that will be returned to the actions that will be process by this event. Specify a static
// value, or javascript objects that will gather data at the time the event fires. The event will appear in the
// ActionParams as the EventValue. By default, this will be the value passed in to the event as event data.
// See on: and trigger: in goradd.js.
// Example: ActionValue(javascript.JsCode{"event.target.id"}) will return the target id from the event object passed in to the event handler.
func (e *Event) ActionValue(r interface{}) *Event {
	e.actionValue = r
	return e
}

// GetActionValue returns the value associated with the action of the event.
func (e *Event) GetActionValue() interface{} {
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

// HasServerAction returns true if at least one of the event's actions is a server action.
func (e *Event) HasServerAction() bool {
	switch a := e.action.(type) {
	case action2.FrameworkCallbackActionI:
		return a.IsServerAction()
	case action2.ActionGroup:
		return a.HasServerAction()
	default:
		return false
	}
}

// HasCallbackAction returns true if at least one of the event's actions is a callback action.
func (e *Event) HasCallbackAction() bool {
	switch a := e.action.(type) {
	case action2.CallbackActionI:
		return true
	case action2.ActionGroup:
		return a.HasCallbackAction()
	default:
		return false
	}
}

func (e *Event) Name() string {
	return e.JsEvent
}

func (e *Event) addAction(a action2.ActionI) {
	e.action = a
}

func (e *Event) renderActions(control ControlI, eventID EventID) string {
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

	var params = action2.RenderParams{control.ID(), control.ActionValue(), uint16(eventID), e.actionValue}

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
			actionJs, e.JsEvent,
		)
	} else {
		actionJs = fmt.Sprintf(
			"goradd.queueAction({f: (function(){\n%s\n}).bind(this.element), name: '%s', d: %d, k: '%s'});\n",
			actionJs, e.JsEvent, e.delay,
			// An event key unique to events used on the page, for preventing the same event from repeating during the delay time
			control.ID()+"-"+strconv.Itoa(int(e.eventID)))
	}

	if e.condition != "" {
		js = fmt.Sprintf("if (%s) {%s%s\n};", e.condition, js, actionJs)
	} else {
		js = js + actionJs
	}

	js = control.WrapEvent(e.JsEvent, e.selector, js, options)

	if !config.Minify {
		// Render a comment
		js = fmt.Sprintf("/*** Event: %s  ControlBase Type: %T, ControlBase Id: %s  ***/\n%s\n", e.JsEvent, control, control.ID(), js)
	}

	return js
}

func (e *Event) getCallbackAction() action2.FrameworkCallbackActionI {
	switch a := e.action.(type) {
	case action2.CallbackActionI:
		return a.(action2.FrameworkCallbackActionI)
	case action2.ActionGroup:
		return a.GetCallbackAction()
	default:
		return nil
	}
}

func (e *Event) isPrivate() bool {
	return e.private
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
	Action                    action2.ActionI
	PreventDefault            bool
	StopPropagation           bool
	ValidationOverride        ValidationType
	ValidationTargetsOverride []string
}

func (e *Event) GobEncode() (data []byte, err error) {
	s := eventEncoded{
		JsEvent:                   e.JsEvent,
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

func (e *Event) GobDecode(data []byte) (err error) {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	s := eventEncoded{}
	if err = dec.Decode(&s); err != nil {
		return
	}
	e.JsEvent = s.JsEvent
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

func init() {
	gob.Register(map[string]interface{}{}) // This is so that action values that use this can be serialized
}
