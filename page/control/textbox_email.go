package control

import (
	"fmt"
	"github.com/spekary/goradd/page"
	"github.com/spekary/goradd/page/control/control_base"
	"net/mail"
	"context"
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
func NewEmailTextbox(parent page.ControlI) *EmailTextbox {
	t := &EmailTextbox{}
	t.Init(t, parent)
	t.maxItems = 1
	t.SetType(TEXTBOX_TYPE_EMAIL)
	return t
}

func (t *EmailTextbox) Init(self control_base.TextboxI, parent page.ControlI) {
	t.Textbox.Init(self, parent)
}

func (t *EmailTextbox) SetMaxCount(max int) {
	t.maxItems = max
	if t.maxItems > 1 {
		t.SetType(TEXTBOX_TYPE_DEFAULT) // Some browsers cannot handle multiple emails in an email type of text input
	}
	t.Refresh()
}

func (t *EmailTextbox) Validate(ctx context.Context) bool {
	ret := t.Textbox.Validate(ctx)

	if ret {
		if t.parseErr != nil {
			t.SetValidationError(t.Translate("Not a valid email address"))
			return false
		} else if len(t.items) > t.maxItems {
			if t.maxItems == 1 {
				t.SetValidationError(t.Translate("Enter only one email address"))
			} else {
				t.SetValidationError(fmt.Sprintf(t.Translate("Enter at most %d email addresses separated by commas"), t.maxItems))
			}

			return false
		}
	}
	return true
}

func (t *EmailTextbox) UpdateFormValues(ctx *page.Context) {
	t.Textbox.UpdateFormValues(ctx)
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
