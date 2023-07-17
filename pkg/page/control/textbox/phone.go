package textbox

import (
	"context"
	"strings"

	"github.com/goradd/goradd/pkg/page"
	strings2 "github.com/goradd/goradd/pkg/strings"
)

type PhoneI interface {
	TextboxI
}

// PhoneTextbox is a Textbox control that validates for phone numbers.
type PhoneTextbox struct {
	Textbox
}

// NewPhoneTextbox creates a new textbox that validates its input as an email address.
// multi will allow the textbox to accept multiple email addresses separated by a comma.
func NewPhoneTextbox(parent page.ControlI, id string) *PhoneTextbox {
	t := &PhoneTextbox{}
	t.Init(t, parent, id)
	return t
}

func (t *PhoneTextbox) Init(self any, parent page.ControlI, id string) {
	t.Textbox.Init(self, parent, id)
	t.SetType(TelType)
}

func (t *PhoneTextbox) this() PhoneI {
	return t.Self().(PhoneI)
}

func (t *PhoneTextbox) Validate(ctx context.Context) bool {
	ret := t.Textbox.Validate(ctx)

	if ret {
		text := strings.TrimSpace(t.Text())
		if text == "" {
			return true // parent will have already checked if its required to have a value
		}
		if text[0:1] == "+" {
			// assume international numbers are valid. Some day maybe we have a per country validator
			return true
		} else {
			n := strings2.ExtractNumbers(text)
			if len(n) == 10 { // us based numbers
				// Reformat to U.S. format
				t.SetText("(" + n[0:3] + ") " + n[3:6] + "-" + n[6:])
				return true
			} else {
				t.SetValidationError(t.GT("Invalid phone number"))
				return false
			}
		}
	}
	return ret
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

func (c PhoneTextboxCreator) Init(ctx context.Context, ctrl PhoneI) {
	// Reuse subclass
	sub := TextboxCreator{
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
