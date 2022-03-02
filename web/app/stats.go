package app

import (
	"fmt"
	"net/http"
	"time"
)

var minDuration, maxDuration, totalDuration time.Duration
var totalServes int64

// StatsHandler records performance statistics on the server.
func (a *Application) StatsHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		end := time.Now()
		d := end.Sub(start)
		if minDuration == 0 || minDuration > d {
			minDuration = d
		}
		if maxDuration == 0 || maxDuration < d {
			maxDuration = d
		}
		totalDuration += d
		totalServes += 1
	}
	return http.HandlerFunc(fn)
}

func GetStats() string {
	return fmt.Sprintf("min: %d us, max: %d us, avg: %d us\n",
		minDuration.Microseconds(),
		maxDuration.Microseconds(),
		totalDuration.Microseconds() / totalServes,
		)
}