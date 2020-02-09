// The datetime package contains utilities for time and date related functions. It is a wrapper for the time.Time package,
// with expanded functionality and routines to help shepard between the types, and to provide localization and better
// timezone support in the context of a web application.
package datetime

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"github.com/goradd/goradd/pkg/javascript"
	"time"
)

const (
	Current           = "now"
	Zero              = "zero"
	UsDate            = "1/2/2006"
	ShortUsDate       = "1/2/06"
	EuroDate          = "2/1/2006"
	UsDateTime        = "1/2/2006 3:04 PM"
	EuroDateTime      = "1/2/2006 15:04"
	UsTime            = "3:04 PM"
	EuroTime          = "15:04"
	UsDateTimeSeconds = "1/2/2006 3:04:05 PM"
	LongDateDOW       = "Monday, January 2, 2006"
	LongDate          = "January 2, 2006"
)

// DateTime is the default time value for the system (as opposed to the built in time.Time).
// DateTime has a number of enhancements over time.Time
// DateTime is designed to interact nicely with dates and times stored in databases.
// In particular, SQL databases can store a date-time as either a DATETIME, or a TIMESTAMP. It
// also has DATE only and TIME only representations.
//
// A DATETIME and TIME are for situations where the time is independent of the timezone. For example,
// in a scheduling application when you have scheduled a future event at 8:00 am, but in between
// now and then, there is a change int daylight savings time. The event should still be at 8:00 am.
// These amounts are stored as UTC time, since UTC is not associated with any real timezone and has
// no DST issues. Convert to GMT if you truly mean a TIMESTAMP at zero offset.
//
// A TIMESTAMP is for recording particular moments in world time. When the event is moved to another
// location, the time will be displayed in the local timezone as opposed to the timezone it was created
// in. Some databases store these internally in UTC time, so that if the database is moved to a machine
// in another location, the actual time is preserved, and then converted to local time when read from
// the database.
//
// Time only and Date only representations are always stored in UTC time. Time only times
// have a date of January 1, year 1. Date only values have zero times at UTC. This means that
// a Time only representation at time zero and zero time are equal, and
// a Date only representation at the zero date, and a DATETIME representation at time zero are also equal, and you
// should resolve that ambiguity by the context you are using these in.
type DateTime struct {
	time.Time
	isTimestamp bool
}

// SetIsTimestamp will change the state of the date time to the given value. Timestamps are communicated to the server
// as a UTC time, whereas non-timestamp times are communicated in whatever value is currently in the DateTime without
// changing the timezone
func (d *DateTime) SetIsTimestamp(t bool) {
	d.isTimestamp = t
}

// IsTimestamp returns whether the DateTime is a timestamp, representing a particular moment in world time.
func (d *DateTime) IsTimestamp() bool {
	return d.isTimestamp
}

// Now returns the current time as a timestamp, but with the time in local time.
func Now() DateTime {
	return DateTime{time.Now(), true}
}

// Return a date-time that represents an empty date and time
func NewZeroDate() DateTime {
	return DateTime{}
}

// Date creates a DateTime with the given information.
// Set hour, min, sec, nsec and loc to zeros for a date-only representation.
// Set year to 0, and month and day to 1's for a time-only representation.
// Pass nil to loc to indicate a non-timestamp value . Otherwise it will create a timestamp
// in the given timezone.
func Date(year int, month Month, day, hour, min, sec, nsec int, loc *time.Location) DateTime {
	if loc == nil {
		t := time.Date(year, time.Month(month), day, hour, min, sec, nsec, time.UTC)
		return DateTime{t, false}
	} else {
		t := time.Date(year, time.Month(month), day, hour, min, sec, nsec, loc)
		return DateTime{t, true}
	}
}

func DateOnly(year int, month Month, day int) DateTime {
	return Date(year, month, day, 0, 0, 0, 0, nil)
}

// Time creates a DateTime that only represents a time of day.
func Time(hour, min, sec, nsec int) DateTime {
	return Date(0, 1, 1, hour, min, sec, nsec, nil)
}

// NewDateTime creates a new date time from the given information. You can give it the following:
// () = zero time
// (DateTime) = copy the given datetime
// (time.Time or *time.Time) = copy the given time
// (string), as follows:
//   datetime.Current - same as calling datetime.Now()
//   datetime.Zero - same as calling datetime.NewZeroDate()
//   anything else - RFC3339 string
// (string, string) = a string representation of the date and time, followed by a time.Parse layout string
func NewDateTime(args ...interface{}) DateTime {
	d := DateTime{}
	if len(args) == 0 || args[0] == nil {
		return d
	}
	v := args[0]

	switch c := v.(type) {
	case DateTime:
		d = c
	case time.Time:
		d.Time = c
	case *time.Time:
		d.Time = *c
	case string:
		if c == Current {
			d = Now()
		} else if c == Zero || c == "" {
			// do nothing, we are already zero'd
		} else {
			if len(args) == 2 {
				d, _ = Parse(args[1].(string), c)
			} else {
				_ = d.UnmarshalText([]byte(c))
			}
		}
	}
	return d
}

// NewTimestamp is the same as NewDateTime, but also sets it as a Timestamp
func NewTimestamp(args ...interface{}) DateTime {
	d := NewDateTime(args...)
	d.isTimestamp = true
	return d
}

// Returns true if the given DateTime object is equal to the current one.
// Timestamps are evaluated as being at the same instant in world time.
// Non-timestamps are compared just for their values ignoring timestamp.
func (d DateTime) Equal(d2 DateTime) bool {
	// non-timestamp values should be stored in UTC, so they can be compared
	return d.Time.Equal(d2.Time)
}

