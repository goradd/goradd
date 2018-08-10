package control_base

import (
	"github.com/spekary/goradd/html"
	"github.com/spekary/goradd/page"
	"github.com/spekary/goradd/page/action"
	localPage "goradd-project/override/page"
	"github.com/spekary/goradd/page/event"
)

type ButtonI interface {
	page.ControlI
}

type Button struct {
	localPage.Control

	isPrimary bool
}


func (b *Button) Init(self page.ControlI, parent page.ControlI, id string) {
	b.Control.Init(self, parent, id)
	b.Tag = "button"
	b.SetValidationType(page.ValidateForm) // default to validate the entire form. Can be changed after creation.
}

// SetLabel is an alias for SetText on buttons. Buttons do not normally have separate labels.
func (b *Button) SetLabel(label string) page.ControlI {
	b.SetText(label)
	return b
}


func (b *Button) On(e page.EventI, actions ...action.ActionI) page.EventI {
	e.Terminating() // prevent default action (override submit)
	b.Control.On(e, actions...)
	return e
}

// DrawingAttributes retrieves the tag's attributes at draw time. You should not normally need to call this, and the
// attributes are disposed of after drawing, so they are essentially read-only.
func (b *Button) DrawingAttributes() *html.Attributes {
	a := b.Control.DrawingAttributes()
	a.SetDataAttribute("grctl", "button")

	a.Set("name", b.ID()) // needed for posts
	if b.isPrimary {
		a.Set("type", "submit")
	} else {
		a.Set("type", "button")
	}

	return a
}

func (b *Button) SetIsPrimary(isPrimary bool) {
	b.isPrimary = isPrimary
	b.Refresh() // redraw
}

func (b *Button) IsPrimary() bool {
	return b.isPrimary
}


// OnClick is a shortcut for adding a click event handler that is particular to buttons and button like objects.
// It debounces the click, to prevent potential accidental multiple form submissions.
func (b *Button) OnClick(actions ...action.ActionI) page.EventI {
	return b.On(event.Click().Terminating().Delay(200).Blocking(), actions...)
}


// Use SetText to set the text of the button
