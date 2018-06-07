package action

import (
	"fmt"
	"github.com/spekary/goradd"
	"github.com/spekary/goradd/javascript"
	"github.com/spekary/goradd/util/types"
	"strings"
)

// CallbackI defines actions that result in a callback to us. Specifically Server and Ajax actions are defined for now.
// Potential for WebSocket action.
type CallbackActionI interface {
	ActionI
	ActionValue(interface{}) CallbackActionI
	DestinationControlID(string) CallbackActionI
	Validator(interface{}) CallbackActionI
	Async() CallbackActionI

	ID() int
	GetDestinationControlID() string
	GetDestinationControlSubID() string
	GetActionValue() interface{}
}

// CallbackAction is a kind of superclass for Ajax and Server actions
type callbackAction struct {
	goradd.Base
	id                 int
	destControlId      string
	subId              string
	actionValue        interface{}
	validationOverride interface{} // overrides the validation setting that is on the control
	async              bool
}

func (a *callbackAction) This() CallbackActionI {
	return a.Self.(CallbackActionI)
}

func (a *callbackAction) ID() int {
	return a.id
}

// SetValue lets you set a value that will be available to the action handler as the GetActionValue in the ActionParam structure
// sent to the handler. This can be any go type, including slices and maps, or a javascript.JavaScripter interface type.
// javascript.Closures will be called immediately with a (this) parameter.
func (a *callbackAction) ActionValue(v interface{}) CallbackActionI {
	a.actionValue = v
	return a.This()
}

func (a *callbackAction) GetActionValue() interface{} {
	return a.actionValue
}

func (a *callbackAction) Validator(v interface{}) CallbackActionI {
	a.validationOverride = v
	return a.This()
}

func (a *callbackAction) Async() CallbackActionI {
	a.async = true
	return a.This()
}

// Assign the destination control id. You can specify a sub id which indicates that the action should be sent to something
// inside the main control by concatenating the controls id with another id that indicates the internal destination,
// separated with an underscore.
func (a *callbackAction) DestinationControlID(id string) CallbackActionI {
	parts := strings.SplitN(id, "_", 2)
	if len(parts) == 2 {
		a.destControlId = parts[0]
		a.subId = parts[1]
	} else {
		a.destControlId = id
	}
	return a.This()
}

func (a *callbackAction) GetDestinationControlID() string {
	return a.destControlId
}

func (a *callbackAction) GetDestinationControlSubID() string {
	return a.subId
}

func (a *callbackAction) RenderScript(params RenderParams) string {
	panic("You need to embed this action and implement RenderScript")
	return ""
}

type serverAction struct {
	callbackAction
}

func Server(destControlId string, actionId int) *serverAction {
	a := &serverAction{
		callbackAction{
			id: actionId,
		},
	}
	a.Self = a
	a.DestinationControlID(destControlId)
	return a
}

func (a *serverAction) RenderScript(params RenderParams) string {
	v := types.NewOrderedMap()
	v.Set("controlID", params.TriggeringControlID)
	v.Set("eventId", params.EventID)
	if a.async {
		v.Set("async", true)
	}

	if eV, aV, cV := params.EventActionValue, a.actionValue, params.ControlActionValue; eV != nil || aV != nil || cV != nil {
		v2 := types.NewOrderedMap()
		if eV != nil {
			v2.Set("event", eV)
		}
		if aV != nil {
			v2.Set("action", aV)
		}
		if cV != nil {
			v2.Set("control", cV)
		}
		v.Set("actionValues", v2)
	}
	return fmt.Sprintf("goradd.postBack(%s);\n", javascript.ToJavaScript(v))
}

type ajaxAction struct {
	callbackAction
}

func Ajax(destControlId string, actionId int) *ajaxAction {
	a := &ajaxAction{
		callbackAction{
			id: actionId,
		},
	}
	a.Self = a
	a.DestinationControlID(destControlId)
	return a
}

func (a *ajaxAction) RenderScript(params RenderParams) string {
	v := types.NewOrderedMap()
	v.Set("controlID", params.TriggeringControlID)
	v.Set("eventId", params.EventID)
	if a.async {
		v.Set("async", true)
	}

	if eV, aV, cV := params.EventActionValue, a.actionValue, params.ControlActionValue; eV != nil || aV != nil || cV != nil {
		v2 := types.NewOrderedMap()
		if eV != nil {
			v2.Set("event", eV)
		}
		if aV != nil {
			v2.Set("action", aV)
		}
		if cV != nil {
			v2.Set("control", cV)
		}
		v.Set("actionValues", v2)
	}
	return fmt.Sprintf("goradd.postAjax(%s);\n", javascript.ToJavaScript(v))
}
