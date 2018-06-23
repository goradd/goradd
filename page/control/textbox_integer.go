package control

import (
	"fmt"
	"github.com/spekary/goradd/page"
	"github.com/spekary/goradd/page/control/control_base"
	"strconv"
)

type IntegerTextbox struct {
	Textbox
}

func NewIntegerTextbox(parent page.ControlI, id string) *IntegerTextbox {
	t := &IntegerTextbox{}
	t.Init(t, parent, id)
	return t
}

func (i *IntegerTextbox) Init(self control_base.TextboxI, parent page.ControlI, id string) {
	i.Textbox.Init(self, parent, id)
	i.ValidateWith(IntValidator{})
}

// SetMinValue creates a validator that makes sure the value of the text box is at least the
// given value. Specify your own error message, or leave the error message blank and a standard error message will
// be presented if the value is not valid.
func (i *IntegerTextbox) SetMinValue(minValue int, invalidMessage string) {
	i.ValidateWith(MinIntValidator{minValue, invalidMessage})
}

// SetMaxValue creates a validator that makes sure the value of the text box is at most the
// given value. Specify your own error message, or leave the error message blank and a standard error message will
// be presented if the value is not valid.
func (i *IntegerTextbox) SetMaxValue(maxValue int, invalidMessage string) {
	i.ValidateWith(MaxIntValidator{maxValue, invalidMessage})
}

func (i *IntegerTextbox) Value() interface{} {
	t := i.Textbox.Text()
	v, _ := strconv.Atoi(t)
	return v
}

type IntValidator struct {
	Message string
}

func (v IntValidator) Validate(t page.Translater, s string) (msg string) {
	if s == "" {
		return "" // empty textbox is checked elsewhere
	}
	if _, err := strconv.Atoi(s); err != nil {
		if v.Message == "" {
			return t.Translate("Please enter an integer.")
		} else {
			return v.Message
		}
	}
	return
}

type MinIntValidator struct {
	MinValue int
	Message  string
}

func (v MinIntValidator) Validate(t page.Translater, s string) (msg string) {
	if s == "" {
		return "" // empty textbox is checked elsewhere
	}
	if val, _ := strconv.Atoi(s); val < v.MinValue {
		if v.Message == "" {
			return fmt.Sprintf(t.Translate("Enter at least %d"), v.MinValue)
		} else {
			return v.Message
		}
	}
	return
}

type MaxIntValidator struct {
	MaxValue int
	Message  string
}

func (v MaxIntValidator) Validate(t page.Translater, s string) (msg string) {
	if s == "" {
		return "" // empty textbox is checked elsewhere
	}
	if val, _ := strconv.Atoi(s); val < v.MaxValue {
		if v.Message == "" {
			return fmt.Sprintf(t.Translate("Enter at most %d"), v.MaxValue)
		} else {
			return v.Message
		}
	}
	return
}
