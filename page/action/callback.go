package action

import (
	"fmt"
	"github.com/spekary/goradd/javascript"
	"github.com/spekary/goradd/page"
	"github.com/spekary/goradd/util/types"
)

// CallbackAction is a kind of superclass for Ajax and Server actions
type callbackAction struct {
	id int
	destControlId string
	actionValue interface{}
	validationOverride interface{}	// overrides the validation setting that is on the control
	async bool
}

// SetValue lets you set a value that will be available to the action handler as the GetActionValue in the ActionParam structure
// sent to the handler. This can be any go type, including slices and maps, or a javascript.JavaScripter interface type.
// javascript.Closures will be called immediately with a (this) parameter.
func (a *callbackAction) ActionValue(v interface{}) *callbackAction {
	a.actionValue = v
	return a
}

func (a *callbackAction) Validator(v interface{}) *callbackAction {
	a.validationOverride = v
	return a
}

func (a *callbackAction) Async() *callbackAction {
	a.async = true
	return a
}

func (a *callbackAction) DestinationControlId(id string) *callbackAction {
	a.destControlId = id
	return a
}

type serverAction struct {
	callbackAction
}

func Server(destControlId string, actionId int) *serverAction {
	return &serverAction{
		callbackAction{
			id: actionId,
			destControlId: destControlId,
		},
	}
}

func (a *serverAction) RenderScript(params page.RenderParams) string {
	v := types.NewOrderedMap()
	v.Set("controlId", a.destControlId)
	v.Set("eventId", params.EventId)
	if a.async {
		v.Set("async", true)
	}

	if eV,aV,cV := params.EventActionValue, a.actionValue, params.TriggeringControl.ActionValue(); eV != nil || aV != nil || cV != nil {
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
	return fmt.Sprintf(`goradd.postBack(%s)`, javascript.ToJavaScript(v))
}



type ajaxAction struct {
	callbackAction
}

func Ajax(destControlId string, actionId int) *ajaxAction {
	return &ajaxAction{
		callbackAction{
			id: actionId,
			destControlId: destControlId,
		},
	}
}

func (a *ajaxAction) RenderScript(params page.RenderParams) string {
	v := types.NewOrderedMap()
	v.Set("controlId", a.destControlId)
	v.Set("eventId", params.EventId)
	if a.async {
		v.Set("async", true)
	}

	if eV,aV,cV := params.EventActionValue, a.actionValue, params.TriggeringControl.ActionValue(); eV != nil || aV != nil || cV != nil {
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
	return fmt.Sprintf(`goradd.postAjax(%s)`, javascript.ToJavaScript(v))
}
