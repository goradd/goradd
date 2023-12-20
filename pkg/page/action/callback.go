package action

import (
	"encoding/gob"
	"fmt"
	"strings"

	"github.com/goradd/goradd/pkg/javascript"
	"github.com/goradd/maps"
)

// DefaultActionId is the default action id for events that do not want to specify a custom action id
const DefaultActionId = 0

// RefreshActionId is the action id that causes a control to refresh.
const RefreshActionId = -1

// DefaultControlId indicates that the control that the action is first sent to is the same as the control that received the event.
const DefaultControlId = ""

// CallbackActionI defines actions that result in a callback to the server through a Post or Ajax call.
//   - ID sets an optional action id, which can be retrieved in the DoAction function.
//   - ActionValue sets the value that will be available to the DoAction handler as the ActionValue() function in the Param structure.
//     This can be any go type, including slices and maps, or a javascript.JavaScripter interface type.
//     javascript.Closures will be called immediately with a (this) parameter.
//   - Async will cause the action to be handled asynchronously. Use this only in special situations where you know that you
//     do not need information from other actions.
//   - ControlID sets the id of the control that will receive the action.
//     Composite controls can specify a sub id which indicates that the action should be sent to
//     a component inside the main control by concatenating the control's id with another id that
//     indicates the internal destination, separated with an underscore.
//     By default, the control will be the control receiving the event that triggered the action.
//   - Post turns the action into a form post action, vs. an Ajax action. You only need this in special situations, when
//     the standard Ajax action will not work.
type CallbackActionI interface {
	ActionI
	ID(id int) CallbackActionI
	ActionValue(v interface{}) CallbackActionI
	Async() CallbackActionI
	ControlID(id string) CallbackActionI
	Post() CallbackActionI
}

// CallbackActionAccessor is used by the framework and the internals of controls to access actions.
// Callback actions must satisfy this interface, as well as the CallbackActionI interface.
type CallbackActionAccessor interface {
	GetActionID() int
	// GetDestinationControlID returns the id that the action was sent to.
	GetDestinationControlID() string
	// GetDestinationControlSubID returns the id of the sub-control that is the destination of the action, if one was assigned.
	GetDestinationControlSubID() string
	// IsPostAction returns true if the action is a server action
	IsPostAction() bool
}

// callbackAction handles Ajax and Post actions through the CallbackActionI interface.
type callbackAction struct {
	ActionID      int
	DestControlID string
	SubID         string
	Value         any
	CallAsync     bool
	IsPost        bool // if false, this is an Ajax call
}

// ID sets the action id of the action, and that can be retrieved from the action parameters in DoAction.
// This value should be positive. Negative values and zero are reserved for goradd internals.
func (a *callbackAction) ID(id int) CallbackActionI {
	a.ActionID = id
	return a
}

// GetActionID returns the action id given to the action when it was created.
func (a *callbackAction) GetActionID() int {
	return a.ActionID
}

// GetActionValue returns the action value given to the action when it was created.
func (a *callbackAction) GetActionValue() any {
	return a.ActionValue
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

// ActionValue sets a value that will be available to the action handler as the ActionValue() function in the
// Params structure sent to the event handler.
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

// ControlID sets the id of the control that will receive the action.
// You can specify a sub id which indicates that the action should be sent to something
// inside the main control by concatenating the control's id with another id that indicates the internal destination,
// separated with an underscore.
//
// By default, this will be the control that triggers the action.
func (a *callbackAction) ControlID(id string) CallbackActionI {
	a.setDestinationControlID(id)
	return a
}

// Post specifies that the action will be invoked by calling an http POST on the form.
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

// IsPostAction will return true if this is a post action.
// This is used by the test harness.
func (a *callbackAction) IsPostAction() bool {
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
//	myControl.On(event.Click(), action.Do("myControl", MyActionIdConst).EventValue("myActionValue").Async())
//
// Deprecated: Use Do() instead.
func Ajax(destControlId string, actionID int) CallbackActionI {
	return Do().ControlID(destControlId).ID(actionID)
}

// Server returns a Server based action that is invoked with a Post http call.
// Deprecated: Use Do instead.
func Server(destControlId string, actionID int) CallbackActionI {
	return Do().ControlID(destControlId).ID(actionID).Post()
}

// Do returns an action that will invoke the DoAction function on the destination control.
//
// By default, the control that will receive the action will be the control that triggered the action.
//
// The returned action uses a Builder pattern to add options, so for example you might call:
//
//	myControl.On(event.Click(), action.Do("myControl", MyActionIdConst).EventValue("myActionValue").Async())
//
// By default, the action uses an ajax process defined in the goradd.js file. If you want to instead use
// a more traditional http POST process, call the Post() function on the returned action, but that would
// rarely be needed.
func Do() CallbackActionI {
	a := &callbackAction{}
	return a
}

// Refresh will cause the given control to redraw
func Refresh(id string) CallbackActionI {
	return Do().ControlID(id).ID(RefreshActionId)
}

func init() {
	// Register actions so they can be serialized
	gob.Register(&callbackAction{})
}
