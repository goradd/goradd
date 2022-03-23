package session

import (
	"context"
	"github.com/alexedwards/scs/v2"
	"github.com/goradd/goradd/pkg/log"
	"net/http"
)

const scsSessionDataKey = "goradd.data"

// ScsManager implements the ManagerI interface for the github.com/alexedwards/scs session manager.
//
// Note that this manager does post-processing on the response writer, including
// writing headers. It therefore relies on the buffered output handler to collect the output
// before writing the headers.
type ScsManager struct {
	*scs.SessionManager
}

func NewScsManager(mgr *scs.SessionManager) ManagerI {
	return ScsManager{mgr}
}

// Use is an http handler that wraps the session management process. It will get and put session data
// into the http context.
func (mgr ScsManager) Use(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var token string

		// get the session. All of our session data is stored in only one key in the session manager.

		cookie, err := r.Cookie(mgr.SessionManager.Cookie.Name)
		if err == nil {
			token = cookie.Value
		}

		ctx, err := mgr.SessionManager.Load(r.Context(), token)

		if err != nil {
			panic("Error loading or unpacking session: " + err.Error())
		}

		var sess *Session
		if d := mgr.SessionManager.Get(ctx, scsSessionDataKey); d != nil {
			sess = d.(*Session)
			log.FrameworkDebug("Found session")
		} else {
			sess = NewSession()
			log.FrameworkDebug("Creating new session")
		}

		ctx = context.WithValue(ctx, sessionContext, sess)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)

		//  Not sure why this is here. See LoadAndSave from SCS.
		/*
			if r.MultipartForm != nil {
				r.MultipartForm.RemoveAll()
			}*/

		if sess.data.Has(sessionResetKey) {
			// Our session requested a reset. We are safe to do this here because the buffered output handler
			// will ensure that we can write headers even if output has been sent.
			sess.data.Delete(sessionResetKey)
			if err = mgr.SessionManager.RenewToken(ctx); err != nil {
				panic("Error renewing session token: %s" + err.Error())
			}
		}

		if sess.data.Len() > 0 {
			mgr.SessionManager.Put(ctx, scsSessionDataKey, sess)
			token2, expiry, _ := mgr.SessionManager.Commit(ctx)
			if err != nil {
				panic("Error marshalling session data: " + err.Error())
				return
			}
			log.FrameworkDebug("Writing session cookie")
			mgr.SessionManager.WriteSessionCookie(ctx, w, token2, expiry)
		} else {
			if err = mgr.SessionManager.Clear(ctx); err != nil {
				panic("Error clearing session: " + err.Error())
			}
		}
	}
	return http.HandlerFunc(fn)
}
