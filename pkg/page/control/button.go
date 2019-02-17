package control

import (
	"github.com/goradd/goradd/pkg/html"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/action"
	"github.com/goradd/goradd/pkg/page/event"
	"reflect"
)

type ButtonI interface {
	page.ControlI
}

// Button is a standard html button. It derives from the button in the override class, allowing you to customize the
// behavior of all buttons in your application.
type Button struct {
	page.Control
	isPrimary bool
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

	a.Set("name", b.ID()) // needed for posts
	if b.isPrimary {
		a.Set("type", "submit")
	} else {
		a.Set("type", "button")
	}

	return a
}

// SetIsPrimary will make this button fire when a return is pressed. It may also style the button
// to show that it is the primary button.
func (b *Button) SetIsPrimary(isPrimary bool) {
	b.isPrimary = isPrimary
	b.Refresh() // redraw
}

// IsPrimary indicates this is the primary button in the form, and will fire when a return is pressed
// on the keyboard.
func (b *Button) IsPrimary() bool {
	return b.isPrimary
}

// Serialize is called by the page serializer. You do not normally need to call this.
func (b *Button) Serialize(e page.Encoder) (err error) {
	if err = b.Control.Serialize(e); err != nil {
		return
	}

	if err = e.Encode(b.isPrimary); err != nil {
		return err
	}
	return
}

// ΩisSerializer is used by the automated control serializer to determine how far down the control chain the control
// has to go before just calling serialize and deserialize
func (b *Button) ΩisSerializer(i page.ControlI) bool {
	return reflect.TypeOf(b) == reflect.TypeOf(i)
}

// Deserialize is called by the page serializer. You do not normally need to call this.
func (b *Button) Deserialize(d page.Decoder, p *page.Page) (err error) {
	if err = b.Control.Deserialize(d, p); err != nil {
		return
	}

	if err = d.Decode(&b.isPrimary); err != nil {
		return
	}
	return
}


// OnSubmit is a shortcut for adding a click event handler that is particular to buttons and button like objects.
// It debounces the click, so that all other events are lost until this event processes. It should generally be used for
// operations that will eventually redirect to a different page. If coupling this with an ajax response, you should
// probably also make the response priority PriorityFinal.
func (b *Button) OnSubmit(actions ...action.ActionI) page.EventI {
	// We delay here to try to make sure any other delayed events are executed first.
	return b.On(event.Click().Terminating().Delay(200).Blocking(), actions...)
}

