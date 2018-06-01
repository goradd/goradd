// The datetime package contains utilities for time and date related functions. It is a wrapper for the time.Time package,
// with expanded functionality and routines to help shepard between the types, and to provide localization and better
// timezone support in the context of a web application.
package datetime

import (
	"time"
	"encoding/json"
	"github.com/spekary/goradd/javascript"
	"fmt"
)


const (
	Current = "now"
	Zero = "zero"
	UsDate = "1/2/2006"
	EuroDate = "2/1/2006"
	LongDateDOW = "Monday, January 2, 2006"
	LongDate = "January 2, 2006"
)

type DateTime struct {
	time.Time
}

func Now() DateTime {
	return DateTime{time.Now()}
}

//
func NewDateTime(args... interface{}) DateTime {
	d := DateTime{}
	v := args[0]
	if v == nil {
		return d
	}

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
		} else if c == Zero {
			// do nothing, we are already zero'd
		} else {
			if len(args) > 1 {
				d.Time, _ = time.Parse(args[1].(string), c)
			} else {
				d.UnmarshalText([]byte(c))
			}
		}
	}
	return d
}

func (d DateTime) Equal(d2 DateTime) bool {
	return d.Time.Equal(d2.Time)
}

// Satisfies the javacript.JavaScripter interface to output the date as a javascript value
func (d DateTime) JavaScript() string {
	return fmt.Sprintf("new Date(%d, %d, %d, %d, %d, %d)", d.Year(), d.Month() - 1, d.Day(), d.Hour(), d.Minute(), d.Second())
}

// Satisfies the json.Marshaller interface to output the date as a value embedded in JSON and that will be unpacked by our javascript file.
func (d DateTime) MarshalJSON() (buf []byte, err error) {
	// We specify numbers explicitly to avoid the warnings about browsers parsing date strings inconsistently
	var obj = map[string]interface{} {
		javascript.JsonObjectType: "datetime",
		"year": d.Year(),
		"month": d.Month() - 1, // javascript is zero based
		"day": d.Day(),
		"hour": d.Hour(),
		"minute": d.Minute(),
		"second": d.Second(),
	}

	buf,err = json.Marshal(obj)
	return

	// TODO: Deal with timezones vs. local time. As of 2017, there still is not a good consistent javascript way of discovering the browser timezone.
}

