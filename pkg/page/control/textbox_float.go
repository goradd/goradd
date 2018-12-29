package control

import (
	"fmt"
	"github.com/goradd/goradd/pkg/page"
	"strconv"
)

// FloatTextbox is a text control that ensures a valid floating point number is entered in the field.
type FloatTextbox struct {
	Textbox
}

func NewFloatTextbox(parent page.ControlI, id string) *FloatTextbox {
	t := &FloatTextbox{}
	t.Init(t, parent, id)
	return t
}

func (i *FloatTextbox) Init(self TextboxI, parent page.ControlI, id string) {
	i.Textbox.Init(self, parent, id)
	i.ValidateWith(FloatValidator{})
}

func (i *FloatTextbox) SetMinValue(minValue float64, invalidMessage string) {
	i.ValidateWith(MinFloatValidator{minValue, invalidMessage})
}

func (i *FloatTextbox) SetMaxValue(maxValue float64, invalidMessage string) {
	i.ValidateWith(MaxFloatValidator{maxValue, invalidMessage})
}

func (i *FloatTextbox) Value() interface{} {
	return i.Float64()
}

func (i *FloatTextbox) Float64() float64 {
	t := i.Textbox.Text()
	v, _ := strconv.ParseFloat(t, 64)
	return v
}

func (i *FloatTextbox) Float32() float32 {
	t := i.Textbox.Text()
	v, _ := strconv.ParseFloat(t, 32)
	return float32(v)
}


func (i *FloatTextbox) SetFloat64(v float64) *FloatTextbox {
	i.Textbox.SetValue(v)
	return i
}

func (i *FloatTextbox) SetFloat32(v float32) *FloatTextbox {
	i.Textbox.SetValue(v)
	return i
}

type FloatValidator struct {
	Message string
}

func (v FloatValidator) Validate(c page.ControlI, s string) (msg string) {
	if s == "" {
		return "" // empty textbox is checked elsewhere
	}
	if _, err := strconv.ParseFloat(s, 64); err != nil {
		if msg == "" {
			return c.ΩT("Please enter a number.")
		} else {
			return v.Message
		}
	}
	return
}

type MinFloatValidator struct {
	MinValue float64
	Message  string
}

func (v MinFloatValidator) Validate(c page.ControlI, s string) (msg string) {
	if s == "" {
		return "" // empty textbox is checked elsewhere
	}
	if val, _ := strconv.ParseFloat(s, 64); val < v.MinValue {
		if msg == "" {
			return fmt.Sprintf(c.ΩT("Enter at least %f"), v.MinValue)
		} else {
			return v.Message
		}
	}
	return
}

type MaxFloatValidator struct {
	MaxValue float64
	Message  string
}

func (v MaxFloatValidator) Validate(c page.ControlI, s string) (msg string) {
	if s == "" {
		return "" // empty textbox is checked elsewhere
	}
	if val, _ := strconv.ParseFloat(s, 64); val < v.MaxValue {
		if msg == "" {
			return fmt.Sprintf(c.ΩT("Enter at most %f"), v.MaxValue)
		} else {
			return v.Message
		}
	}
	return
}
