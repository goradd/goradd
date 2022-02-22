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

// InOffsetTime converts the time.Time to a fixed locale with the given timezone offset.
// tzOffset is minutes from UTC.
// See page.GetContext(ctx).ClientTimezoneOffset.
func InOffsetTime(t time.Time, tzOffset int) time.Time {
	return t.In(time.FixedZone("", tzOffset * 60))
}

// DayDiff returns the number of days between the two dates, taking into consideration only the day of the month.
// So 1 am on the 2nd and 11 pm on the 1st are a day apart.
// If dt1 is before dt2, the result will be negative. If they are the same day, the result is zero.
// The dates should be set to the timezone that will determine what day the datetime falls on.
func DayDiff(dt1, dt2 time.Time) int {
	d1 := dt1.YearDay() + NumLeaps(dt1.Year() - 1) + dt1.Year() * 365
	d2 := dt2.YearDay() + NumLeaps(dt2.Year() - 1) + dt2.Year() * 365
	return d1 - d2
}

func IsLeap(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

func NumLeaps(year int) int {
	return year / 4 - year / 100 + year / 400
}

// NewDateTime creates a time.Time in UTC.
func NewDateTime(year int, month time.Month, day, hour, min, sec, nsec int) time.Time {
	t := time.Date(year, time.Month(month), day, hour, min, sec, nsec, time.UTC)
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