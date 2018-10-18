package page

//go:generate hero -source template

import (
	"bytes"
	"context"
	"fmt"
	"github.com/spekary/gengen/maps"
	"github.com/spekary/goradd/goradd"
	"github.com/spekary/goradd/html"
	"github.com/spekary/goradd/log"
	action2 "github.com/spekary/goradd/page/action"
	"github.com/spekary/goradd/page/session"
	"goradd-project/config"
	gohtml "html"
	"reflect"
)

const PrivateActionBase = 1000
const sessionControlStates string = "goradd.controlStates"
const sessionControlTypeState string = "goradd.controlType"

const RequiredErrorMessage string = "A value is required"

type ValidationType int

const (
	ValidateDefault ValidationType = iota // This is used by the event to indicate it is not overriding.
	ValidateNone                          // Force no validation.
	ValidateForm

	ValidateSiblingsAndChildren // container validations. Will include validating the current control.
	ValidateSiblingsOnly
	ValidateChildrenOnly

	// ValidateContainer will use the validation setting of a parent control with ValidateSiblingsAndChildren, ValidateSiblingsOnly,
	// ValidateChildrenOnly, or ValidateTarget as the stopping point for validation.
	ValidateContainer

	// ValidateTargetsOnly will only validate the specified targets
	ValidateTargetsOnly
)

type ValidationState int

const (
	NotValidated ValidationState = iota
	Valid
	Invalid
)

type ControlTemplateFunc func(ctx context.Context, control ControlI, buffer *bytes.Buffer) error

type ControlWrapperFunc func(ctx context.Context, control ControlI, ctrl string, buffer *bytes.Buffer)

var DefaultCheckboxLabelDrawingMode = html.LabelAfter // Setting used by checkboxes and radio buttons to default how they draw labels.


type ControlI interface {
	ID() string
	control() *Control
	DrawI

	// Drawing support

	DrawTag(context.Context) string
	DrawInnerHtml(context.Context, *bytes.Buffer) error
	DrawTemplate(context.Context, *bytes.Buffer) error
	PreRender(context.Context, *bytes.Buffer) error
	PostRender(context.Context, *bytes.Buffer) error
	ShouldAutoRender() bool
	SetShouldAutoRender(bool)
	DrawAjax(ctx context.Context, response *Response) error
	DrawChildren(ctx context.Context, buf *bytes.Buffer) error
	DrawText(ctx context.Context, buf *bytes.Buffer)
	With(w WrapperI) ControlI
	HasWrapper() bool
	Wrapper() WrapperI

	// Hierarchy functions

	Parent() ControlI
	Children() []ControlI
	SetParent(parent ControlI)
	Remove()
	RemoveChild(id string)
	RemoveChildren()
	Page() *Page
	ParentForm() FormI
	Child(string) ControlI

	// hmtl and css

	SetAttribute(name string, val interface{}) ControlI
	SetWrapperAttribute(name string, val interface{}) ControlI
	Attribute(string) string
	HasAttribute(string) bool
	DrawingAttributes() *html.Attributes
	WrapperAttributes() *html.Attributes
	AddClass(class string) ControlI
	RemoveClass(class string) ControlI
	AddWrapperClass(class string) ControlI
	SetStyles(*html.Style)
	SetStyle(name string, value string) ControlI
	SetWidthStyle(w interface{}) ControlI
	SetHeightStyle(w interface{}) ControlI

	PutCustomScript(ctx context.Context, response *Response)

	HasFor() bool
	SetHasFor(bool) ControlI

	Label() string
	SetLabel(n string) ControlI
	TextIsLabel() bool
	Text() string
	SetText(t string) ControlI
	ValidationMessage() string
	SetValidationError(e string)
	Instructions() string
	SetInstructions(string) ControlI

	WasRendered() bool
	IsRendering() bool
	IsVisible() bool
	SetVisible(bool)

	Refresh()

	Action(context.Context, ActionParams)
	PrivateAction(context.Context, ActionParams)
	SetActionValue(interface{}) ControlI
	ActionValue() interface{}
	On(e EventI, a ...action2.ActionI) EventI
	Off()
	WrapEvent(eventName string, selector string, eventJs string) string

	UpdateFormValues(*Context)

	Validate(ctx context.Context) bool
	ValidationState() ValidationState
	ValidationType(EventI) ValidationType

	// SaveState tells the control whether to save the basic state of the control, so that when the form is reentered, the
	// data in the control will remain the same. This is particularly useful if the control is used as a filter for the
	// contents of another control.
	SaveState(context.Context, bool)
	MarshalState(m maps.Setter)
	UnmarshalState(m maps.Loader)

	// Shortcuts for translation
	T(string) string
	Translate(string) string

	// Serialization helpers
	Restore(self ControlI)

	// API
	SetIsRequired(r bool) ControlI

}