// After reports whether the time instant d is after u.
func (d DateTime) After(u DateTime) bool {
	return d.Time.After(u.Time)
}

// Before reports whether the time instant d is before u.
func (d DateTime) Before(u DateTime) bool {
	return d.Time.Before(u.Time)
}



// Satisfies the javacript.JavaScripter interface to output the date as a javascript value.
// TIMESTAMPS are converted to the local time corresponding to the given world time.
// Non-timestamps are transmitted as if they were in the browser's local time.
func (d DateTime) JavaScript() string {
	if d.IsZero() {
		return "null"
	} else if d.isTimestamp {
		t := d.Time.UTC()
		return fmt.Sprintf("new Date(Date.UTC(%d, %d, %d, %d, %d, %d, %d))", t.Year(), t.Month()-1, t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond()/1000000)
	} else {
		return fmt.Sprintf("new Date(%d, %d, %d, %d, %d, %d, %d)", d.Year(), d.Month()-1, d.Day(), d.Hour(), d.Minute(), d.Second(), d.Nanosecond()/1000000)
	}
}

func (d *DateTime) UnmarshalJSON(data []byte) (err error) {
	err = d.Time.UnmarshalJSON(data)
	if err != nil {
		err = fmt.Errorf("JSON dates and times must be ISO8601 formatted AND in UTC: %s", err.Error());
	}
	d.isTimestamp = true // json times are always utc
	return
}

// MarshalJSON satisfies the json.Marshaller interface to output the date as a value embedded in JSON and that
// will be unpacked by our javascript file.
//
// The JSON "standard" behavior is to output an ISO8601 string in UTC time. There are a bunch of problems with this.
// 1) You cannot use this to represent a display date and time in local time. It will get converted to world time.
//	  even if you try to output a time without timezone information, or with timezone information that is not UTC,
//    which is what the built-in GO Time object does, most browsers will ignore that part of the string and completely
//	  botch the time transfer. We would need to create our own interpreter to make that work.
// 2) Its a string. Unless you know in advance you want a date, or you look for an ISO8601 pattern in every string,
//    you cannot just tell the other side its a date.
//
// So, we have taken the approach of outputting an object that self-identifies as a goradd date, and already
// breaks apart the time components so its easily reassembled on the other side.
func (d DateTime) MarshalJSON() (buf []byte, err error) {
	// We specify numbers explicitly to avoid the warnings about browsers parsing date strings inconsistently
	isTimestamp := d.IsTimestamp()
	var t time.Time
	if isTimestamp {
		t = d.Time.UTC()
	} else {
		t = d.GoTime()
	}
	var obj = map[string]interface{}{
		javascript.JsonObjectType: "date",
		"y":   t.Year(),
		"mo":  t.Month() - 1, // javascript is zero based
		"d":    t.Day(),
		"h":   t.Hour(),
		"m": t.Minute(),
		"s": t.Second(),
		"ms": t.Nanosecond()/1000000,
		"t": isTimestamp,
		"z": d.IsZero(),
	}

	buf, err = json.Marshal(obj)
	return
}

// GoTime returns a GO time.Time value. Be careful, as this returns the time in whatever timestamp it currently
// is set to. Consider converting to UTC time before getting this.
func (d DateTime) GoTime() time.Time {
	return d.Time
}

// GetTime returns a new DateTime object set to only the time portion of the given DateTime object.
func (d DateTime) GetTime() DateTime {
	return Time(d.Hour(), d.Minute(), d.Second(), d.Nanosecond())
}

// GetDate returns a new DateTime object set to only the date portion of the given DateTime object.
func (d DateTime) GetDate() DateTime {
	return Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, nil)
}

// Month returns the month value of the datetime
func (d DateTime) Month() Month {
	return Month(d.Time.Month())
}

// Local returns a new DateTime with the date and time converted to local values in the server's timezone.
func (d DateTime) Local() DateTime {
	return DateTime{d.Time.Local(), d.isTimestamp}
}

// UTC returns a new datetime in UTC time.
func (d DateTime) UTC() DateTime {
	return DateTime{d.Time.UTC(), d.isTimestamp}
}

// In converts the datetime to the given locale.
func (d DateTime) In(location *time.Location) DateTime {
	return DateTime{d.Time.In(location), d.isTimestamp}
}

func (d DateTime) AddDate(years int, months int, days int) DateTime {
	return DateTime{d.Time.AddDate(years, months, days), d.isTimestamp}
}


// String returns the datetime in string form, suitable for sending to the NewDateTime function.
func (d DateTime) String() string {
	b, err := d.MarshalText()
	_ = err // would only happen if year was greater than 10000
	return string(b)
}

// DebugString returns the datetime in human readable form.
func (d DateTime) DebugString() string {
	return d.Time.String()
}

// ToJsonString outputs the date in ISO8601, which is the generally excepted representation of a
// date. The problem with this is that the receiver must know to look for a date instead of a string,
// so it is not actually the best way of transmitting a date to the other side. This is here for
// use by the ORM, which will send standard dates in the expected JSON format for potential consumption
// by REST services.
func (d DateTime) ToJsonString() string {
	isTimestamp := d.IsTimestamp()
	var t time.Time
	if isTimestamp {
		t = d.Time.UTC() // should put a trailing Z to indicate a UTC time
	} else {
		t = d.GoTime() // TODO: need to change this so it definitely does not transmit timezone info
	}

	b,err := t.MarshalJSON()
	if err != nil {
		return ""
	}
	return string(b)
}

func init() {
	gob.Register(time.Time{})
	gob.Register(DateTime{})
}
