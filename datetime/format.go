package datetime

import (
	"log"
	"strings"
	"time"
)

const (
	// taken from time.Format

	ANSIC       = "Mon Jan _2 15:04:05 2006"
	UnixDate    = "Mon Jan _2 15:04:05 MST 2006"
	RubyDate    = "Mon Jan 02 15:04:05 -0700 2006"
	RFC822      = "02 Jan 06 15:04 MST"
	RFC822Z     = "02 Jan 06 15:04 -0700" // RFC822 with numeric zone
	RFC850      = "Monday, 02-Jan-06 15:04:05 MST"
	RFC1123     = "Mon, 02 Jan 2006 15:04:05 MST"
	RFC1123Z    = "Mon, 02 Jan 2006 15:04:05 -0700" // RFC1123 with numeric zone
	RFC3339     = "2006-01-02T15:04:05Z07:00"
	RFC3339Nano = "2006-01-02T15:04:05.999999999Z07:00"
	Kitchen     = "3:04PM"
	// Handy time stamps.
	Stamp      = "Jan _2 15:04:05"
	StampMilli = "Jan _2 15:04:05.000"
	StampMicro = "Jan _2 15:04:05.000000"
	StampNano  = "Jan _2 15:04:05.000000000"

	// Our additional formatters

	SqlDateTime = "2006-01-02 15:04:05.000000"
)

func Parse(layout, value string) (DateTime, error) {
	t, err := time.Parse(layout, value)

	return DateTime{t}, err
}

// FromSqlDateTime will receive a Date, Time, DateTime or Timestamp type of string that is typically output by SQL and
// convert it to our own DateTime object.
func FromSqlDateTime(s string) DateTime {
	var t DateTime
	var err error
	var hasDate, hasTime, hasNano bool

	if strings.ToUpper(s) == "CURRENT_TIMESTAMP" {
		return NewDateTime(Current)
	}

	if strings.Contains(s, ".") {
		hasNano = true
	}
	if strings.Contains(s, "-") {
		hasDate = true
	}
	if strings.Contains(s, ":") {
		hasTime = true
	}

	if hasNano {
		t, err = Parse("2006-01-02 15:04:05.999999", s)
	} else if hasDate && hasTime {
		t, err = Parse("2006-01-02 15:04:05", s)
	} else if hasDate {
		t, err = Parse("2006-01-02", s)
	} else {
		t, err = Parse("15:04:05", s)
	}
	if err != nil {
		log.Panic(err)
	}
	return t

}
