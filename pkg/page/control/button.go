package control

import (
	"github.com/goradd/goradd/pkg/html"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/action"
	"github.com/goradd/goradd/pkg/page/event"
)

type ButtonI interface {
	page.ControlI
}

// Button is a standard html form submit button. It corresponds to a <button> tag in html.
//
// The default behavior of a submit button is to submit a form. If you have text boxes on your
// form, pressing enter will submit the FIRST button in the form, and so this essentially becomes
// your default button. If you have more than one button, and you want the default button to
// NOT be the first button on the screen, you can handle this in one of two ways:
// - Make sure your default button comes out first in the html, but then use css to change
// the visible order of the buttons. Be sure to also set the tab order if you do this to reflect
// the visible arrangement of the buttons. Or,
// - Create another button as the first button, and give it a display attribute
// of none so that it is not visible. Set its action to be the default action you want.
//
// If you want the button to display an image, simple create an Image control as a child of the button.
type Button struct {
	page.Control
}

// NewButton creates a new standard html button
func NewButton(parent page.ControlI, id string) *Button {
	b := &Button{}
	b.Init(b, parent, id)
	return b
}

// Init is called by subclasses of Button to initialize the button control structure.
func (b *Button) Init(self page.ControlI, parent page.ControlI, id string) {
	b.Control.Init(self, parent, id)
	b.Tag = "button"
	b.SetValidationType(page.ValidateForm) // default to validate the entire form. Can be changed after creation.
}

// SetLabel is an alias for SetText on buttons. Standard buttons do not normally have separate labels.
// Subclasses can redefine this if they use separate labels.
func (b *Button) SetLabel(label string) page.ControlI {
	b.SetText(label)
	return b
}

// On causes the given actions to execute when the given event is triggered.
func (b *Button) On(e page.EventI, actions ...action.ActionI) page.EventI {
	e.Terminating() // prevent default action (override submit)
	b.Control.On(e, actions...)
	return e
}

// ΩDrawingAttributes retrieves the tag's attributes at draw time. You should not normally need to call this, and the
// attributes are disposed of after drawing, so they are essentially read-only.
func (b *Button) ΩDrawingAttributes() *html.Attributes {
	a := b.Control.ΩDrawingAttributes()
	a.SetDataAttribute("grctl", "button")

	a.Set("name", page.HtmlVarAction) // needed for non-javascript posts
	a.Set("value", b.ID())
	a.Set("type", "submit")

	return a
}


// OnSubmit is a shortcut for adding a click event handler that is particular to buttons and button like objects.
// It debounces the click, so that all other events are lost until this event processes. It should generally be used for
// operations that will eventually redirect to a different page. If coupling this with an ajax response, you should
// probably also make the response priority PriorityFinal.
func (b *Button) OnSubmit(actions ...action.ActionI) page.EventI {
	// We delay here to try to make sure any other delayed events are executed first.
	return b.On(event.Click().Terminating().Delay(200).Blocking(), actions...)
}

