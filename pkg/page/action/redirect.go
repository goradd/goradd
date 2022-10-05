package action

import (
	"encoding/gob"
	"fmt"
)

type redirectAction struct {
	Location string
}

// Redirect will navigate to the given page.
// TODO: If javascript is turned off, this should still work. We would need to detect the presence of javascript,
// and then emit a server action instead
func Redirect(url string) ActionI {
	return redirectAction{Location: url}
}

// RenderScript is called by the framework to output the action as JavaScript.
func (a redirectAction) RenderScript(params RenderParams) string {
	return fmt.Sprintf(`goradd.redirect("%s");`, a.Location)
}

func init() {
	gob.Register(redirectAction{})
}
