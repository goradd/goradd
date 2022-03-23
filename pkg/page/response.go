package page

import (
	"encoding/json"
	"fmt"
	"github.com/goradd/gengen/pkg/maps"
	"github.com/goradd/html5tag"
	"github.com/goradd/goradd/pkg/javascript"
	"github.com/goradd/html5tag"
	"github.com/goradd/maps"
	"sync"
)

const (
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
	PriorityFinal // TODO: Note that this currently requires a preliminary ajax command, or it will not fire. Should fix that, but its tricky.
)

// responseCommand is a response packet that leads to execution of a javascript function
type responseCommand struct {
	Script   string        `json:"script,omitempty"` // if just straight javascript
	Selector string        `json:"selector,omitempty"`
	Id       string        `json:"id,omitempty"`
	JqueryId string        `json:"jqueryId,omitempty"`
	Function string        `json:"func,omitempty"`
	Args     []interface{} `json:"params,omitempty"`
	Final    bool          `json:"final,omitempty"`
}

// responseControl is the response packet that leads to the manipulation or replacement of an html object
type responseControl struct {
	Html       string            `json:"html,omitempty"`       // replaces the entire control's html
	Attributes map[string]string `json:"attributes,omitempty"` // replace only specific attributes of the control
	Value      string            `json:"value,omitempty"`      // sets the control's value. See goradd.js val:
}

type attributeMap = maps.SliceMap[string, html5tag.Attributes]

