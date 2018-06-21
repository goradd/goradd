package page

import (
	"github.com/spekary/goradd/page/action"
	"github.com/spekary/goradd/javascript"
	"github.com/spekary/goradd/log"
)

// ActionParams are sent to the Action() function in controls in response to a user action.
type ActionParams struct {
	// Id is the id assigned when the action is created
	ID int
	// Action is an interface to the action itself
	Action action.ActionI
	// Values are the action values returned by the application. Note that if you are expecting arrays, maps or structs
	// to be returned that
	Values ActionValues
	// ControlID is the control that originated the action
	ControlId string
}


// ActionValues is the structure representing the values sent in an ActionParam. Note that all numeric values are returned as
// a json.Number type, use the helper functions to extract the proper value. If you are returning a javascript object, it
// will come through as a string map, and you will have to call the javascript.Number* helper functions directly.
type ActionValues struct {
	Event   interface{} `json:"event"`
	Control interface{} `json:"control"`
	Action  interface{} `json:"action"`
}

// String returns the action's value as a string. If you are expecting values from
// more than one place, you should call the more specific helper function. If will log a warning if more than one
// value comes through. The precedence is Control over Action over Event.
func (a ActionValues) String() (ret string) {
	var count int

	if a.Event != nil {
		ret = a.EventString()
		count++
	}
	if a.Action != nil {
		ret = a.ActionString()
		count++
	}
	if a.Control != nil {
		ret = a.ControlString()
		count++
	}
	if count > 1 {
		log.Warning("The action had more than one value. Call the specific action helper.")
	}
	return ret
}

// Int returns the action's value as a int. If you are expecting values from
// more than one place, you should call the more specific helper function. If will log a warning if more than one
// value comes through. The precedence is Control over Action over Event.
func (a ActionValues) Int() (ret int) {
	var count int

	if a.Event != nil {
		ret = a.EventInt()
		count++
	}
	if a.Action != nil {
		ret = a.ActionInt()
		count++
	}
	if a.Control != nil {
		ret = a.ControlInt()
		count++
	}
	if count > 1 {
		log.Warning("The action had more than one value. Call the specific action helper.")
	}
	return ret
}

// Float returns the action's value as a float64. If you are expecting values from
// more than one place, you should call the more specific helper function. If will log a warning if more than one
// value comes through. The precedence is Control over Action over Event.
func (a ActionValues) Float() (ret float64) {
	var count int

	if a.Event != nil {
		ret = a.EventFloat()
		count++
	}
	if a.Action != nil {
		ret = a.ActionFloat()
		count++
	}
	if a.Control != nil {
		ret = a.ControlFloat()
		count++
	}
	if count > 1 {
		log.Warning("The action had more than one value. Call the specific action helper.")
	}
	return ret
}

// Value returns the action's value as an interface{}. If you are expecting values from
// more than one place, you should call the more specific helper function. If will log a warning if more than one
// value comes through. The precedence is Control over Action over Event.
func (a ActionValues) Value() (ret interface{}) {
	var count int

	if a.Event != nil {
		ret = a.Event
		count++
	}
	if a.Action != nil {
		ret = a.Action
		count++
	}
	if a.Control != nil {
		ret = a.Control
		count++
	}
	if count > 1 {
		log.Warning("The action had more than one value. Call the specific action helper.")
	}
	return ret
}

func (a ActionValues) EventString() string {
	return javascript.NumberString(a.Event)
}

func (a ActionValues) EventInt() int {
	return javascript.NumberInt(a.Event)
}

func (a ActionValues) EventFloat() float64 {
	return javascript.NumberFloat(a.Event)
}

func (a ActionValues) ActionString() string {
	return javascript.NumberString(a.Action)
}

func (a ActionValues) ActionInt() int {
	return javascript.NumberInt(a.Action)
}

func (a ActionValues) ActionFloat() float64 {
	return javascript.NumberFloat(a.Action)
}

func (a ActionValues) ControlString() string {
	return javascript.NumberString(a.Control)
}

func (a ActionValues) ControlInt() int {
	return javascript.NumberInt(a.Control)
}

func (a ActionValues) ControlFloat() float64 {
	return javascript.NumberFloat(a.Control)
}

