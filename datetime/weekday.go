package datetime

import (
	"time"
	"context"
	"goradd/i18n"	// project specific translation
)

type Weekday int

const (
	Sunday Weekday = iota
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
)

func (d Weekday) String() string { return time.Weekday(d).String() }

// Translate implements the i18n.Translater interface
func (d Weekday) Translate(ctx context.Context) string {
	return i18n.T(ctx, d.String())
}