type Control struct {
	goradd.Base

	id   string
	page *Page // Page This control is part of

	parent   ControlI   // Parent control
	children []ControlI // Child controls

	Tag            string
	IsVoidTag      bool                  // tag does not get wrapped with a terminating tag, but just ends instead
	hasNoSpace     bool                  // For special situations where we want no space between This and next tag. Spans in particular may need This.
	attributes     *html.Attributes      // a collection of attributes to apply to the control
	text           string                // multi purpose, can be button text, inner text inside of tags, etc.
	textLabelMode  html.LabelDrawingMode // describes how to draw the internal label
	htmlEscapeText bool                  // whether to escape the text output, or send straight text

	attributeScripts []*[]interface{} // commands to send to our javascript to redraw portions of this control via ajax. Allows us to skip drawing the entire control.

	isRequired       bool
	isHidden         bool
	isOnPage         bool
	shouldAutoRender bool

	// internal status functions. Do not serialize.
	isModified  bool
	isRendering bool
	wasRendered bool

	isBlockLevel      bool // true to use a div for the wrapper, false for a span
	wrapper           WrapperI
	wrapperAttributes *html.Attributes
	label             string // the given label, often used as a label tag. Not drawn by default, but the wrapper drawing function uses it. Can also get controls by label.

	hasFor       bool   // When drawing the label, should it use a for attribute? This is helpful for screen readers and navigation on certain kinds of tags.
	instructions string // Instructions, if the field needs extra explanation. You could also try adding a tooltip to the wrapper.

	// ErrorForRequired is the error that will display if a control value is required but not set.
	ErrorForRequired string
	// ValidMessage is the message to display if the control has successfully been validated. Leave blank if you don't want a message to show when valid. Can be useful to contrast between invalid and valid controls in a busy form.
	ValidMessage          string
	validationMessage     string // The message to display when showing the validation condition
	validationState       ValidationState
	validationType        ValidationType
	validationTargets     []string // List of control IDs to target validation
	blockParentValidation bool     // This blocks a parent from validating this control. Useful for dialogs, or situations where multiple panels control their own space.

	actionValue interface{}

	events        EventMap
	privateEvents EventMap
	eventCounter  EventID

	shouldSaveState bool
}

// Init should be called immediately after a control is created and is responsible for setting up the initial state of a
// new control. Your subclasses should have their own Init function, and
// should call the superclass function. This Init function sets up a parent-child relationship with the given parent
// control, and sets up data structures to use the control in object-oriented ways with virtual functions.
// The id is the control id that will appear as the id in html. Leave blank for the system to create a unique id for you.
func (c *Control) Init(self ControlI, parent ControlI, id string) {
	c.Base.Init(self)
	c.attributes = html.NewAttributes()
	c.wrapperAttributes = html.NewAttributes()
	if parent != nil {
		c.page = parent.Page()
		c.id = c.page.GenerateControlID(id)
	}
	self.SetParent(parent)
	c.htmlEscapeText = true // default to encoding the text portion. Explicitly turn this off if you need something else
}


// this() supports object oriented features by giving easy access to the virtual function interface
// Subclasses should provide a duplicate. Calls that implement chaining should return the result of this function.
func (c *Control) this() ControlI {
	return c.Self.(ControlI)
}

// Restore is called after the control has been deserialized. It creates any required data structures
// that are not saved in serialization.
func (c *Control) Restore(self ControlI) {
	c.Base.Init(self)
	if c.attributes == nil {
		c.attributes = html.NewAttributes()
	}
	if c.wrapperAttributes == nil {
		c.wrapperAttributes = html.NewAttributes()
	}
}

func (c *Control) ID() string {
	return c.id
}

// Extract the control from an interface. This is for private use, when called through the interface
func (c *Control) control() *Control {
	return c
}

func (c *Control) PreRender(ctx context.Context, buf *bytes.Buffer) error {
	form := c.ParentForm()
	if c.Page() == nil ||
		form == nil ||
		c.Page() != form.Page() {

		return NewError(ctx, "The control can not be drawn because it is not a member of a form that is on the override.")
	}

	if c.wasRendered || c.isRendering {
		return NewError(ctx, "This control has already been drawn.")
	}

	// Because we may be rerendering a parent control, we need to make sure all "child" controls are marked as NOT being on the override.
	if c.children != nil {
		for _, child := range c.children {
			child.control().markOnPage(false)
		}
	}

	// Finally, let's specify that we have begun rendering this control
	c.isRendering = true

	return nil
}

// Draws the default control structure into the given buffer.
func (c *Control) Draw(ctx context.Context, buf *bytes.Buffer) (err error) {
	// TODO: Capture errors and panics, writing what we can to the buffer on error

	if err = c.this().PreRender(ctx, buf); err != nil {
		return err
	}

	var h string

	if c.isHidden {
		// We are invisible, but not using a wrapper. This creates a problem, in that when we go visible, we do not know what to replace
		// To fix this, we create an empty, invisible control in the place where we would normally draw
		h = "<span id=\"" + c.this().ID() + "\" style=\"display:none;\" data-grctl></span>\n"
	} else {
		h = c.this().DrawTag(ctx)
	}

	if !config.Minify {
		s := html.Comment(fmt.Sprintf("Control Type:%s, Id:%s", c.Type(), c.ID())) + "\n"
		buf.WriteString(s)
	}

	if c.wrapper != nil && !c.isHidden {
		c.wrapper.Wrap(ctx, c.this(), h, buf)
	} else {
		buf.WriteString(h)
	}

	response := c.ParentForm().Response()
	c.this().PutCustomScript(ctx, response)
	c.GetActionScripts(response)
	c.this().PostRender(ctx, buf)
	return
}

// PutCustomScript is the place where you add javascript that transforms the html into a custom javascript control.
// Do this by calling functions on the response object.
// This implementation is a stub.
func (c *Control) PutCustomScript(ctx context.Context, response *Response) {

}

