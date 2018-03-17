package page

// ActionI is an interface that defines actions that can be triggered by events
type ActionI interface {
	// RenderScript returns the action's javascript
	RenderScript(params RenderParams) string
}

// ActionParams are sent to the Action() function in controls in response to a user action.
type ActionParams struct {
	// Id is the id assigned when the action is created
	Id int
	// Action is an interface to the action itself
	Action ActionI
	// Values will contain one or more of the following:
	//  event: The event's action value, if one is provided.
	//  action: The action's action value, if one is provided.
	//  control: The control's action value, if one is provided.
	Values map[string]interface{}
	// ControlId is the control that originated the action
	ControlId string
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

	Id() 						  int
	GetDestinationControlId()	  string
}
