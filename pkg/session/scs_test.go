package session_test

import (
	"github.com/alexedwards/scs"
	"github.com/alexedwards/scs/stores/memstore"
	"github.com/goradd/goradd/pkg/session"
	"testing"
	"time"
)


func TestSetGet(t *testing.T) {
	// setup the ScsSession
	interval, _ := time.ParseDuration("24h")
	session.SetSessionManager(session.NewScsManager(scs.NewManager(memstore.New(interval))))

	// run the session tests
	runRequestTest(t, setRequestHandler(), testRequestHandler(t))
	runRequestTest(t, setupStackRequestHandler(), testStackRequestHandler(t))
}