// Response contains the various commands you can send to the client in response to a goradd event.
// These commands are packed as JSON (for an Ajax response) or JavaScript (for a Server response),
// sent to the client, unpacked by JavaScript code in the goradd.js file, and then acted upon.
type Response struct {
	sync.RWMutex // This was inserted here for very rare situations of simultaneous access, like in the test harness.

	// exclusiveCommand is a single command that is sent by itself, overriding all other commands
	exclusiveCommand *responseCommand
	// highPriorityCommands are sent first
	highPriorityCommands []*responseCommand
	// mediumPriorityCommands are sent after high priority commands
	mediumPriorityCommands []*responseCommand
	// lowPriorityCommands are sent after medium priority commands
	lowPriorityCommands []*responseCommand
	// finalCommands are acted on after all other commands have been processed
	finalCommands []*responseCommand
	// jsFiles are JavaScript files that should be inserted into the page. This should rarely be used,
	// but is needed in case the programmer inserts a control widget in response to an Ajax event,
	// and that control depends on javascript that has not yet been sent to the client.
	jsFiles *attributeMap // Use slicemap to preserve the order
	// styleSheets are css files that should be inserted into the page.
	styleSheets *attributeMap
	// alerts are strings that should be shown to the user in a javascript alert
	alerts []string
	// newLocation is a URL that the client should be redirected to.
	newLocation string
	// winClose directs the browser to close the current window.
	winClose bool
	// controls are goraddControls that should be inserted or replaced
	controls map[string]responseControl
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

func (r *Response) AddClass(id string, class string, priorities ...Priority) {
	r.ExecuteControlCommand(id, "class", "+"+class, priorities)
}

func (r *Response) RemoveClass(id string, class string, priorities ...Priority) {
	r.ExecuteControlCommand(id, "class", "-"+class, priorities)
}

func (r *Response) SetClass(id string, class string, priorities ...Priority) {
	r.ExecuteControlCommand(id, "class", class, priorities)
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
	c := responseCommand{Script: js}
	r.postCommand(&c, priority)

}

// ExecuteControlCommand executes the named command on the given control. Possible commands are defined
// by the goradd widget class in the javascript file.
func (r *Response) ExecuteControlCommand(controlID string, functionName string, args ...interface{}) {
	args2, priority := r.extractPriority(args...)
	c := responseCommand{Id: controlID, Function: functionName, Args: args2}
	r.postCommand(&c, priority)
}

// ExecuteJqueryCommand executes the named jquery command on the given jquery control.
func (r *Response) ExecuteJqueryCommand(controlID string, functionName string, args ...interface{}) {
	args2, priority := r.extractPriority(args...)
	c := responseCommand{JqueryId: controlID, Function: functionName, Args: args2}
	r.postCommand(&c, priority)
}

// ExecuteSelectorFunction calls a goradd function on a group of objects defined by a selector.
func (r *Response) ExecuteSelectorFunction(selector string, functionName string, args ...interface{}) {
	args2, priority := r.extractPriority(args...)
	c := responseCommand{Selector: selector, Function: functionName, Args: args2}

	r.postCommand(&c, priority)
}

// ExecuteJsFunction calls the given JavaScript function with the given arguments.
// If the function name has a dot(.) in it, the items preceeding the dot will be considered global objects
// to call the function on. If the named function just a function label, then the function is called on the window object.
func (r *Response) ExecuteJsFunction(functionName string, args ...interface{}) {
	args2, priority := r.extractPriority(args...)
	c := responseCommand{Function: functionName, Args: args2}
	r.postCommand(&c, priority)
}

func (r *Response) postCommand(c *responseCommand, priority Priority) {
	r.Lock()
	switch priority {
	case PriorityExclusive:
		r.exclusiveCommand = c
	case PriorityHigh:
		r.highPriorityCommands = append(r.highPriorityCommands, c)
	case PriorityStandard:
		r.mediumPriorityCommands = append(r.mediumPriorityCommands, c)
	case PriorityLow:
		r.lowPriorityCommands = append(r.lowPriorityCommands, c)
	case PriorityFinal:
		c.Final = true
		r.finalCommands = append(r.finalCommands, c)
	}
	r.Unlock()
}

func (r *Response) extractPriority(args ...interface{}) (args2 []interface{}, priority Priority) {
	for i, a := range args {
		if p, ok := a.(Priority); ok {
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
func (r *Response) addStyleSheet(path string, attributes html5tag.Attributes) {
	if r.styleSheets == nil {
		r.styleSheets = new(attributeMap)
	}
	r.styleSheets.Set(path, attributes)
}

// Add javascript files to the response.
func (r *Response) addJavaScriptFile(path string, attributes html5tag.Attributes) {
	if r.jsFiles == nil {
		r.jsFiles = new(attributeMap)
	}
	r.jsFiles.Set(path, attributes)
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

func (r *Response) renderCommandArray(commands []*responseCommand) string {
	var script string
	for _, command := range commands {
		if command.Script != "" {
			script += command.Script + ";\n"
		} else {
			com := make(map[string]interface{})

			if command.Selector != "" {
				com["selector"] = command.Selector
			}
			if command.Id != "" {
				com["id"] = command.Id
			}
			if command.JqueryId != "" {
				com["jqueryId"] = command.Id
			}

			if command.Function != "" {
				com["func"] = command.Function
			}

			if command.Args != nil {
				com["params"] = command.Args
			}
			script += fmt.Sprintf("goradd.processCommand(%s);\n", javascript.ToJavaScript(com))
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
		var commands []*responseCommand
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
			reply[ResponseJavaScripts] = r.jsFiles
			r.jsFiles = nil
		}

		if r.styleSheets != nil {
			reply[ResponseStyleSheets] = r.styleSheets
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
	if v, ok := r.controls[id]; ok && v.Html != "" {
		r.Unlock()
		panic("Setting ajax html twice on same control: " + id)
	}
	r.controls[id] = responseControl{Html: html}
	r.Unlock()
}

// SetControlAttribute sets the named html attribute on the control to the given value.
func (r *Response) SetControlAttribute(id string, attribute string, value string) {
	r.Lock()
	if r.controls == nil {
		r.controls = map[string]responseControl{}
	}
	if v, ok := r.controls[id]; ok {
		if v.Html == "" { // only do attributes if whole control is not being redrawn
			if v.Attributes != nil {
				v.Attributes[attribute] = value
			} else {
				v.Attributes = map[string]string{attribute: value}
			}
		}
	} else {
		r.controls[id] = responseControl{Attributes: map[string]string{attribute: value}}
	}
	r.Unlock()
}

// SetControlValue calls the ".val()" function on the given control, passing it the given value.
func (r *Response) SetControlValue(id string, value string) {
	r.Lock()
	if r.controls == nil {
		r.controls = map[string]responseControl{}
	}
	r.controls[id] = responseControl{Value: value}
	r.Unlock()
}

// use an encoder since some fields could be nil
type responseEncoded struct {
	ExclusiveCommand       *responseCommand
	HighPriorityCommands   []*responseCommand
	MediumPriorityCommands []*responseCommand
	LowPriorityCommands    []*responseCommand
	FinalCommands          []*responseCommand
	JsFiles                *attributeMap
	StyleSheets            *attributeMap
	Alerts                 []string
	NewLocation            string
	WinClose               bool
	Controls               map[string]responseControl
}

// Serialize encodes the response for the pagestate. Currently, serialization of the response is only
// used by the testing framework.
func (r *Response) Serialize(e Encoder) {
	enc := responseEncoded{
		ExclusiveCommand:       r.exclusiveCommand,
		HighPriorityCommands:   r.highPriorityCommands,
		MediumPriorityCommands: r.mediumPriorityCommands,
		LowPriorityCommands:    r.lowPriorityCommands,
		FinalCommands:          r.finalCommands,
		JsFiles:                r.jsFiles,
		StyleSheets:            r.styleSheets,
		Alerts:                 r.alerts,
		NewLocation:            r.newLocation,
		WinClose:               r.winClose,
		Controls:               r.controls,
	}
	if err := e.Encode(enc); err != nil {
		panic(err)
	}
}

// Deserialize unpacks the response from the pagestate. Currently the response is only serialized
// in the testing framework.
func (r *Response) Deserialize(d Decoder) {
	enc := responseEncoded{}
	if err := d.Decode(&enc); err != nil {
		panic(err)
	}
	r.exclusiveCommand = enc.ExclusiveCommand
	r.highPriorityCommands = enc.HighPriorityCommands
	r.mediumPriorityCommands = enc.MediumPriorityCommands
	r.lowPriorityCommands = enc.LowPriorityCommands
	r.finalCommands = enc.FinalCommands
	r.jsFiles = enc.JsFiles
	r.styleSheets = enc.StyleSheets
	r.alerts = enc.Alerts
	r.newLocation = enc.NewLocation
	r.winClose = enc.WinClose
	r.controls = enc.Controls
}
