package page

import (
	"fmt"
	"github.com/spekary/goradd/html"
	action2 "github.com/spekary/goradd/page/action"
	"goradd-project/config"
)

type EventI interface {
	Condition(string) EventI
	Delay(int) EventI
	Selector(string) EventI
	Blocking() EventI
	Terminating() EventI
	Validate(v ValidationType) EventI
	ValidationTargets(targets ...string) EventI
	ActionValue(interface{}) EventI
	GetActionValue() interface{}
	AddActions(a ...action2.ActionI)
	RenderActions(control ControlI, eventID EventID) string
	GetActions() []action2.ActionI
	String() string
	event() *Event
}

type EventID uint16
type EventMap map[EventID]EventI

// Event is a base class for other events. You should not call it directly. See the event package for implementations.
type Event struct {
	JsEvent                   string
	condition                 string
	delay                     int
	selector                  string
	blocking                  bool
	actionValue               interface{} // A static value, or js to get a dynamic value when the action returns to us.
	actions                   []action2.ActionI
	preventDefault            bool
	stopPropagation           bool
	validationOverride        ValidationType
	validationTargetsOverride []string
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

// Call NotBubbling to cause the event to not bubble to enclosing objects.
func (e *Event) NotBubbling() EventI {
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

func (e *Event) AddActions(actions ...action2.ActionI) {
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

func (e *Event) RenderActions(control ControlI, eventID EventID) string {
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

	var params = action2.RenderParams{control.ID(), control.ActionValue(), uint16(eventID), e.actionValue}

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

func (e *Event) GetActions() []action2.ActionI {
	return e.actions
}