/**
* DrawAjax will be called during an Ajax rendering of the controls. Every control gets called. Each control
* is responsible for rendering itself. Some objects automatically render their child objects, and some don't,
* so we detect whether the parent is being rendered, and assume the parent is taking care of rendering for
* us if so.
*
* Override if you want more control over ajax drawing, like if you detect parts of your control that have changed
* and then want to draw only those parts. This will get called on every control on every ajax draw request.
* It is up to you to test the blnRendered flag of the control to know whether the control was already rendered
* by a parent control before drawing here.
*
 */

func (c *Control) DrawAjax(ctx context.Context, response *Response) (err error) {

	if c.isModified {
		// simply re-render the control and assume rendering will handle rendering its children

		func() {
			// wrap in a function to get deferred PutBuffer to execute immediately after drawing
			buf := GetBuffer()
			defer PutBuffer(buf)

			err = c.this().Draw(ctx, buf)
			response.SetControlHtml(c.ID(), buf.String())
		}()
	} else {
		// add attribute changes
		if c.attributeScripts != nil {
			for _, scripts := range c.attributeScripts {
				response.ExecuteControlCommand(c.ID(), (*scripts)[0].(string), (*scripts)[1:]...)
			}
			c.attributeScripts = nil
		}

		// ask the child controls to potentially render, since This control doesn't need to
		for _, child := range c.children {
			err = child.DrawAjax(ctx, response)
			if err != nil {
				return
			}
		}
	}
	return
}

func (c *Control) PostRender(ctx context.Context, buf *bytes.Buffer) (err error) {
	// Update watcher
	//if ($This->objWatcher) {
	//$This->objWatcher->makeCurrent();
	//}

	c.isRendering = false
	c.wasRendered = true
	c.isOnPage = true
	c.isModified = false
	c.attributeScripts = nil // Entire control was redrawn, so don't need these

	return
}

// Draw the control tag itself. Override to draw the tag in a different way, or draw more than one tag if
// drawing a compound control.
func (c *Control) DrawTag(ctx context.Context) string {
	var ctrl string

	attributes := c.this().DrawingAttributes()
	if c.wrapper == nil {
		if a := c.this().WrapperAttributes(); a != nil {
			attributes.Merge(a)
		}
	}

	if c.IsVoidTag {
		ctrl = html.RenderVoidTag(c.Tag, attributes)
	} else {
		buf := GetBuffer()
		defer PutBuffer(buf)
		if err := c.this().DrawInnerHtml(ctx, buf); err != nil {
			panic(err)
		}
		if err := c.RenderAutoControls(ctx, buf); err != nil {
			panic(err)
		}
		if c.hasNoSpace {
			ctrl = html.RenderTagNoSpace(c.Tag, attributes, buf.String())

		} else {
			ctrl = html.RenderTag(c.Tag, attributes, buf.String())
		}
	}
	return ctrl
}

// RenderAutoControls is an internal function to draw controls marked to autoRender. These are generally used for hidden controls
// that can be shown without impacting layout, or that are scripts only.
func (c *Control) RenderAutoControls(ctx context.Context, buf *bytes.Buffer) (err error) {
	// Figuring out where to draw these controls can be difficult.

	for _, ctrl := range c.children {
		if ctrl.ShouldAutoRender() &&
			!ctrl.WasRendered() {

			err = ctrl.Draw(ctx, buf)

			if err != nil {
				break
			}
		}
	}
	return
}

// Controls that use templates should use this function signature for the template. That will override this one, and
// we will then detect that the template was drawn. Otherwise, we detect that no template was defined and it will move
// on to drawing the controls without a template, or just the text if text is defined.
func (c *Control) DrawTemplate(ctx context.Context, buf *bytes.Buffer) (err error) {
	// Don't change this to use some kind of function injection, as such things are not serializable
	return NewFrameworkError(FrameworkErrNoTemplate)
}

// Returns the inner text of the control, if the control is not a self terminating (void) control. Sub-controls can
// override this.
func (c *Control) DrawInnerHtml(ctx context.Context, buf *bytes.Buffer) (err error) {
	if err = c.this().DrawTemplate(ctx, buf); err == nil {
		return
	} else if appErr, ok := err.(FrameworkError); !ok || appErr.Err != FrameworkErrNoTemplate {
		return
	}

	err = nil

	if c.children != nil && len(c.children) > 0 {
		err = c.this().DrawChildren(ctx, buf)
		return
	}

	c.this().DrawText(ctx, buf)

	return
}

func (c *Control) DrawChildren(ctx context.Context, buf *bytes.Buffer) (err error) {
	if c.children != nil {
		for _, child := range c.children {
			err = child.Draw(ctx, buf)
			if err != nil {
				break
			}
		}
	}
	return
}

// Draws the text of the control, escaping if needed
func (c *Control) DrawText(ctx context.Context, buf *bytes.Buffer) {
	if c.text != "" {
		text := c.text

		if c.htmlEscapeText {
			text = gohtml.EscapeString(text)
		}
		buf.WriteString(text)
	}
}

// With sets the wrapper style for the control, essentially setting the wrapper template function that will be used.
func (c *Control) With(w WrapperI) ControlI {
	c.wrapper = w
	return c.this() // for chaining
}

func (c *Control) HasWrapper() bool {
	return c.wrapper != nil
}

func (c *Control) Wrapper() WrapperI {
	return c.wrapper
}


