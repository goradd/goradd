package time

import (
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/goradd/goradd/pkg/log"
	strings2 "github.com/goradd/goradd/pkg/strings"
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
	KitchenLc   = "3:04pm"

	Stamp      = "Jan _2 15:04:05"
	StampMilli = "Jan _2 15:04:05.000"
	StampMicro = "Jan _2 15:04:05.000000"
	StampNano  = "Jan _2 15:04:05.000000000"

	DateOnlyFormat = "2006-01-02"
	TimeOnlyFormat = "15:04:05"

	UsDate              = "1/2/2006"
	ShortUsDate         = "1/2/06"
	EuroDate            = "2/1/2006"
	UsDateTime          = "1/2/2006 3:04 PM"
	UsDateTimeLc        = "1/2/2006 3:04 pm"
	EuroDateTime        = "2/1/2006 15:04"
	UsTime              = "3:04 PM"
	UsTimeLc            = "3:04 pm"
	EuroTime            = "15:04"
	UsDateTimeSeconds   = "1/2/2006 3:04:05 PM"
	UsDateTimeSecondsLc = "1/2/2006 3:04:05 pm"
	EuroDateTimeSeconds = "2/1/2006 15:04:00"
	LongDateDOW         = "Monday, January 2, 2006"
	LongDate            = "January 2, 2006"
	SqlDate             = "2006-01-02 15:04:05.000000-07"
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

// ParseInOffset is like time.ParseInLocation but uses the given timezone name and offset in minutes from UTC to
// be the location of the parsed time. If the timezone name is blank, the offset will only be used.
func ParseInOffset(layout, value string, tz string, tzOffset int) (t time.Time, err error) {
	if t, err = ParseForgiving(layout, value); err != nil {
		return
	}
	var loc *time.Location
	if tz != "" {
		loc, err = time.LoadLocation(tz)
	}
	if loc != nil && err == nil {
		t = As(t, loc)
	} else {
		t = As(t, time.FixedZone(tz, tzOffset*60))
		err = nil
	}
	return t, err
}

// ParseForgiving will parse a value allowing for extra spaces or a difference in case for am/pm.
func ParseForgiving(layout, value string) (time.Time, error) {
	layout = strings.ReplaceAll(layout, " ", "")
	value = strings.ReplaceAll(value, " ", "")
	layout = strings.ReplaceAll(layout, "pm", "PM")
	value = strings.ReplaceAll(value, "pm", "PM")
	value = strings.ReplaceAll(value, "am", "AM")
	return time.Parse(layout, value)
}

// FromSqlDateTime will convert a SQL Date, Time, DateTime or Timestamp string to a time.Time. Various SQL databases
// express their times slightly differently, and this tries to interpret what is
// attempting to be expressed. It can also handle unix time strings that are +- from
// the 1970 epoch, including fractional times up to the microsecond level.
//
// If the SQL date time string does not have timezone information,
// the resulting value will be in UTC time.
// If an error occurs, the returned value will be the zero date.
func FromSqlDateTime(s string) (t time.Time) {
	var form string

	// First check for a unix time
	if u, e := strconv.ParseFloat(s, 32); e == nil {
		i, f := math.Modf(u)
		t = time.Unix(int64(i), int64(f*1000000)*1000)
		return
	}

	if len(s) > 10 && s[10] == 'T' {
		form = time.RFC3339
	} else {
		var hasDate, hasTime, hasTZ, hasLocale bool
		if strings.Contains(s, "-") {
			hasDate = true
		}
		if strings.Contains(s, ":") {
			hasTime = true

			if strings.LastIndexAny(s, "+-") > strings.LastIndex(s, ":") {
				hasTZ = true

				if s[len(s)-1] == 'T' {
					hasLocale = true
				}
			}
		}
		if hasDate {
			form = "2006-01-02"
			if hasTime {
				form += " 15:04:05"
				if hasTZ {
					form += " -0700"
					if hasLocale {
						form += " MST"
					}
				}
			}
		} else {
			form = "15:04:05"
		}
	}
	t, err := time.Parse(form, s)
	if err == nil {
		t = t.UTC()
	} else {
		// We can't return the error, but we can log it
		log.Warning(err)
	}
	return
}

// ToSqlDateTime returns a SQL string representing the give time.Time value.
func ToSqlDateTime(t time.Time) (s string) {
	return t.Format(SqlDate)
}

// LayoutHasDate returns true if the given parse layout indicates a date.
func LayoutHasDate(layout string) bool {
	return strings2.ContainsAnyStrings(layout, "6", "2", "Jan")
}

// LayoutHasTime returns true if the given parse layout indicates a time.
func LayoutHasTime(layout string) bool {
	return strings2.ContainsAnyStrings(layout, "15", "4", "5", "3")
}
