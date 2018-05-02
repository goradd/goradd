package page

import "github.com/spekary/goradd/page/action"


// ActionParams are sent to the Action() function in controls in response to a user action.
type ActionParams struct {
	// Id is the id assigned when the action is created
	Id int
	// Action is an interface to the action itself
	Action action.ActionI
	// Values are the action values returned by the application. Note that if you are expecting arrays, maps or structs
	// to be returned that
	Values ActionValues
	// ControlID is the control that originated the action
	ControlId string
}
