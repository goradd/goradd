package control

import (
	"strings"
	"time"

	"github.com/goradd/goradd/pkg/datetime"
	"github.com/goradd/goradd/pkg/page"
)


// DateTextbox is a textbox that only permits dates and/or times to be entered into it.
type DateTextbox struct {
	Textbox
	format string			 // same as what time.format expects
	dt     datetime.DateTime // Converting from text to a datetime is expensive.
							 // We maintain a copy of the conversion to prevent duplication of effort.
}

// NewDateTextbox creates a new DateTextbox.
func NewDateTextbox(parent page.ControlI, id string) *DateTextbox {
	d := &DateTextbox{}
	d.Init(d, parent, id)
	return d
}

func (d *DateTextbox) Init(self TextboxI, parent page.ControlI, id string) {
	d.Textbox.Init(self, parent, id)
	d.ValidateWith(DateValidator{ctrl: d})
	d.format = datetime.UsDateTime
}

// SetFormat sets the format of the text allowed. The format is any allowable format
// that datetime or time can convert.
func (d *DateTextbox) SetFormat(format string) {
	d.format = format
}

// SetValue will set the DateTextbox to the given value if possible.
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
	return d.format
}

// SetText sets the DateTime to the given text. If you attempt set the text to something that is not
// convertible to a date, an empty string will be entered.
func (d *DateTextbox) SetText(s string) page.ControlI {
	l := d.layout()
	v, err := datetime.Parse(l, s)
	if err == nil {
		d.Textbox.SetText(v.Format(l))
		d.dt = v
	} else {
		d.Textbox.SetText("")
		d.dt = datetime.NewZeroDate()
	}
	return d
}

// SetDate will set the textbox to the give datetime value
func (d *DateTextbox) SetDate(dt datetime.DateTime) {
	s := dt.Format(d.layout())
	d.Textbox.SetText(s)
	d.dt = dt
}

// Value returns the value as an interface, but the underlying value will be a datetime.
// If a bad value was entered into the textbox, it will return an empty datetime.
func (d *DateTextbox) Value() interface{} {
	return d.dt
}

// Date returns the value as a DateTime value based on the format.
// If a bad value was entered into the textbox, it will return an empty datetime.
func (d *DateTextbox) Date() datetime.DateTime {
	return d.dt
}

func (d *DateTextbox) ΩUpdateFormValues(ctx *page.Context) {
	d.Textbox.ΩUpdateFormValues(ctx)

	t := d.Text()

	if t == "" {
		d.dt = datetime.NewZeroDate()
		return
	}

	l := d.layout()
	t = strings.Replace(t, "-", "/", -1) // sometimes dashes are used as a date separator. Allow this.

	v, err := datetime.Parse(l, t)
	if err == nil {
		d.dt = v
	}
}


type DateValidator struct {
	ctrl    *DateTextbox
	Message string
}

func (v DateValidator) Validate(c page.ControlI, s string) (msg string) {
	if s == "" {
		return "" // empty string is valid
	}

	// By the time the validator fires, we will have already parsed and validated the value.
	// We just need to check to see if we were successful.
	if v.ctrl.dt.IsZero() {
		if v.Message == "" {
			return c.ΩT("Enter the format ") + v.ctrl.format
		} else {
			return v.Message
		}
	}
	return
}
