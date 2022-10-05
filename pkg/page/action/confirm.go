package action

import (
	"encoding/gob"
	"fmt"
	"github.com/goradd/goradd/pkg/javascript"
)

type confirmAction struct {
	Message interface{}
}

// Confirm will put up a standard browser confirmation dialog box, and will cancel any following actions if the
// user does not agree.
func Confirm(m interface{}) ActionI {
	return confirmAction{Message: m}
}

// RenderScript is called by the framework to output the action as JavaScript.
func (a confirmAction) RenderScript(params RenderParams) string {
	return fmt.Sprintf("if (!window.confirm(%s)) return false;", javascript.ToJavaScript(a.Message))
}

func init() {
	gob.Register(confirmAction{})
}