func (c *Control) SetAttribute(name string, val interface{}) ControlI {
	if name == "id" {
		panic("You can only set the 'id' attribute of a control when it is created")
	}

	changed, err := c.attributes.SetChanged(name, html.AttributeString(val))
	if err != nil {
		panic(err)
	}

	if changed {
		// The val passed in might be a calculation, so we need to get the ultimate new value
		v2 := c.attributes.Get(name)
		c.AddRenderScript("attr", name, v2)
	}
	return c.this()
}

func (c *Control) SetWrapperAttribute(name string, val interface{}) ControlI {
	if name == "id" {
		panic("You cannot set the 'id' attribute of a wrapper")
	}

	changed, err := c.wrapperAttributes.SetChanged(name, html.AttributeString(val))
	if err != nil {
		panic(err)
	}

	if changed {
		// TODO: Make This an attribute script instead of redrawing the whole control. Will prevent having to redraw the whole control
		c.isModified = true
	}
	return c.this()
}

func (c *Control) Attribute(name string) string {
	return c.attributes.Get(name)
}

func (c *Control) HasAttribute(name string) bool {
	return c.attributes.Has(name)
}

// Returns a set of attributes that should override those set by the user. This allows controls to set attributes
// just before drawing that should take precedence over other attributes, and that are critical to drawing the
// tag of the control. This function is designed to only be called by overriding functions, changes will not be remembered.
func (c *Control) DrawingAttributes() *html.Attributes {
	a := html.NewAttributesFrom(c.attributes)
	a.SetID(c.id)                   // make sure the control id is set at a minimum
	a.SetDataAttribute("grctl", "") // make sure control is registered. Overriding controls can put a control name here.

	if c.HasWrapper() {
		c.wrapper.ModifyDrawingAttributes(c.this(), a)
	}

	return a
}

// WrapperAttributes returns the actual attributes for the wrapper. Changes WILL be remembered so that subsequent ajax
// drawing will draw the wrapper correctly. However, it is up to you to refresh the control if you change anything.
func (c *Control) WrapperAttributes() *html.Attributes {
	return c.wrapperAttributes
}


func (c *Control) SetDataAttribute(name string, val interface{}) {
	var v string
	var ok bool

	if v, ok = val.(string); !ok {
		v = fmt.Sprintf("%v", v)
	}

	changed, err := c.attributes.SetDataAttributeChanged(name, v)
	if err != nil {
		panic(err)
	}

	if changed {
		c.AddRenderScript("data", name, v) // Use the jQuery data method to set the data during ajax requests
	}
}

func (c *Control) AddClass(class string) ControlI {
	if changed := c.attributes.AddClassChanged(class); changed {
		v2 := c.attributes.Class()
		c.AddRenderScript("attr", "class", v2)
	}
	return c.this()
}

func (c *Control) RemoveClass(class string) ControlI {
	if changed := c.attributes.RemoveClass(class); changed {
		c.isModified = true
	}
	return c.this()
}

func (c *Control) AddWrapperClass(class string) ControlI {
	if changed := c.wrapperAttributes.AddClassChanged(class); changed {
		c.isModified = true
	}
	return c.this()
}

// Adds a variadic parameter list to the renderScripts array, which is an array of javascript commands to send to the
// browser the next time control drawing happens. These commands allow javascript to change an aspect of the control without
// having to redraw the entire control. This should primarily be used by control implementations only.
func (c *Control) AddRenderScript(params ...interface{}) {
	c.attributeScripts = append(c.attributeScripts, &params)
}

// Control Hierarchy Functions
func (c *Control) Parent() ControlI {
	return c.parent
}

func (c *Control) Children() []ControlI {
	return c.children
}

// Remove removes the current control from its parent. After this is done, the control and all its child items will
// not be part of the override, but the child items will still be accessible through the control itself.
func (c *Control) Remove() {
	if c.parent != nil {
		c.parent.control().removeChild(c.this().ID(), true)
		if !c.shouldAutoRender {
			//c.Refresh() // TODO: Do this through ajax
		}
	} else {
		c.page.removeControl(c.this().ID())
	}
}

// RemoveChild removes the given child control from both the control and the override.
func (c *Control) RemoveChild(id string) {
	c.removeChild(id, true)
}


// removeChild is a private function that will remove a child control from the current control
func (c *Control) removeChild(id string, fromPage bool) {
	for i, v := range c.children {
		if v.ID() == id {
			c.children = append(c.children[:i], c.children[i+1:]...) // remove found item from list
			if fromPage {
				v.control().removeChildrenFromPage()
				c.page.removeControl(id)
			}
			v.control().parent = nil
			break
		}
	}
}

func (c *Control) removeChildrenFromPage() {
	for _, v := range c.children {
		v.control().removeChildrenFromPage()
		c.page.removeControl(v.ID())
	}
}

// RemoveChildren removes all the child controls from this control and the override
func (c *Control) RemoveChildren() {
	for _, child := range c.children {
		child.control().removeChildrenFromPage()
		c.page.removeControl(child.ID())
		child.control().parent = nil
	}
	c.children = nil
}

func (c *Control) SetParent(newParent ControlI) {
	if c.parent == nil {
		c.control().addChildControlsToPage()
	} else {
		c.parent.control().removeChild(c.ID(), newParent == nil)
		if !c.shouldAutoRender {
			//c.parent.Refresh()
		}
	}
	c.parent = newParent
	if c.parent != nil {
		c.parent.control().addChildControl(c.this())
		if !c.shouldAutoRender {
			// TODO: insert into DOM  instead of c.parent.Refresh()
		}
	}
	c.page.addControl(c.this())

	if c.shouldAutoRender && newParent != nil {
		//c.Refresh()
	}

	// TODO: Refresh as needed, but without refreshing the form
}

