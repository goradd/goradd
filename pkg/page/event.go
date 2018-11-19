package page

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/spekary/goradd/pkg/html"
	action2 "github.com/spekary/goradd/pkg/page/action"
	"goradd-project/config"
)

// EventI defines the interface for an event. Many of the routines implement a Builder pattern, allowing you to
// add options by chaining function calls on the event.
type EventI interface {
	// Condition sets a javascript condition that will be used to decide whether to fire the event.  The condition
	// will be evaluated when the event fires, and if its true, the associated action will be queued.
	Condition(string) EventI
	// Delay sets a delay for when the action happens after the event is fired, in milliseconds.
	Delay(int) EventI
	// Selector is a jQuery CSS selector that will filter bubbled events, only allowing events from the specified sub-items.
	Selector(string) EventI
	// Blocking will create an event that prevents all other events from firing before it is handled. Events that get
	// fired during this time are lost. This is particularly useful for Ajax or Server events that should not be
	// interrupted.
	Blocking() EventI
	// Terminating creates an event that does not bubble, nor will it do the browser default for the event.
	Terminating() EventI
	// Validate overrides the default validation specified by the control.
	Validate(v ValidationType) EventI
	// ValidationTargets overrides the validation targets specified by the control.
	ValidationTargets(targets ...string) EventI
	// ActionValue lets you specify a value that will be sent with the resulting action, and that will be accessible
	// through the EventValue() call on the ActionParams structure in the action handler. This is useful for Ajax and
	// Server actions primarily. This can be a static value, or a javascript.* type object to get dynamic values from javascript.
	ActionValue(interface{}) EventI

	// String returns a description of the event, primarily for debugging
	String() string


	addActions(a ...action2.ActionI)
	renderActions(control ControlI, eventID EventID) string
	getActions() []action2.ActionI
	event() *Event
}

type EventID uint16
type EventMap map[EventID]EventI

// Event represents a javascript event that triggers an action. Create it with a call to NewEvent(), or one of the
// predefined events in the event package, like event.Click()
type Event struct {
	JsEvent                   string
	condition                 string
	delay                     int
	selector                  string
	blocking                  bool
	actionValue               interface{} // A static value, or a javascript.* to get a dynamic value when the action returns to us.
	actions                   []action2.ActionI
	preventDefault            bool
	stopPropagation           bool
	validationOverride        ValidationType
	validationTargetsOverride []string
}

// NewEvent creates an event that triggers on the given event type. Use the builder pattern functions from EventI to
// add delays, conditions, etc.
func NewEvent(name string) EventI {
	return &Event{JsEvent: name}
}


func (e *Event) String() string {
	return fmt.Sprintf("Event: %s, Condition: %s, Delay: %d, Selector: %s, Blocking: %t, ActionCount: %d",
		e.JsEvent, e.condition, e.delay, e.selector, e.blocking, len(e.actions))
}

// returns underlying event structure for private access within package
func (e *Event) event() *Event {
	return e
}

// Condition specifies a javascript condition to check before triggering the event. The given string should be javascript
// code that evaluates to a boolean value.
func (e *Event) Condition(javascript string) EventI {
	e.condition = javascript
	return e
}

// Delay is the time in milliseconds to wait before triggering the actions.
func (e *Event) Delay(d int) EventI {
	e.delay = d
	return e
}

// Selector specifies a CSS filter that is used to check for bubbled events. This allows the event to be fired from
// child controls.
func (e *Event) Selector(s string) EventI {
	e.selector = s
	return e
}

// Call Blocking to cause this event to prevent other events from firing after this fires, but before it processes.
// If another event fires between the time when this event fires and when a response is received, it will be lost.
func (e *Event) Blocking() EventI {
	e.blocking = true
	return e
}


// Call Terminating to cause the event not to bubble or do the default action.
func (e *Event) Terminating() EventI {
	e.preventDefault = true
	e.stopPropagation = true
	return e
}

// Call PreventingDefault to cause the event not to do the default action.
func (e *Event) PreventingDefault() EventI {
	e.preventDefault = true
	return e
}

// Call PreventBubbling to cause the event to not bubble to enclosing objects.
func (e *Event) PreventBubbling() EventI {
	e.stopPropagation = true
	return e
}

