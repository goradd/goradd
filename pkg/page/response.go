package page

import (
	"encoding/json"
	"fmt"
	"github.com/goradd/gengen/pkg/maps"
	"github.com/goradd/goradd/pkg/javascript"
	"strings"
	"sync"
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

// Priority orders the various responses to an Ajax request so that the framework can control the order they are processed,
// and not necessarily order the responses in the order they are sent.
type Priority int

const (
	PriorityExclusive Priority = iota
	PriorityHigh
	PriorityStandard
	PriorityLow
	PriorityFinal	// TODO: Note that this currently requires a preliminary ajax command, or it will not fire. Should fix that, but its tricky.
)

// responseCommand is a response packet that leads to execution of a javascript function
type responseCommand struct {
	script   string // if just straight javascript
	selector string
	function string
	args     []interface{}
	final    bool
}

// MarshalJSON is used to form the Ajax response.
func (r responseCommand) MarshalJSON() (buf []byte, err error) {
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

// responseControl is the response packet that leads to the manipulation or replacement of an html object
type responseControl struct {
	id         string
	html       string            // replaces the entire control's html
	attributes map[string]string // replace only specific attributes of the control
	value      string            // call the jQuery .val function with This value
}

// MarshalJSON is used to form the Ajax response.
func (r responseControl) MarshalJSON() (buf []byte, err error) {
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

// Response contains the various commands you can send to the client in response to a goradd event.
// These commands are packed as JSON (for an Ajax response) or JavaScript (for a Server response),
// sent to the client, unpacked by JavaScript code in the goradd.js file, and then acted upon.
type Response struct {
	sync.RWMutex // This was inserted here for very rare situations of simultaneous access, like in the test harness.

	// exclusiveCommand is a single command that is sent by itself, overriding all other commands
	exclusiveCommand       *responseCommand
	// highPriorityCommands are sent first
	highPriorityCommands   []responseCommand
	// mediumPriorityCommands are sent after high priority commands
	mediumPriorityCommands []responseCommand
	// lowPriorityCommands are sent after medium priority commands
	lowPriorityCommands    []responseCommand
	// finalCommands are acted on after all other commands have been processed
	finalCommands          []responseCommand
	// jsFiles are JavaScript files that should be inserted into the page. This should rarely be used,
	// but is needed in case the programmer inserts a control widget in response to an Ajax event,
	// and that control depends on javascript that has not yet been sent to the client.
	jsFiles                *maps.StringSliceMap
	// styleSheets are css files that should be inserted into the page.
	styleSheets            *maps.StringSliceMap
	// alerts are strings that should be shown to the user in a javascript aler
	alerts                 []string
	// newLocation is a URL that the client should be redirected to.
	newLocation            string
	// winClose directs the browser to close the current window.
	winClose               bool
	// controls are goraddControls that should be inserted or replaced
	controls               map[string]responseControl
	// profileHtml is the html sent from the database profiling tool to display in a special window
	// TODO: This is not used currently, and is here for future ajax db profiling
	profileHtml			   string
}

// NewResponse creates a new event response.
func NewResponse() Response {
	return Response{}
}

func (r *Response) displayAlert(message string) {
	r.Lock()
	r.alerts = append(r.alerts, message)
	r.Unlock()
}

// ExecuteJavaScript will execute the given code with the given priority. Note that all javascript code is run in
// strict mode.
func (r *Response) ExecuteJavaScript(js string, priorities ...Priority) {
	var priority = PriorityStandard
	if priorities != nil {
		if len(priorities) == 1 {
			priority = priorities[0]
		} else {
			panic("Don't call ExecuteJavaScript with arguments")
		}
	}
	r.Lock()
	switch priority {
	case PriorityExclusive:
		r.exclusiveCommand = &responseCommand{script: js}
	case PriorityHigh:
		r.highPriorityCommands = append(r.highPriorityCommands, responseCommand{script: js})
	case PriorityStandard:
		r.mediumPriorityCommands = append(r.mediumPriorityCommands, responseCommand{script: js})
	case PriorityLow:
		r.lowPriorityCommands = append(r.lowPriorityCommands, responseCommand{script: js})
	case PriorityFinal:
		r.finalCommands = append(r.finalCommands, responseCommand{script: js})
	}
	r.Unlock()
}

// ExecuteControlCommand executes the named command on the given goradd control.
func (r *Response) ExecuteControlCommand(controlID string, functionName string, args ...interface{}) {
	r.ExecuteSelectorFunction("#"+controlID, functionName, args...)
}

// ExecuteSelectorFunction calls a function on a jQuery selector
func (r *Response) ExecuteSelectorFunction(selector string, functionName string, args ...interface{}) {
	args2,priority := r.extractPriority(args...)
	c := responseCommand{selector: selector, function: functionName, args: args2}

	r.Lock()
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
	r.Unlock()
}

// ExecuteJsFunction calls the given JavaScript function with the given arguments.
// If the function name has a dot(.) in it, the items preceeding the dot will be considered global objects
// to call the function on. If the named function just a function label, then the function is called on the window object.
func (r *Response) ExecuteJsFunction(functionName string, args ...interface{}) {
	args2,priority := r.extractPriority(args...)
	c := responseCommand{function: functionName, args: args2}

	r.Lock()
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
	r.Unlock()
}

func (r *Response) extractPriority (args ...interface{}) (args2 []interface{}, priority Priority) {
	for i,a := range args {
		if p,ok := a.(Priority); ok {
			priority = p
			args2 = append(args[:i], args[i+1:]...)
			return
		}
	}
	priority = PriorityStandard
	args2 = args
	return
}

// One time add of style sheets, to be used by FormBase only for last minute style sheet injection.
func (r *Response) addStyleSheets(styleSheets ...string) {
	if r.styleSheets == nil {
		r.styleSheets = maps.NewStringSliceMap()
	}
	for _, s := range styleSheets {
		r.styleSheets.Set(s, s)
	}
}

// Add javascript files to the response.
func (r *Response) addJavaScriptFiles(files ...string) {
	if r.jsFiles == nil {
		r.jsFiles = maps.NewStringSliceMap()
	}
	for _, f := range files {
		r.jsFiles.Set(f, f)
	}
}

// JavaScript renders the Response object as JavaScript that will be inserted into the page sent back to the
// client in response to a Server action.
func (r *Response) JavaScript() (script string) {
	r.Lock()
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
		script += fmt.Sprintf(`goradd.redirect("%s");`+"\n", r.newLocation)
		r.newLocation = ""
	}

	// A window close
	if r.winClose {
		script += "window.close();\n"
		r.winClose = false
	}
	r.Unlock()
	return script
}


func (r *Response) renderCommandArray(commands []responseCommand) string {
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
			script += fmt.Sprintf("jQuery('%s').%s(%s);\n", command.selector, command.function, args)
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

// GetAjaxResponse returns the JSON for use by the form ajax response.
// It will also reset the response
func (r *Response) GetAjaxResponse() (buf []byte, err error) {
	var reply = map[string]interface{}{}

	r.Lock()

	if r.exclusiveCommand != nil {
		// only render This one;
		reply[ResponseCommandsMedium] = []responseCommand{*r.exclusiveCommand}
		r.exclusiveCommand = nil
	} else {
		var commands []responseCommand
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

		if commands != nil && len(commands) > 0 {
			reply["commands"] = commands
		}

		if r.jsFiles != nil {
			reply[ResponseJavaScripts] = strings.Join(r.jsFiles.Values(), ",")
			r.jsFiles = nil
		}

		if r.styleSheets != nil {
			reply[ResponseStyleSheets] = strings.Join(r.styleSheets.Values(), ",")
			r.styleSheets = nil
		}

		// alerts
		if r.alerts != nil {
			reply[ResponseAlert] = r.alerts
			r.alerts = nil
		}

		if r.controls != nil {
			reply[ResponseControls] = r.controls
			r.controls = nil
		}

		if r.newLocation != "" {
			reply[ResponseLocation] = r.newLocation
			r.newLocation = ""
		}

		if r.winClose {
			reply[ResponseClose] = 1
			r.winClose = false
		}
	}

	r.Unlock()
	return json.Marshal(reply)
}


// Call SetLocation to change the url of the browser.
func (r *Response) SetLocation(newLocation string) {
	r.Lock()
	r.newLocation = newLocation
	r.Unlock()
}

// Call CloseWindow to close the current window.
func (r *Response) CloseWindow() {
	r.Lock()
	r.winClose = true
	r.Unlock()
}

func (r *Response) hasExclusiveCommand() bool {
	r.RLock()
	v := r.exclusiveCommand != nil
	r.RUnlock()
	return v
}

// SetControlHtml will cause the given control's html to be completely replaced by the given HTML.
func (r *Response) SetControlHtml(id string, html string) {
	r.Lock()
	if r.controls == nil {
		r.controls = map[string]responseControl{}
	}
	if v, ok := r.controls[id]; ok && v.html != "" {
		r.Unlock()
		panic("Setting ajax html twice on same control: " + id)
	}
	r.controls[id] = responseControl{html: html}
	r.Unlock()
}

// SetControlAttribute sets the named html attribute on the control to the given value.
func (r *Response) SetControlAttribute(id string, attribute string, value string) {
	r.Lock()
	if r.controls == nil {
		r.controls = map[string]responseControl{}
	}
	if v, ok := r.controls[id]; ok {
		if v.html == "" { // only do attributes if whole control is not being redrawn
			if v.attributes != nil {
				v.attributes[attribute] = value
			} else {
				v.attributes = map[string]string{attribute: value}
			}
		}
	} else {
		r.controls[id] = responseControl{attributes: map[string]string{attribute: value}}
	}
	r.Unlock()
}

// SetControlValue calls the jQuery ".val()" function on the given control, passing it the given value.
func (r *Response) SetControlValue(id string, value string) {
	r.Lock()
	if r.controls == nil {
		r.controls = map[string]responseControl{}
	}
	r.controls[id] = responseControl{value: value}
	r.Unlock()
}


func (r *Response) setProfileInfo(info string) {
	r.Lock()
	r.profileHtml = info
	r.Unlock()
}

