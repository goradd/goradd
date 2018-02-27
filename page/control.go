package page

//go:generate hero -source template

import (
	"context"
	"bytes"
	"github.com/spekary/goradd/html"
	"fmt"
	 gohtml "html"
	"github.com/spekary/goradd"
)


const  (
	ValidateNone int = iota
	ValidateAll
	ValidateSiblingsAndChildren
	ValidateSiblingsOnly
)

type ControlTemplateFunc func(ctx context.Context, control ControlI, buffer *bytes.Buffer) error

type ControlWrapperFunc func(ctx context.Context, control ControlI, ctrl string, buffer *bytes.Buffer)



type ControlI interface {
	Id() string
	control() *Control
	DrawI

	// Drawing support
	DrawTag(context.Context, *bytes.Buffer) string
	DrawInnerHtml(context.Context, *bytes.Buffer) error
	DrawTemplate(context.Context, *bytes.Buffer) error
	PreRender(context.Context, *bytes.Buffer) error
	PostRender(context.Context, *bytes.Buffer) error
	ShouldAutoRender() bool

	// Hierarchy functions
	Parent() ControlI
	Children() []ControlI
	SetParent(parent ControlI)
	Remove()
	RemoveChild(id string)
	RemoveChildren()
	Page() PageI
	Form() FormI

	SetAttribute(name string, val interface{})
	Attribute(string) string
	Attributes() *html.Attributes
	WrapperAttributes() *html.Attributes

	HasFor() bool
	SetHasFor(bool) ControlI

	Name() string
	SetName(n string) ControlI
	Text() string
	SetText(t string) ControlI
	ValidationError() string
	SetValidationError(e string)
	Instructions() string
	SetInstructions(string) ControlI

	WasRendered() bool
	IsRendering() bool

	Refresh()

	Action(*ActionParams)
	SetActionValue(interface{})
	ActionValue() interface{}
	On(e EventI, a ...ActionI)
	Off()
	wrapEvent(eventName string, eventJs string) string
	getCustomScript(response *Response)
	getScripts(r *Response)
	resetFlags()

	addChildControlsToPage()
	addChildControl(ControlI)
	markOnPage(bool)
}

type Control struct {
	goradd.Base

	id string
	page PageI							// Page this control is part of
	parent   ControlI					// Parent control
	children []ControlI					// Child controls

	Tag         string
	IsVoidTag   bool                        // tag does not get wrapped with a terminating tag, but just ends instead
	hasNoSpace  bool                        // For special situations where we want no space between this and next tag. Spans in particular may need this.
	attributes  *html.Attributes            // a collection of attributes to apply to the control
	labelMode	html.LabelDrawingMode		// how to draw the label when wrapper draws the name as a label
	text        string                      // multi purpose, can be button text, inner text inside of tags, etc.
	textIsLabel bool                        // special situation like checkboxes where the text should be wrapped in a label as part of the control
	textLabelMode	 html.LabelDrawingMode // describes how to draw this special label

	htmlEncodeText bool                // whether to encode the text output, or send straight text

	attributeScripts []*[]string // commands to send to our javascript to redraw portions of this control via ajax. Allows us to skip drawing the entire control.

	isRequired bool
	isVisible  bool
	isOnPage   bool
	shouldAutoRender bool

	// internal status functions. Do not serialize.
	isModified  bool
	isRendering bool
	wasRendered bool

	isBlockLevel      bool           // true to use a div for the wrapper, false for a span
	wrapper           WrapperI
	wrapperAttributes *html.Attributes
	name              string              // the given name, often used as a label. Not drawn by default, but the wrapper drawing function uses it. Can also get controls by name.

	hasFor	 			bool			// When drawing the label, should it use a for attribute? This is helpful for screen readers and navigation on certain kinds of tags.
	instructions		string			// Instructions, if the field needs extra explanation. You could also try adding a tooltip to the wrapper.
	validationError		string			// The message to display if there was a validation error
	warning				string			// Warning message

	actionValue			interface{}

	events	EventMap
	privateEvents EventMap
	eventCounter EventId
}

