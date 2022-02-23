package session

import (
	"context"
	"net/http"
)

type Mock struct {
}

func NewMock() *Mock {
	return new(Mock)
}

func (mgr Mock) Use(next http.Handler) http.Handler {
	return nil
}

// With inserts the mock session into the current session
func (mgr Mock) With(ctx context.Context) context.Context {
	sessionData := NewSession()

	ctx = context.WithValue(ctx, sessionContext, sessionData)
	return ctx
}
