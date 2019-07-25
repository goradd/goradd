package control

import (
	"context"
	"fmt"
	"github.com/goradd/goradd/pkg/page"
	"net/mail"
)

// EmailTextbox is a Text control that validates for email addresses.
// EmailTextbox can accept multiple addresses separated by commas, and can accept any address in RFC5322 format (Barry Gibbs <bg@example.com>)
// making it useful for people copying addresses out of an email client and pasting into the field.
type EmailTextbox struct {
	Textbox
	maxItems int
	items    []*mail.Address
	parseErr error
}

// NewEmailTextbox creates a new textbox that validates its input as an email address.
// multi will allow the textbox to accept multiple email addresses separated by a comma.
func NewEmailTextbox(parent page.ControlI, id string) *EmailTextbox {
	t := &EmailTextbox{}
	t.Init(t, parent, id)
	t.maxItems = 1
	t.SetType(TextboxTypeEmail)
	return t
}

func (t *EmailTextbox) Init(self TextboxI, parent page.ControlI, id string) {
	t.Textbox.Init(self, parent, id)
}

func (t *EmailTextbox) SetMaxCount(max int) {
	t.maxItems = max
	if t.maxItems > 1 {
		t.SetType(TextboxTypeDefault) // Some browsers cannot handle multiple emails in an email type of text input
	}
	t.Refresh()
}

func (t *EmailTextbox) Validate(ctx context.Context) bool {
	ret := t.Textbox.Validate(ctx)

	if ret {
		if t.parseErr != nil {
			t.SetValidationError(t.ΩT("Not a valid email address"))
			return false
		} else if len(t.items) > t.maxItems {
			if t.maxItems == 1 {
				t.SetValidationError(t.ΩT("Enter only one email address"))
			} else {
				t.SetValidationError(fmt.Sprintf(t.ΩT("Enter at most %d email addresses separated by commas"), t.maxItems))
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
func (t *EmailTextbox) Addresses() []string {
	ret := []string{}

	for _, item := range t.items {
		ret = append(ret, item.Address)
	}
	return ret
}

type EmailTextboxCreator struct {
	ID string
	Placeholder string
	Type string
	MinLength int
	MaxLength int
	RowCount int
	ReadOnly bool
	SaveState bool
	// Text is the initial value of the textbox. Often its best to load the value in a separate Load step after creating the control.
	Text string

	page.ControlOptions
}

func (t EmailTextboxCreator) Create(ctx context.Context, parent page.ControlI) page.ControlI {
	ctrl := NewEmailTextbox(parent, t.ID)
	if t.Placeholder != "" {
		ctrl.SetPlaceholder(t.Placeholder)
	}
	if t.Type != "" {
		ctrl.typ = t.Type
	}
	ctrl.minLength = t.MinLength
	ctrl.maxLength = t.MaxLength
	ctrl.rowCount = t.RowCount
	ctrl.readonly = t.ReadOnly
	if t.Text != "" {
		ctrl.SetText(t.Text)
	}

	ctrl.ApplyOptions(t.ControlOptions)
	if t.SaveState {
		ctrl.SaveState(ctx, t.SaveState)
	}
	return ctrl
}

// GetEmailTextbox is a convenience method to return the control with the given id from the page.
func GetEmailTextbox(c page.ControlI, id string) *EmailTextbox {
	return c.Page().GetControl(id).(*EmailTextbox);
}
