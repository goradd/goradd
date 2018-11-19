package control

import (
	"github.com/spekary/goradd/pkg/page"
	"github.com/spekary/goradd/pkg/page/control/control_base"
	"github.com/spekary/goradd/pkg/datetime"
	"time"
	"strings"
)

type DateTextboxFormat int


// Go does not yet have internal support for international dates so we will need to do this ourselves
const (
	DateTextboxUS DateTextboxFormat = iota // default to US format M/D/Y
	DateTextboxEuro // D/M/Y
)

type DateTextbox struct {
	Textbox
	format DateTextboxFormat
	value datetime.DateTime
}

func NewDateTextbox(parent page.ControlI, id string) *DateTextbox {
	t := &DateTextbox{}
	t.Init(t, parent, id)
	return t
}

func (i *DateTextbox) Init(self control_base.TextboxI, parent page.ControlI, id string) {
	i.Textbox.Init(self, parent, id)
	i.ValidateWith(DateValidator{ctrl:i})
}

// SetMinValue creates a validator that makes sure the value of the text box is at least the
// given value. Specify your own error message, or leave the error message blank and a standard error message will
// be presented if the value is not valid.
func (i *DateTextbox) SetFormat(format DateTextboxFormat) {
	i.format = format
}

func (i *DateTextbox) SetValue(val interface{}) *DateTextbox {
	switch v := val.(type) {
	case string:
		i.SetText(v)
	case datetime.DateTime:
		i.SetDate(v)
	case time.Time:
		d := datetime.NewDateTime(v)
		i.SetDate(d)
	}
	return i
}

func (i *DateTextbox) layout() string {
	switch i.format{
	case DateTextboxEuro:
		return datetime.EuroDate
	default:
		return datetime.UsDate
	}
}

func (i *DateTextbox) SetText(s string) page.ControlI {
	l := i.layout()
	s = strings.Replace(s, "-", "/", -1)
	v, err := datetime.Parse(l, s)
	if err != nil {
		i.Textbox.SetText(v.Format(l))
		i.value = v
	}
	return i
}

func (i *DateTextbox) SetDate(d datetime.DateTime) {
	s := d.Format(i.layout())
	i.Textbox.SetText(s)
	i.value = d
}

func (i *DateTextbox) Value() interface{} {
	return i.value
}

func (i *DateTextbox) Date() datetime.DateTime {
	return i.value
}


type DateValidator struct {
	ctrl *DateTextbox
	Message string
}

func (v DateValidator) Validate(t page.Translater, s string) (msg string) {
	if s == "" {
		return "" // empty textbox is checked elsewhere
	}
	if _, err := datetime.Parse(v.ctrl.layout(), s); err != nil {
		if v.Message == "" {
			return t.Translate("Enter a date.")
		} else {
			return v.Message
		}
	}
	return
}

