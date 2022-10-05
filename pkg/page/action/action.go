// Package action defines actions that you can trigger using events.
// Normally you would do this with the On() function that all goradd controls have.
//
// # Defining Your Own Actions
//
// You can define your own actions by creating a class that implements the ActionI interface, AND that is
// encodable by gob.Serialize, meaning it either implements the gob.Encoder interface or exports its structure, AND
// registers itself with gob.Register() so that the gob.Decoder knows how to deserialize it into an interface.
package action

// ActionI is an interface that defines actions that can be triggered by events
type ActionI interface {
	// RenderScript is called by the framework to return the action's javascript.
	RenderScript(params RenderParams) string
}

// RenderParams is used by the framework to give extra parameters that help with some actions.
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
