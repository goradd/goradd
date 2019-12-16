package control

import (
	"bytes"
	"context"
	"github.com/goradd/goradd/pkg/config"
	"github.com/goradd/goradd/pkg/datetime"
	"github.com/goradd/goradd/pkg/html"
	"github.com/goradd/goradd/pkg/page"
	"time"
)

// DateTimeSpan is a span that displays a datetime value as static text. This is a typical default control to use
// for a timestamp in the database.
type DateTimeSpan struct {
	Span
	format string
	value  datetime.DateTime
}

// NewDateTimeSpan create a new DateTimeSpan.
func NewDateTimeSpan(parent page.ControlI, id string) *DateTimeSpan {
	s := new(DateTimeSpan)
	s.Self = s
	s.Init(parent, id)
	return s
}

// Init is called by subclasses to initialize the parent. You do not normally need to call this.
func (s *DateTimeSpan) Init(parent page.ControlI, id string) {
	s.Span.Init(parent, id)
	s.format = config.DefaultDateTimeFormat
}

// SetValue sets the value display. You can set the value to a datetime.DateTime, a time.Time,
// or a string that can be parsed by the format string.
func (s *DateTimeSpan) SetValue(v interface{}) {
	switch v2 := v.(type) {
	case datetime.DateTime:
		s.SetDateTime(v2)
	case time.Time:
		s.SetDateTime(datetime.NewDateTime(v2))
	case string:
		d, err := datetime.Parse(s.format, v2)
		if err != nil {
			panic(err)
		}
		s.SetDateTime(d)
	}
}

// SetDateTime sets the value to a datetime.DateTime.
func (s *DateTimeSpan) SetDateTime(d datetime.DateTime) {
	s.value = d
	s.Refresh()
}

// Value returns the value as a datetime.DateTime object.
func (s *DateTimeSpan) Value() datetime.DateTime {
	return s.value
}

// SetFormat sets the format string. This should be a time.TimeFormat string described at
// https://golang.org/pkg/time/#Time.Format
func (s *DateTimeSpan) SetFormat(format string) *DateTimeSpan {
	s.format = format
	s.Refresh()
	return s
}

// DrawInnerHtml is called by the framework to draw the inner html of the span.
func (s *DateTimeSpan) DrawInnerHtml(ctx context.Context, buf *bytes.Buffer) error {
	buf.WriteString(s.value.Format(s.format))
	// TODO: Internationalize this. Golang does not currently have date/time localization support,
	// as in translation of month and weekday strings, but it does support arbitrary timezones.
	// The other option is to let JavaScript format it, though that limits you to formatting in
	// local time or UTC. JavaScript does not have a means to specify the timezone that is well supported.
	// However, JavaScript will translate month and weekday names to the local language.
	return nil
}
func (s *DateTimeSpan) DrawingAttributes(ctx context.Context) html.Attributes {
	return s.ControlBase.DrawingAttributes(ctx)
}

func (s *DateTimeSpan) Serialize(e page.Encoder) (err error) {
	if err = s.ControlBase.Serialize(e); err != nil {
		return
	}

	if err = e.Encode(s.format); err != nil {
		return
	}

	if err = e.Encode(s.value); err != nil {
		return
	}

	return
}

func (s *DateTimeSpan) Deserialize(dec page.Decoder) (err error) {
	if err = s.ControlBase.Deserialize(dec); err != nil {
		return
	}

	if err = dec.Decode(&s.format); err != nil {
		return
	}

	if err = dec.Decode(&s.value); err != nil {
		return
	}

	return
}



type DateTimeSpanCreator struct {
	ID string
	Format string
	Value datetime.DateTime
	page.ControlOptions
}

func (c DateTimeSpanCreator) Create(ctx context.Context, parent page.ControlI) page.ControlI {
	ctrl := NewDateTimeSpan(parent, c.ID)
	ctrl.value = c.Value
	if c.Format != "" {
		ctrl.format = c.Format
	}
	ctrl.ApplyOptions(ctx, c.ControlOptions)
	return ctrl
}

func init() {
	page.RegisterControl(&DateTimeSpan{})
}
