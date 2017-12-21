package datetime

import (
	"time"
	"grlocal/i18n"
	"context"
)

// A Month specifies a month of the year (January = 1, ...).
type Month int

const (
	January Month = 1 + iota
	February
	March
	April
	May
	June
	July
	August
	September
	October
	November
	December
)

func (m Month) String() string { return time.Month(m).String() }

// Translate implements the i18n.Translater interface
func (m Month) Translate(ctx context.Context) string {
	return i18n.T(ctx, m.String())
}
