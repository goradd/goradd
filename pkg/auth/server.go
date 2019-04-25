// package auth provides a default authentication framework based on user name and password.
// To use it, call RegisterAuthenticationService and pass it a structure that will handle the various
// routines for authentication.
package auth

import (
	"context"
	"github.com/goradd/goradd/web/app"
	"net/http"
)

// TODO: Implement Open ID with OAuth mechanisms too. Remember when doing this that you need both. OAuth is
// only for authorization and NOT authentication!


// The approach here:
// The framework here can be used in a few different ways. You can implement basic authentication, or use a bearer
// token as a refresh token. You can return these as headers, or in the body of the response. Its up to you.
//
// However, once the user successfully logs in, either after creating a new account, or using login credentials,
// you should save the user's id or some other kind of identity token in the session, and then from that point
// on the session essentially becomes the access token. This way, you do not have to pass an access token
// every time, but you could choose to do that approach if you wanted.


import (
	"encoding/json"
	"fmt"
	"github.com/goradd/goradd/pkg/goradd"
	"github.com/goradd/goradd/pkg/session"
	"io/ioutil"
	"time"
)

func MakeAuthApiServer(a app.ApplicationI) http.Handler {
	// the handler chain gets built in the reverse order of getting called
	// These handlers are called in reverse order
	h := serveAuthApi()
	h = a.SessionHandler(h)

	return h
}

// serveAuthHandler serves up the auth api.
func serveAuthApi() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		authApiHandler(w,r)
	}
	return http.HandlerFunc(fn)
}

const (
	AuthHello string = "hello" // Requires a session in order to create a new user. Helps us rate limit new user requests.
	AuthNewUser string = "new"
	AuthOpLogin string = "login"
	AuthOpTokenLogin string = "token"
	AuthOpRevoke string = "logout"
	AuthOpRecover string = "recover"
)

// authMessage is the message sent by the server.
type authMessage struct {
	Operation string `json:"op"` // Required. Possibilities are AuthOp* options above

	UserName string `json:"user"` 	// Only present for login
	Password string `json:"pw"`		// Only for login
	Token string `json:"token"`		// Only for token and logout
	RecoverMethod string `json:"method"`		// Only for recover
}

