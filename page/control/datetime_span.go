package control

import (
	"github.com/spekary/goradd/page"
	"goradd/config"
	"bytes"
	"context"
	"github.com/spekary/goradd/datetime"
	"time"
)

// DateTimeSpan is a span that displays a datetime value as static text. This is a typical default control to use
// for a timestamp in the database.
type DateTimeSpan struct {
	Span
	format string
	value datetime.DateTime
}

func NewDateTimeSpan(parent page.ControlI, id string) *DateTimeSpan {
	s := new(DateTimeSpan)
	s.Init(s, parent, id)
	return s
}

func (s *DateTimeSpan) Init(self page.ControlI, parent page.ControlI, id string) {
	s.Span.Init(self, parent, id)
	s.format = config.DefaultDateTimeFormat
}

func (s *DateTimeSpan) SetValue(v interface{}) {
	switch v2 := v.(type) {
	case datetime.DateTime:
		s.SetDateTime(v2)
	case time.Time:
		s.SetDateTime(datetime.NewDateTime(v2))
	case string:
		d,err := datetime.Parse(s.format, v2)
		if err != nil {
			panic(err)
		}
		s.SetDateTime(d)
	}
}

func (s *DateTimeSpan) SetDateTime(d datetime.DateTime){
	s.value = d
	s.Refresh()
}

func (s *DateTimeSpan) Value() datetime.DateTime {
	return s.value
}

func (s *DateTimeSpan) SetFormat(format string)  *DateTimeSpan {
	s.format = format
	s.Refresh()
	return s
}


func (s *DateTimeSpan) DrawInnerHtml(ctx context.Context, buf *bytes.Buffer) error {
	buf.WriteString(s.value.Format(s.format))
	// TODO: Internationalize this. Golang does not currently have date/time localization support
	return nil
}