func (c *Control) Init (self ControlI, parent ControlI, id string) {
	c.Base.Init(self)
	c.attributes = html.NewAttributes()
	c.wrapperAttributes = html.NewAttributes()
	if parent != nil {
		c.page = parent.Page()
		id = c.page.GenerateControlId(id)
	}
	c.id = id
	self.SetParent(parent)
	c.htmlEncodeText = true // default to encoding the text portion. Explicitly turn this off if you need something else
	c.isVisible = true
	c.labelMode = html.LABEL_BEFORE

}

func (c *Control) this() ControlI {
	return c.Self.(ControlI)
}

func (c *Control) Id() string {
	return c.id
}

// Extract the control from an interface. This is for private use, when called through the interface
func (c *Control)  control() *Control {
	return c
}


func (c *Control) PreRender(ctx context.Context, buf *bytes.Buffer) error {
	form := c.Form()
	if c.Page() == nil ||
		form == nil ||
		c.Page() != form.Page() {

		return NewError(ctx, "The control can not be drawn because it is not a member of a form that is on the page.")
	}

	if c.wasRendered || c.isRendering {
		return NewError(ctx, "This control has already been drawn.")
	}

	// Because we may be re-isRendering a parent control, we need to make sure all "children" controls are marked as NOT being on the page.
	if c.children != nil {
		for _,child := range c.children {
			child.markOnPage(false)
		}
	}

	// Finally, let's specify that we have begun isRendering this control
	c.isRendering = true

	return nil
}

// Draws the default control structure into the given buffer.
func (c *Control) Draw(ctx context.Context, buf *bytes.Buffer) (err error) {
	// TODO: Capture errors and panics, writing what we can to the buffer on error

	if err = c.this().PreRender(ctx, buf); err != nil {
		return err
	}

	var ctrl string

	if !c.isVisible {
		if c.wrapper == nil {
			// We are invisible, but not using a wrapper. This creates a problem, in that when we go visible, we do not know what to replace
			// To fix this, we create an empty, invisible control in the place where we would normally draw
			ctrl = "<span id=\"" + c.this().Id() + "\" style=\"display:none;\"></span>"
		} else {
			ctrl = "" // when going visible, we will redraw the inner text of the wrapper
		}
	} else {
		ctrl = c.this().DrawTag(ctx, buf)
	}

	if c.wrapper != nil {
		if GetContext(ctx).AppContext.requestMode == Ajax {
			if c.parent == nil {
				if !c.isOnPage {
					c.wrapper.Wrap(ctx, ctrl, buf)
				}
			} else {
				if c.parent.WasRendered() || c.parent.IsRendering() {
					c.wrapper.Wrap(ctx, ctrl, buf)
				}
				// otherwise RenderAjax will handle the drawing
			}
		} else {
			c.wrapper.Wrap(ctx, ctrl, buf)
			// TODO: Comment if in debug mode
		}
	} else {
		buf.WriteString(ctrl)
	}

	c.this().PostRender(ctx, buf)
	return
}

func (c *Control) PostRender(ctx context.Context, buf *bytes.Buffer) (err error){
	// Update watcher
	//if ($this->objWatcher) {
	//$this->objWatcher->makeCurrent();
	//}

	c.isRendering = false
	c.wasRendered = true
	c.isOnPage = true
	return
}

// Draw the control tag itself. Override to draw the tag in a different way, or draw more than one tag if
// drawing a compound control.
func (c *Control) DrawTag(ctx context.Context, buf *bytes.Buffer) string {
	var ctrl string

	attributes := html.NewAttributesFrom(c.this().Attributes())
	if c.wrapper == nil {
		if a := c.this().WrapperAttributes(); a != nil {
			attributes.Merge(a)
		}
	}
	attributes.SetId(c.this().Id())

	if c.IsVoidTag {
		ctrl = html.RenderVoidTag(c.Tag, attributes)
	} else {
		buf := new(bytes.Buffer) // TODO: Get this buffer from the buffer pool, or better simply render the tag manually straight to the buffer.
		c.this().DrawInnerHtml(ctx, buf)
		if c.hasNoSpace {
			ctrl = html.RenderTagNoSpace(c.Tag, attributes, buf.String())

		} else {
			ctrl = html.RenderTag(c.Tag, attributes, buf.String())
		}
	}
	return ctrl
}

