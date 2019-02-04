package session

import (
	"context"
	"github.com/alexedwards/scs"
	"net/http"
)

const scsSessionDataKey = "goradd.data"

// SCS_Manager satisfies the ManagerI interface for the github.com/alexedwards/scs session manager. You can use it as an example
// of how to incorporate a different session manager into your app.
type SCS_Manager struct {
	*scs.Manager
}

func NewSCSManager(mgr *scs.Manager) ManagerI {
	return SCS_Manager{mgr}
}

// Use is an http handler that wraps the session management process. It will get and put session data
// into the http context.
func (mgr SCS_Manager) Use(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var data []byte
		var temp string
		// get the session. All of our session data is stored in only one key in the session manager.
		session := mgr.Manager.Load(r)
		session.Touch(w) // Make sure to get a cookie in our header if we don't have one
		data, _ = session.GetBytes(scsSessionDataKey)
		sessionData := NewSession()
		if data != nil {
			sessionData.UnmarshalBinary(data)
		}

		if sessionData.Has(sessionResetKey) {
			// Our previous session requested a reset. We can't reset after writing, so we reset here at the start of the next request.
			sessionData.Delete(sessionResetKey)
			session.RenewToken(w)
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, sessionContext, sessionData)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)

		// write out the changed session. The below will attempt to write a cookie, but it can't because headers have already been written.
		// That is OK, because of our Touch above.
		if sessionData.Len() > 0 {
			var err error
			data, err = sessionData.MarshalBinary()
			if err != nil {
				w.Write([]byte(err.Error()))
			}
			temp = string(data)
			_ = temp
			session.PutBytes(w, scsSessionDataKey, data)
		} else {
			session.Clear(w)
		}
	}
	return mgr.Manager.Use(http.HandlerFunc(fn))
}

