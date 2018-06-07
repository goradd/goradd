package page

import (
	"encoding/json"
	"fmt"
	"github.com/spekary/goradd/javascript"
	"github.com/spekary/goradd/util/types"
	"strings"
)

const (
	ResponseWatcher        = "watcher"
	ResponseControls       = "controls"
	ResponseCommandsHigh   = "commandsHigh"
	ResponseCommandsMedium = "commands"
	ResponseCommandsLow    = "commandsLow"
	ResponseCommandsFinal  = "commandsFinal"
	ResponseRegC           = "regc" // register control list
	ResponseHtml           = "html"
	ResponseValue          = "value"
	ResponseId             = "id"
	ResponseAttributes     = "attributes"
	ResponseCss            = "css"
	ResponseClose          = "winclose"
	ResponseLocation       = "loc"
	ResponseAlert          = "alert"
	ResponseStyleSheets    = "ss"
	ResponseJavaScripts    = "js"
)

type Priority int

const (
	PriorityExclusive Priority = iota
	PriorityHigh
	PriorityStandard
	PriorityLow
	PriorityFinal
)

// ResponseCommand is a response packet that leads to execution of a javascript function
type ResponseCommand struct {
	script   string // if just straight javascript
	selector string
	function string
	args     []interface{}
	final    bool
}

func (r ResponseCommand) MarshalJSON() (buf []byte, err error) {
	var reply = map[string]interface{}{}

	if r.script != "" {
		reply["script"] = r.script
	} else if r.selector != "" {
		reply["selector"] = r.selector
		reply["func"] = r.function
		reply["params"] = javascript.Arguments(r.args)
		if r.final {
			reply["final"] = true
		}
	} else {
		reply["func"] = r.function
		reply["params"] = javascript.Arguments(r.args)
		if r.final {
			reply["final"] = true
		}
	}
	return json.Marshal(reply)
}

// A response packet that leads to the manipulation or replacement of an html object
type ResponseControl struct {
	id         string
	html       string            // replaces the entire control's html
	attributes map[string]string // replace only specific attributes of the control
	value      string            // call the jQuery .val function with This value
}

func (r ResponseControl) MarshalJSON() (buf []byte, err error) {
	var reply = map[string]interface{}{}

	if r.html != "" {
		reply["html"] = r.html
	} else if r.attributes != nil {
		reply["attributes"] = r.attributes
	} else {
		reply["value"] = r.value
	}

	return json.Marshal(reply)
}

type Response struct {
	exclusiveCommand       *ResponseCommand
	highPriorityCommands   []ResponseCommand
	mediumPriorityCommands []ResponseCommand
	lowPriorityCommands    []ResponseCommand
	finalCommands          []ResponseCommand
	jsFiles                *types.OrderedStringMap
	alerts                 []string
	styleSheets            *types.OrderedStringMap
	newLocation            string
	winClose               bool
	controls               map[string]ResponseControl
}

func NewResponse() Response {
	return Response{}
}

func (r *Response) displayAlert(message string) {
	r.alerts = append(r.alerts, message)
}

// ExecuteJavaScript will execute the given code with the given priority. Note that all javascript code is run in
// strict mode.
func (r *Response) ExecuteJavaScript(js string, priority Priority) {
	switch priority {
	case PriorityExclusive:
		r.exclusiveCommand = &ResponseCommand{script: js}
	case PriorityHigh:
		r.highPriorityCommands = append(r.highPriorityCommands, ResponseCommand{script: js})
	case PriorityStandard:
		r.mediumPriorityCommands = append(r.mediumPriorityCommands, ResponseCommand{script: js})
	case PriorityLow:
		r.lowPriorityCommands = append(r.lowPriorityCommands, ResponseCommand{script: js})
	case PriorityFinal:
		r.finalCommands = append(r.finalCommands, ResponseCommand{script: js})
	}
}

