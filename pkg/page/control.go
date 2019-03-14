package page

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"github.com/goradd/gengen/pkg/maps"
	"github.com/goradd/goradd/pkg/base"
	"github.com/goradd/goradd/pkg/config"
	"github.com/goradd/goradd/pkg/html"
	"github.com/goradd/goradd/pkg/i18n"
	"github.com/goradd/goradd/pkg/log"
	action2 "github.com/goradd/goradd/pkg/page/action"
	"github.com/goradd/goradd/pkg/session"
	gohtml "html"
	"reflect"
)

const PrivateActionBase = 1000
const sessionControlStates string = "goradd.controlStates"
const sessionControlTypeState string = "goradd.controlType"

const RequiredErrorMessage string = "A value is required"

// ValidationType is used by active controls, like buttons, to determine what other items on the form will get validated
// when the button is pressed. You can set the ValidationType for a control, but you can also set it for individual events
// and override the control's validation setting.
type ValidationType int

const (
	// ValidateDefault is used by events to indicate they are not overriding a control validation. You should not need to use this.
	ValidateDefault ValidationType = iota
	// ValidateNone indicates the control will not validate the form
	ValidateNone
	// ValidateForm is the default validation for buttons, and indicates the entire form and all controls will validate.
	ValidateForm
	// ValidateSiblingsAndChildren will validate the current control, and all siblings of the control and all
	// children of the siblings and current control.
	ValidateSiblingsAndChildren
	// ValidateSiblingOnly will validate only the siblings of the current control, but not any child controls.
	ValidateSiblingsOnly
	// ValidateChildrenOnly will validate only the children of the current control.
 	ValidateChildrenOnly
	// ValidateContainer will use the validation setting of a parent control with ValidateSiblingsAndChildren, ValidateSiblingsOnly,
	// ValidateChildrenOnly, or ValidateTargetsOnly as the stopping point for validation.
	ValidateContainer
	// ValidateTargetsOnly will only validate the specified targets
	ValidateTargetsOnly
)

// ValidationState is used internally by the framework to determine how the control's wrapper handles drawing validation error
// messages. Different wrappers use it to set classes or attributes of the error message or the overall control.
type ValidationState int

const (
	// ValidationWaiting is the default for controls that accept validation. It means that the control expects to be validated,
	// but has not yet been validated. Wrappers should save a spot for the error message of this control so that if
	// an error appears, it will not change the layout of the form.
	ValidationWaiting ValidationState = iota
	// ValidationNever indicates that the control will never fail validation. Essentially it indicates that the wrapper does not
	// need to save a spot for an error message for this control.
	ValidationNever
	// ValidationValid indicates the control has been validated. This state gets entered if some control on the form has failed validation, but
	// this control passed validation. You can choose to display a special message, or a special color, etc., to
	// indicate to the user that this is not the source of the validation problem, or do nothing.
	ValidationValid
	// ValidationInvalid indicates the control has failed validation, and the wrapper should somehow call that out to the user. The error message
	// should be displayed at a minimum, but likely other things should happen as well, like a special color, and
	// aria attributes should be set.
	ValidationInvalid
)

// ControlTemplateFunc is the type of function control templates should create
type ControlTemplateFunc func(ctx context.Context, control ControlI, buffer *bytes.Buffer) error

// ControlWrapperFunc is a template function that specifies how wrappers will draw
type ControlWrapperFunc func(ctx context.Context, control ControlI, ctrl string, buffer *bytes.Buffer)


// DefaultCheckboxLabelDrawingMode is a setting used by checkboxes and radio buttons to default how they draw labels.
// Some CSS framworks are very picky about whether checkbox labels wrap the control, or sit next to the control,
// and whether the label is before or after the control
var DefaultCheckboxLabelDrawingMode = html.LabelAfter

// ControlI is the interface that all controls must support. The functions are implemented by the
// Control methods. See the Control method implementation for a description of each method.
type ControlI interface {
	ID() string
	control() *Control
	DrawI

	// Drawing support

	ΩDrawTag(context.Context) string
	ΩDrawInnerHtml(context.Context, *bytes.Buffer) error
	DrawTemplate(context.Context, *bytes.Buffer) error
	ΩPreRender(context.Context, *bytes.Buffer) error
	ΩPostRender(context.Context, *bytes.Buffer) error
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
	ΩDrawingAttributes() *html.Attributes
	WrapperAttributes() *html.Attributes
	AddClass(class string) ControlI
	RemoveClass(class string) ControlI
	AddWrapperClass(class string) ControlI
	SetStyles(*html.Style)
	SetStyle(name string, value string) ControlI
	SetWidthStyle(w interface{}) ControlI
	SetHeightStyle(w interface{}) ControlI

	ΩPutCustomScript(ctx context.Context, response *Response)

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
	IsOnPage() bool

	Refresh()

	Action(context.Context, ActionParams)
	PrivateAction(context.Context, ActionParams)
	SetActionValue(interface{}) ControlI
	ActionValue() interface{}
	On(e EventI, a ...action2.ActionI) EventI
	Off()
	WrapEvent(eventName string, selector string, eventJs string) string
	HasServerAction(eventName string) bool

	ΩUpdateFormValues(*Context)

	Validate(ctx context.Context) bool
	ValidationState() ValidationState
	ValidationType(EventI) ValidationType

	// SaveState tells the control whether to save the basic state of the control, so that when the form is reentered, the
	// data in the control will remain the same. This is particularly useful if the control is used as a filter for the
	// contents of another control.
	SaveState(context.Context, bool)
	ΩMarshalState(m maps.Setter)
	ΩUnmarshalState(m maps.Loader)

	// Shortcuts for translation

	ΩT(format string) string
	T(format string, params... interface{}) string
	TPrintf(format string, params... interface{}) string

	// Serialization helpers

	Restore(self ControlI)

	// API

	SetIsRequired(r bool) ControlI

	Serialize(e Encoder) (err error)
	Deserialize(d Decoder, p *Page) (err error)
	ΩisSerializer(i ControlI) bool

}

