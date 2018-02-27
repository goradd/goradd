package action

import (
	"fmt"
	"github.com/spekary/goradd/javascript"
	"github.com/spekary/goradd/page"
)

type messageAction struct {
	message interface{}
}

// Message returns an action that will display a standard browser alert message. Specify a string, or one of the
// javascript.* types.
func Message(m interface{}) *messageAction {
	return &messageAction{message: m}
}

func (a *messageAction) RenderScript(params page.RenderParams) string {
	return fmt.Sprintf(`goradd.msg(%s)`, javascript.ToJavaScript(a.message))
}

type confirmAction struct {
	message interface{}
}


func Confirm(m interface{}) *confirmAction {
	return &confirmAction{message: m}
}

func (a *confirmAction) RenderScript(params page.RenderParams) string {
	return fmt.Sprintf(`if (!goradd.confirm(%s)) return false;`, javascript.ToJavaScript(a.message))
}


type blurAction struct {
	controlId string
}

func Blur(controlId string) *blurAction {
	return &blurAction{controlId: controlId}
}

func (a *blurAction) RenderScript(params page.RenderParams) string {
	return fmt.Sprintf(`goradd.blur('%s');`, a.controlId)
}

type focusAction struct {
	controlId string
}

func Focus(controlId string) *focusAction {
	return &focusAction{controlId: controlId}
}

func (a *focusAction) RenderScript(params page.RenderParams) string {
	return fmt.Sprintf(`goradd.focus('%s');`, a.controlId)
}

type selectAction struct {
	controlId string
}

func Select(controlId string) *selectAction {
	return &selectAction{controlId: controlId}
}

func (a *selectAction) RenderScript(params page.RenderParams) string {
	return fmt.Sprintf(`goradd.select('%s');`, a.controlId)
}


type cssPropertyAction struct {
	property string
	value interface{}
	controlId string
}

func SetCssProperty(controlId string, property string, value interface{}) *cssPropertyAction {
	return &cssPropertyAction{controlId: controlId}
}

func (a *cssPropertyAction) RenderScript(params page.RenderParams) string {
	return fmt.Sprintf(`goradd.css('%s', '%s', '%s');`, a.controlId, a.property, a.value)
}

type cssAddClassAction struct {
	classes string
	controlId string
}

func AddClass(controlId string, addClasses string) *cssAddClassAction {
	return &cssAddClassAction{controlId: controlId, classes:addClasses}
}

func (a *cssAddClassAction) RenderScript(params page.RenderParams) string {
	return fmt.Sprintf(`goradd.addClass('%s', '%s');`, a.controlId, a.classes)
}

type cssToggleClassAction struct {
	classes string
	controlId string
}

func ToggleClass(controlId string, classes string) *cssToggleClassAction {
	return &cssToggleClassAction{controlId: controlId, classes: classes}
}

func (a *cssToggleClassAction) RenderScript(params page.RenderParams) string {
	return fmt.Sprintf(`goradd.toggleClass('%s', '%s');`, a.controlId, a.classes)
}


type redirectAction struct {
	location string
}

func Redirect(l string) *redirectAction {
	return &redirectAction{location: l}
}

func (a *redirectAction) RenderScript(params page.RenderParams) string {
	return fmt.Sprintf(`goradd.redirect(%s)`, a.location)
}


