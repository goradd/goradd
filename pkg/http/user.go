package http

import "net/http"

// User is the interface for http managers that can be injected into the
// handler stack.
type User interface {
	// Use wraps the given handler.
	Use(http.Handler) http.Handler
}
