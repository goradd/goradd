package action

import "encoding/gob"

type javascriptAction struct {
	JavaScript string
}

// Javascript will execute the given javascript
func Javascript(js string) ActionI {
	if js != "" {
		if js[len(js)-1:] != ";" {
			js += ";"
		}
	}
	return javascriptAction{JavaScript: js}
}

// RenderScript is called by the framework to output the action as JavaScript.
func (a javascriptAction) RenderScript(params RenderParams) string {
	return a.JavaScript
}

func init() {
	// Register actions so they can be serialized
	gob.Register(javascriptAction{})
}