// Controls that use templates should use this function signature for the template. That will override this one, and
// we will then detect that the template was drawn. Otherwise, we detect that no template was defined and it will move
// on to drawing the controls without a template, or just the text if text is defined.
func (c *Control) DrawTemplate(ctx context.Context, buf *bytes.Buffer) (err error) {
	return NewAppErr(AppErrNoTemplate)
}

// Returns the inner text of the control, if the control is not a self terminating (void) control. Sub-controls can
// override this.
func (c *Control) DrawInnerHtml(ctx context.Context, buf *bytes.Buffer) (err error) {
	if err = c.this().DrawTemplate(ctx, buf); err == nil {
		return
	} else if appErr,ok := err.(AppErr); !ok || appErr.Err != AppErrNoTemplate {
		return
	} else if c.children != nil && len(c.children) > 0 {
		for _, child := range c.children {
			err = child.Draw(ctx, buf)
			if err != nil {
				break
			}
		}
		return
	} else {
		err = nil
		text := c.text

		if c.htmlEncodeText {
			text = gohtml.EscapeString(text)
		}
		buf.WriteString(text)
	}
	return
}


// With sets the wrapper style for the control, essentially setting the wrapper template function that will be used.
func (c *Control) With(w WrapperI) ControlI {
	// We must use a string here, because we may need to serialize the control, and so know how to restore the wrapper function upon
	// deserializing the control.

	c.wrapper = w
	return c.this() // for chaining

}


func (c *Control) SetAttribute(name string, val interface{}) {
	var v string
	var ok bool

	if v,ok = val.(string); !ok {
		v = fmt.Sprintf("%v", v)
	}


	changed, err := c.attributes.SetChanged(name, v)
	if err != nil {
		panic (err)
	}

	if changed {
		c.addRenderScript("attr", name, v)
	}
}

func (c *Control) Attribute(name string) string {
	return c.attributes.Get(name)
}

// Returns the set of attributes. Subclasses can override this and add isRequired attributes. This gets called just before drawing.
func (c *Control) Attributes() *html.Attributes {
	return c.attributes
}

func (c *Control) WrapperAttributes() *html.Attributes {
	return c.wrapperAttributes
}

func (c *Control) SetDataAttribute(name string, val interface{}) {
	var v string
	var ok bool

	if v,ok = val.(string); !ok {
		v = fmt.Sprintf("%v", v)
	}

	changed, err := c.attributes.SetDataAttributeChanged(name, v)
	if err != nil {
		panic (err)
	}

	if changed {
		c.addRenderScript("data", name, v)	// Use the jQuery data method to set the data during ajax requests
	}
}

// Adds a variadic parameter list to the renderScripts array, which is an array of javascript commands to send to the
// browser the next time control drawing happens. These commands allow javascript to change an aspect of the control without
// having to redraw the entire control
func (c *Control) addRenderScript(params ...string) {
	c.attributeScripts = append (c.attributeScripts, &params)
}

// Control Hierarchy Functions
func (c *Control) Parent() ControlI {
	return c.parent
}

func (c *Control) Children() [] ControlI {
	return c.children
}

func (c *Control) Remove() {
	if c.parent != nil {
		c.parent.RemoveChild(c.this().Id())
	} else {
		c.RemoveChildren()
		c.page.removeControl(c.this().Id())
	}
}

func (c *Control) RemoveChild(id string) {
	for i,v := range c.children {
		if v.Id() == id {
			c.children = append(c.children[:i], c.children[i+1:]...) // remove found item from list
			break
		}
	}
}

func (c *Control) RemoveChildren() {
	for _,child := range c.children {
		child.RemoveChildren()
		c.page.removeControl(child.Id())
	}
	c.children = nil
}

func (c *Control) SetParent(newParent ControlI) {
	if c.parent == nil {
		c.addChildControlsToPage()
	}
	c.parent = newParent
	if c.parent != nil {
		c.parent.addChildControl(c.this())
	}
	c.page.addControl(c.this())
}

func (c *Control) addChildControlsToPage() {
	for _,child := range c.children {
		child.addChildControlsToPage()
		c.page.addControl(child)
	}
}

// Private function called by setParent on parent function
func (c *Control) addChildControl(child ControlI) {
	if c.children == nil {
		c.children = make([]ControlI,0)
	}
	c.children = append(c.children, child)
}