type attributeScriptEntry struct {
	id string	// id of the object to execute the command on. This should be the id of the control, or a a related html object.
	f string	// the jquery function to call
	commands []interface{}	// parameters to the jquery function
}

// A Control is a basic UI widget in goradd. It corresponds to a standard html form object or tag, or a custom javascript
// widget. The Control renders a tag and everything inside of the tag, but can also include a wrapper which associates
// a label, instructions and error messages with the tag. A Control can also associate javascript
// with itself to make sure the javascript is loaded on the page when the control is drawn, and can render
// javascript that will initialize a custom javascript widget.
//
// A Control can have child Controls. It
// can either allow the framework to automatically draw the child Controls as part of the inner-html of
// the Control, can use a template to draw the Child controls, or manually draw them. The Control is part
// of a hierarchical tree structure, with the Form being the root of the tree.
//
// A Control is part of a system that will reflect the state of the control between the client and server.
// When a user updates a control in the browser and performs an action that requires a response from the
// server, the goradd javascript will gather up all the changes in the form and send those to the server.
// The control can read those values and update its own internal state, so that from the perspective
// of the programmer referring to the control, the values in the Control are the same as what the user sees in a browser.
//
// This Control struct is a mixin that all controls should use. You would not normally create a Control directly,
// but rather create one of the "subclasses" of Control. See the control package for Controls that implement
// standard html widgets.
type Control struct {
	base.Base

	// id is the id passed to the control when it is created, or assigned automatically if empty.
	id   string
	// page is a pointer to the page that encloses the entire control tree.
	page *Page

	// parent is the immediate parent control of this control. Only the form object will not have a parent.
	parent   ControlI
	// children are the child controls that belong to this control
	children []ControlI // Child controls

	// Tag is text of the tag that will enclose the control, like "div" or "input"
	Tag            string
	// IsVoidTag should be true if the tag should not have a closing tag, like "img"
	IsVoidTag      bool
	// hasNoSpace is for special situations where we want no space between this and the next tag. Spans in particular may need this.
	hasNoSpace     bool
	// attributes are the collection of custom attributes to apply to the control. This does not include all the
	// attributes that will be drawn, as some are added temporarily just before drawing by GetDrawingAttributes()
	attributes     *html.Attributes
	// test is a multi purpose string that can be button text, inner text inside of tags, etc. depending on the control.
	text           string
	// textLabelMode describes how to draw the internal label
	textLabelMode  html.LabelDrawingMode
	// htmlEscapeText tells us whether to escape the text output, or send straight text
	htmlEscapeText bool

	// attributeScripts are commands to send to our javascript to redraw portions of the control via ajax.
	attributeScripts []attributeScriptEntry

	// isRequired indicates that we will require a value during validation
	isRequired       bool
	// isHidden indicates that we will not draw the control, but rather an invisible placeholder for the control.
	isHidden         bool
	// isOnPage indicates we have drawn the control at some point in the past
	isOnPage         bool
	// shouldAutoRender indicates that we will eventually draw the control even if it is not drawn directly.
	shouldAutoRender bool

	// internal status functions. Do not serialize.

	// isModified will cause the control to redraw as part of the response.
	isModified  bool
	// isRendering is true when we are in the middle of rendering the control.
	isRendering bool
	// wasRendered indicates that the page was drawn during the current response.
	wasRendered bool

	// isBlock is true to use a div for the wrapper, false for a span
	isBlock           bool
	// wrapper is the wrapper object the control will use to draw the label, instructions and error message for the control.
	wrapper           WrapperI
	// wrapperAttributes are the attributes to add to the wrapper tag.
	wrapperAttributes *html.Attributes
	// label is the test to use for the label tag. Not drawn by default, but the wrapper drawing function uses it. Can also get controls by label.
	label             string
	// hasFor tells us if we should draw a for attribute in the label tag. This is helpful for screen readers and navigation on certain kinds of tags.
	hasFor       bool
	// instructions is text associated with the control for extra explanation. You could also try adding a tooltip to the wrapper.
	instructions string

	// ErrorForRequired is the error that will display if a control value is required but not set.
	ErrorForRequired string

	// ValidMessage is the message to display if the control has successfully been validated.
	// Leave blank if you don't want a message to show when valid.
	// Can be useful to contrast between invalid and valid controls in a busy form.
	ValidMessage          string
	// validationMessage is the current validation message that will display when drawing the control
	// This gets copied from ValidMessage at drawing time if the control is in an invalid state
	validationMessage     string
	// validationState is the current validation state of the control, and will effect how the control is drawn.
	validationState       ValidationState
	// validationType indicates how the control will validate itself. See ValidationType for a description.
	validationType        ValidationType
	// validationTargets is the list of control IDs to target validation
	validationTargets     []string
	// This blocks a parent from validating this control. Useful for dialogs, and other situations where sub-controls should control their own space.
	blockParentValidation bool

	// actionValue is the value that will be provided as the ControlValue for any actions that are triggered by this control.
	actionValue interface{}
	// events are all the events added by the control user that the control might trigger
	events        EventMap
	// privateEvents are events that are private to the control and that should not be allowed to be canceled by a control's user.
	privateEvents EventMap
	// eventCounter is used to generate a unique id for an event to help us route the event through the system.
	eventCounter  EventID
	// shouldSaveState indicates that we should save parts of our state into a session variable so that if
	// the client should come back to the form, we will attempt to restore the state of the control. The state
	// in this situation would be the user's input, so text in a textbox, or the selection from a list.
	shouldSaveState bool
	// encoded is used during the serialization process to prevent encoding a control multiple times.
	encoded bool

	// anything added here needs to be also added to the GOB encoder!
}

