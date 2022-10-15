package textbox

import (
	"context"

	"github.com/goradd/goradd/pkg/page"
)

// PasswordI is the interface that defines a Password
type PasswordI interface {
	TextboxI
}

// Password is a Textbox for passwords. It has the "password" type attribute, and it is specially
// controlled so that the password value is never stored in cleartext, either through the pagestate store
// or through the state store.
type Password struct {
	Textbox
}

// NewPasswordTextbox creates a new Password
func NewPasswordTextbox(parent page.ControlI, id string) *Password {
	t := &Password{}
	t.Self = t
	t.Init(parent, id)
	return t
}

// Init is called by the framework to initialize the control. Only subclasses need to call it.
func (t *Password) Init(parent page.ControlI, id string) {
	t.Textbox.Init(parent, id)
	t.SetAttribute("autocomplete", "off")
	t.SetType(PasswordType)
	// Indicate to goradd.js to always post this ajax value
	// We need this since we are not storing the value in the pagestate
	t.SetDataAttribute("grPost", true)
}

// Serialize is used by the framework to serialize the textbox into the pagestate.
//
// This special override prevents the value of the password from ever getting put into the pagestate store.
func (t *Password) Serialize(e page.Encoder) {
	t.value = ""
	t.Textbox.Serialize(e)
	return
}

// SaveState normally is used to save the text of the control to restore it if the page is returned to.
// This version panics, so that you never SaveState on a password text box.
func (t *Password) SaveState(_ context.Context, _ bool) {
	panic("do not call SaveState on a password textbox as it would be a security risk")
}

// PasswordCreator creates a Password.
// Pass it to AddControls of a control, or as a Child of
// a FormFieldWrapper.
type PasswordCreator struct {
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
	// Text is the initial value of the textbox. Often it is best to load the value in a separate Load step after creating the control.
	Text string

	page.ControlOptions
}

// Create is called by the framework to turn the PasswordCreator into a control. You do not
// normally need to call it.
func (c PasswordCreator) Create(ctx context.Context, parent page.ControlI) page.ControlI {
	ctrl := NewPasswordTextbox(parent, c.ID)
	c.Init(ctx, ctrl)
	return ctrl
}

// Init is called by the framework to initialize a newly created Password. You do not
// normally need to call it.
func (c PasswordCreator) Init(ctx context.Context, ctrl PasswordI) {
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
func GetPasswordTextbox(c page.ControlI, id string) *Password {
	return c.Page().GetControl(id).(*Password)
}

func init() {
	page.RegisterControl(&Password{})
}
