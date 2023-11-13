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
// There are two types of actions:
//   - Javascript Actions
//   - Callback Actions
//
// Javascript Actions execute javascript code that is handled immediately by the browser. Goradd provides
// a number of standard Javascript actions, like [Redirect]. Also, the [Javascript] action can execute any javascript.
//
// Callback Actions invoke the DoAction() function on GoRADD controls. There are currently two kinds of
// Callback Actions:
//   - Ajax and
//   - Post
//
// Create a callback action by calling action.Do. By default, this will create an Ajax action. Call Post
// on the resulting action to turn it into a Post action.
// You specify the control id of the receiving control, and an integer representing the action. The DoAction
// function of the control can then do whatever is needed on the server side.
//
// Ajax actions use the javascript XMLHttpRequest mechanism to pass data without a complete refresh.
// In response to an Ajax action, a
// control can call its Refresh() function to redraw just that control. If a Post action is working, but not
// an Ajax action, one thing to check is to make sure Refresh() is called if anything in a control changes.
//
// Post actions use the standard http Post method of html forms, and cause a page to completely refresh.
// This is the method used in the early days of the web, and while it still works, it isn't as efficient
// as Ajax actions. However, there are still specific times to use a Post action:
//   - When starting a file upload. Currently, you must use a Post action for file uploads.
//   - When debugging a problem with an Ajax action, it may be helpful to see if it works as a Post action.
//
// In addition to the action id, a Callback action can receive data that is sent with the action, data that is sent by
// the control, and also data that is sent by the event that triggered the action. This data can be static data
// when the action is created, or dynamic data that is gathered by javascript when the action is invoked.
// See the action.Params structure for a description of what is supplied to the DoAction() function.
//
// To execute multiple actions in response to an event, put the actions in an actionGroup. The actionGroup
// can have multiple javascript actions, but only one callback action.
package action

// ActionI is an interface that defines actions that can be triggered by events
type ActionI interface {
	// RenderScript is called by the framework to return the action's javascript.
	// You do not normally need to call this function.
	RenderScript(params RenderParams) string
}

// RenderParams is used by the framework to give information that helps with rendering actions.
type RenderParams struct {
	// TriggeringControlID is the id of the control that triggered the action
	TriggeringControlID string
	// ControlActionValue is the control action value that will be received by the DoAction() function.
	ControlActionValue interface{}
	// EventID is the event that triggered the action
	EventID uint16
	// EventActionValue is the event action value that will be received by the DoAction() function.
	EventActionValue interface{}
}
