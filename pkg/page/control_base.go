package page

import (
	"context"
	"fmt"
	"github.com/goradd/gengen/pkg/maps"
	"github.com/goradd/goradd/pkg/base"
	"github.com/goradd/goradd/pkg/config"
	"github.com/goradd/goradd/pkg/html"
	"github.com/goradd/goradd/pkg/i18n"
	"github.com/goradd/goradd/pkg/javascript"
	"github.com/goradd/goradd/pkg/log"
	"github.com/goradd/goradd/pkg/orm/query"
	"github.com/goradd/goradd/pkg/page/action"
	buf2 "github.com/goradd/goradd/pkg/pool"
	"github.com/goradd/goradd/pkg/session"
	"github.com/goradd/goradd/pkg/stringmap"
	"github.com/goradd/goradd/pkg/watcher"
	gohtml "html"
	"io"
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
type ControlTemplateFunc func(ctx context.Context, control ControlI, w io.Writer) error

// ControlWrapperFunc is a template function that specifies how wrappers will draw
type ControlWrapperFunc func(ctx context.Context, control ControlI, ctrl string, w io.Writer) error

// DefaultCheckboxLabelDrawingMode is a setting used by checkboxes and radio buttons to default how they draw labels.
// Some CSS framworks are very picky about whether checkbox labels wrap the control, or sit next to the control,
// and whether the label is before or after the control
var DefaultCheckboxLabelDrawingMode = html.LabelAfter

// The DataConnector moves data between the control and the database model. It is a thin view-model controller
// that can be customized on a per-control basis.
type DataConnector interface {
	// Refresh reads from the model, and puts it into the control
	Refresh(i ControlI, model interface{})
	// Update reads data from the control, and puts it into the model
	Update(i ControlI, model interface{})
}

// DataLoader is an optional interface that DataConnectors can use if they need to load data from the database
// to present a choice of items to the user to select from. The Load method will be called whenever the entire control
// gets redrawn.
type DataLoader interface {
	Load(ctx context.Context) []interface{}
}


// ControlI is the interface that all controls must support. The functions are implemented by the
// ControlBase methods. See the ControlBase method implementation for a description of each method.
type ControlI interface {
	ID() string
	control() *ControlBase
	DrawI

	// Drawing support

	DrawTag(context.Context) string
	DrawInnerHtml(context.Context, io.Writer)
	DrawTemplate(context.Context, io.Writer) error
	PreRender(context.Context, io.Writer)
	PostRender(context.Context, io.Writer)
	ShouldAutoRender() bool
	SetShouldAutoRender(bool)
	DrawAjax(ctx context.Context, response *Response)
	DrawChildren(ctx context.Context, w io.Writer)
	DrawText(ctx context.Context, w io.Writer)

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
	RangeAllChildren(func(ControlI))
	RangeSelfAndAllChildren(func(ControlI))

	// hmtl and css

	SetAttribute(name string, val interface{}) ControlI
	Attribute(string) string
	HasAttribute(string) bool
	ProcessAttributeString(s string) ControlI
	DrawingAttributes(context.Context) html.Attributes
	AddClass(class string) ControlI
	RemoveClass(class string) ControlI
	HasClass(class string) bool
	SetStyles(html.Style)
	SetStyle(name string, value string) ControlI
	SetWidthStyle(w interface{}) ControlI
	SetHeightStyle(w interface{}) ControlI
	Attributes() html.Attributes
	SetDisplay(d string) ControlI
	SetDisabled(d bool)
	IsDisabled() bool

	PutCustomScript(ctx context.Context, response *Response)

	TextIsLabel() bool
	Text() string
	SetText(t string) ControlI
	ValidationMessage() string
	SetValidationError(e string)
	ResetValidation()


	WasRendered() bool
	IsRendering() bool
	IsVisible() bool
	SetVisible(bool)
	IsOnPage() bool

	Refresh()
	NeedsRefresh() bool

	Action(context.Context, ActionParams)
	PrivateAction(context.Context, ActionParams)
	SetActionValue(interface{}) ControlI
	ActionValue() interface{}
	On(e *Event, a action.ActionI) ControlI
	Off()
	WrapEvent(eventName string, selector string, eventJs string, options map[string]interface{}) string
	HasServerAction(eventName string) bool
	HasCallbackAction(eventName string) bool

	// UpdateFormValues is used by the framework to cause the control to retrieve its values from the form
	UpdateFormValues(context.Context)

	Validate(ctx context.Context) bool
	ValidationState() ValidationState
	ValidationType(*Event) ValidationType
	SetValidationType(typ ValidationType) ControlI
	ChildValidationChanged()

	// SaveState tells the control whether to save the basic state of the control, so that when the form is reentered, the
	// data in the control will remain the same. This is particularly useful if the control is used as a filter for the
	// contents of another control.
	SaveState(context.Context, bool)
	MarshalState(m maps.Setter)
	UnmarshalState(m maps.Loader)

	// Shortcuts for translation

	// GT translates strings using the Goradd dictionary.
	GT(format string) string
	// T translates strings using the application provided dictionary.
	T(message string, params ...interface{}) string
	TPrintf(format string, params ...interface{}) string

	// Serialization helpers

	Restore()
	Cleanup()

	// API

	SetIsRequired(r bool) ControlI

	Serialize(e Encoder) (err error)
	Deserialize(d Decoder) (err error)

	ApplyOptions(ctx context.Context, o ControlOptions)
	AddControls(ctx context.Context, creators ...Creator)

	DataConnector() DataConnector
	SetDataConnector(d DataConnector) ControlI
	RefreshData(data interface{})
	UpdateData(data interface{})

	WatchDbTables(ctx context.Context, nodes... query.NodeI)
	WatchDbRecord(ctx context.Context, n query.NodeI, pk string)
	WatchChannel(ctx context.Context, channel string)
}

type attributeScriptEntry struct {
	id       string        // id of the object to execute the command on. This should be the id of the control, or a a related html object.
	f        string        // the  function to call
	commands []interface{} // parameters to the function
}

// ControlBase is the basis for UI controls and widgets in goradd. It corresponds to a standard html form object or tag, or a custom javascript
// widget. A Control renders a tag and everything inside of the tag, but can also include a wrapper which associates
// a label, instructions and error messages with the tag. A Control can also associate javascript
// with itself to make sure the javascript is loaded on the page when the control is drawn, and can render
// javascript that will initialize a custom javascript widget.
//
// A Control can have child Controls. It
// can either allow the framework to automatically draw the child Controls as part of the inner-html of
// the ControlBase, can use a template to draw the Child controls, or manually draw them. The ControlBase is part
// of a hierarchical tree structure, with the Form being the root of the tree.
//
// A Control is part of a system that will reflect the state of the control between the client and server.
// When a user updates a control in the browser and performs an action that requires a response from the
// server, the goradd javascript will gather up all the changes in the form and send those to the server.
// The control can read those values and update its own internal state, so that from the perspective
// of the programmer referring to the control, the values in the ControlBase are the same as what the user sees in a browser.
//
// This ControlBase struct is a mixin that all controls should use. You would not normally create a ControlBase directly,
// but rather create one of the "subclasses" of ControlBase. See the control package for Controls that implement
// standard html widgets.
type ControlBase struct {
	base.Base

	// id is the id passed to the control when it is created, or assigned automatically if empty.
	id string
	// page is a pointer to the page that encloses the entire control tree.
	page *Page

	// parentId is the id of the immediate parent control of this control. Only the form object will not have a parent.
	// We use the id here to prevent a memory leak if we remove the control from the form.
	parentId string
	// children are the child controls that belong to this control. They are cached for speed, and to allow
	// children of controls to be accessed even when the control is not part of the form.
	children []ControlI // Child controls

	// Tag is text of the tag that will enclose the control, like "div" or "input"
	Tag string
	// IsVoidTag should be true if the tag should not have a closing tag, like "img"
	IsVoidTag bool
	// hasNoSpace is for special situations where we want no space between this and the next tag. Spans in particular may need this.
	hasNoSpace bool
	// attributes are the collection of custom attributes to apply to the control. This does not include all the
	// attributes that will be drawn, as some are added temporarily just before drawing by GetDrawingAttributes()
	attributes html.Attributes
	// text is a multi purpose string that can be button text, inner text inside of tags, etc. depending on the control.
	text string
	// textLabelMode describes how to draw the internal label
	textLabelMode html.LabelDrawingMode
	// textIsHtml will prevent the text output from being escaped
	textIsHtml bool

	// attributeScripts are commands to send to our javascript to redraw portions of the control via ajax.
	attributeScripts []attributeScriptEntry

	// isRequired indicates that we will require a value during validation
	isRequired bool
	// isHidden indicates that we will not draw the control, but rather an invisible placeholder for the control.
	isHidden bool
	// isOnPage indicates we have drawn the control at some point in the past
	isOnPage bool
	// shouldAutoRender indicates that we will eventually draw the control even if it is not drawn directly.
	shouldAutoRender bool

	// internal status functions. Do not serialize.

	// isModified will cause the control to redraw as part of the response.
	isModified bool
	// isRendering is true when we are in the middle of rendering the control.
	isRendering bool
	// wasRendered indicates that the page was drawn during the current response.
	wasRendered bool

	// isBlock is true to use a div for the wrapper, false for a span
	isBlock bool

	// ErrorForRequired is the error that will display if a control value is required but not set.
	ErrorForRequired string

	// ValidMessage is the message to display if the control has successfully been validated.
	// Leave blank if you don't want a message to show when valid.
	// Can be useful to contrast between invalid and valid controls in a busy form.
	ValidMessage string
	// validationMessage is the current validation message that will display when drawing the control
	// This gets copied from ValidMessage at drawing time if the control is in an invalid state
	validationMessage string
	// validationState is the current validation state of the control, and will effect how the control is drawn.
	validationState ValidationState
	// validationType indicates how the control will validate itself. See ValidationType for a description.
	validationType ValidationType
	// validationTargets is the list of control IDs to target validation
	validationTargets []string
	// This blocks a parent from validating this control. Useful for dialogs, and other situations where sub-controls should control their own space.
	blockParentValidation bool

	// actionValue is the value that will be provided as the ControlValue for any actions that are triggered by this control.
	actionValue interface{}
	// events are all the events added by the control user that the control might trigger
	events EventMap
	// eventCounter is used to generate a unique id for an event to help us route the event through the system.
	eventCounter EventID
	// shouldSaveState indicates that we should save parts of our state into a session variable so that if
	// the client should come back to the form, we will attempt to restore the state of the control. The state
	// in this situation would be the user's input, so text in a textbox, or the selection from a list.
	shouldSaveState bool
	// encoded is used during the serialization process to prevent encoding a control multiple times.
	encoded bool

	dataConnector DataConnector

	watchedKeys map[string]string

	// anything added here needs to be also added to the GOB encoder!
}

// Init is used by ControlBase implementations to initialize the standard control structure. You would only call this if you
// are subclassing one of the standard controls.
// ControlBase implementations should call this immediately after a control is created.
// The ControlBase subclasses should have their own Init function that
// call this superclass function. This Init function sets up a parent-child relationship with the given parent
// control, and sets up data structures to use the control in object-oriented ways with virtual functions.
// The id is the control id that will appear as the id in html. Leave blank for the system to create a unique id for you.
func (c *ControlBase) Init(parent ControlI, id string) {
	c.attributes = html.NewAttributes()
	if parent != nil {
		c.page = parent.Page()
		c.id = c.page.GenerateControlID(id)
	}
	c.this().SetParent(parent)
	c.isModified = true
}

// this supports object oriented features by giving easy access to the virtual function interface.
// Subclasses should provide a duplicate. Calls that implement chaining should return the result of this function.
func (c *ControlBase) this() ControlI {
	return c.Self.(ControlI)
}

// Restore is called after the control has been deserialized. It notifies the control tree so that it
// can restore internal pointers.
// TODO: Serialization is not yet implemented
func (c *ControlBase) Restore() {
}

// ID returns the id assigned to the control. If you do not provide an ID when the control is created,
// the framework will give the control a unique id.
func (c *ControlBase) ID() string {
	return c.id
}

// Extract the control from an interface. This is for package private use, when called through the interface.
func (c *ControlBase) control() *ControlBase {
	return c
}

// PreRender is called by the framework to notify the control that it is about to be drawn. If you
// override it, be sure to also call this parent function as well.
func (c *ControlBase) PreRender(ctx context.Context, w io.Writer) {
	form := c.ParentForm()
	if c.Page() == nil ||
		form == nil ||
		c.Page() != form.Page() {

		panic (fmt.Sprintf("Control %s can not be drawn because it is not a member of a form that is on the override.", c.ID()))
	}

	if c.wasRendered || c.isRendering {
		panic(fmt.Sprintf("Control %s has already been drawn.", c.ID()))
	}

	// Because we may be rerendering a parent control, we need to make sure all "child" controls are marked as NOT being on the form
	// before rendering it again.
	for _, child := range c.children {
		child.control().markOnPage(false)
	}

	// Finally, let's specify that we have begun rendering this control
	c.isRendering = true
}

// Draw renders the control structure into the given buffer.
func (c *ControlBase) Draw(ctx context.Context, w io.Writer)  {
	c.this().PreRender(ctx, w)

	var h string

	if c.isHidden {
		// We are invisible, but not using a wrapper. This creates a problem, in that when we go visible, we do not know what to replace
		// To fix this, we create an empty, invisible control in the place where we would normally draw
		h = "<span id=\"" + c.this().ID() + "\" style=\"display:none;\" data-grctl></span>\n"
	} else {
		h = c.this().DrawTag(ctx)
	}

	if !config.Minify && GetContext(ctx).RequestMode() != Ajax {
		s := html.Comment(fmt.Sprintf("ControlBase Type:%s, Id:%s", c.Type(), c.ID())) + "\n"
		if _,err := io.WriteString(w, s); err != nil {panic(err)}
	}

	if _,err := io.WriteString(w, h); err != nil {panic(err)}

	response := c.ParentForm().Response()
	c.this().PutCustomScript(ctx, response)
	c.GetActionScripts(response)
	c.this().PostRender(ctx, w)
	return
}

// PutCustomScript is called by the framework to ask the control to inject any javascript it needs into the form.
// In particular, this is the place where Controls add javascript that transforms the html into a custom javascript control.
// A ControlBase implementation does this by calling functions on the response object.
// This implementation is a stub.
func (c *ControlBase) PutCustomScript(ctx context.Context, response *Response) {

}

// NeedsRefresh returns true if the control needs to be completely redrawn. Generally you control
// this by calling Refresh(), but subclasses can implement other ways of detecting this.
func (c *ControlBase) NeedsRefresh() bool {
	return c.isModified
}

// DrawAjax will be called by the framework during an Ajax rendering of the ControlBase. Every ControlBase gets called. Each ControlBase
// is responsible for rendering itself. Some objects automatically render their child objects, and some don't,
// so we detect whether the parent is being rendered, and assume the parent is taking care of rendering for
// us if so.
//
// Override if you want more control over ajax drawing, like if you detect parts of your control that have changed
// and then want to draw only those parts. This will get called on every control on every ajax draw request.
// It is up to you to test the blnRendered flag of the control to know whether the control was already rendered
// by a parent control before drawing here.
func (c *ControlBase) DrawAjax(ctx context.Context, response *Response) {

	if c.this().NeedsRefresh() {
		// simply re-render the control and assume rendering will handle rendering its children

		func() {
			// wrap in a function to get deferred PutBuffer to execute immediately after drawing
			buf := buf2.GetBuffer()
			defer buf2.PutBuffer(buf)

			c.this().Draw(ctx, buf)
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

		// ask the child controls to potentially render, since this control doesn't need to
		for _, child := range c.children {
			if child.IsOnPage() || child.ShouldAutoRender() {
				child.DrawAjax(ctx, response)
			}
		}
	}
	return
}

// PostRender is called by the framework at the end of drawing, and is the place where controls
// do any post-drawing cleanup needed.
func (c *ControlBase) PostRender(ctx context.Context, w io.Writer) {
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

// DrawTag is responsible for drawing the ControlBase's tag itself.
// ControlBase implementations can override this to draw the tag in a different way, or draw more than one tag if
// drawing a compound control.
func (c *ControlBase) DrawTag(ctx context.Context) string {
	// TODO: Implement this with a buffer to reduce string allocations
	var ctrl string

	log.FrameworkDebug("Drawing tag: " + c.ID())

	attributes := c.this().DrawingAttributes(ctx)

	if c.IsVoidTag {
		ctrl = html.RenderVoidTag(c.Tag, attributes)
	} else {
		buf := buf2.GetBuffer()
		defer buf2.PutBuffer(buf)
		c.this().DrawInnerHtml(ctx, buf)
		c.RenderAutoControls(ctx, buf)
		if c.Tag == "" {
			ctrl = buf.String() // a wrapper with no tag. Just inserts functionality and draws its children.
		} else if c.hasNoSpace {
			ctrl = html.RenderTagNoSpace(c.Tag, attributes, buf.String())

		} else {
			ctrl = html.RenderTag(c.Tag, attributes, buf.String())
		}
	}
	return ctrl
}

// RenderAutoControls is an internal function to draw controls marked to autoRender. These are generally used for hidden controls
// that can be shown without impacting layout, or that are scripts only. ControlBase implementations that need to
// put these controls in particular locations on the form can override this.
func (c *ControlBase) RenderAutoControls(ctx context.Context, w io.Writer) {
	// Figuring out where to draw these controls can be difficult.

	for _, ctrl := range c.children {
		if ctrl.ShouldAutoRender() &&
			!ctrl.WasRendered() {

			ctrl.Draw(ctx, w)
		}
	}
	return
}

// DrawTemplate is used by the framework to draw the ControlBase with a template.
// Controls that use templates should use this function signature for the template. That will override this one, and
// we will then detect that the template was drawn. Otherwise, we detect that no template was defined and it will move
// on to drawing the controls without a template, or just the text if text is defined.
func (c *ControlBase) DrawTemplate(ctx context.Context, w io.Writer) (err error) {
	// Don't change this to use some kind of function injection, as such things are not serializable
	return NewFrameworkError(FrameworkErrNoTemplate)
}

// DrawInnerHtml is used by the framework to draw just the inner html of the control, if the control is not a self
// terminating (void) control. Sub-controls can override this.
func (c *ControlBase) DrawInnerHtml(ctx context.Context, w io.Writer) {
	if err := c.this().DrawTemplate(ctx, w); err == nil {
		return
	} else if appErr, ok := err.(FrameworkError); !ok || appErr.Err != FrameworkErrNoTemplate {
		panic(err)
	}
	// No template found, so draw children instead

	if c.children != nil && len(c.children) > 0 {
		c.this().DrawChildren(ctx, w)
		return
	}

	c.this().DrawText(ctx, w)

	return
}

// DrawChildren renders the child controls that have not yet been drawn into the buffer.
func (c *ControlBase) DrawChildren(ctx context.Context, w io.Writer) {
	for _, child := range c.children {
		if !child.WasRendered() {
			child.Draw(ctx, w)
		}
	}
	return
}

// DrawText renders the text of the control, escaping if needed.
func (c *ControlBase) DrawText(ctx context.Context, w io.Writer) {
	if c.text != "" {
		text := c.text

		if !c.textIsHtml {
			text = gohtml.EscapeString(text)
		}
		if _,err := io.WriteString(w, text); err != nil {panic(err)}
	}
	return
}

// SetAttribute sets an html attribute of the control. You can manually set most any attribute, but be careful
// not to set the id attribute, or any attribute that is managed by the control itself. If you are setting
// a data-* attribute, use SetDataAttribute instead. If you are adding a class to the control, use AddAttributeValue.
func (c *ControlBase) SetAttribute(name string, val interface{}) ControlI {
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

// Return the value of a custom attribute. Note that this will not return values that are set only during
// drawing and that are managed by the ControlBase implementation.
func (c *ControlBase) Attribute(name string) string {
	return c.attributes.Get(name)
}

// HasAttribute returns true if the control has the indicated custom attribute defined.
func (c *ControlBase) HasAttribute(name string) bool {
	return c.attributes.Has(name)
}

// DrawingAttributes is called by the framework just before drawing a control, and should
// return a set of attributes that should override those set by the user. This allows controls to set attributes
// that should take precedence over other attributes, and that are critical to drawing the
// tag of the control. This function is designed to only be called by ControlBase implementations.
func (c *ControlBase) DrawingAttributes(ctx context.Context) html.Attributes {
	a := html.NewAttributesFrom(c.attributes)
	a.SetID(c.id)                   // make sure the control id is set at a minimum
	a.SetDataAttribute("grctl", "") // make sure control is registered. Overriding controls can put a control name here.

	if c.isRequired {
		a.Set("aria-required", "true")
	}

	channels := stringmap.JoinStrings(c.watchedKeys, "=", ";")

	if channels != "" {
		a.SetDataAttribute("grWatch", channels)
	}

	return a
}

// SetDataAttribute will set a data-* attribute. You do not need to include the "data-" in the name, it will be added
// automatically.
func (c *ControlBase) SetDataAttribute(name string, val interface{}) {
	var v string
	var ok bool

	if v, ok = val.(string); !ok {
		v = fmt.Sprint(v)
	}

	changed, err := c.attributes.SetDataAttributeChanged(name, v)
	if err != nil {
		panic(err)
	}

	if changed {
		c.AddRenderScript("data", name, v) // Use the data method to set the data during ajax requests
	}
}

func (c *ControlBase) MergeAttributes(a html.Attributes) ControlI {
	c.attributes.Merge(a)
	return c.this()
}

// ProcessAttributeString is used by the drawing template to let you set attributes in the draw tag.
// Attributes are of the form `name="value"`.
func (c *ControlBase) ProcessAttributeString(s string) ControlI {
	if s != "" {
		c.attributes.Merge(s)
	}
	return c.this()
}


// AddAttributeValue will add a class or classes to the control. If adding multiple classes at once, separate them with
// a space.
func (c *ControlBase) AddClass(class string) ControlI {
	if changed := c.attributes.AddClassChanged(class); changed {
		// Note here. We cannot just draw the class, because DrawingAttributes might return
		// a class, and DrawingAttributes requires a context. So we coordinate with goradd.js
		// to be able to add and remove a class.
		c.AddRenderScript("class", "+" + class)
	}
	return c.this()
}

// RemoveClass will remove the named class from the control.
func (c *ControlBase) RemoveClass(class string) ControlI {
	if changed := c.attributes.RemoveClass(class); changed {
		c.AddRenderScript("class", "-" + class)
	}
	return c.this()
}

// HasClass returns true if the class has been assigned to the control from the GO side. We do not currently detect
// class changes done in javascript.
func (c *ControlBase)HasClass(class string) bool {
	return c.attributes.HasClass(class)
}


// Attributes returns a pointer to the attributes of the control. Use this with caution.
// Some controls setup attributes at initialization time, so you could potentially write over those.
// Also, if you change attributes during an ajax call, the changes will not be reflected unless you redraw
// the control. The primary use for this function is to allow controls to set up attributes during initialization.
func (c *ControlBase) Attributes() html.Attributes {
	return c.attributes
}

// AddRenderScript adds a javascript command to be executed on the next ajax draw.
// These commands allow javascript to change an aspect of the control without
// having to redraw the entire control. This should be used by ControlBase implementations only.
func (c *ControlBase) AddRenderScript(f string, params ...interface{}) {
	c.attributeScripts = append(c.attributeScripts, attributeScriptEntry{id: c.ID(), f: f, commands: params})
}

// AddRelatedRenderScript adds a render script for a related html object. This is primarily used by control implementations.
func (c *ControlBase) AddRelatedRenderScript(id string, f string, params ...interface{}) {
	c.attributeScripts = append(c.attributeScripts, attributeScriptEntry{id: id, f: f, commands: params})
}

// Parent returns the parent control of the control. All controls have a parent, except the Form control.
func (c *ControlBase) Parent() ControlI {
	if c.Page().HasControl(c.parentId) {
		return c.Page().GetControl(c.parentId)
	}
	return nil
}

// Children returns the child controls of the control.
func (c *ControlBase) Children() []ControlI {
	return c.children
}

// RangeAllChildren recursively calls the given function on every child control and subcontrol.
// It calls the function on the child controls of each control first, and then on the control itself.
func (c *ControlBase) RangeAllChildren(f func(child ControlI)) {
	for _, child := range c.children {
		child.RangeAllChildren(f)
		f(child)
	}
}

// RangeSelfAndAllChildren recursively calls the given function on this control and every child control and subcontrol.
// It calls the function on the child controls of each control first, and then on the control itself.
func (c *ControlBase) RangeSelfAndAllChildren(f func(ctrl ControlI)) {
	c.RangeAllChildren(f)
	f(c.this())
}

// Remove removes the current control from its parent. After this is done, the control and all its child items will
// not be part of the drawn form, but the child items will still be accessible through the control itself.
func (c *ControlBase) Remove() {
	if c.parentId != "" {
		c.Parent().control().removeChild(c.this().ID(), true)
		if !c.shouldAutoRender {
			//c.Refresh() // TODO: Do this through ajax
		}
	} else {
		c.page.removeControl(c.this().ID())
	}
}

// RemoveChild removes the given child control from both the control and the form.
func (c *ControlBase) RemoveChild(id string) {
	c.removeChild(id, true)
}

// removeChild is a private function that will remove a child control from the current control
func (c *ControlBase) removeChild(id string, fromPage bool) {
	for i, v := range c.children {
		if v.ID() == id {
			c.children = append(c.children[:i], c.children[i+1:]...) // remove found item from list
			if fromPage {
				v.control().removeChildrenFromPage()
				c.page.removeControl(id)
			}
			v.control().parentId = ""
			break
		}
	}
}

func (c *ControlBase) removeChildrenFromPage() {
	c.RangeAllChildren(func(child ControlI) {
		c.page.removeControl(child.ID())
	})
}

// RemoveChildren removes all the child controls from this control and the form so that the memory manager can delete them.
func (c *ControlBase) RemoveChildren() {
	for _, child := range c.children {
		child.control().removeChildrenFromPage()
		c.page.removeControl(child.ID())
		child.control().parentId = ""
	}
	c.children = nil
}

// SetParent sets the parent of the control. Use this primarily if you are responding to some kind of user
// interface that will move a child ControlBase from one parent ControlBase to another.
func (c *ControlBase) SetParent(newParent ControlI) {
	if c.parentId == "" {
		c.control().addChildControlsToPage()
	} else {
		c.Parent().control().removeChild(c.ID(), newParent == nil)
		if !c.shouldAutoRender {
			//c.parent.Refresh()
		}
	}
	if newParent != nil {
		c.parentId = newParent.ID()
		c.Parent().control().addChildControl(c.this())
		if !c.shouldAutoRender {
			// TODO: insert into DOM  instead of c.parent.Refresh()
		}
	} else {
		c.parentId = ""
	}
	c.page.addControl(c.this())

	if c.shouldAutoRender && newParent != nil {
		//c.Refresh()
	}

	// TODO: Refresh as needed, but without refreshing the form
}

// Child returns the child control with the given id.
// TODO: This should be a map, both to speed it up, and add the ability to sort it
func (c *ControlBase) Child(id string) ControlI {
	for _, c := range c.children {
		if c.ID() == id {
			return c
		}
	}
	return nil
}

func (c *ControlBase) addChildControlsToPage() {
	for _, child := range c.children {
		child.control().addChildControlsToPage()
		c.page.addControl(child)
	}
}

// Private function called by setParent on parent function
func (c *ControlBase) addChildControl(child ControlI) {
	if c.children == nil {
		c.children = make([]ControlI, 0)
	}
	c.children = append(c.children, child)
}

// ParentForm returns the form object that encloses this control.
func (c *ControlBase) ParentForm() FormI {
	return c.page.Form()
}

// Page returns the page object associated with the control.
func (c *ControlBase) Page() *Page {
	return c.page
}

// Refresh will force the control to be completely redrawn on the next update.
func (c *ControlBase) Refresh() {
	c.isModified = true
}

// SetIsRequired will set whether the control requires a value from the user. Setting it to true
// will cause the ControlBase to check this during validation, and show an appropriate error message if the user
// did not enter a value.
func (c *ControlBase) SetIsRequired(r bool) ControlI {
	c.isRequired = r
	return c.this()
}

// IsRequired returns true if the control requires input from the user to pass validation.
func (c *ControlBase) IsRequired() bool {
	return c.isRequired
}

// ValidationMessage is the currently set validation message that will print with the control. Normally this only
// gets set when a validation error occurs.
func (c *ControlBase) ValidationMessage() string {
	return c.validationMessage
}

// SetValidationError sets the validation error to the given string. It will also handle setting the wrapper class
// to indicate an error. Override if you have a different way of handling errors.
func (c *ControlBase) SetValidationError(e string) {
	if c.validationMessage != e {
		c.validationMessage = e

		if e == "" {
			c.validationState = ValidationWaiting
			c.AddRenderScript("removeAttr", "aria-invalid")
		} else {
			c.validationState = ValidationInvalid
			c.AddRenderScript("attr", "aria-invalid", "true")
		}
		if c.Parent() != nil {
			c.Parent().ChildValidationChanged() // notify parent wrappers
		}
	}
}

func (f *ControlBase) ResetValidation() {
	f.RangeSelfAndAllChildren(func(ctrl ControlI) {
		c := ctrl.control()
		var changed bool
		if c.validationMessage != "" {
			c.validationMessage = ""
			changed = true
		}
		if c.validationState != ValidationWaiting {
			c.validationState = ValidationWaiting
			changed = true
		}
		if changed {
			if p := c.Parent(); p != nil {
				p.ChildValidationChanged()
			}
		}
	})
}


// ChildValidationChanged is sent by the framework when a child control's validation message
// has changed. Parent controls can use this to change messages or attributes in response.
func (c *ControlBase) ChildValidationChanged() {
	if c.Parent() != nil {
		c.Parent().ChildValidationChanged()
	}
}

// ValidationState returns the current ValidationState value.
func (c *ControlBase) ValidationState() ValidationState {
	return c.validationState
}

// SetText sets the text of the control. Not all controls use this value.
func (c *ControlBase) SetText(t string) ControlI {
	if t != c.text {
		c.text = t
		c.Refresh()
	}
	return c.this()
}

// Text returns the text of the control.
func (c *ControlBase) Text() string {
	return c.text
}

// TextIsLabel is used by the drawing routines to determine if the control's text should be wrapped with a label tag.
// This is normally used by checkboxes and radio buttons that use the label tag in a special way.
func (c *ControlBase) TextIsLabel() bool {
	return false
}

func (c *ControlBase) markOnPage(v bool) {
	c.isOnPage = v
}

func (c *ControlBase) IsOnPage() bool {
	return c.isOnPage
}

// WasRendered returns true if the control has been rendered.
func (c *ControlBase) WasRendered() bool {
	return c.wasRendered
}

// IsRendering returns true if we are in the process of rendering the control.
func (c *ControlBase) IsRendering() bool {
	return c.isRendering
}

// SetHasNoSpace tells the control to draw its inner html with no space around it.
// This should generally only be called by control implementations. If this is not set, spaces
// might be added to make the HTML more readable, which can affect some html control types.
func (c *ControlBase) SetHasNoSpace(v bool) ControlI {
	c.hasNoSpace = v
	return c
}

// SetShouldAutoRender sets whether this control will automatically render. AutoRendered controls are drawn
// by the form automatically, after all other controls are drawn, if the control was not drawn in
// some other way. An example of an auto-rendered control would be a dialog box that starts out hidden,
// but then is shown by some user response. Such controls are normally shown by javascript, and are
// absolutely positioned so that they do not effect the layout of the rest of the form.
func (c *ControlBase) SetShouldAutoRender(r bool) {
	c.shouldAutoRender = r
}

// ShouldAutoRender returns true if the control is set up to auto-render.
func (c *ControlBase) ShouldAutoRender() bool {
	return c.shouldAutoRender
}

// On adds an event listener to the control that will trigger the given actions.
// To have a single event fire multiple actions, use action.Group() to combine the actions into one.
func (c *ControlBase) On(e *Event, a action.ActionI) ControlI {
	c.Refresh() // completely redraw the control. The act of redrawing will turn off old scripts.
	// TODO: Adding scripts should instead just redraw the associated script block. We will need to
	// implement a script block with every control connected by id
	e.addAction(a)
	c.eventCounter++

	// Get a new event id
	for {
		if _, ok := c.events[c.eventCounter]; ok {
			c.eventCounter++
		} else {
			break
		}
	}

	if c.events == nil {
		c.events = map[EventID]*Event{}
	}
	c.events[c.eventCounter] = e
	e.eventID = c.eventCounter
	return c.this()
}

// Off removes all event handlers from the control
func (c *ControlBase) Off() {
	for id,e := range c.events {
		if !e.isPrivate() {
			delete(c.events, id)
		}
	}
}

// HasServerAction returns true if one of the actions attached to the given event is a Server action.
func (c *ControlBase) HasServerAction(eventName string) bool {
	for _, e := range c.events {
		if e.Name() == eventName && e.HasServerAction() {
			return true
		}
	}
	return false
}

// HasCallbackAction returns true if one of the actions attached to the given event is a Server action.
func (c *ControlBase) HasCallbackAction(eventName string) bool {
	for _, e := range c.events {
		if e.Name() == eventName && e.HasCallbackAction() {
			return true
		}
	}
	return false
}


// GetEvent returns the event associated with the eventName, which corresponds to the javascript
// trigger name.
func (c *ControlBase) GetEvent(eventName string) *Event {
	for _, e := range c.events {
		if e.Name() == eventName {
			return e
		}
	}
	return nil
}

// SetActionValue sets a value that is provided to actions when they are triggered. The value can be a static value
// or one of the javascript.* objects that can dynamically generate values. The value is then sent back to the action
// handler after the action is triggered.
func (c *ControlBase) SetActionValue(v interface{}) ControlI {
	c.actionValue = v
	return c.this()
}

// ActionValue returns the control's action value
func (c *ControlBase) ActionValue() interface{} {
	return c.actionValue
}

// Action processes actions. Typically, the Action function will first look at the id to know how to handle it.
// This is just an empty implemenation. Sub-controls should implement this.
func (c *ControlBase) Action(ctx context.Context, a ActionParams) {
}

// PrivateAction processes actions that a control sets up for itself, and that it does not want to give the opportunity
// for users of the control to manipulate or remove those actions. Generally, private actions should call their superclass
// PrivateAction function too.
func (c *ControlBase) PrivateAction(ctx context.Context, a ActionParams) {
}

// GetActionScripts is an internal function called during drawing to gather up all the event related
// scripts attached to the control and send them to the response.
func (c *ControlBase) GetActionScripts(r *Response) {
	// Render actions
	if c.events != nil {
		for id, e := range c.events {
			s := e.renderActions(c.this(), id)
			r.ExecuteJavaScript(s, PriorityStandard)
		}
	}
}

// WrapEvent is an internal function to allow the control to customize its treatment of event processing.
func (c *ControlBase) WrapEvent(eventName string, selector string, eventJs string, options map[string]interface{}) string {
	if selector != "" {
		return fmt.Sprintf("g$('%s').on('%s', '%s', function(event, eventData){%s}, %s);", c.ID(), eventName, selector, eventJs, javascript.ToJavaScript(options))
	} else {
		return fmt.Sprintf("g$('%s').on('%s', function(event, eventData){%s}, %s);", c.ID(), eventName, eventJs, javascript.ToJavaScript(options))
	}
}

// UpdateFormValues is used by the framework to cause the control to retrieve its values from the form
func (c *ControlBase) UpdateFormValues(ctx context.Context) {

}



// doAction is an internal function that the form manager uses to send callback actions to controls.
func (c *ControlBase) doAction(ctx context.Context) {
	var e *Event
	var ok bool
	var isPrivate bool
	var grCtx = GetContext(ctx)

	if e, ok = c.events[grCtx.eventID]; ok {
		isPrivate = e.isPrivate()
	}

	if !ok {
		// This is the situation where we are submitting a form using a button in a browser
		// where javascript has been turned off. We assume we only have a click event on the button
		// and so just grab it.
		var id EventID
		for id, e = range c.events {
			break
		}
		if id == 0 {
			return
		}
	}

	if (e.validationOverride != ValidateNone && e.validationOverride != ValidateDefault) ||
		(e.validationOverride == ValidateDefault && c.this().ValidationType(e) != ValidateNone) {
		c.ParentForm().ResetValidation()
	}

	if c.passesValidation(ctx, e) {
		log.FrameworkDebug("doAction - triggered event: ", e.String())
		if callbackAction := e.getCallbackAction(); callbackAction != nil {
			p := ActionParams{
				ID:        callbackAction.ID(),
				Action:    callbackAction,
				ControlId: c.ID(),
			}

			// grCtx.actionValues is a json representation of the action values. We extract the json, but since json does
			// not differentiate between float and int, we will leave all numbers as json.Number types so we can extract later.
			// use javascript.NumberInt() to easily convert numbers in interfaces to int values.
			p.values = grCtx.actionValues

			if c.Page().HasControl(callbackAction.GetDestinationControlID()) {
				dest := c.Page().GetControl(callbackAction.GetDestinationControlID())
				if isPrivate {
					if log.HasLogger(log.FrameworkDebugLog) {
						log.FrameworkDebugf("doAction - PrivateAction, DestId: %s, ActionId: %d, Action: %s, TriggerId: %s",
							dest.ID(), p.ID, reflect.TypeOf(p.Action).String(), p.ControlId)
					}
					dest.PrivateAction(ctx, p)
				} else {
					if log.HasLogger(log.FrameworkDebugLog) {
						log.FrameworkDebugf("doAction - Action, DestId: %s, ActionId: %d, Action: %s, TriggerId: %s",
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
func (c *ControlBase) SetBlockParentValidation(block bool) {
	c.blockParentValidation = block
}

// SetValidationType specifies how this control validates other controls. Typically its either ValidateNone or ValidateForm.
// ValidateForm will validate all the controls on the form.
// ValidateSiblingsAndChildren will validate the immediate siblings of the target controls and their children
// ValidateSiblingsOnly will validate only the siblings of the target controls
// ValidateTargetsOnly will validate only the specified target controls
func (c *ControlBase) SetValidationType(typ ValidationType) ControlI {
	c.validationType = typ
	return c.this()
}

// ValidationType is an internal function to return the validation type. It allows subclasses to override it.
func (c *ControlBase) ValidationType(e *Event) ValidationType {
	if c.validationType == ValidateNone || c.validationType == ValidateDefault {
		return ValidateNone
	} else {
		return c.validationType
	}
}

// SetValidationTargets specifies which controls to validate, in conjunction with the ValidationType setting,
// giving you very fine-grained control over validation. The default
// is to use just this control as the target.
func (c *ControlBase) SetValidationTargets(controlIDs ...string) {
	c.validationTargets = controlIDs
}

// passesValidation checks to see if the event requires validation, and if so, if it passes the required validation
func (c *ControlBase) passesValidation(ctx context.Context, event *Event) (valid bool) {
	validation := c.this().ValidationType(event)

	if v := event.validationOverride; v != ValidateDefault {
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
			if targets == nil {
				// Target is the form
				targets = []ControlI{c.ParentForm()}
				validation = ValidateForm
			}
		} else {
			targets = []ControlI{c}
		}
	} else {
		if c.validationType == ValidateForm ||
			c.validationType == ValidateContainer {
			panic("Unsupported validation type and target combo.")
		}
		for _, id := range c.validationTargets {
			if c.Page().HasControl(id) {
				targets = append(targets, c.Page().GetControl(id))
			}
		}
	}

	valid = true

	switch validation {
	case ValidateForm:
		valid = c.ParentForm().control().validateSelfAndChildren(ctx)
	case ValidateSiblingsAndChildren:
		for _, t := range targets {
			valid = t.control().validateSiblingsAndChildren(ctx) && valid
		}
	case ValidateSiblingsOnly:
		for _, t := range targets {
			valid = t.control().validateSelfAndSiblings(ctx) && valid
		}
	case ValidateChildrenOnly:
		for _, t := range targets {
			valid = t.control().validateSelfAndChildren(ctx) && valid
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
// It is designed to be overridden by ControlBase implementations.
// Overriding controls should call the parent version before doing their own validation.
func (c *ControlBase) Validate(ctx context.Context) bool {
	if c.validationState != ValidationNever {

		if c.validationMessage != c.ValidMessage {
			c.validationMessage = c.ValidMessage
		}
		if c.validationState != ValidationValid {
			c.validationState = ValidationValid
		}
	}
	return true
}

// validateSelfAndSiblings will validate self and siblings
func (c *ControlBase) validateSelfAndSiblings(ctx context.Context) bool {

	if c.parentId == "" {
		// the one and only form
		return true
	}

	p := c.Parent().control()
	siblings := p.children

	var valid = true
	for _, sibling := range siblings {
		if sibling.IsOnPage() {
			valid = sibling.Validate(ctx) && valid
		}
	}
	return valid
}

func (c *ControlBase) validateSelfAndChildren(ctx context.Context) bool {
	if !c.IsOnPage() {
		return true
	}

	if c.children == nil || len(c.children) == 0 {
		return c.this().Validate(ctx)
	}

	var isValid = true
	for _, child := range c.children {
		if !child.control().blockParentValidation && child.IsOnPage() {
			isValid = child.control().validateSelfAndChildren(ctx) && isValid
		}
	}
	// validate self after validating all children, because self might want to invalidate child items
	// also make sure we validate the parent even if the children are invalid in case the parent is looking at the validation state of the children
	isValid = c.this().Validate(ctx) && isValid

	return isValid
}

func (c *ControlBase) validateSiblingsAndChildren(ctx context.Context) bool {

	if c.parentId == "" {
		return true
	}

	p := c.Parent().control()
	siblings := p.children

	var isValid = true
	for _, sibling := range siblings {
		if !sibling.IsOnPage() {
			continue
		}
		isValid = sibling.control().validateSelfAndChildren(ctx) && isValid
	}

	return isValid
}

// SaveState sets whether the control should save its value and other state information so that if the form is redrawn,
// the value can be restored. Call this during control initialization to cause the control to remember what it
// is set to, so that if the user returns to the form, it will keep its value.
// This function is also responsible for restoring the previously saved state of the control,
// so call this only after you have set the default state of a control during creation or initialization.
func (c *ControlBase) SaveState(ctx context.Context, saveIt bool) {
	c.shouldSaveState = saveIt
	c.readState(ctx)
}

// writeState is an internal function that will write out the state of itself
// This state is used by controls to restore the visual state of the control if the page is returned to. This is helpful
// in situations where a control is used to filter what is shown on the page, you zoom into an item, and then return to
// the parent control. In this situation, you want to see things in the same state they were in, and not have to set up
// the filter all over again.
func (c *ControlBase) writeState(ctx context.Context) {
	var stateStore *maps.Map
	var state *maps.Map
	var ok bool

	if c.shouldSaveState {
		state = maps.NewMap()
		c.this().MarshalState(state)
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
}

// readState is an internal function that will read the state of itself
func (c *ControlBase) readState(ctx context.Context) {
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
}

/* I think to do this you would just reset the control itself.

func (c *ControlBase) ResetSavedState(ctx context.Context) {
	c.resetState(ctx)
}

func (c *ControlBase) resetState(ctx context.Context) {
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

// MarshalState is a helper function for controls to save their basic state, so that if the form is reloaded, the
// value that the user entered will not be lost. Implementing controls should add items to the given map.
// Note that the control id is used as a key for the state,
// so that if you are dynamically adding controls, you should make sure you give a specific, non-changing control id
// to the control, or the state may be lost.
func (c *ControlBase) MarshalState(m maps.Setter) {
}

// UnmarshalState is a helper function for controls to get their state from the stateStore. To implement it, a control
// should read the data out of the given map. If needed, implemet your own version checking scheme. The given map will
// be guaranteed to have been written out by the same kind of control as the one reading it. Be sure to call the super-class
// version too.
func (c *ControlBase) UnmarshalState(m maps.Loader) {
}

// GT translates strings using the Goradd dictionary.
func (c *ControlBase) GT(message string) string {
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
//	 textbox.SetText(textbox.T("S", i18n.ID("South")));
func (c *ControlBase) T(message string, params ...interface{}) string {
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
func (c *ControlBase) TPrintf(message string, params ...interface{}) string {
	builder, args := i18n.ExtractBuilderFromArguments(params)

	return builder.
		Lang(c.page.LanguageCode()).
		Sprintf(message, args...)
}

// SetDisable will set the "disabled" attribute of the control.
func (c *ControlBase) SetDisabled(d bool) {
	c.attributes.SetDisabled(d)
	c.Refresh()
}

// IsDisabled returns true if the disabled attribute is true.
func (c *ControlBase) IsDisabled() bool {
	return c.attributes.IsDisabled()
}

// SetDisplay sets the "display" property of the style attribute of the html control to the given value.
// Also consider using SetVisible. If you use SetDisplay to hide a control, the control will still be
// rendered in html, but the browser will not show it.
func (c *ControlBase) SetDisplay(d string) ControlI {
	c.attributes.SetDisplay(d)
	c.Refresh()
	return c.this()
}

// IsDisplayed returns true if the control will be displayed.
func (c *ControlBase) IsDisplayed() bool {
	return c.attributes.IsDisplayed()
}

// IsVisible returns whether the control will be drawn.
func (c *ControlBase) IsVisible() bool {
	return !c.isHidden
}

// SetVisible controls whether the ControlBase will be drawn. Controls that are not visible are not rendered in
// html, but rather a hidden stub is rendered as a placeholder in case the control is made visible again.
func (c *ControlBase) SetVisible(v bool) {
	if c.isHidden == v { // these are opposite in meaning
		c.isHidden = !v
		c.Refresh()
	}
}

// SetStyles sets the style attribute of the control to the given values.
func (c *ControlBase) SetStyles(s html.Style) {
	c.attributes.SetStyles(s)
	c.Refresh() // TODO: Do this with javascript
}

// SetStyle sets a particular property of the style attribute on the control.
func (c *ControlBase) SetStyle(name string, value string) ControlI {
	if changed, _ := c.attributes.SetStyleChanged(name, value); changed {
		c.Refresh() // TODO: Do this with javascript
	}
	return c.this()
}

// RemoveClassesWithPrefix will remove the classes on a control that start with the given string.
// Some CSS frameworks use prefixes to as a kind of namespace for their class tags, and this can
// make it easier to remove a group of classes with this kind of prefix.
func (c *ControlBase) RemoveClassesWithPrefix(prefix string) {
	if c.attributes.RemoveClassesWithPrefix(prefix) {
		c.Refresh() // TODO: Do this with javascript
	}
}

// SetWidthStyle sets the width style property
func (c *ControlBase) SetWidthStyle(w interface{}) ControlI {
	v := html.StyleString(w)
	c.attributes.SetStyle("width", v)
	c.AddRenderScript("css", "width", v) // use javascript to set this value
	return c.this()
}

// SetHeightStyle sets the height style property
func (c *ControlBase) SetHeightStyle(h interface{}) ControlI {
	v := html.StyleString(h)
	c.attributes.SetStyle("height", v)
	c.AddRenderScript("css", "height", v) // use javascript to set this value
	return c.this()
}

// SetTextIsHtml to true to turn off html escaping of the text output.
func (c *ControlBase) SetTextIsHtml(h bool) ControlI {
	c.textIsHtml = h
	return c.this()
}

// ExecuteWidgetFunction will execute the given JavaScript function on the matching client object, with the given
// parameters. The function is a widget function of the goradd widget wrapper or similar type of object.
func (c *ControlBase) ExecuteWidgetFunction(command string, params ...interface{}) {
	c.ParentForm().Response().ExecuteControlCommand(c.ID(), command, params...)
}

// SetWillBeValidated indicates to the wrapper whether to save a spot for a validation message or not.
func (c *ControlBase) SetWillBeValidated(v bool) {
	if v {
		c.validationState = ValidationWaiting
	} else {
		c.validationState = ValidationNever
	}
}

// DataConnector returns the data connector.
func (c *ControlBase) DataConnector() DataConnector {
	return c.dataConnector
}

// SetDataConnector sets the data connector. The connector must be registered with Gob to be serializable.
func (c *ControlBase) SetDataConnector(d DataConnector) ControlI {
	c.dataConnector = d
	return c.this()
}

func (c *ControlBase) RefreshData(data interface{}) {
	if c.dataConnector != nil {
		c.dataConnector.Refresh(c.this(), data)
	}
}

func (c *ControlBase) UpdateData(data interface{}) {
	if c.dataConnector != nil && c.IsOnPage() {
		c.dataConnector.Update(c.this(), data)
	}
}

// WatchDbTables will add the table nodes to the list of database tables that the control is watching.
// It also adds all the parents of those nodes.
// For example, WatchDbTables(ctx, node.Project().Manager()) will watch the project table and the person table.
func (c *ControlBase) WatchDbTables(ctx context.Context, nodes... query.NodeI) {
	if c.watchedKeys == nil {
		c.watchedKeys = make(map[string]string)
	}
	for _,n := range nodes {
		for {
			c.watchedKeys[watcher.MakeKey(ctx, query.NodeDbKey(n), query.NodeTableName(n), "")] = ""
			n = query.ParentNode(n)
			if n == nil {
				break
			}
		}
	}
}

// WatchDbRecord will watch a specific record. Specify a table node to watch all fields in the record, or a column node
// to watch the changes to a particular field of the table.
func (c *ControlBase) WatchDbRecord(ctx context.Context, n query.NodeI, pk string) {
	if c.watchedKeys == nil {
		c.watchedKeys = make(map[string]string)
	}
	channel := watcher.MakeKey(ctx, query.NodeDbKey(n), query.NodeTableName(n), pk)
	if cn, ok := n.(*query.ColumnNode); ok {
		c.watchedKeys[channel] = query.ColumnNodeDbName(cn)
	} else {
		c.watchedKeys[channel] = ""
	}
}


// WatchChannel allows you to specify any channel to watch that will cause a redraw
func (c *ControlBase) WatchChannel(ctx context.Context, channel string) {
	if c.watchedKeys == nil {
		c.watchedKeys = make(map[string]string)
	}
	c.watchedKeys[channel] = ""
}

// MockFormValue will mock the process of getting a form value from an http response for
// testing purposes. This includes calling UpdateFormValues and Validate on the control.
// It returns the result of the Validate function.
func (c *ControlBase) MockFormValue(value string) bool {
	ctx := NewMockContext()

	grctx := GetContext(ctx)
	grctx.formVars.Set(c.ID(), value)
	c.this().UpdateFormValues(ctx)
	return c.this().Validate(ctx)
}

// Cleanup is called by the framework when a control is being removed from the page cache. It is an opportunity to remove
// any potential circular references in your controls that would prevent the garbage collector from removing the
// control from memory. In particular, references to parent objects would be a problem.
func (c *ControlBase) Cleanup() {
	c.page = nil
}

func (c *ControlBase) MarshalJSON() (data []byte, err error) {
	return
}

func (c *ControlBase) UnmarshalJSON(data []byte) (err error) {
	return
}

type controlEncoding struct {
	Id                    string
	ParentID              string
	ChildIDs              []string
	Tag                   string
	IsVoidTag             bool
	HasNoSpace            bool
	Attributes            html.Attributes
	Text                  string
	TextLabelMode         html.LabelDrawingMode
	TextIsHtml        	  bool
	IsRequired            bool
	IsHidden              bool
	IsOnPage              bool
	ShouldAutoRender      bool
	IsBlock               bool
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
	EventCounter          EventID
	ShouldSaveState       bool
	IsModified			  bool		// For testing framework
	DataConnector 		  DataConnector
	WatchedKeys			  map[string]string
}

// Serialize is used by the framework to serialize a control to be saved in the pagestate.
// It is overridable, and control implementations should call this function first before their
// own serializer.
func (c *ControlBase) Serialize(e Encoder) (err error) {
	s := controlEncoding{
		Id:					   c.id,
		Tag:                   c.Tag,
		IsVoidTag:             c.IsVoidTag,
		HasNoSpace:            c.hasNoSpace,
		Attributes:            c.attributes,
		Text:                  c.text,
		TextLabelMode:         c.textLabelMode,
		TextIsHtml:        	   c.textIsHtml,
		IsRequired:            c.isRequired,
		IsHidden:              c.isHidden,
		IsOnPage:              c.isOnPage,
		ShouldAutoRender:      c.shouldAutoRender,
		IsBlock:               c.isBlock,
		ErrorForRequired:      c.ErrorForRequired,
		ValidMessage:          c.ValidMessage,
		ValidationMessage:     c.validationMessage,
		ValidationState:       c.validationState,
		ValidationType:        c.validationType,
		ValidationTargets:     c.validationTargets,
		BlockParentValidation: c.blockParentValidation,
		ActionValue:           c.actionValue,
		Events:                c.events,
		EventCounter:          c.eventCounter,
		ShouldSaveState:       c.shouldSaveState,
		ParentID:			   c.parentId,
		IsModified:				c.isModified,
		DataConnector:			c.dataConnector,
		WatchedKeys:			c.watchedKeys,
	}

	for _,child := range c.children {
		s.ChildIDs = append(s.ChildIDs, child.ID())
	}

	if err = e.Encode(s); err != nil {
		panic(err)
	}

	return
}

// Deserialize is called by GobDecode to deserialize the control.  It is overridable, and control implementations
// should call this first before calling their own version. However, after deserialization, the control will
// not be ready for use, since its parent, form or child controls still need to be deserialized.
// The Decoded function should be called to fix up the necessary internal pointers.
func (c *ControlBase) Deserialize(d Decoder) (err error) {
	var s controlEncoding

	if err = d.Decode(&s); err != nil {
		panic(err)
	}

	c.id = s.Id
	c.parentId = s.ParentID
	c.Tag = s.Tag
	c.IsVoidTag = s.IsVoidTag
	c.hasNoSpace = s.HasNoSpace
	c.attributes = s.Attributes
	c.text = s.Text
	c.textLabelMode = s.TextLabelMode
	c.textIsHtml = s.TextIsHtml
	c.isRequired = s.IsRequired
	c.isHidden = s.IsHidden
	c.isOnPage = s.IsOnPage
	c.shouldAutoRender = s.ShouldAutoRender
	c.isBlock = s.IsBlock
	c.ErrorForRequired = s.ErrorForRequired
	c.ValidMessage = s.ValidMessage
	c.validationMessage = s.ValidationMessage
	c.validationState = s.ValidationState
	c.validationType = s.ValidationType
	c.validationTargets = s.ValidationTargets
	c.blockParentValidation = s.BlockParentValidation
	c.actionValue = s.ActionValue
	c.events = s.Events
	c.eventCounter = s.EventCounter
	c.shouldSaveState = s.ShouldSaveState
	c.isModified = s.IsModified
	c.dataConnector = s.DataConnector
	c.watchedKeys = s.WatchedKeys

	// This relies on the children being deserialized first, which is taken care of by the page serializer
	for _,id := range s.ChildIDs {
		c.children = append(c.children, c.page.GetControl(id))
	}
	return
}

// EventList is a list of event and action pairs. Use action.Group as the Action to assign multiple actions to
// an event.
type EventList []struct {
	Event *Event
	Action action.ActionI
}

type DataAttributeMap map[string]interface{}

func Nodes(n ...query.NodeI) []query.NodeI {
	return n
}

// ControlOptions are options common to all controls
type ControlOptions struct {
	// Attributes will set the attributes of the control. Use DataAttributes to set data attributes, Styles to set styles, and Class to set the class
	Attributes html.Attributes
	// Attributes will set the attributes of the control. Use DataAttributes to set data attributes, Styles to set styles, and Class to set the class
	DataAttributes DataAttributeMap
	// Styles sets the styles of the control's tag
	Styles html.Style
	// Class sets the class of the control's tag. Prefix a class with "+" to add a class, or "-" to remove a class.
	Class string
	// IsDisabled initializes the control in the disabled state, with a "disabled" attribute
	IsDisabled bool
	// IsRequired is used by the validator. If a value is required, and the control is empty, it will not pass validation.
	IsRequired bool
	// IsHidden initializes this control as hidden. A place holder will be sent in the html so that when the control is shown through ajax, we will know where to put it.
	IsHidden bool
	// On adds events with actions to the control
	On EventList
	// DataConnector is the ViewModel layer that moves data between the control and an attached model.
	DataConnector DataConnector
	// WatchedDbTables lets you specify database nodes to watch for changes. When a record in the table is altered, added or deleted,
	// this control will automatically redraw. To watch a specific record, call WatchDbRecord when you load the control's data.
	WatchedDbTables []query.NodeI
}


func (c *ControlBase) ApplyOptions(ctx context.Context, o ControlOptions) {
	if o.Attributes != nil {
		c.MergeAttributes(o.Attributes)
	}
	for k,v := range o.DataAttributes {
		c.SetDataAttribute(k, v)
	}
	for k,v := range o.Styles {
		c.SetStyle(k, v)
	}
	for _,a := range o.On {
		c.On(a.Event, a.Action)
	}
	if o.Class != "" {
		c.attributes.AddClass(o.Class) // Responds to add and remove class commands
	}
	if o.IsDisabled {
		c.attributes.SetDisabled(o.IsDisabled)
	}
	if o.IsRequired {
		c.isRequired = true
	}
	if o.IsHidden {
		c.isHidden = true
	}
	c.dataConnector = o.DataConnector

	if o.WatchedDbTables != nil {
		c.WatchDbTables(ctx, o.WatchedDbTables...)
	}
}

// Creator is the interface all declarative helpers need to implement
type Creator interface {
	Create(ctx context.Context, parent ControlI) ControlI
}

// AddControls adds subcontrols to a control using a Create function
func (c *ControlBase) AddControls(ctx context.Context, creators ...Creator) {
	for _,creator := range creators {
		creator.Create(ctx, c)
	}
}

// FireTestMarker sends a marker signal to the browser test runner. You would normally send this from some place
// in your application if you want to wait until your app has gotten to that spot. Call WaitMarker on the test
// form to wait for the marker.
func (c *ControlBase) FireTestMarker(marker string) {
	if config.Debug {
		log.FrameworkDebug("Firing test marker: ", marker)
		c.ParentForm().Response().ExecuteJsFunction("goradd.postMarker", marker)
	}
}


