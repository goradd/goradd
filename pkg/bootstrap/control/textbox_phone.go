package control

import (
	"context"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/control/textbox"
	"github.com/goradd/html5tag"
)

type PhoneTextboxI interface {
	textbox.PhoneI
}

type PhoneTextbox struct {
	textbox.Phone
}

func NewPhoneTextbox(parent page.ControlI, id string) *PhoneTextbox {
	t := new(PhoneTextbox)
	t.Self = t
	t.Init(parent, id)
	return t
}

func (t *PhoneTextbox) DrawingAttributes(ctx context.Context) html5tag.Attributes {
	a := t.Phone.DrawingAttributes(ctx)
	a.AddClass("form-control")
	return a
}

// PhoneTextboxCreator creates an phone textbox.
// Pass it to AddControls of a control, or as a Child of
// a FormFieldWrapper.
type PhoneTextboxCreator struct {
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
	// ReadOnly sets the readonly attribute of the textbox, which prevents it from being changed by the user.
	ReadOnly bool
	// SaveState will save the text in the textbox, to be restored if the user comes back to the page.
	// It is particularly helpful when the textbox is being used to filter the results of a query, so that
	// when the user comes back to the page, he does not have to type the filter text again.
	SaveState bool
	// Text is the initial value of the textbox. Often its best to load the value in a separate Load step after creating the control.
	Text string

	page.ControlOptions
}

func (c PhoneTextboxCreator) Create(ctx context.Context, parent page.ControlI) page.ControlI {
	ctrl := NewPhoneTextbox(parent, c.ID)
	c.Init(ctx, ctrl)
	return ctrl
}

func (c PhoneTextboxCreator) Init(ctx context.Context, ctrl PhoneTextboxI) {
	// Reuse subclass
	sub := textbox.PhoneCreator{
		Placeholder:    c.Placeholder,
		Type:           c.Type,
		MinLength:      c.MinLength,
		MaxLength:      c.MaxLength,
		ColumnCount:    c.ColumnCount,
		ReadOnly:       c.ReadOnly,
		ControlOptions: c.ControlOptions,
		SaveState:      c.SaveState,
		Text:           c.Text,
	}
	sub.Init(ctx, ctrl)
}

// GetPhoneTextbox is a convenience method to return the control with the given id from the page.
func GetPhoneTextbox(c page.ControlI, id string) *PhoneTextbox {
	return c.Page().GetControl(id).(*PhoneTextbox)
}

func init() {
	page.RegisterControl(&PhoneTextbox{})
}