func (c *Control) Child(id string) ControlI {
	for _, c := range c.children {
		if c.ID() == id {
			return c
		}
	}
	return nil
}

func (c *Control) addChildControlsToPage() {
	for _, child := range c.children {
		child.control().addChildControlsToPage()
		c.page.addControl(child)
	}
}

// Private function called by setParent on parent function
func (c *Control) addChildControl(child ControlI) {
	if c.children == nil {
		c.children = make([]ControlI, 0)
	}
	c.children = append(c.children, child)
}

// ParentForm returns the form object that encloses this control.
func (c *Control) ParentForm() FormI {
	return c.page.Form()
}

func (c *Control) Page() *Page {
	return c.page
}

// Drawing aids
func (c *Control) Refresh() {
	c.isModified = true
}

func (c *Control) SetIsRequired(r bool) ControlI {
	c.isRequired = r
	return c.this()
}

func (c *Control) IsRequired() bool {
	return c.isRequired
}

func (c *Control) ValidationMessage() string {
	return c.validationMessage
}

// SetValidationError sets the validation error to the given string. It will also handle setting the wrapper class
// to indicate an error. Override if you have a different way of handling errors.
func (c *Control) SetValidationError(e string) {
	if !c.HasWrapper() {
		panic("Only controls with wrappers can get a validation error.")
	}
	if c.validationMessage != e {
		c.validationMessage = e
		c.Refresh() // TODO: Set a response attribute instead to only update the inner text of the ctlId_err div and possibly the div class. Tricky because bootstrap has multiple options to set the div class.

		if e == "" {
			c.validationState = NotValidated
			c.SetWrapperAttribute("class", "- error")
		} else {
			c.validationState = Invalid
			c.SetWrapperAttribute("class", "+ error")
		}
	}
}

func (c *Control) ValidationState() ValidationState {
	return c.validationState
}

func (c *Control) SetText(t string) ControlI {
	if t != c.text {
		c.text = t
		c.Refresh()
	}
	return c.this()
}

func (c *Control) Text() string {
	return c.text
}

func (c *Control) SetLabel(n string) ControlI {
	if n != c.label {
		c.label = n
		c.Refresh()
	}
	return c.this()
}

func (c *Control) Label() string {
	return c.label
}

// TextIsLabel is used by the drawing routines to determine if the control's text should be wrapped with a label tag,
// Meaning the Label itself
func (c *Control) TextIsLabel() bool {
	return false
}

func (c *Control) SetInstructions(i string) ControlI {
	if i != c.instructions {
		c.instructions = i
		c.Refresh()
	}
	return c.this()
}

func (c *Control) Instructions() string {
	return c.instructions
}

func (c *Control) markOnPage(v bool) {
	c.isOnPage = v
}

func (c *Control) WasRendered() bool {
	return c.wasRendered
}

func (c *Control) IsRendering() bool {
	return c.isRendering
}

func (c *Control) HasFor() bool {
	return c.hasFor
}

func (c *Control) SetHasFor(v bool) ControlI {
	if v != c.hasFor {
		c.hasFor = v
		c.Refresh()
	}
	return c.this()
}

func (c *Control) SetShouldAutoRender(r bool) {
	c.shouldAutoRender = r
}

func (c *Control) ShouldAutoRender() bool {
	return c.shouldAutoRender
}

// On adds an event listener to the control that will trigger the given actions
func (c *Control) On(e EventI, actions ...action2.ActionI) EventI {
	var isPrivate bool
	c.Refresh() // completely redraw the control. The act of redrawing will turn off old scripts.
	// TODO: Adding scripts should instead just redraw the associated script block. We will need to
	// implement a script block with every control connected by id
	e.AddActions(actions...)
	c.eventCounter++
	for _, action := range actions {
		if _, ok := action.(action2.PrivateAction); ok {
			isPrivate = true
			break
		}
	}

	// Get a new event id
	for {
		if _, ok := c.events[c.eventCounter]; ok {
			c.eventCounter++
		} else if _, ok := c.privateEvents[c.eventCounter]; ok {
			c.eventCounter++
		} else {
			break
		}
	}

	if isPrivate {
		if c.privateEvents == nil {
			c.privateEvents = map[EventID]EventI{}
		}
		c.privateEvents[c.eventCounter] = e
	} else {
		if c.events == nil {
			c.events = map[EventID]EventI{}
		}
		c.events[c.eventCounter] = e
	}
	return e
}

// Off removes all event handlers from the control
func (c *Control) Off() {
	c.events = nil
}

// SetActionValue sets a value that is provided to actions when they are triggered. The value can be a static value
// or one of the javascript.* objects that can dynamically generated values. The value is then sent back to the action
// handler after the action is triggered.
func (c *Control) SetActionValue(v interface{}) ControlI {
	c.actionValue = v
	return c.this()
}

// ActionValue returns the control's action value
func (c *Control) ActionValue() interface{} {
	return c.actionValue
}

// Action processes actions. Typically, the Action function will first look at the id to know how to handle it.
// This is just an empty implemenation. Sub-controls should implement this.
func (c *Control) Action(ctx context.Context, a ActionParams) {
}

