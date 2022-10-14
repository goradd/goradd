package action

import (
	"encoding/gob"
	"fmt"
	"strings"

	"github.com/goradd/goradd/pkg/javascript"
	"github.com/goradd/maps"
)

// CallbackActionI defines actions that result in a callback to the server.
// Specifically Server and Ajax actions are defined for now.
// There is potential for a Message action, like through WebSocket, PubHub, etc.
type CallbackActionI interface {
	ActionI
	ActionValue(v interface{}) CallbackActionI
	Async() CallbackActionI
	DestinationControlID(id string) CallbackActionI
}

// G_CallbackActionI is used by the framework to access actions. Callback actions must satisfy
// this interface, as well as the CallbackActionI interface.
//
// This is separated here so that IDEs will not pick up these functions for framework users.
type G_CallbackActionI interface {
	// ID returns the id assigned to the action when the action was created.
	ID() int
	// GetDestinationControlID returns the id that the action was sent to.
	GetDestinationControlID() string
	// GetDestinationControlSubID returns the id of the sub-control that is the destination of the action, if one was assigned.
	GetDestinationControlSubID() string
	// GetActionValue returns the action value that was assigned to the action when the action was fired.
	GetActionValue() interface{}
	// IsServerAction returns true if this is a server action.
	IsServerAction() bool
}

// callbackAction is a kind of superclass for Ajax and Server actions.
type callbackAction struct {
	ActionID      int
	DestControlID string
	SubID         string
	Value         interface{}
	CallAsync     bool
}

// G_CB is a private helper for encoding callback actions
type G_CB = callbackAction

// ID returns the action id that was defined when the action was created.
func (a *callbackAction) ID() int {
	return a.ActionID
}

// GetActionValue returns the action value given to the action when it was created.
func (a *callbackAction) GetActionValue() interface{} {
	return a.Value
}

// Assign the destination control id. You can specify a sub id which indicates that the action should be sent to something
// inside the main control by concatenating the control's id with another id that indicates the internal destination,
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

// GetDestinationControlID returns the control that the action will operate on.
func (a *callbackAction) GetDestinationControlID() string {
	return a.DestControlID
}

// GetDestinationControlSubID returns the sub id so that a composite control can send the
// action to a sub control.
func (a *callbackAction) GetDestinationControlSubID() string {
	return a.SubID
}

type serverAction struct {
	G_CB
}

// Server creates a server action, which is an action that will use a POST submission mechanism to trigger the action.
// Generally, with modern browsers, server actions are not that useful, since they cause an entire page to reload, while
// Ajax actions do not, and therefore Ajax actions are quicker to process.
// However, there are special cases where a server action might be
// useful, like:
//   - You are moving to a new page anyway.
//   - You are having trouble making an Ajax action work for some reason, and a Server action might get around the problem.
//   - You are submitting a multipart form, like when uploading a file.
//
// When the action fires, the Action() function of the GoRADD control identified by the
// destControlId will be called, and the given actionID will be the ID passed in the Params of the call.
//
// You send an action to a sub-control  specify a sub id which indicates that the action should be sent to something
// inside the main control by concatenating the control's id with another id that indicates the internal destination,
// separated with an underscore.
//
// The returned action uses a Builder pattern to add options, so for example you might call:
//
//	myControl.On(event.Click(), action.Server("myControl", MyActionIdConst).ActionValue("myActionValue").Async())
func Server(destControlId string, actionId int) CallbackActionI {
	a := &serverAction{
		callbackAction{
			ActionID: actionId,
		},
	}
	a.DestinationControlID(destControlId)
	return a
}

// RenderScript is called by the framework to render the script as javascript.
func (a *serverAction) RenderScript(params RenderParams) string {
	v := new(maps.SliceMap[string, any])
	v.Set("controlID", params.TriggeringControlID)
	v.Set("eventId", params.EventID)
	if a.CallAsync {
		v.Set("async", true)
	}

	eV, aV, cV := params.EventActionValue, a.Value, params.ControlActionValue
	v2 := new(maps.SliceMap[string, any])
	if eV != nil {
		v2.Set("event", eV)
	} else {
		v2.Set("event", javascript.JsCode("eventData"))
	}
	if aV != nil {
		v2.Set("action", aV)
	}
	if cV != nil {
		v2.Set("control", cV)
	}
	v.Set("actionValues", v2)
	return fmt.Sprintf("goradd.postBack(%s);", javascript.ToJavaScript(v))
}

