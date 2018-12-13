package control

import (
	"strings"
	"time"

	"github.com/spekary/goradd/pkg/datetime"
	"github.com/spekary/goradd/pkg/page"
)

type DateTextboxFormat int

// Go does not yet have internal support for international dates so we will need to do this ourselves
const (
	DateTextboxUS   DateTextboxFormat = iota // default to US format M/D/Y
	DateTextboxEuro                          // D/M/Y
)

type DateTextbox struct {
	Textbox
	format DateTextboxFormat
	value  datetime.DateTime
}

func NewDateTextbox(parent page.ControlI, id string) *DateTextbox {
	t := &DateTextbox{}
	t.Init(t, parent, id)
	return t
}

func (d *DateTextbox) Init(self TextboxI, parent page.ControlI, id string) {
	d.Textbox.Init(self, parent, id)
	d.ValidateWith(DateValidator{ctrl: d})
}

// SetMinValue creates a validator that makes sure the value of the text box is at least the
// given value. Specify your own error message, or leave the error message blank and a standard error message will
// be presented if the value is not valid.
func (d *DateTextbox) SetFormat(format DateTextboxFormat) {
	d.format = format
}

func (d *DateTextbox) SetValue(val interface{}) *DateTextbox {
	switch v := val.(type) {
	case string:
		d.SetText(v)
	case datetime.DateTime:
		d.SetDate(v)
	case time.Time:
		dt := datetime.NewDateTime(v)
		d.SetDate(dt)
	}
	return d
}

func (d *DateTextbox) layout() string {
	switch d.format {
	case DateTextboxEuro:
		return datetime.EuroDate
	default:
		return datetime.UsDate
	}
}

func (d *DateTextbox) SetText(s string) page.ControlI {
	l := d.layout()
	s = strings.Replace(s, "-", "/", -1)
	v, err := datetime.Parse(l, s)
	if err != nil {
		d.Textbox.SetText(v.Format(l))
		d.value = v
	}
	return d
}

func (d *DateTextbox) SetDate(dt datetime.DateTime) {
	s := dt.Format(d.layout())
	d.Textbox.SetText(s)
	d.value = dt
}

func (d *DateTextbox) Value() interface{} {
	return d.value
}

func (d *DateTextbox) Date() datetime.DateTime {
	return d.value
}

type DateValidator struct {
	ctrl    *DateTextbox
	Message string
}

func (v DateValidator) Validate(c page.ControlI, s string) (msg string) {
	if s == "" {
		return "" // empty textbox is checked elsewhere
	}
	if _, err := datetime.Parse(v.ctrl.layout(), s); err != nil {
		if v.Message == "" {
			return c.ΩT("Enter a date.")
		} else {
			return v.Message
		}
	}
	return
}