// PrivateAction processes actions that a control sets up for itself, and that it does not want to give the opportunity
// for users of the control to manipulate or remove those actions. Generally, private actions should call their superclass
// PrivateAction function too.
func (c *Control) PrivateAction(ctx context.Context, a ActionParams) {
}

// GetActionScripts is an internal function called during drawing to recursively gather up all the event related
// scripts attached to the control and send them to the response.
func (c *Control) GetActionScripts(r *Response) {
	// Render actions
	if c.privateEvents != nil {
		for id, e := range c.privateEvents {
			s := e.RenderActions(c.this(), id)
			r.ExecuteJavaScript(s, PriorityStandard)
		}
	}

	if c.events != nil {
		for id, e := range c.events {
			s := e.RenderActions(c.this(), id)
			r.ExecuteJavaScript(s, PriorityStandard)
		}
	}
}

// Recursively reset the drawing flags
func (c *Control) resetDrawingFlags() {
	c.wasRendered = false
	c.isModified = false

	if children := c.this().Children(); children != nil {
		for _, child := range children {
			child.control().resetDrawingFlags()
		}
	}
}

// Recursively reset the validation state
func (c *Control) resetValidation() {
	if c.HasWrapper() &&
		(c.validationMessage != "" || c.validationState != NotValidated) {
		c.validationMessage = ""
		c.validationState = NotValidated
		c.wrapperAttributes.RemoveClass("error")
		c.Refresh() // TODO: Handle the above with javascript calls so base control does not have to be redrawn
	}

	if children := c.this().Children(); children != nil {
		for _, child := range children {
			child.control().resetValidation()
		}
	}
}

// WrapEvent is an internal function to allow the control to customize its treatment of event processing.
func (c *Control) WrapEvent(eventName string, selector string, eventJs string) string {
	if selector != "" {
		return fmt.Sprintf("$j('#%s').on('%s', '%s', function(event, ui){%s});", c.ID(), eventName, selector, eventJs)
	} else {
		return fmt.Sprintf("$j('#%s').on('%s', function(event, ui){%s});", c.ID(), eventName, eventJs)
	}
}

// updateValues is called by the form during event handling. It reflexively updates the values in each of its child controls
func (c *Control) updateValues(ctx *Context) {
	c.this().UpdateFormValues(ctx)
	children := c.Children()
	if children != nil {
		for _, child := range children {
			child.control().updateValues(ctx)
		}
	}
}

// UpdateFormValues should be implemented by controls to get their values from the context.
func (c *Control) UpdateFormValues(ctx *Context) {

}

// doAction is an internal function that the override manager uses to send actions to controls.
func (c *Control) doAction(ctx context.Context) {
	var e EventI
	var ok bool
	var isPrivate bool
	var grCtx = GetContext(ctx)

	if e, ok = c.events[grCtx.eventID]; !ok {
		e, ok = c.privateEvents[grCtx.eventID]
		isPrivate = true
	}

	if !ok {
		log.FrameworkDebug("doAction - event not found: ", grCtx.eventID)
		return
	}

	if (e.event().validationOverride != ValidateNone && e.event().validationOverride != ValidateDefault) ||
		(e.event().validationOverride == ValidateDefault && c.this().ValidationType(e) != ValidateNone) {
		c.ParentForm().control().resetValidation()
	}

	if c.passesValidation(ctx, e) {
		log.FrameworkDebug("doAction - triggered event: ", e.String())
		for _, a := range e.GetActions() {
			callbackAction := a.(action2.CallbackActionI)
			p := ActionParams{
				ID:        callbackAction.ID(),
				Action:    a,
				ControlId: c.ID(),
			}

			// grCtx.actionValues is a json representation of the action values. We extract the json, but since json does
			// not differentiate between float and int, we will leave all numbers as json.Number types so we can extract later.
			// use javascript.NumberInt() to easily convert numbers in interfaces to int values.
			p.Values = grCtx.actionValues
			dest := c.Page().GetControl(callbackAction.GetDestinationControlID())

			if dest != nil {
				if isPrivate {
					if log.HasLogger(log.FrameworkDebugLog) {
						log.FrameworkDebugf("doAction - PrivateAction, DestId: %s, action2.ActionId: %d, Action: %s, TriggerId: %s",
							dest.ID(), p.ID, reflect.TypeOf(p.Action).String(), p.ControlId)
					}
					dest.PrivateAction(ctx, p)
				} else {
					if log.HasLogger(log.FrameworkDebugLog) {
						log.FrameworkDebugf("doAction - Action, DestId: %s, action2.ActionId: %d, Action: %s, TriggerId: %s",
							dest.ID(), p.ID, reflect.TypeOf(p.Action).String(), p.ControlId)
					}
					dest.Action(ctx, p)
				}
			}
		}
	} else {
		log.FrameworkDebug("doAction - failed validation: ", e.String())
	}
}

// SetBlockParentValidation will prevent a parent from validating this control. This is generally useful for panels and
// other containers of controls that wish to have their own validation scheme. Dialogs in particular need this since
// they essentially act as a separate form, even though technically they are included in a form.
func (c *Control) SetBlockParentValidation(block bool) {
	c.blockParentValidation = block
}

// SetValidationType specifies how this control validates other controls. Typically its either ValidateNone or ValidateForm.
// ValidateForm will validate all the controls on the form.
// ValidateSiblingsAndChildren will validate the immediate siblings of the target controls and their children
// ValidateSiblingsOnly will validate only the siblings of the target controls
// ValidateTargetsOnly will validate only the specified target controls
func (c *Control) SetValidationType(typ ValidationType) {
	c.validationType = typ
}

