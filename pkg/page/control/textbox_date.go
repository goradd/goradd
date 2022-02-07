package control

import (
	"context"
	"encoding/gob"
	"github.com/goradd/goradd/pkg/session"
	"strings"
	"time"

	"github.com/goradd/goradd/pkg/datetime"
	"github.com/goradd/goradd/pkg/page"
)

type DateTextboxI interface {
	TextboxI
	SetFormats(formats []string) DateTextboxI
	Date() datetime.DateTime
	Formats() []string
}

// DateTextbox is a textbox that only permits dates and/or times to be entered into it.
//
// Dates and times will be converted to Browser local time.
type DateTextbox struct {
	Textbox
	formats []string         // Variety of formats it will accept. Same as what time.format expects.
	dt     datetime.DateTime // Converting from text to a datetime is expensive.
							 // We maintain a copy of the conversion to prevent duplication of effort.
}

// NewDateTextbox creates a new DateTextbox.
func NewDateTextbox(parent page.ControlI, id string) *DateTextbox {
	d := &DateTextbox{}
	d.Self = d
	d.Init(parent, id)
	return d
}

func (d *DateTextbox) Init(parent page.ControlI, id string) {
	d.Textbox.Init(parent, id)
	d.ValidateWith(DateValidator{})
	d.formats = []string{datetime.UsDateTime}
}

// SetFormats sets the format of the text allowed. The format is any allowable format
// that datetime or time can convert.
func (d *DateTextbox) SetFormats(formats []string) DateTextboxI {
	d.formats = formats
	return d
}

// Formats returns the format string specified previously
func (d *DateTextbox) Formats() []string {
	return d.formats
}


// SetValue will set the DateTextbox to the given value if possible.
func (d *DateTextbox) SetValue(val interface{}) page.ControlI {
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

func (d *DateTextbox) layouts() []string {
	return d.formats
}

func (d *DateTextbox) parseDate(ctx context.Context, s string) (result datetime.DateTime, layoutUsed string, err error) {
	var grctx *page.Context

	if ctx != nil {
		grctx = page.GetContext(ctx)
	}
	for _,layoutUsed = range d.layouts() {
		if grctx != nil && datetime.LayoutHasDate(layoutUsed) && datetime.LayoutHasTime(layoutUsed){
			result, err = datetime.ParseInOffset(layoutUsed, s, session.ClientTimezoneOffset(ctx))
		} else {
			result, err = datetime.Parse(layoutUsed, s)
		}
		if err == nil {
			break
		}
	}
	return
}

// SetText sets the DateTime to the given text. If you attempt set the text to something that is not
// convertible to a date, an empty string will be entered. The resulting datetime will be in UTC time.
// Use SetDate if you want to make sure the date is in a certain timezone.
func (d *DateTextbox) SetText(s string) page.ControlI {
	v, layout, err := d.parseDate(nil, s)

	if err == nil {
		d.Textbox.SetText(v.Format(layout))
		d.dt = v
	} else {
		d.Textbox.SetText("")
		d.dt = datetime.NewZeroDate()
	}
	return d
}

// SetDate will set the textbox to the give datetime value
func (d *DateTextbox) SetDate(dt datetime.DateTime) {
	s := dt.Format(d.layouts()[0])
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

func (d *DateTextbox) UpdateFormValues(ctx context.Context) {
	d.Textbox.UpdateFormValues(ctx)

	if d.readonly {
		// This would happen if someone was attempting to hack the browser.
		return
	}
	if _, ok := page.GetContext(ctx).FormValue(d.ID()); !ok {
		return
	}
	t := d.Text()
	if t == "" {
		d.dt = datetime.NewZeroDate()
		return
	}

	v, _, err := d.parseDate(ctx, t)

	if err == nil {
		d.dt = v
	} else {
		d.dt = datetime.DateTime{} // set to zero value to indicate an error
	}
}

func (d *DateTextbox) Serialize(e page.Encoder) (err error) {
	if err = d.Textbox.Serialize(e); err != nil {
		return
	}
	if err = e.Encode(d.formats); err != nil {
		return
	}
	if err = e.Encode(d.dt); err != nil {
		return
	}

	return
}

func (d *DateTextbox) Deserialize(dec page.Decoder) (err error) {
	if err = d.Textbox.Deserialize(dec); err != nil {
		return
	}
	if err = dec.Decode(&d.formats); err != nil {
		return
	}
	if err = dec.Decode(&d.dt); err != nil {
		return
	}

	return
}

type DateValidator struct {
	Message string
}

func (v DateValidator) Validate(c page.ControlI, s string) (msg string) {
	if s == "" {
		return "" // empty string is valid
	}

	// By the time the validator fires, we will have already parsed and validated the value.
	// We just need to check to see if we were successful.
	ctrl := c.(DateTextboxI)
	if ctrl.Date().IsZero() {
		if v.Message == "" {
			return c.GT("Enter one of these formats: ") + strings.Join(ctrl.Formats(), ", ")
		} else {
			return v.Message
		}
	}
	return
}

// DateTextboxCreator creates an date textbox.
// Pass it to AddControls of a control, or as a Child of
// a FormFieldWrapper.
type DateTextboxCreator struct {
	// ID is the control id of the html widget and must be unique to the page
	ID string
	// Placeholder is the placeholder attribute of the textbox and shows as help text inside the field
	Placeholder string
	// Type is the type attribute of the textbox
	Type string
	// MinLength is the minimum number of characters that the user is required to enter. If the
	// length is less than this number, a validation error will be shown.
	MinLength int
	// MaxLength is the maximum number of characters that the user is required to enter. If the
	// length is more than this number, a validation error will be shown.
	MaxLength int
	// ColumnCount is the number of characters wide the textbox will be, and becomes the width attribute in the tag.
	// The actual width is browser dependent. For better control, use a width style property.
	ColumnCount int
	// ReadOnly sets the readonly attribute of the textbox, which prevents it from being changed by the user.
	ReadOnly bool
	// SaveState will save the text in the textbox, to be restored if the user comes back to the page.
	// It is particularly helpful when the textbox is being used to filter the results of a query, so that
	// when the user comes back to the page, he does not have to type the filter text again.
	SaveState bool
	// Text is the initial value of the textbox. Often its best to load the value in a separate Load step after creating the control.
	Text string
	// Formats is the time.format strings to use to decode the text into a date or to display the date. By default it is datetime.UsDateTime.
	Formats []string

	page.ControlOptions
}

func (c DateTextboxCreator) Create(ctx context.Context, parent page.ControlI) page.ControlI {
	ctrl := NewDateTextbox(parent, c.ID)
	c.Init(ctx, ctrl)
	return ctrl
}

func (c DateTextboxCreator) Init(ctx context.Context, ctrl DateTextboxI) {
	if c.Formats != nil {
		ctrl.SetFormats(c.Formats)
	}
	// Reuse subclass
	sub := TextboxCreator{
		Placeholder:    c.Placeholder,
		MinLength:      c.MinLength,
		MaxLength:      c.MaxLength,
		Type:           c.Type,
		ColumnCount:    c.ColumnCount,
		ReadOnly:       c.ReadOnly,
		ControlOptions: c.ControlOptions,
		SaveState:      c.SaveState,
		Text: c.Text,
	}
	sub.Init(ctx, ctrl)
}

// GetDateTextbox is a convenience method to return the control with the given id from the page.
func GetDateTextbox(c page.ControlI, id string) *DateTextbox {
	return c.Page().GetControl(id).(*DateTextbox)
}

func init() {
	gob.Register(DateValidator{})
	page.RegisterControl(&DateTextbox{})
}