// Init is used by Control implementations to initialize the standard control structure. You would only call this if you
// are subclassing one of the standard controls.
// Control implementations should call this immediately after a control is created.
// The Control subclasses should have their own Init function that
// call this superclass function. This Init function sets up a parent-child relationship with the given parent
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


// this supports object oriented features by giving easy access to the virtual function interface.
// Subclasses should provide a duplicate. Calls that implement chaining should return the result of this function.
func (c *Control) this() ControlI {
	return c.Self.(ControlI)
}

// Restore is called after the control has been deserialized. It creates any required data structures
// that are not saved in serialization.
// TODO: Serialization is not yet implemented
func (c *Control) Restore(self ControlI) {
	c.Base.Init(self)
	if c.attributes == nil {
		c.attributes = html.NewAttributes()
	}
	if c.wrapperAttributes == nil {
		c.wrapperAttributes = html.NewAttributes()
	}
}

// ID returns the id assigned to the control. If you do not provide an ID when the control is created,
// the framework will give the control a unique id.
func (c *Control) ID() string {
	return c.id
}

// Extract the control from an interface. This is for package private use, when called through the interface.
func (c *Control) control() *Control {
	return c
}

// ΩPreRender is called by the framework to notify the control that it is about to be drawn. If you
// override it, be sure to also call this parent function as well.
func (c *Control) ΩPreRender(ctx context.Context, buf *bytes.Buffer) error {
	form := c.ParentForm()
	if c.Page() == nil ||
		form == nil ||
		c.Page() != form.Page() {

		return NewError(ctx, "The control can not be drawn because it is not a member of a form that is on the override.")
	}

	if c.wasRendered || c.isRendering {
		return NewError(ctx, "This control has already been drawn.")
	}

	// Because we may be rerendering a parent control, we need to make sure all "child" controls are marked as NOT being on the form
	// before rendering it again.
	if c.children != nil {
		for _, child := range c.children {
			child.control().markOnPage(false)
		}
	}

	// Finally, let's specify that we have begun rendering this control
	c.isRendering = true

	return nil
}

// Draw renders the default control structure into the given buffer. Call this function from your templates
// to draw the control.
func (c *Control) Draw(ctx context.Context, buf *bytes.Buffer) (err error) {
	if err = c.this().ΩPreRender(ctx, buf); err != nil {
		return err
	}

	var h string

	if c.isHidden {
		// We are invisible, but not using a wrapper. This creates a problem, in that when we go visible, we do not know what to replace
		// To fix this, we create an empty, invisible control in the place where we would normally draw
		h = "<span id=\"" + c.this().ID() + "\" style=\"display:none;\" data-grctl></span>\n"
	} else {
		h = c.this().ΩDrawTag(ctx)
	}

	if !config.Minify && GetContext(ctx).RequestMode() != Ajax {
		s := html.Comment(fmt.Sprintf("Control Type:%s, Id:%s", c.Type(), c.ID())) + "\n"
		buf.WriteString(s)
	}

	if c.wrapper != nil && !c.isHidden {
		c.wrapper.ΩWrap(ctx, c.this(), h, buf)
	} else {
		buf.WriteString(h)
	}

	response := c.ParentForm().Response()
	c.this().ΩPutCustomScript(ctx, response)
	c.GetActionScripts(response)
	c.this().ΩPostRender(ctx, buf)
	return
}

// ΩPutCustomScript is called by the framework to ask the control to inject any javascript it needs into the form.
// In particular, this is the place where Controls add javascript that transforms the html into a custom javascript control.
// A Control implementation does this by calling functions on the response object.
// This implementation is a stub.
func (c *Control) ΩPutCustomScript(ctx context.Context, response *Response) {

}

// DrawAjax will be called by the frameowkr during an Ajax rendering of the Control. Every Control gets called. Each Control
// is responsible for rendering itself. Some objects automatically render their child objects, and some don't,
// so we detect whether the parent is being rendered, and assume the parent is taking care of rendering for
// us if so.
//
// Override if you want more control over ajax drawing, like if you detect parts of your control that have changed
// and then want to draw only those parts. This will get called on every control on every ajax draw request.
// It is up to you to test the blnRendered flag of the control to know whether the control was already rendered
// by a parent control before drawing here.
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
			for _, entry := range c.attributeScripts {
				response.ExecuteControlCommand(entry.id, entry.f, entry.commands...)
			}
			c.attributeScripts = nil
		}

		if c.wrapper != nil {
			c.wrapper.ΩAjaxRender(ctx, response, c)
		}

		// ask the child controls to potentially render, since this control doesn't need to
		for _, child := range c.children {
			err = child.DrawAjax(ctx, response)
			if err != nil {
				return
			}
		}
	}
	return
}

