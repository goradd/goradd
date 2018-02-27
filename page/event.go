package page

import (
	"fmt"
	"goradd/config"
)

type EventI interface {
	Condition(string) EventI
	Delay(int) EventI
	Selector(string) EventI
	Blocking() EventI
	ActionValue(interface{}) EventI
	GetActionValue() interface{}
	AddActions(a... ActionI)
	RenderActions(control ControlI, eventId EventId) string
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
}

func (e *Event) Condition (c string) EventI {
	e.condition = c
	return e
}

func (e *Event) Delay (d int) EventI {
	e.delay = d
	return e
}

func (e *Event) Selector (s string) EventI {
	e.selector = s
	return e
}

func (e *Event) Blocking () EventI {
	e.blocking = true
	return e
}

func (e *Event) Terminating() EventI {
	e.actionsMustTerminate = true
	return e
}

func (e *Event) ActionValue(r interface{}) EventI {
	e.actionValue = r
	return e
}

func (e *Event) GetActionValue() interface{} {
	return e.actionValue
}

func (e *Event) AddActions(a... ActionI) {
	e.actions = append (e.actions, a...)
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

	if e.delay != 0 {
		js = fmt.Sprintf("goradd.setTimeout('%s', $j.proxy(function(){\n%s\n},this), %d);\n", control.Id(), js, e.delay)
	}

	if e.condition != "" {
		js = fmt.Sprintf("\nif (%s) {\n%s\n};", e.condition, js)
	}

	js = control.wrapEvent(e.JsEvent, js)

	if config.Mode == config.Dev {
		// Render a comment
		js = fmt.Sprintf("/*** Event: %s  Control Type: %T, Control Name: %s, Control Id: %s  ***/\n%s\n", e.JsEvent, control, control.Name(), control.Id(), js)
	}

	return js
}
