// The datetime package contains utilities for time and date related functions. It is a wrapper for the time.Time package,
// with expanded functionality and routines to help shepard between the types, and to provide localization and better
// timezone support in the context of a web application.
package datetime

import "time"


const (
	Current = "now"
	Zero = "zero"
	UsDate     = "1/2/2006"

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

