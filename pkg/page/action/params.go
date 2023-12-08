package action

import (
	"encoding/json"
	"strconv"
	"strings"
)

// Params are sent to the control.DoAction() function in response to a user action.
// Each action might contain one of three kinds of values:
//
//   - Action value, a value assigned to the action when it is created,
//   - Control value, a value returned by certain custom controls
//   - Event value, a value assigned by certain event types
//
// Use the accessor functions like [Params.ActionValueInt] or [Params.ControlValueString] to extract an appropriately
// typed value from the action.
type Params struct {
	// Event is the javascript name of the event that was triggered.
	Event string
	// ID is the action id assigned when the action was created
	ID int
	// Action is an interface to the action itself
	Action ActionI
	// ControlID is the control that originated the action
	ControlId string

	// values, to be accessed with the helper functions
	values RawActionValues
}

// RawActionValues is the structure representing the json values sent in an ActionParam.
// Use the helper functions to get to the values.
type RawActionValues struct {
	Event   json.RawMessage `json:"event"`
	Control json.RawMessage `json:"control"`
	Action  json.RawMessage `json:"action"`
}

// NewActionParams is used internally by the framework to create new action parameters.
// You should not normally need to call this function.
func NewActionParams(event string, id int, action ActionI, controlId string, rawValues RawActionValues) Params {
	return Params{
		Event:     event,
		ID:        id,
		Action:    action,
		ControlId: controlId,
		values:    rawValues,
	}
}

// EventValue will attempt to put the Event value into the given object using json.Unmarshal.
// You should primarily use it to get object or array values out of the Action value. If you are
// expecting a basic type, use one of the EventValue* helper functions instead.
// The given value should be a pointer to an object or variable that is the expected type for the data.
// ok will be false if no data was found.
// It will return an error if the data was found, but could not be converted to the given type.
func (a *Params) EventValue(i interface{}) (ok bool, err error) {
	if a.values.Event != nil {
		ok = true
		err = json.Unmarshal(a.values.Event, i)
	}
	return
}

// ActionValue will return the action value into the given object using json.Unmarshal.
// You should primarily use it to get object or array values out of the action value. If you are
// expecting a basic type, use one of the ActionValue* helper functions instead.
// The given value should be a pointer to an object or variable that is the expected type for the data.
// ok will be false if no data was found.
// It will return an error if the data was found, but could not be converted to the given type.
func (a *Params) ActionValue(i interface{}) (ok bool, err error) {
	if a.values.Action != nil {
		ok = true
		err = json.Unmarshal(a.values.Action, i)
	}
	return
}

// ControlValue will return the control value into the given object using json.Unmarshal.
// You should primarily use it to get object or array values out of the control value. If you are
// expecting a basic type, use one of the ControlValue* helper functions instead.
// The given value should be a pointer to an object or variable that is the expected type for the data.
// ok will be false if no data was found.
// It will return an error if the data was found, but could not be converted to the given type.
func (a *Params) ControlValue(i interface{}) (ok bool, err error) {
	if a.values.Control != nil {
		ok = true
		err = json.Unmarshal(a.values.Control, i)
	}
	return
}

// EventValueString returns the event value as a string. It will convert to a string, even if the value
// is not a string.
func (a *Params) EventValueString() string {
	v := string(a.values.Event)
	if len(v) > 1 && v[0] == '"' && v[len(v)-1] == '"' {
		// It is surrounded by quotes, so remove the quotes
		v = v[1 : len(v)-1]
	}
	return string(v)
}

// EventValueInt returns the event value as an integer. If the value was a floating point value at the client,
// it will be truncated to an integer. If the value is not numeric, will return 0.
func (a *Params) EventValueInt() int {
	return int(a.EventValueFloat()) // json is always sent as float
}

// EventValueFloat returns the event value as a float64. If the value was not numeric, it will return 0.
func (a *Params) EventValueFloat() float64 {
	f, _ := strconv.ParseFloat(a.EventValueString(), 64)
	return f
}

// EventValueBool returns the event value as a bool. To be false, the value should be a boolean false,
// a numeric 0, an empty string, null or undefined on the client side. Otherwise, will return true.
func (a *Params) EventValueBool() bool {
	return actionValueToBool(a.EventValueString())
}

// EventValueStringMap returns the event value as a map[string]string.
func (a *Params) EventValueStringMap() (m map[string]string) {
	_, _ = a.EventValue(&m)
	return m
}

// ActionValueString returns the action value as a string. It will convert to a string, even if the value
// is not a string.
func (a *Params) ActionValueString() string {
	v := string(a.values.Action)
	if len(v) > 1 && v[0] == '"' && v[len(v)-1] == '"' {
		// It is surrounded by quotes, so remove the quotes
		v = v[1 : len(v)-1]
	}
	return string(v)
}

// ActionValueInt returns the action value as an integer.
func (a *Params) ActionValueInt() int {
	return int(a.ActionValueFloat()) // json is always sent as float
}

// ActionValueFloat returns the action value as a float64.
func (a *Params) ActionValueFloat() float64 {
	f, _ := strconv.ParseFloat(a.ActionValueString(), 64)
	return f
}

// ActionValueBool returns the action value as a bool.
func (a *Params) ActionValueBool() bool {
	return actionValueToBool(a.ActionValueString())
}

// ControlValueString returns the control value as a string. It will convert to a string, even if the value
// is not a string.
func (a *Params) ControlValueString() string {
	v := string(a.values.Control)
	if len(v) > 1 && v[0] == '"' && v[len(v)-1] == '"' {
		// It is surrounded by quotes, so remove the quotes
		v = v[1 : len(v)-1]
	}
	return string(v)
}

// ControlValueInt returns the control value as an int.
func (a *Params) ControlValueInt() int {
	return int(a.ControlValueFloat()) // json is always sent as float
}

// ControlValueFloat returns the control value as a float64.
func (a *Params) ControlValueFloat() float64 {
	f, _ := strconv.ParseFloat(a.ControlValueString(), 64)
	return f
}

// ControlValueBool returns the control value as a bool.
func (a *Params) ControlValueBool() (ret bool) {
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