// ValidationType is an internal function to return the validation type. It allows subclasses to override it.
func (c *Control) ValidationType(e EventI) ValidationType {
	if c.validationType == ValidateNone || c.validationType == ValidateDefault {
		return ValidateNone
	} else {
		return c.validationType
	}
}

// SetValidationTargets specifies which controls to validate, in conjunction with the ValidationType setting,
// giving you very fine-grained control over validation. The default
// is to use just this control as the target.
func (c *Control) SetValidationTargets(controlIDs ...string) {
	c.validationTargets = controlIDs
}

// passesValidation checks to see if the event requires validation, and if so, if it passes the required validation
func (c *Control) passesValidation(ctx context.Context, event EventI) (valid bool) {
	validation := c.this().ValidationType(event)

	if v := event.event().validationOverride; v != ValidateDefault {
		validation = v
	}

	if validation == ValidateDefault || validation == ValidateNone {
		return true
	}

	var targets []ControlI

	if c.validationTargets == nil {
		if c.validationType == ValidateForm {
			targets = []ControlI{c.ParentForm()}
		} else if c.validationType == ValidateContainer {
			for target := c.Parent(); target != nil; target = target.Parent() {
				switch target.control().validationType {
				case ValidateChildrenOnly:
					fallthrough
				case ValidateSiblingsAndChildren:
					fallthrough
				case ValidateSiblingsOnly:
					fallthrough
				case ValidateTargetsOnly:
					validation = target.control().validationType
					targets = []ControlI{target}
					break
				}
			}
			// Target is the form
			targets = []ControlI{c.ParentForm()}
			validation = ValidateForm
		} else {
			targets = []ControlI{c}
		}
	} else {
		if c.validationType == ValidateForm ||
			c.validationType == ValidateContainer {
			panic("Unsupported validation type and target combo.")
		}
		for _, id := range c.validationTargets {
			if c2 := c.Page().GetControl(id); c2 != nil {
				targets = append(targets, c2)
			}
		}
	}

	valid = true

	switch validation {
	case ValidateForm:
		valid = c.ParentForm().control().validateChildren(ctx)
	case ValidateSiblingsAndChildren:
		for _, t := range targets {
			valid = t.control().validateSiblingsAndChildren(ctx) && valid
		}
	case ValidateSiblingsOnly:
		for _, t := range targets {
			valid = t.control().validateSiblings(ctx) && valid
		}
	case ValidateChildrenOnly:
		for _, t := range targets {
			valid = t.control().validateChildren(ctx) && valid
		}

	case ValidateTargetsOnly:
		var valid bool
		for _, t := range targets {
			valid = t.Validate(ctx) && valid
		}
	}
	return valid
}

// Validate is designed to be overridden. Overriding controls should call the parent version before doing their own validation.
func (c *Control) Validate(ctx context.Context) bool {
	if c.HasWrapper() && c.validationMessage != c.ValidMessage {
		c.validationMessage = c.ValidMessage
		c.Refresh() // TODO: Do this in javascript so whole control does not need to refresh
	}
	c.validationState = Valid
	return true
}

func (c *Control) validateSiblings(ctx context.Context) bool {

	if c.parent == nil {
		return true
	}

	p := c.parent.control()
	siblings := p.children

	var valid = true
	for _, child := range siblings {
		valid = child.Validate(ctx) && valid
	}
	return valid
}

func (c *Control) validateChildren(ctx context.Context) bool {

	if c.children == nil || len(c.children) == 0 {
		return c.this().Validate(ctx)
	}

	var isValid = true
	for _, child := range c.children {
		if !child.control().blockParentValidation {
			isValid = child.control().validateChildren(ctx) && isValid
		}
	}
	if isValid {
		isValid = c.this().Validate(ctx)	// validate self after validating all children, because self might want to invalidate child items
	}

	return isValid
}

func (c *Control) validateSiblingsAndChildren(ctx context.Context) bool {

	if c.parent == nil {
		return true
	}

	p := c.parent.control()
	siblings := p.children

	var isValid = true
	for _, child := range siblings {
		if child.ID() != c.ID() {
			isValid = child.control().validateChildren(ctx) && isValid
		} else {
			// validate self and children
			var childrenValid = true
			if c.children != nil {
				for _, child := range c.children {
					if !child.control().blockParentValidation {
						childrenValid = child.Validate(ctx) && childrenValid
					}
				}
			}
			if childrenValid {
				isValid = c.this().Validate(ctx) // only validate self if children validate
			} else {
				isValid = false
			}
		}
	}
	return isValid
}

// SaveState sets whether the control should save its value and other state information so that if the form is redrawn,
// the value can be restored. This function is also responsible for restoring the previously saved state of the control,
// so call This only after you have set the default state of a control during creation or initialization.
func (c *Control) SaveState(ctx context.Context, saveIt bool) {
	c.shouldSaveState = saveIt
	c.readState(ctx)
}

