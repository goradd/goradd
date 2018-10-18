package page

import (
	"encoding/json"
	"github.com/spekary/goradd/log"
	"github.com/spekary/goradd/page/action"
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
	Event   json.RawMessage `json:"event"`
	Control json.RawMessage `json:"control"`
	Action  json.RawMessage `json:"action"`
}

// String returns the action's value as a string. If you are expecting values from
// more than one place, you should call the more specific helper function. If will log a warning if more than one
// value comes through. The precedence is Control over Action over Event.
func (a ActionValues) String() (ret string) {
	var count int

	if a.Event != nil {
		a.EventValue(&ret)
		count++
	}
	if a.Action != nil {
		a.ActionValue(&ret)
		count++
	}
	if a.Control != nil {
		a.ControlValue(&ret)
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
		a.EventValue(&ret)
		count++
	}
	if a.Action != nil {
		a.ActionValue(&ret)
		count++
	}
	if a.Control != nil {
		a.ControlValue(&ret)
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
		a.EventValue(&ret)
		count++
	}
	if a.Action != nil {
		a.ActionValue(&ret)
		count++
	}
	if a.Control != nil {
		a.ControlValue(&ret)
		count++
	}
	if count > 1 {
		log.Warning("The action had more than one value. Call the specific action helper.")
	}
	return ret
}

// Float returns the action's value as a bool. If you are expecting values from
// more than one place, you should call the more specific helper function. If will log a warning if more than one
// value comes through. The precedence is Control over Action over Event.
func (a ActionValues) Bool() (ret bool) {
	var count int

	if a.Event != nil {
		a.EventValue(&ret)
		count++
	}
	if a.Action != nil {
		a.ActionValue(&ret)
		count++
	}
	if a.Control != nil {
		a.ControlValue(&ret)
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
func (a ActionValues) Value(ret interface{}) (ok bool, err error) {
	var count int

	if a.Event != nil {
		ok,err = a.EventValue(ret)
		count++
	}
	if a.Action != nil {
		ok,err = a.ActionValue(ret)
		count++
	}
	if a.Control != nil {
		ok,err = a.ControlValue(ret)
		count++
	}
	if count > 1 {
		log.Warning("The action had more than one value. Call the specific action helper.")
	}
	return
}

// EventValue will attempt to put the Event value into the given object. The given value should be a pointer
// to an object or variable that is the expected type for the data. ok will be false if no data was found.
// It will return an error if the data was found, but could not be converted to the given type.
func (a ActionValues) EventValue(i interface{}) (ok bool, err error) {
	if a.Event != nil {
		ok = true
		err = json.Unmarshal(a.Event, i)
	}
	return
}

// ActionValue will attempt to put the Action value into the given object. The given value should be a pointer
// to an object or variable that is the expected type for the data. ok will be false if no data was found.
// It will return an error if the data was found, but could not be converted to the given type.
func (a ActionValues) ActionValue(i interface{}) (ok bool, err error) {
	if a.Action != nil {
		ok = true
		err = json.Unmarshal(a.Action, i)
	}
	return
}

// ControlValue will attempt to put the Control value into the given object. The given value should be a pointer
// to an object or variable that is the expected type for the data. ok will be false if no data was found.
// It will return an error if the data was found, but could not be converted to the given type.
func (a ActionValues) ControlValue(i interface{}) (ok bool, err error) {
	if a.Control != nil {
		ok = true
		err = json.Unmarshal(a.Control, i)
	}
	return
}


func (a ActionValues) EventString() (ret string) {
	a.EventValue(&ret)
	return
}

func (a ActionValues) EventInt() (ret int) {
	a.EventValue(&ret)
	return
}

func (a ActionValues) EventFloat() (ret float64) {
	a.EventValue(&ret)
	return
}

func (a ActionValues) EventBool() (ret bool) {
	a.EventValue(&ret)
	return
}


func (a ActionValues) ActionString() (ret string) {
	a.ActionValue(&ret)
	return
}

func (a ActionValues) ActionInt() (ret int) {
	a.ActionValue(&ret)
	return
}

func (a ActionValues) ActionFloat() (ret float64) {
	a.ActionValue(&ret)
	return
}

func (a ActionValues) ActionBool() (ret bool) {
	a.ActionValue(&ret)
	return
}


func (a ActionValues) ControlString() (ret string) {
	a.ControlValue(&ret)
	return
}

func (a ActionValues) ControlInt() (ret int) {
	a.ControlValue(&ret)
	return
}

func (a ActionValues) ControlFloat() (ret float64) {
	a.ControlValue(&ret)
	return
}

func (a ActionValues) ControlBool() (ret bool) {
	a.ControlValue(&ret)
	return
}


