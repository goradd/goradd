package page

// ActionI is an interface that defines actions that can be triggered by events
type ActionI interface {
	// RenderScript returns the action's javascript
	RenderScript(params RenderParams) string
}

// ActionsParams are sent to the Action() function in controls in response to a user action.
type ActionParams struct {
	// Id is the id assigned when the action is created
	Id int
	// Action is an interface to the action itself
	Action ActionI
	// GetActionValue is the value set by the Value() function when the action was created.
	ActionValue interface{}
	// JsValues are the values coming from javascript, and possibly further altered by the control. If a javascript object, it will be a string map to interface objects.
	JsValues interface{}
	// OriginalParams  are the original javascript parameters. Controls can alter the parameters coming from javascript.
	OriginalJsValues interface{}
	// ControlId is the control that originated the action
	ControlId string
	// Form is the form that sent the action, in case your action handler responds to multiple forms.
	Form FormI
}

// RenderParams are sent to the action rendering function.
type RenderParams struct {
	TriggeringControl ControlI
	EventId EventId
	EventActionValue interface{}
}

// CallbackI defines actions that result in a callback to us. Specifically Server and Ajax actions are defined for now.
// Potential for WebSocket action.
type CallbackActionI interface {
	ActionI
	ActionValue(interface{})      CallbackActionI
	DestinationControlId(string)  CallbackActionI
	Validator(interface{})        CallbackActionI
	Async()                       CallbackActionI
}
