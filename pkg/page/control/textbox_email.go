package control

import (
	"context"
	"fmt"
	"github.com/goradd/goradd/pkg/page"
	"net/mail"
)

type EmailTextboxI interface {
	TextboxI
	SetMaxItemCount(count int) EmailTextboxI
}

// EmailTextbox is a Textbox control that validates for email addresses.
// EmailTextbox can accept multiple addresses separated by commas, and can accept any address in RFC5322 format (Barry Gibbs <bg@example.com>)
// making it useful for people copying addresses out of an email client and pasting into the field.
type EmailTextbox struct {
	Textbox
	maxItemCount int
	items        []*mail.Address
	parseErr     error
}

// NewEmailTextbox creates a new textbox that validates its input as an email address.
// multi will allow the textbox to accept multiple email addresses separated by a comma.
func NewEmailTextbox(parent page.ControlI, id string) *EmailTextbox {
	t := &EmailTextbox{}
	t.Init(t, parent, id)
	t.maxItemCount = 1
	t.SetType(TextboxTypeEmail)
	return t
}

func (t *EmailTextbox) Init(self TextboxI, parent page.ControlI, id string) {
	t.Textbox.Init(self, parent, id)
}

func (t *EmailTextbox) this() EmailTextboxI {
	return t.Self.(EmailTextboxI)
}

func (t *EmailTextbox) SetMaxItemCount(max int) EmailTextboxI {
	t.maxItemCount = max
	if t.maxItemCount > 1 {
		t.SetType(TextboxTypeDefault) // Some browsers cannot handle multiple emails in an email type of text input
	}
	t.Refresh()
	return t.this()
}

func (t *EmailTextbox) Validate(ctx context.Context) bool {
	ret := t.Textbox.Validate(ctx)

	if ret {
		if t.parseErr != nil {
			t.SetValidationError(t.ΩT("Not a valid email address"))
			return false
		} else if len(t.items) > t.maxItemCount {
			if t.maxItemCount == 1 {
				t.SetValidationError(t.ΩT("Enter only one email address"))
			} else {
				t.SetValidationError(fmt.Sprintf(t.ΩT("Enter at most %d email addresses separated by commas"), t.maxItemCount))
			}

			return false
		}
	}
	return true
}

func (t *EmailTextbox) ΩUpdateFormValues(ctx *page.Context) {
	t.Textbox.ΩUpdateFormValues(ctx)
	if t.Text() == "" {
		t.items = nil
		t.parseErr = nil
		return
	}
	t.items, t.parseErr = mail.ParseAddressList(t.Text())
}

// Addresses returns a slice of the individual addresses entered, stripped of any extra text entered.
func (t *EmailTextbox) Addresses() (ret []string) {
	for _, item := range t.items {
		ret = append(ret, item.Address)
	}
	return ret
}

// Use EmailTextboxCreator to create an email textbox.
// Pass it to AddControls of a control, or as a Child of
// a FormFieldWrapper.
type EmailTextboxCreator struct {
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
	// MaxItemCount is the maximum number of email addresses allowed to be entered, separated by commas
	// By default it allows only 1.
	MaxItemCount int

	page.ControlOptions
}

func (c EmailTextboxCreator) Create(ctx context.Context, parent page.ControlI) page.ControlI {
	ctrl := NewEmailTextbox(parent, c.ID)
	c.Init(ctx, ctrl)
	return ctrl
}

func (c EmailTextboxCreator) Init(ctx context.Context, ctrl EmailTextboxI) {
	if c.MaxItemCount > 1 {
		ctrl.SetMaxItemCount(c.MaxItemCount)
	}
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

// GetEmailTextbox is a convenience method to return the control with the given id from the page.
func GetEmailTextbox(c page.ControlI, id string) *EmailTextbox {
	return c.Page().GetControl(id).(*EmailTextbox)
}
