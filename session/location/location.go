// The location package implements a location queue built on top of the session service.
// You can use it to keep track of the path a user took to get to a particular location,
// and then allow that user to back out of that path. This is helpful in situations where
// there are multiple paths to get to a page, but the user presses a Save or Cancel button
// indicating the user wants to go back to whatever location they were previously at.
package location

import (
	"context"
	"github.com/spekary/goradd/session"
)

const key = "goradd.locations"
func Push(ctx context.Context, loc string) {
	var locations []string
	if session.Has(ctx, key) {
		locations,_ = session.Get(ctx, key).([]string)
	}
	locations = append(locations, loc)
	session.Set(ctx, key, locations)
}

func Pop(ctx context.Context) (loc string) {
	var locations []string
	if session.Has(ctx, key) {
		locations,_ = session.Get(ctx, key).([]string)
	}
	if len(locations) > 0 {
		loc = locations[len(locations)-1]
		locations = locations[:len(locations)-1]
	}

	if len(locations) == 0 {
		session.Remove(ctx, key)
	} else {
		session.Set(ctx, key, locations)
	}
	return
}

func Clear(ctx context.Context) {
	session.Remove(ctx, key)
}

