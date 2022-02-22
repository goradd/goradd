package time

import (
	"github.com/goradd/goradd/pkg/log"
	strings2 "github.com/goradd/goradd/pkg/strings"
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

	DateOnlyFormat = "2006-01-02"
	TimeOnlyFormat = "15:04:05"

	UsDate            = "1/2/2006"
	ShortUsDate       = "1/2/06"
	EuroDate          = "2/1/2006"
	UsDateTime        = "1/2/2006 3:04 PM"
	EuroDateTime      = "2/1/2006 15:04"
	UsTime            = "3:04 PM"
	EuroTime          = "15:04"
	UsDateTimeSeconds = "1/2/2006 3:04:05 PM"
	LongDateDOW       = "Monday, January 2, 2006"
	LongDate          = "January 2, 2006"
)


// Parse parses the given layout string to turn a string int a DateTime
// Since the time.Time doc is not that great about really describing the format, here it is:
//
// Day of the week		Mon   Monday
// Day					02   2   _2   (width two, right justified)
// Month				01   1   Jan   January
// Year					06   2006
// Hour					03   3   15
// Minute				04   4
// Second				05   5
// ms μs ns				.000   .000000   .000000000
// ms μs ns				.999   .999999   .999999999   (trailing zeros removed)
// am/pm				PM   pm
// Timezone				MST
// Offset				-0700   -07   -07:00   Z0700   Z07:00


// ParseInOffset is like time.ParseInLocation but uses the given timezone offset in minutes from UTC to
// be the location of the parsed time.
func ParseInOffset(layout, value string, tzOffset int) (time.Time, error) {
	t, err := time.Parse(layout, value)
	if err == nil {
		t = As(t, time.FixedZone("", tzOffset * 60))
	}
	return t, err
}

// FromSqlDateTime will receive a Date, Time, DateTime or Timestamp type of string that is typically output by SQL and
// convert it to a time.Time. If the SQL date time string does not have timezone information,
// the resulting value will be in UTC time.
// If an error occurs, the returned value will be the zero date.
func FromSqlDateTime(s string) (t time.Time) {
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

	t, err := time.Parse(form, s)
	if err == nil {
		t = t.UTC()
	} else {
		// We can't return it, but we can log it
		log.Warning(err)
	}
	return
}

// LayoutHasDate returns true if the given parse layout indicates a date.
func LayoutHasDate(layout string) bool {
	return strings2.ContainsAnyStrings(layout, "6", "2", "Jan")
}

// LayoutHasTime returns true if the give parse layout indicates a time.
func LayoutHasTime(layout string) bool {
	return strings2.ContainsAnyStrings(layout, "15", "4", "5", "3")
}