func (r *Response) ExecuteControlCommand(controlID string, functionName string, priority Priority, args ...interface{}) {
	r.ExecuteSelectorFunction("#"+controlID, functionName, priority, args...)
}

// Calls a function on a jQuery selector
func (r *Response) ExecuteSelectorFunction(selector string, functionName string, priority Priority, args ...interface{}) {
	c := ResponseCommand{selector: selector, function: functionName, args: args}

	switch priority {
	case PriorityExclusive:
		r.exclusiveCommand = &c
	case PriorityHigh:
		r.highPriorityCommands = append(r.highPriorityCommands, c)
	case PriorityStandard:
		r.mediumPriorityCommands = append(r.mediumPriorityCommands, c)
	case PriorityLow:
		r.lowPriorityCommands = append(r.lowPriorityCommands, c)
	case PriorityFinal:
		c.final = true
		r.finalCommands = append(r.finalCommands, c)
	}

}

// Call the given function with the given arguments. If just a function label, then the window object is searched.
// The function can be inside an object accessible from the global namespace by separating with periods.
func (r *Response) ExecuteJsFunction(functionName string, priority Priority, args ...interface{}) {
	c := ResponseCommand{function: functionName, args: args}

	switch priority {
	case PriorityExclusive:
		r.exclusiveCommand = &c
	case PriorityHigh:
		r.highPriorityCommands = append(r.highPriorityCommands, c)
	case PriorityStandard:
		r.mediumPriorityCommands = append(r.mediumPriorityCommands, c)
	case PriorityLow:
		r.lowPriorityCommands = append(r.lowPriorityCommands, c)
	case PriorityFinal:
		c.final = true
		r.finalCommands = append(r.finalCommands, c)
	}
}

// One time add of style sheets, to be used by Form only for last minute style sheet injection.
func (r *Response) addStyleSheets(styleSheets ...string) {
	if r.styleSheets == nil {
		r.styleSheets = types.NewOrderedStringMap()
	}
	for _, s := range styleSheets {
		r.styleSheets.Set(s, s)
	}
}

// Add javascript files to the response.
func (r *Response) addJavaScriptFiles(files ...string) {
	if r.jsFiles == nil {
		r.jsFiles = types.NewOrderedStringMap()
	}
	for _, f := range files {
		r.jsFiles.Set(f, f)
	}
}

/**
 * Function renders all the Javascript commands as output to the client browser. This is a mirror of what
 * occurs in the success function in the qcubed.js ajax code.
 *
 * @param bool $blnBeforeControls True to only render the javascripts that need to come before the controls are defined.
 *                                This is used to break the commands issued into two groups.
 * @static
 * @return string
 */
func (r *Response) JavaScript() (script string) {
	// Style sheet injection by a control. Not very common, as other ways of adding style sheets would normally be done first.
	if r.styleSheets != nil {
		for _, s := range r.styleSheets.Keys() {
			script += `goradd.loadStyleSheetFile("` + s + `", "all);\n"`
		}
		r.styleSheets = nil
	}

	// alerts
	if r.alerts != nil {
		for _, a := range r.alerts {
			b, err := json.Marshal(a)
			if err != nil {
				panic(err)
			}
			script += fmt.Sprintf("goradd.msg(%s);\n", b[:])
		}
		r.alerts = nil
	}

	if r.highPriorityCommands != nil {
		script += r.renderCommandArray(r.highPriorityCommands)
		r.highPriorityCommands = nil
	}

	if r.mediumPriorityCommands != nil {
		script += r.renderCommandArray(r.mediumPriorityCommands)
		r.mediumPriorityCommands = nil
	}
	if r.lowPriorityCommands != nil {
		script += r.renderCommandArray(r.lowPriorityCommands)
		r.lowPriorityCommands = nil
	}

	// A redirect
	if r.newLocation != "" {
		script += fmt.Sprintf(`goradd.redirect(%s);`+"\n", r.newLocation)
		r.newLocation = ""
	}

	// A window close
	if r.winClose {
		script += "window.close();\n"
		r.winClose = false
	}

	return script
}