// ActionValue lets you set a value that will be available to the action handler as the ActionValue() function in the ActionParam structure
// sent to the event handler. This can be any go type, including slices and maps, or a javascript.JavaScripter interface type.
// javascript.Closures will be called immediately with a (this) parameter.
func (a *serverAction) ActionValue(v interface{}) CallbackActionI {
	a.Value = v
	return a
}

// Async will cause the action to be handled asynchronously. Use this only in special situations where you know that you
// do not need information from other actions.
func (a *serverAction) Async() CallbackActionI {
	a.CallAsync = true
	return a
}

// DestinationControlID sets the id of the control that will receive the action.
//
// Implementations of composite controls can specify a sub id which indicates that the action should be sent to
// an item inside a control by concatenating the control's id with another id that indicates the internal destination,
// separated with an underscore.
func (a *serverAction) DestinationControlID(id string) CallbackActionI {
	a.setDestinationControlID(id)
	return a
}

// IsServerAction returns true if this is a server action
func (a *serverAction) IsServerAction() bool {
	return true
}

/*
func (a *serverAction) GobEncode() (data []byte, err error) {
	var buf bytes.Buffer
	e := gob.NewEncoder(&buf)

	if err = a.callbackAction.encode(e); err != nil {
		return nil, err
	}
	data = buf.Bytes()
	return
}

func (a *serverAction) GobDecode(data []byte) (err error) {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	return a.callbackAction.decode(dec)
}*/

type ajaxAction struct {
	G_CB
}

// Ajax creates an ajax action. When the action fires, the Action() function of the GoRADD control identified by the
// destControlId will be called, and the given actionID will be the ID passed in the Params of the call.
// You can specify a sub id which indicates that the action should be sent to something
// inside the main control by concatenating the control's id with another id that indicates the internal destination,
// separated with an underscore.
//
// The returned action uses a Builder pattern to add options, so for example you might call:
//
//	myControl.On(event.Click(), action.Ajax("myControl", MyActionIdConst).ActionValue("myActionValue").Async())
func Ajax(destControlId string, actionID int) CallbackActionI {
	a := &ajaxAction{
		callbackAction{
			ActionID: actionID,
		},
	}
	a.DestinationControlID(destControlId)
	return a
}

// RenderScript renders the script as javascript.
func (a *ajaxAction) RenderScript(params RenderParams) string {
	v := new(maps.SliceMap[string, any])
	v.Set("controlID", params.TriggeringControlID)
	v.Set("eventId", params.EventID)
	if a.CallAsync {
		v.Set("async", true)
	}

	eV, aV, cV := params.EventActionValue, a.Value, params.ControlActionValue
	v2 := new(maps.SliceMap[string, any])
	if eV != nil {
		v2.Set("event", eV)
	} else {
		v2.Set("event", javascript.JsCode("eventData"))
	}
	if aV != nil {
		v2.Set("action", aV)
	}
	if cV != nil {
		v2.Set("control", cV)
	}
	v.Set("actionValues", v2)
	return fmt.Sprintf("goradd.postAjax(%s);", javascript.ToJavaScript(v))
}

// ActionValue lets you set a value that will be available to the action handler as the ActionValue() function in the ActionParam structure
// sent to the event handler. This can be any go type, including slices and maps, or a javascript.JavaScripter interface type.
// javascript.Closures will be called immediately with a (this) parameter.
func (a *ajaxAction) ActionValue(v interface{}) CallbackActionI {
	a.Value = v
	return a
}

// Async will cause the action to be handled asynchronously. Use this only in special situations where you know that you
// do not need information from other actions.
func (a *ajaxAction) Async() CallbackActionI {
	a.CallAsync = true
	return a
}

// DestinationControlID sets the id of the control that will receive the action.
// You can specify a sub id which indicates that the action should be sent to something
// inside the main control by concatenating the control's id with another id that indicates the internal destination,
// separated with an underscore.
func (a *ajaxAction) DestinationControlID(id string) CallbackActionI {
	a.setDestinationControlID(id)
	return a
}

// IsServerAction will return false if this is not a server action.
func (a *ajaxAction) IsServerAction() bool {
	return false
}

/*
	func (a *ajaxAction) GobEncode() (data []byte, err error) {
		var buf bytes.Buffer
		e := gob.NewEncoder(&buf)

		if err = a.callbackAction.encode(e); err != nil {
			return nil, err
		}
		data = buf.Bytes()
		return
	}

	func (a *ajaxAction) GobDecode(data []byte) (err error) {
		buf := bytes.NewBuffer(data)
		dec := gob.NewDecoder(buf)
		return a.callbackAction.decode(dec)
	}
*/
func init() {
	// Register actions so they can be serialized
	gob.Register(&ajaxAction{})
	gob.Register(&serverAction{})
}