// ΩPostRender is called by the framework at the end of drawing, and is the place where controls
// do any post-drawing cleanup needed.
func (c *Control) ΩPostRender(ctx context.Context, buf *bytes.Buffer) (err error) {
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

// ΩDrawTag is responsible for drawing the Control's tag itself.
// Control implementations can override this to draw the tag in a different way, or draw more than one tag if
// drawing a compound control.
func (c *Control) ΩDrawTag(ctx context.Context) string {
	// TODO: Implement this with a buffer to reduce string allocations
	var ctrl string

	attributes := c.this().ΩDrawingAttributes()
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
		if err := c.this().ΩDrawInnerHtml(ctx, buf); err != nil {
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
// that can be shown without impacting layout, or that are scripts only. Control implementations that need to
// put these controls in particular locations on the form can override this.
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

// DrawTemplate is used by the framework to draw the Control with a template.
// Controls that use templates should use this function signature for the template. That will override this one, and
// we will then detect that the template was drawn. Otherwise, we detect that no template was defined and it will move
// on to drawing the controls without a template, or just the text if text is defined.
func (c *Control) DrawTemplate(ctx context.Context, buf *bytes.Buffer) (err error) {
	// Don't change this to use some kind of function injection, as such things are not serializable
	return NewFrameworkError(FrameworkErrNoTemplate)
}

// ΩDrawInnerHtml is used by the framework to draw just the inner html of the control, if the control is not a self
// terminating (void) control. Sub-controls can override this.
func (c *Control) ΩDrawInnerHtml(ctx context.Context, buf *bytes.Buffer) (err error) {
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

// DrawChildren renders the child controls into the buffer.
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

// DrawText renders the text of the control, escaping if needed.
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

// HasWrapper returns true if the control has a wrapper.
func (c *Control) HasWrapper() bool {
	return c.wrapper != nil
}

// Wrapper returns the controls wrapper, or nil if the control does not have a wrapper defined.
func (c *Control) Wrapper() WrapperI {
	return c.wrapper
}

// SetAttribute sets an html attribute of the control. You can manually set most any attribute, but be careful
// not to set the id attribute, or any attribute that is managed by the control itself. If you are setting
// a data-* attribute, use SetDataAttribute instead. If you are adding a class to the control, use AddAttributeValue.
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
		// We are recording here that the attribute intends to change. If we are responding to an ajax
		// request, we will send back a command to only change the attribute on the control if the
		// control does not get completely redrawn. If the control is completely redrawn, the new
		// attribute will automatically be drawn, so there would be no need to also send an attribute change command.
		c.AddRenderScript("attr", name, v2)
	}
	return c.this()
}

// SetWrapperAttribute sets an attribute for the tag that wraps the control. Obviously this only works if
// you have defined a wrapper for the control.
func (c *Control) SetWrapperAttribute(name string, val interface{}) ControlI {
	if name == "id" {
		panic("You cannot set the 'id' attribute of a wrapper")
	}

	changed, err := c.wrapperAttributes.SetChanged(name, html.AttributeString(val))
	if err != nil {
		panic(err)
	}

	if changed {
		// The val passed in might be a calculation, so we need to get the ultimate new value
		v2 := c.wrapperAttributes.Get(name)
		c.AddRelatedRenderScript(c.ID() + "_ctl", "attr", name, v2)
	}
	return c.this()
}

// Return the value of a custom attribute. Note that this will not return values that are set only during
// drawing and that are managed by the Control implementation.
func (c *Control) Attribute(name string) string {
	return c.attributes.Get(name)
}

// HasAttribute returns true if the control has the indicated custom attribute defined.
func (c *Control) HasAttribute(name string) bool {
	return c.attributes.Has(name)
}

// ΩDrawingAttributes is called by the framework just before drawing a control, and should
// return a set of attributes that should override those set by the user. This allows controls to set attributes
// that should take precedence over other attributes, and that are critical to drawing the
// tag of the control. This function is designed to only be called by Control implementations.
func (c *Control) ΩDrawingAttributes() *html.Attributes {
	a := html.NewAttributesFrom(c.attributes)
	a.SetID(c.id)                   // make sure the control id is set at a minimum
	a.SetDataAttribute("grctl", "") // make sure control is registered. Overriding controls can put a control name here.

	if c.HasWrapper() {
		c.wrapper.ΩModifyDrawingAttributes(c.this(), a)
	}

	if c.isRequired {
		a.Set("aria-required", "true")
	}

	return a
}

// WrapperAttributes returns the actual attributes for the wrapper. Changes WILL be remembered so that subsequent ajax
// drawing will draw the wrapper correctly. However, it is up to you to refresh the control if you change anything.
func (c *Control) WrapperAttributes() *html.Attributes {
	return c.wrapperAttributes
}

// SetDataAttribute will set a data-* attribute. You do not need to include the "data-" in the name, it will be added
// automatically.
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

// AddAttributeValue will add a class or classes to the control. If adding multiple classes at once, separate them with
// a space.
func (c *Control) AddClass(class string) ControlI {
	if changed := c.attributes.AddClassChanged(class); changed {
		v2 := c.attributes.Class()
		c.AddRenderScript("attr", "class", v2)
	}
	return c.this()
}

// RemoveClass will remove the named class from the control.
func (c *Control) RemoveClass(class string) ControlI {
	if changed := c.attributes.RemoveClass(class); changed {
		v2 := c.attributes.Class()
		c.AddRenderScript("attr", "class", v2)
	}
	return c.this()
}

// AddWrapperClass will add a class or classes to the control's wrapper, if one is defined. Separate multiple
// classes with a space.
func (c *Control) AddWrapperClass(class string) ControlI {
	if changed := c.wrapperAttributes.AddClassChanged(class); changed {
		v2 := c.wrapperAttributes.Class()
		c.AddRelatedRenderScript(c.ID() + "_ctl", "attr", "class", v2)
	}
	return c.this()
}

// AddRenderScript adds a jQuery command to be executed on the next ajax draw.
// These commands allow javascript to change an aspect of the control without
// having to redraw the entire control. This should be used by Control implementations only.
func (c *Control) AddRenderScript(f string, params ...interface{}) {
	c.attributeScripts = append(c.attributeScripts, attributeScriptEntry{id: c.ID(), f: f, commands: params})
}

// AddRelatedRenderScript adds a render script for a related html object. This is primarily used by control implementations.
func (c *Control) AddRelatedRenderScript(id string, f string, params ...interface{}) {
	c.attributeScripts = append(c.attributeScripts, attributeScriptEntry{id: id, f: f, commands: params})
}


// Parent returns the parent control of the control. All controls have a parent, except the Form control.
func (c *Control) Parent() ControlI {
	return c.parent
}

// Children returns the child controls of the control.
func (c *Control) Children() []ControlI {
	return c.children
}

// Remove removes the current control from its parent. After this is done, the control and all its child items will
// not be part of the drawn form, but the child items will still be accessible through the control itself.
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

// RemoveChild removes the given child control from both the control and the form.
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

// RemoveChildren removes all the child controls from this control and the form
func (c *Control) RemoveChildren() {
	for _, child := range c.children {
		child.control().removeChildrenFromPage()
		c.page.removeControl(child.ID())
		child.control().parent = nil
	}
	c.children = nil
}

// SetParent sets the parent of the control. Use this primarily if you are responding to some kind of user
// interface that will move a child Control from one parent Control to another.
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

// Child returns the child control with the given id.
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

// Page returns the page object associated with the control.
func (c *Control) Page() *Page {
	return c.page
}

// Refresh will force the control to be completely redrawn on the next update.
func (c *Control) Refresh() {
	c.isModified = true
}

// SetIsRequired will set whether the control requires a value from the user. Setting it to true
// will cause the Control to check this during validation, and show an appropriate error message if the user
// did not enter a value.
func (c *Control) SetIsRequired(r bool) ControlI {
	c.isRequired = r
	return c.this()
}

// IsRequired returns true if the control requires input from the user to pass validation.
func (c *Control) IsRequired() bool {
	return c.isRequired
}

// ValidationMessage is the currently set validation message that will print with the control. Normally this only
// gets set when a validation error occurs.
func (c *Control) ValidationMessage() string {
	return c.validationMessage
}

// SetValidationError sets the validation error to the given string. It will also handle setting the wrapper class
// to indicate an error. Override if you have a different way of handling errors.
func (c *Control) SetValidationError(e string) {
/*
	Keeping this here to show that these have been considered and rejected.
	We can still set the aria state in validation situations, even if we are not showing a message
	and subclasses might have a special need for validation without a wrapper.

	if !c.HasWrapper() {
		return // Validation only applies if you have a wrapper to show the message
	}
	if c.validationState == ValidationNever {
		panic(fmt.Errorf("control %s has been set to never validate, so you cannot set a validation error message for it", c.ID()))
	}
*/

	if c.validationMessage != e {
		c.validationMessage = e
		if c.wrapper != nil {
			c.wrapper.ΩSetValidationMessageChanged()
			c.wrapper.ΩSetValidationStateChanged()
		}

		if e == "" {
			c.validationState = ValidationWaiting
			c.AddRenderScript("removeAttr", "aria-invalid")
		} else {
			c.validationState = ValidationInvalid
			c.AddRenderScript("attr", "aria-invalid", "true")
		}
	}
}

// ValidationState returns the current ValidationState value.
func (c *Control) ValidationState() ValidationState {
	return c.validationState
}

// SetText sets the text of the control. Not all controls use this value.
func (c *Control) SetText(t string) ControlI {
	if t != c.text {
		c.text = t
		c.Refresh()
	}
	return c.this()
}

// Text returns the text of the control.
func (c *Control) Text() string {
	return c.text
}

// SetLabel sets the text of the label that will be associated with the control. Labels only get rendered by
// wrappers, so if there is no wrapper with the control, no label will be printed.
func (c *Control) SetLabel(n string) ControlI {
	if n != c.label {
		c.label = n
		c.Refresh()
	}
	return c.this()
}

// Label returns the text of the label associated with the control.
func (c *Control) Label() string {
	return c.label
}

// TextIsLabel is used by the drawing routines to determine if the control's text should be wrapped with a label tag.
// This is normally used by checkboxes and radio buttons that use the label tag in a special way.
func (c *Control) TextIsLabel() bool {
	return false
}

// SetInstructions sets the instructions that will be printed with the control. Instructions only get rendered
// by wrappers, so if there is no wrapper, or the wrapper does not render  the instructions, this will not appear.
func (c *Control) SetInstructions(i string) ControlI {
	if i != c.instructions {
		c.instructions = i
		c.Refresh()
	}
	return c.this()
}

// Instructions returns the instructions to be printed with the control
func (c *Control) Instructions() string {
	return c.instructions
}

func (c *Control) markOnPage(v bool) {
	c.isOnPage = v
}

func (c *Control) IsOnPage() bool {
	return c.isOnPage
}


// WasRendered returns true if the control has been rendered.
func (c *Control) WasRendered() bool {
	return c.wasRendered
}

// IsRendering returns true if we are in the process of rendering the control.
func (c *Control) IsRendering() bool {
	return c.isRendering
}

// HasFor is true if the label should have a "for" attribute. Most browsers respond to this by allowing the
// label to be clicked in order to give focus to the control. Not all controls use this.
func (c *Control) HasFor() bool {
	return c.hasFor
}

// SetHasFor sets whether the control's label should have a "for" attribute that points to the Control.
func (c *Control) SetHasFor(v bool) ControlI {
	if v != c.hasFor {
		c.hasFor = v
		c.Refresh()
	}
	return c.this()
}

// SetHasNoSpace tells the control to draw its inner html with no space around it.
// This should generally only be called by control implementations. If this is not set, spaces
// might be added to make the HTML more readable, which can affect some html control types.
func (c *Control) SetHasNoSpace(v bool) ControlI {
	c.hasNoSpace = v
	return c
}

// SetShouldAutoRender sets whether this control will automatically render. AutoRendered controls are drawn
// by the form automatically, after all other controls are drawn, if the control was not drawn in
// some other way. An example of an auto-rendered control would be a dialog box that starts out hidden,
// but then is shown by some user response. Such controls are normally shown by javascript, and are
// absolutely positioned so that they do not effect the layout of the rest of the form.
func (c *Control) SetShouldAutoRender(r bool) {
	c.shouldAutoRender = r
}

// ShouldAutoRender returns true if the control is set up to auto-render.
func (c *Control) ShouldAutoRender() bool {
	return c.shouldAutoRender
}

// On adds an event listener to the control that will trigger the given actions.
// It returns the event for chaining.
func (c *Control) On(e EventI, actions ...action2.ActionI) EventI {
	var isPrivate bool
	c.Refresh() // completely redraw the control. The act of redrawing will turn off old scripts.
	// TODO: Adding scripts should instead just redraw the associated script block. We will need to
	// implement a script block with every control connected by id
	e.addActions(actions...)
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
	e.event().eventID = c.eventCounter
	return e
}

// Off removes all event handlers from the control
func (c *Control) Off() {
	c.events = nil
}

// HasServerAction returns true if one of the actions attached to the given event is a Server action.
func (c *Control) HasServerAction(eventName string) bool {
	for _,e := range c.events {
		if e.Name() == eventName && e.HasServerAction() {
			return true
		}
	}
	return false
}

// GetEvent returns the event associated with the eventName, which corresponds to the javascript
// trigger name.
func (c *Control) GetEvent(eventName string) EventI {
	for _,e := range c.events {
		if e.Name() == eventName {
			return e
		}
	}
	return nil
}



// SetActionValue sets a value that is provided to actions when they are triggered. The value can be a static value
// or one of the javascript.* objects that can dynamically generate values. The value is then sent back to the action
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
			s := e.renderActions(c.this(), id)
			r.ExecuteJavaScript(s, PriorityStandard)
		}
	}

	if c.events != nil {
		for id, e := range c.events {
			s := e.renderActions(c.this(), id)
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
	if c.validationMessage != "" {
		if c.wrapper != nil {
			c.wrapper.ΩSetValidationMessageChanged()
		}
		c.validationMessage = ""
	}
	if c.validationState != ValidationWaiting {
		if c.wrapper != nil {
			c.wrapper.ΩSetValidationStateChanged()
		}
		c.validationState = ValidationWaiting
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
	children := c.Children()
	if children != nil {
		for _, child := range children {
			child.control().updateValues(ctx)
		}
	}
	// Parent is updated after children so that parent can read the state of the children
	// to update any internal caching of the state. Parent can then delete or recreate children
	// as needed.
	c.this().ΩUpdateFormValues(ctx)
}

// ΩUpdateFormValues should be implemented by Control implementations to get their values from the context.
// This is where a Control updates its internal state based on actions by the client.
func (c *Control) ΩUpdateFormValues(ctx *Context) {

}

// doAction is an internal function that the form manager uses to send actions to controls.
func (c *Control) doAction(ctx context.Context) {
	var e EventI
	var ok bool
	var isPrivate bool
	var grCtx = GetContext(ctx)

	if e, ok = c.events[grCtx.eventID]; !ok {
		if e, ok = c.privateEvents[grCtx.eventID]; ok {
			isPrivate = true
		}
	}

	if !ok {
		// This is the situation where we are submitting a form using a button in a browser
		// where javascript has been turned off. We assume we only have a click event on the button
		// and so just grab it.
		var id EventID
		for id,e = range c.events {
			break
		}
		if id == 0 {
			return
		}

	}

	if (e.event().validationOverride != ValidateNone && e.event().validationOverride != ValidateDefault) ||
		(e.event().validationOverride == ValidateDefault && c.this().ValidationType(e) != ValidateNone) {
		c.ParentForm().control().resetValidation()
	}

	if c.passesValidation(ctx, e) {
		log.FrameworkDebug("doAction - triggered event: ", e.String())
		for _, a := range e.getActions() {
			callbackAction := a.(action2.CallbackActionI)
			p := ActionParams{
				ID:        callbackAction.ID(),
				Action:    a,
				ControlId: c.ID(),
			}

			// grCtx.actionValues is a json representation of the action values. We extract the json, but since json does
			// not differentiate between float and int, we will leave all numbers as json.Number types so we can extract later.
			// use javascript.NumberInt() to easily convert numbers in interfaces to int values.
			p.values = grCtx.actionValues
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

// Validate is called by the framework to validate a control, but not the control's children.
// It is designed to be overridden by Control implementations.
// Overriding controls should call the parent version before doing their own validation.
func (c *Control) Validate(ctx context.Context) bool {
	if c.validationState != ValidationNever {

		if c.validationMessage != c.ValidMessage {
			c.validationMessage = c.ValidMessage
			if c.wrapper != nil {
				c.wrapper.ΩSetValidationMessageChanged()
			}
		}
		if c.validationState != ValidationValid {
			c.validationState = ValidationValid
			if c.wrapper != nil {
				c.wrapper.ΩSetValidationStateChanged()
			}
		}
	}
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
// the value can be restored. Call this during control initialization to cause the control to remember what it
// is set to, so that if the user returns to the form, it will keep its value.
// This function is also responsible for restoring the previously saved state of the control,
// so call this only after you have set the default state of a control during creation or initialization.
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
		c.this().ΩMarshalState(state)
		stateKey := c.ParentForm().ID() + ":" + c.ID()
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
			stateStore.Set(stateKey, state)
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

			c.this().ΩUnmarshalState(state)
		}
	}

	if c.children == nil || len(c.children) == 0 {
		return
	}

	for _, child := range c.children {
		child.control().readState(ctx)
	}
}

/* I think to do this you would just reset the control itself.

func (c *Control) ResetSavedState(ctx context.Context) {
	c.resetState(ctx)
}

func (c *Control) resetState(ctx context.Context) {
	var stateStore *maps.Map
	var ok bool

	if c.shouldSaveState {
		i := session.Get(ctx, sessionControlStates)
		if stateStore, ok = i.(*maps.Map); ok {
			key := c.ParentForm().ID() + ":" + c.ID()
			stateStore.Set(key, nil) // we need to notify writeState to remove it, or writeState will just stomp on it
		}
	}

	if c.children == nil || len(c.children) == 0 {
		return
	}

	for _, child := range c.children {
		child.control().resetState(ctx)
	}
}
 */

// ΩMarshalState is a helper function for controls to save their basic state, so that if the form is reloaded, the
// value that the user entered will not be lost. Implementing controls should add items to the given map.
// Note that the control id is used as a key for the state,
// so that if you are dynamically adding controls, you should make sure you give a specific, non-changing control id
// to the control, or the state may be lost.
func (c *Control) ΩMarshalState(m maps.Setter) {
}

// ΩUnmarshalState is a helper function for controls to get their state from the stateStore. To implement it, a control
// should read the data out of the given map. If needed, implemet your own version checking scheme. The given map will
// be guaranteed to have been written out by the same kind of control as the one reading it. Be sure to call the super-class
// version too.
func (c *Control) ΩUnmarshalState(m maps.Loader) {
}

// ΩT is a shortcut for the translator that uses the internal Goradd domain for translations.
func (c *Control) ΩT(message string) string {
	// at this point, there is no need for comments or disambiguation, so we go right to translation

	return i18n.
		Build().
		Domain(i18n.GoraddDomain).
		Lang(c.page.LanguageCode()).
		T(message)
}


// T sends strings to the translator for translation, and returns the translated string. The language is taken from the
// session. See the i18n package for more info on that mechanism.
// Additionally, you can add an i18n.ID() call to add an id to the translation to disambiguate it from similar strings, and
// you can add a i18n.Comment() call to add an extracted comment for the translators. The message string should be a literal
// string and not a variable, so that an extractor can extract it from your source to put it into a translation file.
// This version passes the literal string.
//
// Examples
//   textbox.T("I have %d things", count, i18n.Comment("This will need multiple translations based on the count value"));
//	 textbox.SetLabel(textbox.T("S", i18n.ID("South")));
func (c *Control) T(message string, params... interface{}) string {
	builder, args := i18n.ExtractBuilderFromArguments(params)
	if len(args) > 0 {
		panic("T() cannot have arguments")
	}

	return builder.
		Lang(c.page.LanguageCode()).
		T(message)
}

// TPrintf is like T(), but works like Sprintf, returning the translated string, but sending the arguments to the message
// as if the message was an Sprintf format string. The go/text extractor has code that can do interesting things with
// this kind of string.
func (c *Control) TPrintf(message string, params... interface{}) string {
	builder, args := i18n.ExtractBuilderFromArguments(params)

	return builder.
		Lang(c.page.LanguageCode()).
		Sprintf(message, args...)
}

// SetDisable will set the "disabled" attribute of the control.
func (c *Control) SetDisabled(d bool) {
	c.attributes.SetDisabled(d)
	c.Refresh()
}

// IsDisabled returns true if the disabled attribute is true.
func (c *Control) IsDisabled() bool {
	return c.attributes.IsDisabled()
}

// SetDisplay sets the "display" property of the style attribute of the html control to the given value.
// Also consider using SetVisible. If you use SetDisplay to hide a control, the control will still be
// rendered in html, but the browser will not show it.
func (c *Control) SetDisplay(d string) {
	c.attributes.SetDisplay(d)
	c.Refresh()
}

// IsDisplayed returns true if the control will be displayed.
func (c *Control) IsDisplayed() bool {
	return c.attributes.IsDisplayed()
}

// IsVisible returns whether the control will be drawn.
func (c *Control) IsVisible() bool {
	return !c.isHidden
}

// SetVisible controls whether the Control will be drawn. Controls that are not visible are not rendered in
// html, but rather a hidden stub is rendered as a placeholder in case the control is made visible again.
func (c *Control) SetVisible(v bool) {
	if c.isHidden == v { // these are opposite in meaning
		c.isHidden = !v
		c.Refresh()
	}
}

// SetStyles sets the style attribute of the control to the given values.
func (c *Control) SetStyles(s *html.Style) {
	c.attributes.SetStyles(s)
	c.Refresh() // TODO: Do this with javascript
}

// SetStyle sets a particular property of the style attribute on the control.
func (c *Control) SetStyle(name string, value string) ControlI {
	if changed,_ := c.attributes.SetStyleChanged(name, value); changed {
		c.Refresh() // TODO: Do this with javascript
	}
	return c.this()
}

// RemoveClassesWithPrefix will remove the classes on a control that start with the given string.
// Some CSS frameworks use prefixes to as a kind of namespace for their class tags, and this can
// make it easier to remove a group of classes with this kind of prefix.
func (c *Control) RemoveClassesWithPrefix(prefix string) {
	if c.attributes.RemoveClassesWithPrefix(prefix) {
		c.Refresh() // TODO: Do this with javascript
	}
}

// SetWidthStyle sets the width style property
func (c *Control) SetWidthStyle(w interface{}) ControlI {
	v := html.StyleString(w)
	c.attributes.SetStyle("width", v)
	c.AddRenderScript("css", "width", v) // use javascript to set this value
	return c.this()
}

// SetHeightStyle sets the height style property
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

// SetWillBeValidated indicates to the wrapper whether to save a spot for a validation message or not.
func (c *Control) SetWillBeValidated(v bool) {
	if v {
		c.validationState = ValidationWaiting
	} else {
		c.validationState = ValidationNever
	}
}

// MockFormValue will mock the process of getting a form value from an http response for
// testing purposes. This includes calling ΩUpdateFormValues and Validate on the control.
// It returns the result of the Validate function.
func (c *Control) MockFormValue(value string) bool {
	ctx := NewMockContext()

	grctx := GetContext(ctx)
	grctx.formVars.Set(c.ID(), value)
	c.this().ΩUpdateFormValues(grctx)
	return c.this().Validate(ctx)
}


// GobEncode here is implemented to intercept the GobSerializer to only encode an empty structure. We use this as part
// of our overall serialization stratgey for forms. Controls still need to be registered with gob.
func (c *Control) GobEncode() (data []byte, err error) {
	return
}

func (c *Control) GobDecode(data []byte) (err error) {
	return
}

func (c *Control) MarshalJSON() (data []byte, err error) {
	return
}

func (c *Control) UnmarshalJSON(data []byte) (err error) {
	return
}

type controlEncoding struct {
	Id                    string
	ParentID              string
	Children              []ControlI
	Tag                   string
	IsVoidTag             bool
	HasNoSpace            bool
	Attributes            *html.Attributes
	Text                  string
	TextLabelMode         html.LabelDrawingMode
	HtmlEscapeText        bool
	IsRequired            bool
	IsHidden              bool
	IsOnPage              bool
	ShouldAutoRender      bool
	IsBlock               bool
	Wrapper               WrapperI
	WrapperAttributes     *html.Attributes
	Label                 string
	HasFor                bool
	Instructions          string
	ErrorForRequired      string
	ValidMessage          string
	ValidationMessage     string
	ValidationState       ValidationState
	ValidationType        ValidationType
	ValidationTargets     []string
	BlockParentValidation bool
	ActionValue           interface{}
	Events                EventMap
	PrivateEvents         EventMap
	EventCounter          EventID
	ShouldSaveState       bool
}

// Serialize is used by the framework to serialize a control to be saved in the formstate.
func (c *Control) Serialize(e Encoder) (err error) {
	// TODO: This is in development and not yet used.
	if err = e.Encode(c.id); err != nil {
		return
	}

	e.Encode(len(c.children))
	c.encoded = true	// Make sure circular references in child controls do not encode twice
	for _,child := range c.children {
		err = e.EncodeControl(child)
		if err != nil {
			return
		}
	}

	s := controlEncoding {
		Tag: c.Tag,
		IsVoidTag: c.IsVoidTag,
		HasNoSpace: c.hasNoSpace,
		Attributes: c.attributes,
		Text: c.text,
		TextLabelMode: c.textLabelMode,
		HtmlEscapeText: c.htmlEscapeText,
		IsRequired: c.isRequired,
		IsHidden: c.isHidden,
		IsOnPage: c.isOnPage,
		ShouldAutoRender: c.shouldAutoRender,
		IsBlock: c.isBlock,
		Wrapper: c.wrapper,
		Label: c.label,
		HasFor: c.hasFor,
		Instructions: c.instructions,
		ErrorForRequired: c.ErrorForRequired,
		ValidMessage: c.ValidMessage,
		ValidationMessage: c.validationMessage,
		ValidationState: c.validationState,
		ValidationType: c.validationType,
		ValidationTargets: c.validationTargets,
		BlockParentValidation: c.blockParentValidation,
		ActionValue: c.actionValue,
		Events: c.events,
		PrivateEvents: c.privateEvents,
		EventCounter: c.eventCounter,
		ShouldSaveState: c.shouldSaveState,
	}

	if c.parent !=  nil {
		s.ParentID = c.parent.ID()
	}

	err = e.Encode(&s)

	return
}

// ΩisSerializer is used by the automated control serializer to determine how far down the control chain the control
// has to go before just calling serialize and deserialize
func (c *Control) ΩisSerializer(i ControlI) bool {
	return reflect.TypeOf(c) == reflect.TypeOf(i)
}

// Deserialize is called by the page serializer.
func (c *Control) Deserialize(d Decoder, p *Page) (err error) {
	if err = d.Decode(&c.id); err != nil {
		return
	}

	var count int
	if err = d.Decode(&count); err != nil {
		return
	}

	for i := 0; i < count; i++ {
		var ci ControlI
		if ci, err = d.DecodeControl(p); err != nil {
			return
		}
		c.children = append(c.children, ci)
	}

	var s controlEncoding

	if err = d.Decode(&s); err != nil {
		return
	}

	c.parent = p.GetControl(s.ParentID)
	c.Tag = s.Tag
	c.IsVoidTag = s.IsVoidTag
	c.hasNoSpace = s.HasNoSpace
	c.attributes = s.Attributes
	c.text = s.Text
	c.textLabelMode = s.TextLabelMode
	c.htmlEscapeText = s.HtmlEscapeText
	c.isRequired = s.IsRequired
	c.isHidden = s.IsHidden
	c.isOnPage = s.IsOnPage
	c.shouldAutoRender = s.ShouldAutoRender
	c.isBlock = s.IsBlock
	c.wrapper = s.Wrapper
	c.label = s.Label
	c.hasFor = s.HasFor
	c.instructions = s.Instructions
	c.ErrorForRequired = s.ErrorForRequired
	c.ValidMessage = s.ValidMessage
	c.validationState = s.ValidationState
	c.validationType = s.ValidationType
	c.validationTargets = s.ValidationTargets
	c.blockParentValidation = s.BlockParentValidation
	c.actionValue = s.ActionValue
	c.events = s.Events
	c.privateEvents = s.PrivateEvents
	c.eventCounter = s.EventCounter
	c.shouldSaveState = s.ShouldSaveState

	return
}


// ControlConnectorParams returns a list of options setable by the connector dialog (not currently implemented)
func ControlConnectorParams() *maps.SliceMap {
	m := maps.NewSliceMap()

	// TODO: Add setable options for all controls
	return m
}

func init() {
	gob.Register(&Control{})
}