/**
 * @param array $commandArray
 * @return string
 */
func (r *Response) renderCommandArray(commands []ResponseCommand) string {
	var script string
	for _, command := range commands {
		if command.script != "" {
			script += command.script + ";\n"
		} else if command.selector != "" {
			if command.function == "" {
				panic("Cannot process a selector without a function")
			}
			var args string

			if command.args != nil {
				args = javascript.Arguments(command.args).JavaScript()
			}
			script += fmt.Sprintf("jQuery(%s).%s(%s);\n", command.selector, command.function, args)
		} else if command.function != "" {
			var args string
			if command.args != nil {
				args = javascript.Arguments(command.args).JavaScript()
			}
			script += fmt.Sprintf("%s(%s);\n", command.function, args)
		}
	}

	return script
}

// Return the JSON for use by the form ajax response. Will essentially do the same thing as
// above, but working in cooperation with the javascript file to process these through an ajax response.
func (r *Response) MarshalJSON() (buf []byte, err error) {
	var reply = map[string]interface{}{}

	if r.exclusiveCommand != nil {
		// only render This one;
		reply[ResponseCommandsMedium] = []ResponseCommand{*r.exclusiveCommand}
		r.exclusiveCommand = nil
	} else {
		var commands []ResponseCommand
		if r.highPriorityCommands != nil {
			commands = append(commands, r.highPriorityCommands...)
			r.highPriorityCommands = nil
		}
		if r.mediumPriorityCommands != nil {
			commands = append(commands, r.mediumPriorityCommands...)
			r.mediumPriorityCommands = nil
		}
		if r.lowPriorityCommands != nil {
			commands = append(commands, r.lowPriorityCommands...)
			r.lowPriorityCommands = nil
		}
		if r.finalCommands != nil {
			commands = append(commands, r.finalCommands...)
			r.finalCommands = nil
		}
		reply["commands"] = commands

		if r.jsFiles != nil {
			reply[ResponseJavaScripts] = strings.Join(r.jsFiles.Values(), ",")
		}

		if r.styleSheets != nil {
			reply[ResponseStyleSheets] = strings.Join(r.styleSheets.Values(), ",")
		}

		// alerts
		if r.alerts != nil {
			reply[ResponseAlert] = r.alerts
		}

		if r.controls != nil {
			reply[ResponseControls] = r.controls
		}

		if r.newLocation != "" {
			reply[ResponseLocation] = r.newLocation
		}

		if r.winClose {
			reply[ResponseClose] = 1
		}
	}

	return json.Marshal(reply)
}

func (r *Response) SetLocation(newLocation string) {
	r.newLocation = newLocation
}

func (r *Response) CloseWindow() {
	r.winClose = true
}

func (r *Response) hasExclusiveCommand() bool {
	return r.exclusiveCommand != nil
}

func (r *Response) SetControlHtml(id string, html string) {
	if r.controls == nil {
		r.controls = map[string]ResponseControl{}
	}
	if v, ok := r.controls[id]; ok && v.html != "" {
		panic("Setting ajax html twice on same control: " + id)
	}
	r.controls[id] = ResponseControl{html: html}
}

func (r *Response) SetControlAttribute(id string, attribute string, value string) {
	if r.controls == nil {
		r.controls = map[string]ResponseControl{}
	}
	if v, ok := r.controls[id]; ok {
		if v.html != "" {
			return // whole control is being redrawn so ignore individual attribute changes
		}
		if v.attributes != nil {
			v.attributes[attribute] = value
		} else {
			v.attributes = map[string]string{attribute: value}
		}
	} else {
		r.controls[id] = ResponseControl{attributes: map[string]string{attribute: value}}
	}
}

func (r *Response) SetControlValue(id string, value string) {
	if r.controls == nil {
		r.controls = map[string]ResponseControl{}
	}
	r.controls[id] = ResponseControl{value: value}
}
