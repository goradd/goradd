package session

import (
	"context"
	"github.com/alexedwards/scs/v2"
	"github.com/goradd/goradd/pkg/log"
	"net/http"
	"time"
)

const scsSessionDataKey = "goradd.data"

// ScsManager satisfies the ManagerI interface for the github.com/alexedwards/scs session manager.
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
			log.Errorf("Error loading or unpacking session: %s", err.Error()) // we can't panic here, because our panic handlers have not been set up
		}
		var sessionData *Session
		if data := mgr.SessionManager.Get(ctx, scsSessionDataKey); data != nil {
			sessionData = data.(*Session)
			log.FrameworkDebug("Found session")
		} else {
			sessionData = NewSession()
			log.FrameworkDebug("Creating new session")
		}

		ctx = context.WithValue(ctx, sessionContext, sessionData)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)

		//  Not sure why this is here. See LoadAndSave from SCS.
		/*
		if r.MultipartForm != nil {
			r.MultipartForm.RemoveAll()
		}*/


		if sessionData.Has(sessionResetKey) {
			// Our session requested a reset. We are safe to do this here because the buffered output handler
			// will ensure that we can write headers even if output has been sent.
			sessionData.Delete(sessionResetKey)
			if err := mgr.SessionManager.RenewToken(ctx); err != nil {
				log.Errorf("Error renewing session token: %s", err.Error())
			}
		}

		if sessionData.Len() > 0 {
			mgr.SessionManager.Put(ctx, scsSessionDataKey, sessionData)
			token, expiry, err := mgr.SessionManager.Commit(ctx)
			if err != nil {
				s := "Error marshalling session data: " + err.Error()
				log.Error(s)
				http.Error(w, s, 500)
				return
			}
			log.FrameworkDebug("Writing session cookie")
			writeSessionCookie(w, mgr.SessionManager.Cookie, token, expiry)
		} else {
			if err := mgr.SessionManager.Clear(ctx); err != nil {
				log.Errorf("Error clearing session: %s", err.Error())
			}
		}
	}
	return http.HandlerFunc(fn)
}


// The following two function come from SCS source code. Copyright (c) 2016 Alex Edwards.

func writeSessionCookie(w http.ResponseWriter, sc scs.SessionCookie, token string, expiry time.Time) {
	cookie := &http.Cookie{
		Name:     sc.Name,
		Value:    token,
		Path:     sc.Path,
		Domain:   sc.Domain,
		Secure:   sc.Secure,
		HttpOnly: sc.HttpOnly,
		SameSite: sc.SameSite,
	}

	if expiry.IsZero() {
		cookie.Expires = time.Unix(1, 0)
		cookie.MaxAge = -1
	} else if sc.Persist {
		cookie.Expires = time.Unix(expiry.Unix()+1, 0)        // Round up to the nearest second.
		cookie.MaxAge = int(time.Until(expiry).Seconds() + 1) // Round up to the nearest second.
	}

	w.Header().Add("Set-Cookie", cookie.String())
	addHeaderIfMissing(w, "Cache-Control", `no-cache="Set-Cookie"`)
	addHeaderIfMissing(w, "Vary", "Cookie")

}

func addHeaderIfMissing(w http.ResponseWriter, key, value string) {
	for _, h := range w.Header()[key] {
		if h == value {
			return
		}
	}
	w.Header().Add(key, value)
}