// writeState is an internal function that will recursively write out the state of itself and its subcontrols
func (c *Control) writeState(ctx context.Context) {
	var stateStore *maps.Map
	var state *maps.Map
	var ok bool

	if c.shouldSaveState {
		state = maps.NewMap()
		c.this().MarshalState(state)
		if state.Len() > 0 {
			state.Set(sessionControlTypeState, c.Type()) // so we can make sure the type is the same when we read, in situations where control Ids are dynamic
			i := session.Get(ctx, sessionControlStates)
			if i == nil {
				stateStore = maps.NewMap()
				session.Set(ctx, sessionControlStates, stateStore)
			} else if _, ok = i.(*maps.Map); !ok {
				stateStore = maps.NewMap()
				session.Set(ctx, sessionControlStates, stateStore)
			} else {
				stateStore = i.(*maps.Map)
			}
			key := c.ParentForm().ID() + ":" + c.ID()
			stateStore.Set(key, state)
		}
	}

	if c.children == nil || len(c.children) == 0 {
		return
	}

	for _, child := range c.children {
		child.control().writeState(ctx)
	}
}

// readState is an internal function that will recursively read the state of itself and its subcontrols
func (c *Control) readState(ctx context.Context) {
	var stateStore *maps.Map
	var state *maps.Map
	var ok bool

	if c.shouldSaveState {
		if i := session.Get(ctx, sessionControlStates); i != nil {
			if stateStore, ok = i.(*maps.Map); !ok {
				return
				// Indicates the entire control state store changed types, so completely ignore it
			}

			key := c.ParentForm().ID() + ":" + c.ID()
			i2 := stateStore.Get(key)
			if state, ok = i2.(*maps.Map); !ok {
				return
				// Indicates This particular item was not stored correctly
			}

			if typ, _ := state.LoadString(sessionControlTypeState); typ != c.Type() {
				return // types are not equal, ids must have changed
			}

			c.this().UnmarshalState(state)
		}
	}

	if c.children == nil || len(c.children) == 0 {
		return
	}

	for _, child := range c.children {
		child.control().readState(ctx)
	}
}

// MarshalState is a helper function for controls to save their basic state, so that if the form is reloaded, the
// value that the user entered will not be lost. Implementing controls should add items to the given map.
// Note that the control id is used as a key for the state,
// so that if you are dynamically adding controls, you should make sure you give a specific, non-changing control id
// to the control, or the state may be lost.
func (c *Control) MarshalState(m maps.Setter) {
}

// UnmarshalState is a helper function for controls to get their state from the stateStore. To implement it, a control
// should read the data out of the given map. If needed, implemet your own version checking scheme. The given map will
// be guaranteed to have been written out by the same kind of control as the one reading it. Be sure to call the super-class
// version too.
func (c *Control) UnmarshalState(m maps.Loader) {
}

// T is a shortcut for the override translator that should only be used by internal goradd code. See Translate() for the
// version to use for your project.
func (c *Control) T(in string) string {
	return c.Page().GoraddTranslator().Translate(in)
}

// Translate is a shortcut to the override translator.
// All static strings that could create output to the user should be wrapped in This. The translator itself is designed
// to be capable of per-override translation, meaning each user of the web service can potentially choose their own language
// and see the web override in that language.
func (c *Control) Translate(in string) string {
	return c.Page().ProjectTranslator().Translate(in)
}

func (c *Control) SetDisabled(d bool) {
	c.attributes.SetDisabled(d)
	c.Refresh()
}

func (c *Control) IsDisabled() bool {
	return c.attributes.IsDisabled()
}

func (c *Control) SetDisplay(d string) {
	c.attributes.SetDisplay(d)
	c.Refresh()
}

func (c *Control) IsDisplayed() bool {
	return c.attributes.IsDisplayed()
}

func (c *Control) IsVisible() bool {
	return !c.isHidden
}

func (c *Control) SetVisible(v bool) {
	if c.isHidden == v { // these are opposite in meaning
		c.isHidden = !v
		c.Refresh()
	}
}

func (c *Control) SetStyles(s *html.Style) {
	c.attributes.SetStyles(s)
	c.Refresh() // TODO: Do this with javascript
}

func (c *Control) SetStyle(name string, value string) ControlI {
	if changed,_ := c.attributes.SetStyleChanged(name, value); changed {
		c.Refresh() // TODO: Do this with javascript
	}
	return c.this()
}

func (c *Control) RemoveClassesWithPrefix(prefix string) {
	if c.attributes.RemoveClassesWithPrefix(prefix) {
		c.Refresh() // TODO: Do this with javascript
	}
}

// SetWidthStyle sets the width css property
func (c *Control) SetWidthStyle(w interface{}) ControlI {
	v := html.StyleString(w)
	c.attributes.SetStyle("width", v)
	c.AddRenderScript("css", "width", v) // use javascript to set this value
	return c.this()
}

// SetHeightStyle sets the height css property
func (c *Control) SetHeightStyle(h interface{}) ControlI {
	v := html.StyleString(h)
	c.attributes.SetStyle("height", v)
	c.AddRenderScript("css", "height", v) // use javascript to set this value
	return c.this()
}


// SetEscapeText to false to turn off html escaping of the text output. It is on by default.
func (c *Control) SetEscapeText(e bool) ControlI {
	c.htmlEscapeText = e
	return c.this()
}

// ExecuteJqueryFunction will execute the given JQuery function on the given command, with the given
// parameters. i.e. jQuery("#id").command(params...); will get executed in javascript.
func (c *Control) ExecuteJqueryFunction(command string, params ...interface{}) {
	c.ParentForm().Response().ExecuteControlCommand(c.ID(), command, params...)
}

// ConnectorParams returns a list of options setable by the connector dialog (not currently implemented)
func ControlConnectorParams() *maps.SliceMap {
	m := maps.NewSliceMap()

	// TODO: Add setable options for all controls
	return m
}