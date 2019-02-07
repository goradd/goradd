package page

import (
	"encoding/json"
	"github.com/goradd/goradd/pkg/page/action"
	"strconv"
	"strings"
)

// ActionParams are sent to the Action() function in controls in response to a user action.
type ActionParams struct {
	// Id is the id assigned when the action is created
	ID int
	// Action is an interface to the action itself
	Action action.ActionI
	// ControlID is the control that originated the action
	ControlId string

	// values, to be accessesed with the Event*, Action* and Control* helper functions
	values actionValues
}


// actionValues is the structure representing the values sent in an ActionParam.
// Use the helper functions to get to the values.
type actionValues struct {
	Event   json.RawMessage `json:"event"`
	Control json.RawMessage `json:"control"`
	Action  json.RawMessage `json:"action"`
}


// EventValue will attempt to put the Event value into the given object using json.Unmarshal.
// You should primarily use it to get object or array values out of the Action value. If you are
// expecting a basic type, use one of the EventValue* helper functions instead.
// The given value should be a pointer to an object or variable that is the expected type for the data.
// ok will be false if no data was found.
// It will return an error if the data was found, but could not be converted to the given type.
func (a *ActionParams) EventValue(i interface{}) (ok bool, err error) {
	if a.values.Event != nil {
		ok = true
		err = json.Unmarshal(a.values.Event, i)
	}
	return
}

// ActionValue will attempt to put the Action value into the given object using json.Unmarshal.
// You should primarily use it to get object or array values out of the Action value. If you are
// expecting a basic type, use one of the ActionValue* helper functions instead.
// The given value should be a pointer to an object or variable that is the expected type for the data.
// ok will be false if no data was found.
// It will return an error if the data was found, but could not be converted to the given type.
func (a *ActionParams) ActionValue(i interface{}) (ok bool, err error) {
	if a.values.Action != nil {
		ok = true
		err = json.Unmarshal(a.values.Action, i)
	}
	return
}

// ControlValue will attempt to put the Control value into the given object using json.Unmarshal.
// You should primarily use it to get object or array values out of the Control value. If you are
// expecting a basic type, use one of the ControlValue* helper functions instead.
// The given value should be a pointer to an object or variable that is the expected type for the data.
// ok will be false if no data was found.
// It will return an error if the data was found, but could not be converted to the given type.
func (a *ActionParams) ControlValue(i interface{}) (ok bool, err error) {
	if a.values.Control != nil {
		ok = true
		err = json.Unmarshal(a.values.Control, i)
	}
	return
}

// EventValueString returns the event value as a string. It will convert to a string, even if the value
// is not a string.
func (a *ActionParams) EventValueString() string {
	return string(a.values.Event)
}

// EventValueInt returns the event value as an integer. If the value was a floating point value at the client,
// it will be truncated to an integer. If the value is not numeric, will return 0.
func (a *ActionParams) EventValueInt() int {
	return int(a.EventValueFloat())	// json is always sent as float
}

// EventValueFloat returns the event value as a float64. If the value was not numeric, it will return 0.
func (a *ActionParams) EventValueFloat() float64 {
	f,_ := strconv.ParseFloat(a.EventValueString(), 64)
	return f
}

// EventValueBool returns the event value as a bool. To be false, the value should be a boolean false,
// a numeric 0, an empty string, null or undefined on the client side. Otherwise, will return true.
func (a *ActionParams) EventValueBool() bool {
	return actionValueToBool(a.EventValueString())
}

// ActionValueString returns the action value as a string. It will convert to a string, even if the value
// is not a string.
func (a *ActionParams) ActionValueString() string {
	return string(a.values.Action)
}

// ActionValueInt returns the action value as an integer.
func (a *ActionParams) ActionValueInt() int {
	return int(a.ActionValueFloat())	// json is always sent as float
}

// ActionValueFloat returns the action value as a float64.
func (a *ActionParams) ActionValueFloat() float64 {
	f,_ := strconv.ParseFloat(a.ActionValueString(), 64)
	return f
}

// ActionValueBool returns the action value as a bool.
func (a *ActionParams) ActionValueBool() bool {
	return actionValueToBool(a.ActionValueString())
}

// ControlValueString returns the control value as a string. It will convert to a string, even if the value
// is not a string.
func (a *ActionParams) ControlValueString() string {
	return string(a.values.Control)
}

// ControlValueInt returns the control value as an int.
func (a *ActionParams) ControlValueInt() int {
	return int(a.ControlValueFloat())	// json is always sent as float
}

// ControlValueFloat returns the control value as a float64.
func (a *ActionParams) ControlValueFloat() float64 {
	f,_ := strconv.ParseFloat(a.ControlValueString(), 64)
	return f
}

// ControlValueBool returns the control value as a bool.
func (a *ActionParams) ControlValueBool() (ret bool) {
	return actionValueToBool(a.ControlValueString())
}

func actionValueToBool(v string) bool {
	v = strings.TrimSpace(v)
	if v == "" ||
		v == "0" ||
		v == "false" {
			return false
	} else {
		return true
	}
}

