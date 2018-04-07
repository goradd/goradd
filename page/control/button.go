package control
import (
	"github.com/spekary/goradd/html"
	"github.com/spekary/goradd/page"
	localPage "goradd/page"
)


type ButtonI interface {
	page.ControlI
}

type Button struct {
	localPage.Control

	isPrimary bool

}

// Creates a new standard html button
func NewButton(parent page.ControlI) *Button {
	b := &Button{}
	b.Init(b, parent)
	return b
}

func (b *Button) Init(self TextboxI, parent page.ControlI) {
	b.Control.Init(self, parent)
	b.Tag = "button"
	b.SetValidationType(page.ValidateForm) // default to validate the entire form. Can be changed after creation.
}


func (b *Button) On(e page.EventI, actions... page.ActionI) {
	e.Terminating()	// prevent default action (page submit)
	b.Control.On(e, actions...)
}

// DrawingAttributes retrieves the tag's attributes at draw time. You should not normally need to call this, and the
// attributes are disposed of after drawing, so they are essentially read-only.
func (b *Button) DrawingAttributes() *html.Attributes {
	a := b.Control.DrawingAttributes()
	a.Set("name", b.Id())	// needed for posts
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

// Use SetText to set the text of the button
