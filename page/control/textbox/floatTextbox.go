package textbox

import (
	"goradd/page/control"
	"strconv"
	"fmt"
	"github.com/spekary/goradd/page"
	grcontrol "github.com/spekary/goradd/page/control"
)

type FloatTextbox struct {
	control.TextBox
}

func (i *FloatTextbox) Init(self grcontrol.TextBoxI, parent page.ControlI, id string) {
	i.TextBox.Init(self, parent, id)
	i.ValidateWith(FloatValidator{})
}

func (i *FloatTextbox) SetMinValue(minValue float64, invalidMessage string) {
	i.ValidateWith(MinFloatValidator{minValue, invalidMessage})
}

func (i *FloatTextbox) SetMaxValue(maxValue float64, invalidMessage string) {
	i.ValidateWith(MaxFloatValidator{maxValue, invalidMessage})
}

func (i *FloatTextbox) Value() interface{} {
	t := i.TextBox.Text()
	v,_ := strconv.ParseFloat(t,64)
	return v
}

type FloatValidator struct {
	Message string
}

func (v FloatValidator) Validate(t page.Translater, s string) (msg string) {
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
	if val,_ := strconv.ParseFloat(s, 64); val < v.MaxValue {
		if msg == "" {
			return fmt.Sprintf (t.Translate("Enter at most %f"), v.MaxValue)
		} else {
			return v.Message
		}
	}
	return
}