// ActionValue is a value that will be returned to the actions that will be process by this event. Specify a static
// value, or javascript objects that will gather data at the time the event fires. The event will appear in the
// ActionParams as the EventValue.
// Example: ActionValue(javascript.ModelName("ui")) will return the "ui" variable that is part of the event call.
func (e *Event) ActionValue(r interface{}) EventI {
	e.actionValue = r
	return e
}

func (e *Event) GetActionValue() interface{} {
	return e.actionValue
}

// Validate overrides the controls validation setting just for this event.
func (e *Event) Validate(v ValidationType) EventI {
	e.validationOverride = v
	return e
}

// ValidationTargets overrides the control's validation targets just for this event.
func (e *Event) ValidationTargets(targets ...string) EventI {
	e.validationTargetsOverride = targets
	return e
}

func (e *Event) addActions(actions ...action2.ActionI) {
	var foundCallback bool
	for _, action := range actions {
		if _, ok := action.(action2.PrivateAction); ok {
			continue
		}
		if _, ok := action.(action2.CallbackActionI); ok {
			if foundCallback {
				panic("You can only associate one callback action with an event, and it must be the last action.")
			}
			foundCallback = true
		} else {
			if foundCallback {
				panic("You can only associate one callback action with an event, and it must be the last action.")
			}
		}

		e.actions = append(e.actions, action)
	}

	// Note, the above could be more robust and allow multiple callback actions, but it would get quite tricky if different
	// kinds of actions were interleaved. We will wait until someone presents a compelling need for something like that.
}

func (e *Event) renderActions(control ControlI, eventID EventID) string {
	if e.actions == nil {
		return ""
	}

	var js string

	if e.preventDefault {
		js += "event.preventDefault();\n"
	}
	if e.stopPropagation {
		js += "event.stopPropagation();\n"
	}

	var params = action2.ΩrenderParams{control.ID(), control.ActionValue(), uint16(eventID), e.actionValue}

	var actionJs string
	for _, a := range e.actions {
		actionJs += a.RenderScript(params)
	}

	if e.blocking {
		actionJs += "goradd.blockEvents = true;\n"
	}

	if !config.Minify {
		actionJs = html.Indent(actionJs)
	}
	actionJs = fmt.Sprintf("goradd.queueAction({f: $j.proxy(function(){\n%s\n},this), d: %d, name: '%s'});\n", actionJs, e.delay, e.JsEvent)

	if e.condition != "" {
		js = fmt.Sprintf("if (%s) {%s%s\n};", e.condition, js, actionJs)
	} else {
		js = js + actionJs
	}

	js = control.WrapEvent(e.JsEvent, e.selector, js)

	if !config.Minify {
		// Render a comment
		js = fmt.Sprintf("/*** Event: %s  Control Type: %T, Control Label: %s, Control Id: %s  ***/\n%s\n", e.JsEvent, control, control.Label(), control.ID(), js)
	}

	return js
}

func (e *Event) getActions() []action2.ActionI {
	return e.actions
}

type eventEncoded struct {
	JsEvent                   string
	Condition                 string
	Delay                     int
	Selector                  string
	Blocking                  bool
	ActionValue               interface{} // A static value, or js to get a dynamic value when the action returns to us.
	Actions                   []action2.ActionI
	PreventDefault            bool
	StopPropagation           bool
	ValidationOverride        ValidationType
	ValidationTargetsOverride []string
}

func (e *Event) GobEncode() (data []byte, err error) {
	s := eventEncoded {
		JsEvent: e.JsEvent,
		Condition: e.condition,
		Delay: e.delay,
		Selector:e.selector,
		Blocking:e.blocking,
		ActionValue:e.actionValue,
		Actions:e.actions,
		PreventDefault:e.preventDefault,
		StopPropagation: e.stopPropagation,
		ValidationOverride: e.validationOverride,
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
	e.actionValue = s.ActionValue
	e.actions = s.Actions
	e.preventDefault = s.PreventDefault
	e.stopPropagation = s.StopPropagation
	e.validationOverride = s.ValidationOverride
	e.validationTargetsOverride = s.ValidationTargetsOverride

	return nil
}