func authApiHandler(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var msg authMessage
	err = json.Unmarshal(b, &msg)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	ctx := r.Context()

	switch msg.Operation {
	case AuthHello:
		// Hello is only for the first time connection in order to establish a session. We use it as a kind of authentication
		// mechanism to frustrate mischief.
		// If the session is already established, it means someone is behaving badly.
		if !session.Has(ctx, goradd.SessionAuthTime) {
			// Valid first time, so we set up a timestamp for rate limiting new account requests
			session.Set(ctx, goradd.SessionAuthTime, time.Now().Unix())
			// TODO: Prevent a DoS attack here by checking for rapid hellos from the same IP address and rate limiting them
			// If we do it, one article suggested returning a 200 in order to prevent an attacker from knowing they were being rate limited
		} else {
			// Someone is sending us a hello when we already have a session. Provided we build our apps to check, this
			// would only be done by someone trying to abuse the endpoint. We simply frustrate them by
			// delaying.
			time.Sleep(HackerDelay)
		}
	case AuthNewUser:
		// Requesting a new user account.
		if !session.Has(ctx, goradd.SessionAuthTime) {
			// Session was not established. Ask them to say hello first.
			w.WriteHeader(404)
			w.Write([]byte("say hello"))
		} else if session.Has(ctx, goradd.SessionAuthSuccess) {
			// This session already has created a new user or has successfully logged in, do not allow this a second time
			//time.Sleep(HackerDelay)
			w.WriteHeader(400)
		} else if authNewUser(ctx, msg.UserName, msg.Password, w) {
			// if the user was successfully created, we mark that so that the same session cannot create another user
			session.Set(ctx, goradd.SessionAuthSuccess, true)
		}
	case AuthOpLogin:
		// Logging in using a user name and password
		if !session.Has(ctx, goradd.SessionAuthTime) {
			// Session was not established. Ask them to say hello first.
			w.WriteHeader(404)
			w.Write([]byte("say hello"))
		} else if session.Has(ctx, goradd.SessionAuthSuccess) {
			// This session already has created a new user or has successfully logged in, do not allow this a second time
			//time.Sleep(HackerDelay)
			w.WriteHeader(400)
		} else {
			lastLogin := session.Get(ctx, goradd.SessionAuthTime).(int64)
			now := time.Now().Unix()
			if now - lastLogin <= LoginRateLimit {
				http.Error(w, fmt.Sprintf("%d", LoginRateLimit - (now - lastLogin) + 1), 425)
			} else {
				if authLogin(ctx, msg.UserName, msg.Password, w) {
					// if the user was successfully logged in, we mark that so that the same session cannot create another user or log in
					session.Set(ctx, goradd.SessionAuthSuccess, true)
				}
				session.Set(ctx, goradd.SessionAuthTime, time.Now().Unix())
			}
		}
	case AuthOpTokenLogin:
		// Logging in using a token
		if !session.Has(ctx, goradd.SessionAuthTime) {
			// Session was not established. Ask them to say hello first.
			w.WriteHeader(404)
			w.Write([]byte("say hello"))
		} else if session.Has(ctx, goradd.SessionAuthSuccess) {
			// This session already has created a new user or has successfully logged in, do not allow this a second time
			//time.Sleep(HackerDelay)
			w.WriteHeader(400)
		} else {
			lastLogin := session.Get(ctx, goradd.SessionAuthTime).(int64)
			now := time.Now().Unix()
			if now - lastLogin <= LoginRateLimit {
				http.Error(w, fmt.Sprintf("%d", LoginRateLimit - (now - lastLogin) + 1), 425)
			} else {
				if authTokenLogin(ctx, msg.Token, w) {
					// if the user was successfully logged in, we mark that so that the same session cannot create another user or log in
					session.Set(ctx, goradd.SessionAuthSuccess, true)
				}
				session.Set(ctx, goradd.SessionAuthTime, time.Now().Unix())
			}
		}
	case AuthOpRevoke:
		if !session.Has(ctx, goradd.SessionAuthTime) {
			// Session was not established. Ask them to say hello first.
			w.WriteHeader(404)
			w.Write([]byte("say hello"))
		} else {
			authRevoke(ctx, msg.Token)
			// kill the session
			session.Clear(ctx)
			session.Reset(ctx)
		}
	case AuthOpRecover:
		if !session.Has(ctx, goradd.SessionAuthTime) {
			// Session was not established. Ask them to say hello first.
			w.WriteHeader(404)
			w.Write([]byte("say hello"))
		} else if session.Has(ctx, goradd.SessionAuthSuccess) {
			// Trying to recover when the person is already logged in. This makes no sense.
			w.WriteHeader(400)
		} else {
			// rate limit recovery attempts
			lastLogin := session.Get(ctx, goradd.SessionAuthTime).(int64)
			now := time.Now().Unix()
			if now - lastLogin <= LoginRateLimit {
				http.Error(w, fmt.Sprintf("%d", LoginRateLimit - (now - lastLogin) + 1), 425)
			} else {
				authRecover(ctx, msg.RecoverMethod)
				session.Set(ctx, goradd.SessionAuthTime, time.Now().Unix())
			}
		}
	default:
		http.Error(w, "Invalid operation: " + msg.Operation, 500)
	}
}

func authNewUser(ctx context.Context, user string, password string, w http.ResponseWriter) bool {
	return authService.NewUser(ctx, user, password, w)
}

func authLogin(ctx context.Context, user string, password string, w http.ResponseWriter) bool {
	return authService.Login(ctx, user, password, w)
}

func authTokenLogin(ctx context.Context, token string, w http.ResponseWriter) bool {
	return authService.TokenLogin(ctx, token, w)
}

func authRevoke(ctx context.Context, token string) {
	authService.RevokeToken(ctx, token)
}

func authRecover(ctx context.Context, method string) {
	authService.Recover(ctx, method)
}