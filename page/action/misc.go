package action

import (
	"fmt"
	"github.com/spekary/goradd/javascript"
)

// ActionI is an interface that defines actions that can be triggered by events
type ActionI interface {
	// RenderScript returns the action's javascript
	RenderScript(params RenderParams) string
}


type RenderParams struct {
	TriggeringControlID string
	ControlActionValue  interface{}
	EventID             uint16
	EventActionValue    interface{}
}



type messageAction struct {
	message interface{}
}

// Note: actions currently depend on a javascript eval if they are introduced to a form during an ajax response.
// One way to fix that would be to register all javascript actions so that they get added to the form at drawing time,
// so that when an event gets attached during an ajax call, the resulting action is already in the browser.

// Message returns an action that will display a standard browser alert message. Specify a string, or one of the
// javascript.* types.
func Message(m interface{}) *messageAction {
	return &messageAction{message: m}
}

func (a *messageAction) RenderScript(params RenderParams) string {
	return fmt.Sprintf(`goradd.msg(%s)`, javascript.ToJavaScript(a.message))
}

type confirmAction struct {
	message interface{}
}


func Confirm(m interface{}) *confirmAction {
	return &confirmAction{message: m}
}

func (a *confirmAction) RenderScript(params RenderParams) string {
	return fmt.Sprintf(`if (!goradd.confirm(%s)) return false;`, javascript.ToJavaScript(a.message))
}


type blurAction struct {
	controlID string
}

func Blur(controlID string) *blurAction {
	return &blurAction{controlID: controlID}
}

func (a *blurAction) RenderScript(params RenderParams) string {
	return fmt.Sprintf(`goradd.blur('%s');`, a.controlID)
}

type focusAction struct {
	controlID string
}

func Focus(controlID string) *focusAction {
	return &focusAction{controlID: controlID}
}

func (a *focusAction) RenderScript(params RenderParams) string {
	return fmt.Sprintf(`goradd.focus('%s');`, a.controlID)
}

type selectAction struct {
	controlID string
}

func Select(controlID string) *selectAction {
	return &selectAction{controlID: controlID}
}

func (a *selectAction) RenderScript(params RenderParams) string {
	return fmt.Sprintf(`goradd.select('%s');`, a.controlID)
}


type cssPropertyAction struct {
	property  string
	value     interface{}
	controlID string
}

func SetCssProperty(controlID string, property string, value interface{}) *cssPropertyAction {
	return &cssPropertyAction{controlID: controlID}
}

func (a *cssPropertyAction) RenderScript(params RenderParams) string {
	return fmt.Sprintf(`goradd.css('%s', '%s', '%s');`, a.controlID, a.property, a.value)
}

type cssAddClassAction struct {
	classes string
	controlID string
}

func AddClass(controlID string, addClasses string) *cssAddClassAction {
	return &cssAddClassAction{controlID: controlID, classes:addClasses}
}

func (a *cssAddClassAction) RenderScript(params RenderParams) string {
	return fmt.Sprintf(`goradd.addClass('%s', '%s');`, a.controlID, a.classes)
}

type cssToggleClassAction struct {
	classes string
	controlID string
}

func ToggleClass(controlID string, classes string) *cssToggleClassAction {
	return &cssToggleClassAction{controlID: controlID, classes: classes}
}

func (a *cssToggleClassAction) RenderScript(params RenderParams) string {
	return fmt.Sprintf(`goradd.toggleClass('%s', '%s');`, a.controlID, a.classes)
}


type redirectAction struct {
	location string
}

func Redirect(l string) *redirectAction {
	return &redirectAction{location: l}
}

func (a *redirectAction) RenderScript(params RenderParams) string {
	return fmt.Sprintf(`goradd.redirect(%s)`, a.location)
}

// PrivateAction is used by control implementations to add a private action to a controls action list. Unless you are
// creating a control, you should not use this.
type PrivateAction struct{}
func (a PrivateAction) RenderScript(params RenderParams) string {
	return ""
}


