// Package time has some utilities for time.Time values
package time

import (
	"encoding/gob"
	"time"
)

// As will express the date and time at a particular location.
// In other words, if the date and time is 4:30, it will be 4:30 in the given timezone.
func As(t time.Time, location *time.Location) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), location)
}

// DayDiff returns the number of days between the two dates, taking into consideration only the day of the month.
// So 1 am on the 2nd and 11 pm on the 1st are a day apart.
//
// If dt1 is before dt2, the result will be negative. If they are the same day, the result is zero.
// The dates should already be set to the timezone that will determine what day the datetime falls on.
// This function also takes into account leap years (assuming that you are dealing with Gregorian dates).
func DayDiff(dt1, dt2 time.Time) int {
	d1 := dt1.YearDay() + NumLeaps(dt1.Year()-1) + (dt1.Year()-1)*365
	d2 := dt2.YearDay() + NumLeaps(dt2.Year()-1) + (dt2.Year()-1)*365
	return d1 - d2
}

// IsLeap returns true if the year is a leap year.
func IsLeap(year int) (isLeap bool) {
	if year <= 1582 {
		return year%4 == 0
	} else {
		return year%4 == 0 && (year%100 != 0 || year%400 == 0)
	}
}

// NumLeaps returns the number of leap years that have occurred since year zero. Note that leap years changed
// slightly with the Gregorian calendar in 1582.
func NumLeaps(year int) (leaps int) {
	leaps = year / 4
	if year > 1582 {
		leaps += -year/100 + year/400
	}
	return
}

// NewDateTime creates a time.Time in UTC.
func NewDateTime(year int, month time.Month, day, hour, min, sec, nsec int) time.Time {
	t := time.Date(year, month, day, hour, min, sec, nsec, time.UTC)
	return t
}

// NewDate creates a time.Time that is treated as a date only.
func NewDate(year int, month time.Month, day int) time.Time {
	return NewDateTime(year, month, day, 0, 0, 0, 0)
}

// NewTime creates a time.Time that only represents a time of day.
func NewTime(hour, min, sec, nsec int) time.Time {
	return NewDateTime(0, 1, 1, hour, min, sec, nsec)
}

// TimeOnly returns a new time.Time object set to only the time portion of the given time.
func TimeOnly(t time.Time) time.Time {
	return NewTime(t.Hour(), t.Minute(), t.Second(), t.Nanosecond())
}

// DateOnly returns a new time.Time object set to only the date portion of the given time.
func DateOnly(t time.Time) time.Time {
	return NewDate(t.Year(), t.Month(), t.Day())
}

func init() {
	gob.Register(time.Time{})
}
