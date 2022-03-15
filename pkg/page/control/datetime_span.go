package control

import (
	"context"
	"github.com/goradd/goradd/pkg/config"
	"github.com/goradd/html5tag"
	"github.com/goradd/goradd/pkg/page"
	"io"
	"time"
)

// DateTimeSpan is a span that displays a datetime value as static text. This is a typical default control to use
// for a timestamp in the database.
type DateTimeSpan struct {
	Span
	format string
	value  time.Time
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
	case time.Time:
		s.SetDateTime(v2)
	case string:
		d, err := time.Parse(s.format, v2)
		if err != nil {
			panic(err)
		}
		s.SetDateTime(d)
	}
}

// SetDateTime sets the value to a datetime.DateTime.
func (s *DateTimeSpan) SetDateTime(d time.Time) {
	s.value = d
	s.Refresh()
}

// Value returns the value as a datetime.DateTime object.
// Also satisfies the Valuer interface
func (s *DateTimeSpan) Value() interface{} {
	return s.value
}

func (s *DateTimeSpan) DateTime() time.Time {
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
func (s *DateTimeSpan) DrawInnerHtml(_ context.Context, w io.Writer) {
	page.WriteString(w, s.value.Format(s.format))
	// TODO: Internationalize this. Golang does not currently have date/time localization support,
	// as in translation of month and weekday strings, but it does support arbitrary timezones.
	// The other option is to let JavaScript format it, though that limits you to formatting in
	// local time or UTC. JavaScript does not have a means to specify the timezone that is well supported.
	// However, JavaScript will translate month and weekday names to the local language.
	return
}
func (s *DateTimeSpan) DrawingAttributes(ctx context.Context) html5tag.Attributes {
	return s.ControlBase.DrawingAttributes(ctx)
}

func (s *DateTimeSpan) Serialize(e page.Encoder) {
	s.ControlBase.Serialize(e)

	if err := e.Encode(s.format); err != nil {
		panic(err)
	}

	if err := e.Encode(s.value); err != nil {
		panic(err)
	}
}

func (s *DateTimeSpan) Deserialize(dec page.Decoder) {
	s.ControlBase.Deserialize(dec)

	if err := dec.Decode(&s.format); err != nil {
		panic(err)
	}

	if err := dec.Decode(&s.value); err != nil {
		panic(err)
	}
}

type DateTimeSpanCreator struct {
	ID string
	Format string
	Value time.Time
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
