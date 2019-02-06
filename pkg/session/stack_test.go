package session_test

import (
	"github.com/goradd/goradd/pkg/session"
	"github.com/goradd/goradd/pkg/session/location"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

const stack = "test.stack"

func setupStackRequestHandler() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()

		session.PushStack(ctx, stack, "A")
		session.PushStack(ctx, stack, "B")
		session.PushStack(ctx, stack, "C")

		location.Push(ctx, "Here")
		location.Push(ctx, "There")
		location.Clear(ctx)
	}
	return http.HandlerFunc(fn)
}

func testStackRequestHandler(t *testing.T) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()

		assert.Equal(t, "C", session.PopStack(ctx, stack))
		assert.Equal(t, "B", session.PopStack(ctx, stack))
		assert.Equal(t, "A", session.PopStack(ctx, stack))
		assert.Equal(t, "", session.PopStack(ctx, stack))

		assert.Equal(t, "", location.Pop(ctx))
	}
	return http.HandlerFunc(fn)
}
