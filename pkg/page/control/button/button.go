// Package button includes button and button-like controls that are clickable, including things
// that toggle.
package button

import (
	"context"

	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/action"
	"github.com/goradd/goradd/pkg/page/event"
	"github.com/goradd/html5tag"
)

type ButtonI interface {
	page.ControlI
	SetLabel(label string) ButtonI
	OnSubmit(action action.ActionI) ButtonI
	OnClick(action action.ActionI) ButtonI
	SetIsPrimary(bool) ButtonI
}

// Button is a standard html form button. It corresponds to a <button> tag in html.
//
// By default, we set the "type" attribute of the button to "button". This will prevent
// the button from submitting the form when the user presses the return key.
// To choose which button will submit on a return, call SetIsPrimary() or set the "type" attribute to "submit".
//
// If multiple "submit" buttons are on the page, the default behavior
// will occur if there are text boxes on the
// form, and pressing enter will submit the FIRST button in the form as encountered in the HTML.
// Using CSS to alter the placement of the buttons (using float for instance), will not change
// which button the browser considers to be the first.
//
// If you want the button to display an image, create an Image control as a child of the button.
type Button struct {
	page.ControlBase
}

// NewButton creates a new standard HTML button.
//
// Detect clicks by assigning an action to the OnClick method.
func NewButton(parent page.ControlI, id string) *Button {
	b := new(Button)
	b.Self = b
	b.Init(parent, id)
	return b
}

// Init is called by subclasses of Button to initialize the button control structure.
func (b *Button) Init(parent page.ControlI, id string) {
	b.ControlBase.Init(parent, id)
	b.Tag = "button"
	b.SetAttribute("type", "button")
	b.SetValidationType(event.ValidateForm) // default to validate the entire form. Can be changed after creation.
}

func (b *Button) this() ButtonI {
	return b.Self.(ButtonI)
}

// SetLabel is an alias for SetText on buttons. Standard buttons do not normally have separate labels.
// Subclasses can redefine this if they use separate labels.
func (b *Button) SetLabel(label string) ButtonI {
	b.SetText(label)
	return b.this()
}

// SetIsPrimary will set this button to be the default button on the form, which is the button clicked when
// the user presses a return. Some browsers only respond to this when there is a textbox on the screen.
func (b *Button) SetIsPrimary(s bool) ButtonI {
	if s {
		b.SetAttribute("type", "submit")
	} else {
		b.SetAttribute("type", "button")
	}
	return b.this()
}

// On causes the given actions to execute when the given event is triggered.
func (b *Button) On(e *event.Event, action action.ActionI) page.ControlI {
	e.Terminating() // prevent default action (override submit)
	b.ControlBase.On(e, action)
	return b.this()
}

// DrawingAttributes retrieves the tag's attributes at draw time. You should not normally need to call this, and the
// attributes are disposed of after drawing, so they are essentially read-only.
func (b *Button) DrawingAttributes(ctx context.Context) html5tag.Attributes {
	a := b.ControlBase.DrawingAttributes(ctx)
	a.SetData("grctl", "button")

	a.Set("name", page.HtmlVarAction) // needed for non-javascript posts
	a.Set("value", b.ID())
	//a.Set("type", "submit")

	return a
}

// OnSubmit is a shortcut for adding a click event handler that is particular to buttons and button like objects.
// It debounces the click, so that all other events are lost until this event processes. It should generally be used for
// operations that will redirect to a different page. If coupling this with an ajax response, you should
// probably also make the response priority PriorityFinal.
func (b *Button) OnSubmit(action action.ActionI) ButtonI {
	// We delay here to try to make sure any other delayed events are executed first.
	b.this().On(event.Click().Delay(200).Blocking(), action)
	return b.this()
}

// OnClick is a shortcut for adding a click event handler.
//
// If your handler causes a new page to load, you should consider using OnSubmit instead.
func (b *Button) OnClick(action action.ActionI) ButtonI {
	b.this().On(event.Click(), action)
	return b.this()
}

// ButtonCreator is the initialization structure for declarative creation of buttons
type ButtonCreator struct {
	// ID is the control id
	ID string
	// Text is the text displayed in the button
	Text string
	// Set IsPrimary to true to make this the default button on the page
	IsPrimary bool
	// OnSubmit is the action to take when the button is submitted. Use this specifically
	// for buttons that move to other pages or processes transactions, as it debounces the button
	// and waits until all other actions complete
	OnSubmit action.ActionI
	// OnClick is an action to take when the button is pressed. Do not specify both
	// a OnClick and OnSubmit.
	OnClick action.ActionI
	// ValidationType controls how the page will be validating when the button is clicked.
	ValidationType event.ValidationType
	// ControlOptions are additional options that are common to all controls.
	page.ControlOptions
}

// Create is called by the framework to create a new control from the Creator. You
// do not normally need to call this.
func (c ButtonCreator) Create(ctx context.Context, parent page.ControlI) page.ControlI {
	ctrl := NewButton(parent, c.ID)

	c.Init(ctx, ctrl)
	return ctrl
}

// Init is called by implementations of Buttons to initialize a control with the
// creator. You do not normally need to call this.
func (c ButtonCreator) Init(ctx context.Context, ctrl ButtonI) {
	ctrl.SetLabel(c.Text)
	if c.OnSubmit != nil {
		ctrl.OnSubmit(c.OnSubmit)
	}
	if c.OnClick != nil {
		ctrl.OnClick(c.OnClick)
	}
	if c.ValidationType != event.ValidateDefault {
		ctrl.SetValidationType(c.ValidationType)
	}
	if c.IsPrimary {
		ctrl.SetIsPrimary(true)
	}
	ctrl.ApplyOptions(ctx, c.ControlOptions)
}

// GetButton is a convenience method to return the button with the given id from the page.
func GetButton(c page.ControlI, id string) *Button {
	return c.Page().GetControl(id).(*Button)
}

func init() {
	page.RegisterControl(&Button{})
}
