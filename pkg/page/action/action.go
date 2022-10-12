// Package action defines actions that you can trigger using events.
// Normally you would do this with the On() function that all GoRADD controls have.
//
// For example:
//
//	button := NewButton(p, "okButton").SetText("OK")
//	button.On(event.Click(), action.Message(javascript.JsCode("event.target.value") + "was clicked"))
//
// will create a Message action that responds to a button click.
//
// There are three types of actions:
//   - [Javascript] Actions
//   - [Server] Action
//   - [Ajax] Action
//
// Javascript Actions execute javascript code that is handled immediately by the browser. Goradd provides
// a number of standard actions, like [Redirect]. However, the [Javascript] action can execute any javascript.
//
// [Server] and [Ajax] actions call the control.Action() function of a GoRADD control when the event is triggered.
// You specify the control id of the receiving control, and an integer representing the action. The Action()
// function can then do whatever is needed.
//
// Server actions cause a page to completely refresh, whereas Ajax actions use the javascript
// XMLHttpRequest mechanism to pass data without a complete refresh. In response to an Ajax action, a
// control can call its Refresh() function to redraw just that control.
package action

// ActionI is an interface that defines actions that can be triggered by events
type ActionI interface {
	// RenderScript is called by the framework to return the action's javascript.
	RenderScript(params RenderParams) string
}

// RenderParams is used by the framework to give information that helps with rendering actions.
type RenderParams struct {
	// TriggeringControlID is the id of the control that triggered the action
	TriggeringControlID string
	// ControlActionValue is the control action value that will be received by the Action() function.
	ControlActionValue interface{}
	// EventID is the event that triggered the action
	EventID uint16
	// EventActionValue is the event action value that will be received by the Action() function.
	EventActionValue interface{}
}
