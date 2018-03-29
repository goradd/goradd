package textbox

import (
	"goradd/page/control"
	"strconv"
	"fmt"
	"github.com/spekary/goradd/page"
	grcontrol "github.com/spekary/goradd/page/control"
)

type IntegerTextbox struct {
	control.TextBox
}

func NewIntegerTextBox(parent page.ControlI, id string) *IntegerTextbox {
	t := &IntegerTextbox{}
	t.Init(t, parent, id)
	return t
}

func (i *IntegerTextbox) Init(self grcontrol.TextBoxI, parent page.ControlI, id string) {
	i.TextBox.Init(self, parent, id)
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
	t := i.TextBox.Text()
	v,_ := strconv.Atoi(t)
	return v
}

type IntValidator struct {
	Message string
}

func (v IntValidator) Validate(t page.Translater, s string) (msg string) {
	if _,err := strconv.Atoi(s); err != nil {
		if msg == "" {
			return t.Translate("Please enter an integer.")
		} else {
			return v.Message
		}
	}
	return
}


type MinIntValidator struct {
	MinValue int
	Message string
}

func (v MinIntValidator) Validate(t page.Translater, s string) (msg string) {
	if val,_ := strconv.Atoi(s); val < v.MinValue {
		if msg == "" {
			return fmt.Sprintf (t.Translate("Enter at least %d"), v.MinValue)
		} else {
			return v.Message
		}
	}
	return
}

type MaxIntValidator struct {
	MaxValue int
	Message string
}

func (v MaxIntValidator) Validate(t page.Translater, s string) (msg string) {
	if val,_ := strconv.Atoi(s); val < v.MaxValue {
		if msg == "" {
			return fmt.Sprintf (t.Translate("Enter at most %d"), v.MaxValue)
		} else {
			return v.Message
		}
	}
	return
}


