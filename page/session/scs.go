package session

import (
	"net/http"
	"github.com/alexedwards/scs"
	"context"
)

// SCS_Manager satisfies the ManagerI interface for the github.com/alexedwards/scs session manager
type SCS_Manager struct {
	*scs.Manager
}

func NewSCSManager(mgr *scs.Manager) ManagerI {
	return SCS_Manager{mgr}
}

func (mgr SCS_Manager) Use (next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var data []byte
		var temp string
		// get the session. All of our session data is stored in only one key in the session manager.
		session := mgr.Manager.Load(r)
		session.Touch(w) // Make sure to get a cookie in our header if we don't have one
		data,_ = session.GetBytes("goradd.data")
		sessionData := NewSession()
		if data != nil {
			sessionData.UnmarshalBinary(data)
		}

		if sessionData.Has(sessionResetKey) {
			// Our previous session requested a reset. We can't reset after writing, so we reset here at the start of the next request.
			sessionData.Remove(sessionResetKey)
			session.RenewToken(w)
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, sessionContext, sessionData)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)

		// write out the changed session. The below will attempt to write a cookie, but it can't because headers have already been written.
		// That is OK, because of our Touch above.
		if sessionData.Len() > 0 {
			data,_ = sessionData.MarshalBinary()
			temp = string(data)
			_ = temp
			session.PutBytes(w, "goradd.data", data)
		} else {
			session.Clear(w)
		}
	}
	return mgr.Manager.Use(http.HandlerFunc(fn))
}

// SCS_Session is a goradd session manager that uses the github.com/alexedwards/scs session manager.
// It implements the SessionI interface
type SCS_Session struct {
	writer http.ResponseWriter
	mgr *scs.Manager
	session *scs.Session
}

func (s *SCS_Session) Load(mgr *scs.Manager, w http.ResponseWriter, r *http.Request) {
	s.writer = w
	s.mgr = mgr
	s.session = mgr.Load(r)
}

