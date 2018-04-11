package control

import (
	"strconv"
	"fmt"
	"github.com/spekary/goradd/page"
)

// FloatTextbox is a text control that ensures a valid floating point number is entered in the field.
type FloatTextbox struct {
	Textbox
}

func NewFloatTextbox(parent page.ControlI) *FloatTextbox {
	t := &FloatTextbox{}
	t.Init(t, parent)
	return t
}

func (i *FloatTextbox) Init(self page.ControlI, parent page.ControlI) {
	i.Textbox.Init(self, parent)
	i.ValidateWith(FloatValidator{})
}

func (i *FloatTextbox) SetMinValue(minValue float64, invalidMessage string) {
	i.ValidateWith(MinFloatValidator{minValue, invalidMessage})
}

func (i *FloatTextbox) SetMaxValue(maxValue float64, invalidMessage string) {
	i.ValidateWith(MaxFloatValidator{maxValue, invalidMessage})
}

func (i *FloatTextbox) Value() interface{} {
	t := i.Textbox.Text()
	v,_ := strconv.ParseFloat(t,64)
	return v
}

type FloatValidator struct {
	Message string
}

func (v FloatValidator) Validate(t page.Translater, s string) (msg string) {
	if s == "" {
		return "" // empty textbox is checked elsewhere
	}
	if _,err := strconv.ParseFloat(s, 64); err != nil {
		if msg == "" {
			return t.Translate("Please enter a number.")
		} else {
			return v.Message
		}
	}
	return
}


type MinFloatValidator struct {
	MinValue float64
	Message string
}

func (v MinFloatValidator) Validate(t page.Translater, s string) (msg string) {
	if s == "" {
		return "" // empty textbox is checked elsewhere
	}
	if val,_ := strconv.ParseFloat(s, 64); val < v.MinValue {
		if msg == "" {
			return fmt.Sprintf (t.Translate("Enter at least %f"), v.MinValue)
		} else {
			return v.Message
		}
	}
	return
}

type MaxFloatValidator struct {
	MaxValue float64
	Message string
}

func (v MaxFloatValidator) Validate(t page.Translater, s string) (msg string) {
	if s == "" {
		return "" // empty textbox is checked elsewhere
	}
	if val,_ := strconv.ParseFloat(s, 64); val < v.MaxValue {
		if msg == "" {
			return fmt.Sprintf (t.Translate("Enter at most %f"), v.MaxValue)
		} else {
			return v.Message
		}
	}
	return
}


