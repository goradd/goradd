package control

import (
	"context"
	"github.com/goradd/goradd/pkg/page"
)

type PasswordTextboxI interface {
	TextboxI
}

type PasswordTextbox struct {
	Textbox
}

func NewPasswordTextbox(parent page.ControlI, id string) *PasswordTextbox {
	t := &PasswordTextbox{}
	t.Self = t
	t.Init(parent, id)
	return t
}

func (t *PasswordTextbox) Init(parent page.ControlI, id string) {
	t.Textbox.Init(parent, id)
	t.SetAttribute("autocomplete", "off")
	t.SetType(TextboxTypePassword)
	// Indicate to goradd.js to always post this ajax value
	// We need this since we are not storing the value in the pagestate
	t.SetDataAttribute("grPost", true)
}

// Serialize is used by the framework to serialize the textbox into the pagestate.
//
// This special override prevents the value of the password from ever getting put into the pagestate store.
func (t *PasswordTextbox) Serialize(e page.Encoder) (err error) {
	t.value = ""
	if err = t.Textbox.Serialize(e); err != nil {
		panic(err)
	}
	return
}

func (t *PasswordTextbox) SaveState(ctx context.Context, saveIt bool) {
	panic("do not call SaveState on a password textbox as it would be a security risk")
}

// PasswordTextboxCreator creates a PasswordTextbox.
// Pass it to AddControls of a control, or as a Child of
// a FormFieldWrapper.
type PasswordTextboxCreator struct {
	// ID is the control id of the html widget and must be unique to the page
	ID string
	// Placeholder is the placeholder attribute of the textbox and shows as help text inside the field
	Placeholder string
	// MinLength is the minimum number of characters that the user is required to enter. If the
	// length is less than this number, a validation error will be shown.
	MinLength int
	// MaxLength is the maximum number of characters that the user is required to enter. If the
	// length is more than this number, a validation error will be shown.
	MaxLength int
	// ColumnCount is the number of characters wide the textbox will be, and becomes the width attribute in the tag.
	// The actual width is browser dependent. For better control, use a width style property.
	ColumnCount int
	// Text is the initial value of the textbox. Often its best to load the value in a separate Load step after creating the control.
	Text string

	page.ControlOptions
}

func (c PasswordTextboxCreator) Create(ctx context.Context, parent page.ControlI) page.ControlI {
	ctrl := NewPasswordTextbox(parent, c.ID)
	c.Init(ctx, ctrl)
	return ctrl
}

func (c PasswordTextboxCreator) Init(ctx context.Context, ctrl PasswordTextboxI) {
	// Reuse subclass
	sub := TextboxCreator{
		Placeholder:    c.Placeholder,
		MinLength:      c.MinLength,
		MaxLength:      c.MaxLength,
		ColumnCount:    c.ColumnCount,
		ControlOptions: c.ControlOptions,
		Text:           c.Text,
	}
	sub.Init(ctx, ctrl)
}


// GetPasswordTextbox is a convenience method to return the control with the given id from the page.
func GetPasswordTextbox(c page.ControlI, id string) *PasswordTextbox {
	return c.Page().GetControl(id).(*PasswordTextbox)
}

func init() {
	page.RegisterControl(&PasswordTextbox{})
}
