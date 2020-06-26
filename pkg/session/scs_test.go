package session_test

import (
	"github.com/alexedwards/scs/v2"
	"github.com/alexedwards/scs/v2/memstore"
	"github.com/goradd/goradd/pkg/session"
	"testing"
	"time"
)

func TestSetGet(t *testing.T) {
	// setup the ScsSession
	store := memstore.NewWithCleanupInterval(24 * time.Hour)
	sm := scs.New()
	sm.Store = store
	session.SetSessionManager(session.NewScsManager(sm))

	// run the session tests
	runRequestTest(t, setRequestHandler(), testRequestHandler(t))
	runRequestTest(t, setupStackRequestHandler(), testStackRequestHandler(t))
}
