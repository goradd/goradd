package datetime

import (
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
)

func Parse(layout, value string) (DateTime, error) {
	var ts bool
	if layout == Stamp || layout == StampMilli || layout == StampMicro || layout == StampNano {
		ts = true
	}

	layout, value = strings.ToUpper(layout), strings.ToUpper(value) // standardize Pm, PM or pm

	t, err := time.Parse(layout, value)
	return DateTime{t, ts}, err
}

// FromSqlDateTime will receive a Date, Time, DateTime or Timestamp type of string that is typically output by SQL and
// convert it to our own DateTime object. If the SQL date time string does not have timezone information,
// the resulting value will be in UTC time, will not be a timestamp, and we are assuming that the SQL itself is in UTC.
// If it DOES have timezone information, we will assume its a timestamp.
// If an error occurs, the returned value will be the zero date.
func FromSqlDateTime(s string) (t DateTime, err error) {
	var hasDate, hasTime, hasNano, hasTZ bool
	var form string

	if strings.Contains(s, ".") {
		hasNano = true
	}
	if strings.Contains(s, "-") {
		hasDate = true
	}
	if strings.Contains(s, ":") {
		hasTime = true

		if strings.LastIndexAny(s, "+-") > strings.LastIndex(s, ":") {
			hasTZ = true
		}
	}

	if hasNano {
		form = "2006-01-02 15:04:05.999999"
		if hasTZ {
			form += "-07"
		}
	} else if hasDate && hasTime {
		form = "2006-01-02 15:04:05"
		if hasTZ {
			form += "-07"
		}
	} else if hasDate {
		form = "2006-01-02"
	} else {
		form = "15:04:05"
	}

	t, err = Parse(form, s)

	if hasTZ {
		t.isTimestamp = true
	}

	return
}

func (d DateTime) Format(layout string) string {
	if d.IsZero() {
		return ""
	}
	return d.Time.Format(layout)
}