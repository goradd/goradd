package session_test

import (
	"context"
	"github.com/goradd/goradd/pkg/session"
	"testing"
)

func TestMockSetGet(t *testing.T) {
	// setup the mock session
	s := session.NewMock()
	session.SetSessionManager(s)
	ctx := s.With(context.Background())

	// run the session tests
	setupTest(ctx)
	runTest(t, ctx)
}
