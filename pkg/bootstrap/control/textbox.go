package control

import (
	"context"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/control/textbox"
	"github.com/goradd/html5tag"
)

type TextboxI interface {
	textbox.TextboxI
}

type Textbox struct {
	textbox.Textbox
}

func NewTextbox(parent page.ControlI, id string) *Textbox {
	t := new(Textbox)
	t.Self = t
	t.Init(parent, id)
	return t
}

func (t *Textbox) DrawingAttributes(ctx context.Context) html5tag.Attributes {
	a := t.Textbox.DrawingAttributes(ctx)
	a.AddClass("form-control")
	return a
}

func init() {
	page.RegisterControl(&Textbox{})
}

// TextboxCreator creates a textbox. Pass it to AddControls of a control, or as a Child of
// a FormGroup.
type TextboxCreator struct {
	// ID is the control id of the html widget and must be unique to the page
	ID string
	// Placeholder is the placeholder attribute of the textbox and shows as help text inside the field
	Placeholder string
	// Type is the type attribute of the textbox
	Type string
	// MinLength is the minimum number of characters that the user is required to enter. If the
	// length is less than this number, a validation error will be shown.
	MinLength int
	// MaxLength is the maximum number of characters that the user is required to enter. If the
	// length is more than this number, a validation error will be shown.
	MaxLength int
	// ColumnCount is the number of characters wide the textbox will be, and becomes the width attribute in the tag.
	// The actual width is browser dependent. For better control, use a width style property.
	ColumnCount int
	// RowCount creates a multi-line textarea with the given number of rows. By default, the
	// textbox will expand vertically by this number of lines. Use a height style property for
	// better control of the height of a textbox.
	RowCount int
	// ReadOnly sets the readonly attribute of the textbox, which prevents it from being changed by the user.
	ReadOnly bool
	// SaveState will save the text in the textbox, to be restored if the user comes back to the page.
	// It is particularly helpful when the textbox is being used to filter the results of a query, so that
	// when the user comes back to the page, he does not have to type the filter text again.
	SaveState bool
	// Text is the initial value of the textbox. Generally you would not use this, but rather load the value in a separate Load step after creating the control.
	Text string

	page.ControlOptions
}

// Create is called by the framework to create a new control from the Creator. You
// do not normally need to call this.
func (c TextboxCreator) Create(ctx context.Context, parent page.ControlI) page.ControlI {
	ctrl := NewTextbox(parent, c.ID)

	c.Init(ctx, ctrl)
	return ctrl
}

func (c TextboxCreator) Init(ctx context.Context, ctrl TextboxI) {
	// Reuse subclass
	sub := textbox.TextboxCreator{
		Placeholder:    c.Placeholder,
		Type:           c.Type,
		MinLength:      c.MinLength,
		MaxLength:      c.MaxLength,
		RowCount:       c.RowCount,
		ColumnCount:    c.ColumnCount,
		ReadOnly:       c.ReadOnly,
		SaveState:      c.SaveState,
		Text:           c.Text,
		ControlOptions: c.ControlOptions,
	}
	sub.Init(ctx, ctrl)
}

// GetTextbox is a convenience method to return the control with the given id from the page.
func GetTextbox(c page.ControlI, id string) *Textbox {
	return c.Page().GetControl(id).(*Textbox)
}
