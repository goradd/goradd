package action

import (
	"encoding/gob"
	"fmt"
	"github.com/spekary/gengen/maps"
	"github.com/spekary/goradd/pkg/javascript"
	"strings"
)

// CallbackI defines actions that result in a callback to us. Specifically Server and Ajax actions are defined for now.
// Potential for Message action, like through WebSocket, PubHub, etc.
type CallbackActionI interface {
	ActionI

	ID() int
	GetDestinationControlID() string
	GetDestinationControlSubID() string
	GetActionValue() interface{}
}

// CallbackAction is a kind of superclass for Ajax and Server actions
type callbackAction struct {
	ActionID           int
	DestControlID      string
	SubID              string
	Value        interface{}
	ValidationOverride interface{} // overrides the validation setting that is on the control
	CallAsync          bool
}

func (a *callbackAction) ID() int {
	return a.ActionID
}

func (a *callbackAction) GetActionValue() interface{} {
	return a.Value
}

// Assign the destination control id. You can specify a sub id which indicates that the action should be sent to something
// inside the main control by concatenating the controls id with another id that indicates the internal destination,
// separated with an underscore.
func (a *callbackAction) setDestinationControlID(id string) {
	parts := strings.SplitN(id, "_", 2)
	if len(parts) == 2 {
		a.DestControlID = parts[0]
		a.SubID = parts[1]
	} else {
		a.DestControlID = id
	}
}

func (a *callbackAction) GetDestinationControlID() string {
	return a.DestControlID
}

func (a *callbackAction) GetDestinationControlSubID() string {
	return a.SubID
}

func (a *callbackAction) RenderScript(params ΩrenderParams) string {
	panic("You need to embed this action and implement RenderScript")
	return ""
}

type ΩserverAction struct {
	callbackAction
}

// Server creates a server action, which is an action that will use a POST submission mechanism to trigger the action.
// Generally, with modern browsers, server actions are not that useful, since they cause an entire page to reload, while
// Ajax actions do not, and so are Ajax actions are quicker to process. However, there are special cases where a server action might be
// useful, like:
// - You are moving to a new page anyway.
// - You are having trouble making an Ajax action work for some reason, and a Server action might get around the problem.
// - You are submitting a multi-part form, like when uploading a file.
// When the action fires, the Action() function of the Goradd control identified by the
// destControlId will be called, and the given actionID will be the ID passed in the ActionParams of the call.
// You can specify a sub id which indicates that the action should be sent to something
// inside the main control by concatenating the controls id with another id that indicates the internal destination,
// separated with an underscore.
//
// The returned action uses a Builder pattern to add options, so for example you might call:
//   myControl.On(event.Click(), action.Server("myControl", MyActionIdConst).ActionValue("myActionValue").Async())
func Server(destControlId string, actionId int) *ΩserverAction {
	a := &ΩserverAction{
		callbackAction{
			ActionID: actionId,
		},
	}
	a.DestinationControlID(destControlId)
	return a
}

func (a *ΩserverAction) RenderScript(params ΩrenderParams) string {
	v := maps.NewSliceMap()
	v.Set("controlID", params.TriggeringControlID)
	v.Set("eventId", params.EventID)
	if a.CallAsync {
		v.Set("async", true)
	}

	if eV, aV, cV := params.EventActionValue, a.ActionValue, params.ControlActionValue; eV != nil || aV != nil || cV != nil {
		v2 := maps.NewSliceMap()
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

// ActionValue lets you set a value that will be available to the action handler as the ActionValue() function in the ActionParam structure
// sent to the event handler. This can be any go type, including slices and maps, or a javascript.JavaScripter interface type.
// javascript.Closures will be called immediately with a (this) parameter.
func (a *ΩserverAction) ActionValue(v interface{}) *ΩserverAction {
	a.Value = v
	return a
}

// Validator lets you override the validation setting for the control that the action is being sent to.
func (a *ΩserverAction) Validator(v interface{}) *ΩserverAction {
	a.ValidationOverride = v
	return a
}

// Aysnc will cause the action to be handled asynchronously. Use this only in special situations where you know that you
// do not need information from other actions.
func (a *ΩserverAction) Async() *ΩserverAction {
	a.CallAsync = true
	return a
}

// DestinationControlID sets the id of the control that will receive the action.
// You can specify a sub id which indicates that the action should be sent to something
// inside the main control by concatenating the controls id with another id that indicates the internal destination,
// separated with an underscore.
func (a *ΩserverAction) DestinationControlID(id string) *ΩserverAction {
	a.setDestinationControlID(id)
	return a
}

type ΩajaxAction struct {
	callbackAction
}


// Ajax creates an ajax action. When the action fires, the Action() function of the Goradd control identified by the
// destControlId will be called, and the given actionID will be the ID passed in the ActionParams of the call.
// You can specify a sub id which indicates that the action should be sent to something
// inside the main control by concatenating the controls id with another id that indicates the internal destination,
// separated with an underscore.
//
// The returned action uses a Builder pattern to add options, so for example you might call:
//   myControl.On(event.Click(), action.Ajax("myControl", MyActionIdConst).ActionValue("myActionValue").Async())
func Ajax(destControlId string, actionID int) *ΩajaxAction {
	a := &ΩajaxAction{
		callbackAction{
			ActionID: actionID,
		},
	}
	a.DestinationControlID(destControlId)
	return a
}

func (a *ΩajaxAction) RenderScript(params ΩrenderParams) string {
	v := maps.NewSliceMap()
	v.Set("controlID", params.TriggeringControlID)
	v.Set("eventId", params.EventID)
	if a.CallAsync {
		v.Set("async", true)
	}

	if eV, aV, cV := params.EventActionValue, a.ActionValue, params.ControlActionValue; eV != nil || aV != nil || cV != nil {
		v2 := maps.NewSliceMap()
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

// ActionValue lets you set a value that will be available to the action handler as the ActionValue() function in the ActionParam structure
// sent to the event handler. This can be any go type, including slices and maps, or a javascript.JavaScripter interface type.
// javascript.Closures will be called immediately with a (this) parameter.
func (a *ΩajaxAction) ActionValue(v interface{}) *ΩajaxAction {
	a.Value = v
	return a
}

// Validator lets you override the validation setting for the control that the action is being sent to.
func (a *ΩajaxAction) Validator(v interface{}) *ΩajaxAction {
	a.ValidationOverride = v
	return a
}

// Aysnc will cause the action to be handled asynchronously. Use this only in special situations where you know that you
// do not need information from other actions.
func (a *ΩajaxAction) Async() *ΩajaxAction {
	a.CallAsync = true
	return a
}

// DestinationControlID sets the id of the control that will receive the action.
// You can specify a sub id which indicates that the action should be sent to something
// inside the main control by concatenating the controls id with another id that indicates the internal destination,
// separated with an underscore.
func (a *ΩajaxAction) DestinationControlID(id string) *ΩajaxAction {
	a.setDestinationControlID(id)
	return a
}

func init() {
	// Register actions so they can be serialized
	gob.Register(&ΩajaxAction{})
	gob.Register(&ΩserverAction{})
}
