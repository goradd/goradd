package page

import (
	"fmt"
	"goradd/config"
	"github.com/spekary/goradd/html"
)

type EventI interface {
	Condition(string) EventI
	Delay(int) EventI
	Selector(string) EventI
	Blocking() EventI
	Terminating() EventI
	ActionValue(interface{}) EventI
	GetActionValue() interface{}
	AddActions(a... ActionI)
	RenderActions(control ControlI, eventId EventId) string
	GetActions() []ActionI
	String() string
	event() *Event
}

type EventId uint16
type EventMap map[EventId]EventI

// Event is a base class for other events. You should not call it directly. See the event package for implementations.
type Event struct {
	JsEvent     string
	condition   string
	delay       int
	selector    string
	blocking    bool
	actionValue interface{} // A static value, or js to get a dynamic value when the action returns to us.
	actions     []ActionI
	actionsMustTerminate bool
	validationOverride ValidationType
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
func (e *Event) Condition (javascript string) EventI {
	e.condition = javascript
	return e
}

// Delay is the time in milliseconds to wait before triggering the actions.
func (e *Event) Delay (d int) EventI {
	e.delay = d
	return e
}

// Selector specifies a CSS filter that is used to check for bubbled events. This allows the event to be fired from
// child controls.
func (e *Event) Selector (s string) EventI {
	e.selector = s
	return e
}

// Call Blocking to cause This event to prevent other events from firing after This fires, but before it processes.
// This is particularly useful to debounce button clicks. (The infamous double-click of the submit button when processing financial transactions for example).
func (e *Event) Blocking () EventI {
	e.blocking = true
	return e
}

// Call Terminating to cause the event not to bubble or do the default action.
func (e *Event) Terminating() EventI {
	e.actionsMustTerminate = true
	return e
}

// ActionValue is a value that will be returned to the actions that will be process by This event. Specify a static
// value, or javascript objects that will gather data at the time the event fires.
// Example: ActionValue(javascript.VarName("ui")) will return the "ui" variable that is part of the event call.
func (e *Event) ActionValue(r interface{}) EventI {
	e.actionValue = r
	return e
}

func (e *Event) GetActionValue() interface{} {
	return e.actionValue
}

// Validate overrides the controls validation setting just for This event.
func (e *Event) Validate(v ValidationType) EventI {
	e.validationOverride = v
	return e
}

// ValidationTargets overrides the control's validation targets just for This event.
func (e *Event) ValidationTargets(targets... string) EventI {
	e.validationTargetsOverride = targets
	return e
}



func (e *Event) AddActions(actions... ActionI) {
	var foundCallback bool
	for _,action := range actions {
		if _,ok := action.(CallbackActionI); ok {
			if foundCallback {
				panic ("You can only associate one callback action with an event, and it must be the last action.")
			}
			foundCallback = true
		} else {
			if foundCallback {
				panic ("You can only associate one callback action with an event, and it must be the last action.")
			}
		}
	}

	// Note, the above could be more robust and allow multiple callback actions, but it would get quite tricky if different
	// kinds of actions were interleaved. We will wait until someone presents a compelling need for something like that.
	e.actions = append (e.actions, actions...)
}

func (e *Event) RenderActions(control ControlI, eventId EventId) string {
	if e.actions == nil {
		return ""
	}

	var js string

	if e.actionsMustTerminate {
		js = "event.preventDefault();\n"
	}

	var params = RenderParams{control, eventId, e.actionValue}

	for _,a := range e.actions {
		js += a.RenderScript(params)
	}

	if e.blocking {
		js += "goradd.blockEvents = true;\n"
	}

	if !config.Minify {
		js = html.Indent(js)
	}
	js = fmt.Sprintf("goradd.queueAction({f: $j.proxy(function(){\n%s\n},this), d: %d});\n", js, e.delay)

	if e.condition != "" {
		js = fmt.Sprintf("\nif (%s) {\n%s\n};", e.condition, js)
	}

	js = control.wrapEvent(e.JsEvent, e.selector, js)

	if !config.Minify {
		// Render a comment
		js = fmt.Sprintf("/*** Event: %s  Control Type: %T, Control Label: %s, Control Id: %s  ***/\n%s\n", e.JsEvent, control, control.Label(), control.Id(), js)
	}

	return js
}

func (e *Event) GetActions() []ActionI {
	return e.actions
}