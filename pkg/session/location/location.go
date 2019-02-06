// The location package implements a location queue built on top of the session service.
// You can use it to keep track of the path a user took to get to a particular location,
// and then allow that user to back out of that path. This is helpful in situations where
// there are multiple paths to get to a page, but the user presses a Save or Cancel button
// indicating the user wants to go back to whatever location they were previously at.
package location

import (
	"context"
	"github.com/goradd/goradd/pkg/session"
)

const key = "goradd.locations"

// Push pushes the given location onto the location stack.
func Push(ctx context.Context, loc string) {
	session.PushStack(ctx, key, loc)
}

// Pop pops the given location off of the location stack and returns it.
func Pop(ctx context.Context) (loc string) {
	return session.PopStack(ctx, key)
}

// Clear removes all locations from the location stack
func Clear(ctx context.Context) {
	session.ClearStack(ctx, key)
}

