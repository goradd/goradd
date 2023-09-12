package action

import (
	"encoding/gob"
	"fmt"
	"strings"

	"github.com/goradd/goradd/pkg/javascript"
	"github.com/goradd/maps"
)

// CallbackActionI defines actions that result in a callback to the server.
// Specifically Post and Ajax actions are defined for now.
// There is potential for a Message action, like through WebSocket, PubHub, etc.
type CallbackActionI interface {
	ActionI
	// ActionValue lets you set a value that will be available to the action handler as the ActionValue() function in the ActionParam structure
	// sent to the event handler. This can be any go type, including slices and maps, or a javascript.JavaScripter interface type.
	// javascript.Closures will be called immediately with a (this) parameter.
	ActionValue(v interface{}) CallbackActionI
	// Async will cause the action to be handled asynchronously. Use this only in special situations where you know that you
	// do not need information from other actions.
	Async() CallbackActionI
	// DestinationControlID sets the id of the control that will receive the action.
	//
	// Composite controls can specify a sub id which indicates that the action should be sent to
	// a component inside the main control by concatenating the control's id with another id that
	// indicates the internal destination, separated with an underscore.
	DestinationControlID(id string) CallbackActionI
	// Post turns the action into a post action
	Post() CallbackActionI
}

// CallbackActionAccessor is used by the framework and the internals of controls to access actions.
// Callback actions must satisfy
// this interface, as well as the CallbackActionI interface.
type CallbackActionAccessor interface {
	// ID returns the id assigned to the action when the action was created.
	ID() int
	// GetDestinationControlID returns the id that the action was sent to.
	GetDestinationControlID() string
	// GetDestinationControlSubID returns the id of the sub-control that is the destination of the action, if one was assigned.
	GetDestinationControlSubID() string
	// IsServerAction returns true if the action is a server action
	IsServerAction() bool
}

// callbackAction handles Ajax and Post actions.
type callbackAction struct {
	ActionID      int
	DestControlID string
	SubID         string
	Value         interface{}
	CallAsync     bool
	IsPost        bool // if false, this is an Ajax call
}

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

// RenderScript renders the script as javascript.
func (a *callbackAction) RenderScript(params RenderParams) string {
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

	if a.IsPost {
		return fmt.Sprintf("goradd.postBack(%s);", javascript.ToJavaScript(v))
	} else {
		return fmt.Sprintf("goradd.postAjax(%s);", javascript.ToJavaScript(v))
	}
}

// ActionValue lets you set a value that will be available to the action handler as the ActionValue() function in the
// ActionParam structure sent to the event handler.
// This can be any go type, including slices and maps, or a javascript.JavaScripter interface type.
// javascript.Closures will be called immediately with a (this) parameter.
func (a *callbackAction) ActionValue(v interface{}) CallbackActionI {
	a.Value = v
	return a
}

// Async will cause the action to be handled asynchronously. Use this only in special situations where you know that you
// do not need information from other actions.
func (a *callbackAction) Async() CallbackActionI {
	a.CallAsync = true
	return a
}

// DestinationControlID sets the id of the control that will receive the action.
// You can specify a sub id which indicates that the action should be sent to something
// inside the main control by concatenating the control's id with another id that indicates the internal destination,
// separated with an underscore.
func (a *callbackAction) DestinationControlID(id string) CallbackActionI {
	a.setDestinationControlID(id)
	return a
}

// Post specifies that the action will be invoked by calling and http POST on the form.
//
// It is unlikely you would need to call this, but there are special cases where a server action might be
// useful, like:
//   - You are moving to a new page anyway.
//   - You are having trouble making an Ajax action work for some reason, and a Post action might get around the problem.
//   - You are submitting a multipart form, like when uploading a file.
func (a *callbackAction) Post() CallbackActionI {
	a.IsPost = true
	return a
}

// IsServerAction will return false if this is not a server action.
// This is used by the test harness.
func (a *callbackAction) IsServerAction() bool {
	return a.IsPost
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

// Ajax creates an ajax action. When the action fires, the DoAction() function of the GoRADD control identified by the
// destControlId will be called, and the given actionID will be the ID passed in the Params of the call.
// You can specify a sub id which indicates that the action should be sent to something
// inside the main control by concatenating the control's id with another id that indicates the internal destination,
// separated with an underscore.
//
// The returned action uses a Builder pattern to add options, so for example you might call:
//
//	myControl.On(event.Click(), action.Do("myControl", MyActionIdConst).ActionValue("myActionValue").Async())
//
// Deprecated: Use Do() instead.
func Ajax(destControlId string, actionID int) CallbackActionI {
	return Do(destControlId, actionID)
}

// Server returns a Server based action that is invoked with a Post http call.
// Deprecated: Use Do instead.
func Server(destControlId string, actionID int) CallbackActionI {
	return Do(destControlId, actionID).Post()
}

// Do returns an action that will invoke the DoAction function on the control with id specified by
// destControlId.
//
// The DoAction function will receive the destControlId as the ControlId value
// in the params parameter of the DoAction function. The given actionID can be accessed
// in the DoAction function by calling one of the ActionValue* functions on the params value of DoAction.
//
// You can specify a sub id which indicates that the action should be sent to something
// inside the main control by concatenating the control's id with another id that indicates the internal destination,
// separated with an underscore.
//
// The returned action uses a Builder pattern to add options, so for example you might call:
//
//	myControl.On(event.Click(), action.Do("myControl", MyActionIdConst).ActionValue("myActionValue").Async())
//
// By default, the action uses an ajax process defined in the goradd.js file. If you want to instead use
// a more traditional http POST process, call the Post() function on the returned action, but that would
// rarely be needed.
func Do(destControlId string, actionID int) CallbackActionI {
	a := &callbackAction{
		ActionID: actionID,
	}
	a.setDestinationControlID(destControlId)
	return a
}

func init() {
	// Register actions so they can be serialized
	gob.Register(&callbackAction{})
}