func (c *Control) Form() FormI {
	p := c.this()
	for p.Parent() != nil {
		p = p.Parent()
	}
	return p.(FormI)
}

func (c *Control) Page() PageI {
	return c.page
}

// Drawing aids
func (c *Control) Refresh() {
	c.isModified = true
}

func (c *Control) SetRequired(r bool) ControlI {
	c.isRequired = r
	return c.this()
}

func (c *Control) Required() bool {
	return c.isRequired
}

func (c *Control) ValidationError() string {
	return c.validationError
}

func (c *Control) SetValidationError(e string) {
	c.validationError = e
	// TODO: use SetAttribute somehow to set just the validation error text through ajax, and possibly class, so we don't have to redraw the entire control
	c.isModified = true
}

func (c *Control) SetText(t string) ControlI {
	if t != c.text {
		c.text = t
		c.isModified = true
	}
	return c.this()
}

func (c *Control) Text() string {
	return c.text
}

func (c *Control) SetName(n string) ControlI {
	if n != c.name {
		c.name = n
		c.isModified = true
	}
	return c.this()
}

func (c *Control) Name() string {
	return c.name
}

func (c *Control) SetInstructions(i string) ControlI {
	if i != c.instructions {
		c.instructions = i
		c.isModified = true
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
		c.isModified = true
	}
	return c.this()
}

func (c *Control) ShouldAutoRender() bool {
	return c.shouldAutoRender
}

// On adds an event listener to the control that will trigger the given actions
func (c *Control) On(e EventI, actions... ActionI) {
	c.isModified = true
	e.AddActions(actions...)
	c.eventCounter++
	if c.events == nil {
		c.events = map[EventId]EventI{}
	}
	for {
		if _,ok := c.events[c.eventCounter]; ok {
			c.eventCounter ++
		}
	}
	c.events[c.eventCounter] = e
}

// Off removes all event handlers from the control
func (c *Control) Off() {
	c.events = nil
}

// SetActionValue sets a value that is provided to actions when they are triggered. The value can be a static value
// or one of the javascript.* objects that can dynamically generated values. The value is then sent back to the action
// handler after the action is triggered.
func (c *Control) SetActionValue(v interface{}) {
	c.actionValue = v
}

// ActionValue returns the control's action value
func (c *Control) ActionValue() interface{} {
	return c.actionValue
}


// Action processes actions. Typically, the Action function will first look at the id to know how to handle it.
// This is just an empty implemenation. Sub-controls should implement this.
func (c *Control) Action(a *ActionParams) {
}


// getScripts is an internal function called by the form to recursively process javascripts needed by the control that should be sent to the
// browser. There are two main ways to send a response, either by calling functions on the response object, or
// by returning javascript to execute. Returned javascript executes at medium priority in the order given.
// This function gets called when the entire control is redrawn.
func (c *Control) getScripts(r *Response) {
	// Render actions
	if c.privateEvents != nil {
		for id,e := range c.privateEvents {
			s := e.RenderActions(c.this(), id)
			r.executeJavaScript(s, PriorityStandard)
		}
	}

	if c.events != nil {
		for id,e := range c.events {
			s := e.RenderActions(c.this(), id)
			r.executeJavaScript(s, PriorityStandard)
		}
	}

	c.this().getCustomScript(r)

	for _,child := range c.this().Children() {
		child.getScripts(r)
	}

	c.attributeScripts = nil // Entire control was redrawn, so don't need these
}


// getCustomScript is the place custom controls will render their javascript that would transform the html object into a javascript control.
// This supports the style of controls jQuery plugins and jQuery UI is based on. It usually involves some kind of
// jQuery call to a function attached to the object. To return a script, call functions on the response object.
func (c *Control) getCustomScript(r *Response) {

}

// Recursively reset the drawing flags
func (c *Control) resetFlags() {
	c.wasRendered = false
	c.isModified = false
	if c.wrapper != nil {
		c.wrapper.SetModified(false)
	}

	if children := c.this().Children(); children != nil {
		for _,child := range children {
			child.resetFlags()
		}
	}
}

// An internal function to allow the control to customize its treatment of event processing.
func (c *Control) wrapEvent(eventName string, eventJs string) string {
	return fmt.Sprintf("$j('#%s').on('%s', function(event, ui){%s});", c.Id(), eventName, eventJs);
}



