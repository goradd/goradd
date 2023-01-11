// Package auth provides on possible authentication framework based on username and password.
// To use it, call RegisterAuthenticationService and pass it a structure that will handle the various
// routines for authentication.
//
// On the client side, you send an auth operation to the auth endpoint by sending 2 form values:
// 1) An "op" value, which is the operation to perform. See OpHello, etc. for possible values
// 2) A "msg" value, which gets passed on the auth service you provide.
//
// You coordinate between your client and your service on how to encode your messages. A common way would be
// to use json, but its up to you to do the encoding and decoding on either end.
//
// See the AuthI interface for details on what each message type should accomplish.
//
// This authentication system can be used to implement a token authorization flow, which is not the recommended way of doing
// authorization in apps, but is an acceptable way of doing it. It requires a client that can store cookies.
// The primary problem with this flow is that the client must gather that username and password, which can be
// problematic for mobile app clients since they must be very careful to not locally store the password accidentally.
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
	"fmt"
	"github.com/goradd/goradd/pkg/goradd"
	"github.com/goradd/goradd/pkg/session"
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
		authApiHandler(w, r)
	}
	return http.HandlerFunc(fn)
}

// These are the operations accepted in the op form variable
const (
	OpHello      = "hello" // Requires a session in order to create a new user. Helps us rate limit new user requests.
	OpNewUser    = "new"
	OpLogin      = "login"
	OpTokenLogin = "token"
	OpRevoke     = "logout"
	OpRecover    = "recover" // When passing this one, you will need to specify your own recovery method in the message.
)

// authMessage is the message sent by the server. It comes through as JSON and gets unpacked into
// this structure. The keys that are reserved and that we use are listed below. You can send other
// keys in the message as well to meet your needs, and the message will get sent on to the auth service.
// The op keyword is always required. The others are required depending on the operation.
//type authMessage map[string]interface{}

const (
	formOperation = "op"  // The operation to perform. Possibilities are AuthOp* options above.
	formMsg       = "msg" // This is the message from your client to your service. Use any format you wish.
)

func authApiHandler(w http.ResponseWriter, r *http.Request) {
	op := r.FormValue(formOperation)
	msg := r.FormValue(formMsg)
	ctx := r.Context()

	switch op {
	case OpHello:
		// Hello is only for the first time connection in order to establish a session. We use it as a kind of authentication
		// mechanism to frustrate mischief.
		// If the session is already established, it means someone is behaving badly.
		if !session.Has(ctx, goradd.SessionAuthTime) {
			// Valid first time, so we set up a timestamp for rate limiting new account requests
			// We subtract LoginRateLimit here because we expect the user to try to login once immediately after
			// saying hello. We rate limit any subsequent attempts.
			session.Set(ctx, goradd.SessionAuthTime, time.Now().Unix()-LoginRateLimit)
			// TODO: Prevent a DoS attack here by checking for rapid hellos from the same IP address and rate limiting them
			// If we do it, one article suggested returning a 200 in order to prevent an attacker from knowing they were being rate limited
		} else {
			// Someone is sending us a hello when we already have a session. Provided we build our apps to check, this
			// would only be done by someone trying to abuse the endpoint. We simply frustrate them by
			// delaying.
			time.Sleep(HackerDelay)
		}
	case OpNewUser:
		// Requesting a new user account.
		if !session.Has(ctx, goradd.SessionAuthTime) {
			// Session was not established. Ask them to say hello first.
			authWriteError(ctx, "say hello", 404, w)
		} else if session.Has(ctx, goradd.SessionAuthSuccess) {
			// This session already has created a new user or has successfully logged in, do not allow this a second time
			//time.Sleep(HackerDelay)
			authWriteError(ctx, "session already established", 400, w)
		} else if authNewUser(ctx, []byte(msg), w) {
			// if the user was successfully created, we mark that so that the same session cannot create another user
			session.Set(ctx, goradd.SessionAuthSuccess, true)
		}
	case OpLogin:
		// Logging in using a user name and password
		if !session.Has(ctx, goradd.SessionAuthTime) {
			// Session was not established. Ask them to say hello first.
			authWriteError(ctx, "say hello", 404, w)

		} else if session.Has(ctx, goradd.SessionAuthSuccess) {
			// This client is already logged in. If logging in again, we will assume the client tried to logout
			// and something failed, like a bad connection at that moment. So, we will logout here and force the
			// client to reestablish a session first
			session.Clear(ctx)
			session.Reset(ctx)
			authWriteError(ctx, "say hello", 404, w)
		} else {
			lastLogin := session.Get(ctx, goradd.SessionAuthTime).(int64)
			now := time.Now().Unix()
			if now-lastLogin < LoginRateLimit {
				authWriteError(ctx, fmt.Sprintf("%d", LoginRateLimit-(now-lastLogin)+1), 425, w)
			} else {
				if authLogin(ctx, []byte(msg), w) {
					// if the user was successfully logged in, we mark that so that the same session cannot create another user or log in
					session.Set(ctx, goradd.SessionAuthSuccess, true)
				}
				session.Set(ctx, goradd.SessionAuthTime, time.Now().Unix())
			}
		}
	case OpTokenLogin:
		// Logging in using a token
		if !session.Has(ctx, goradd.SessionAuthTime) {
			// Session was not established. Ask them to say hello first.
			authWriteError(ctx, "say hello", 404, w)
		} else if session.Has(ctx, goradd.SessionAuthSuccess) {
			// This session already has created a new user or has successfully logged in, do not allow this a second time
			//time.Sleep(HackerDelay)
			authWriteError(ctx, "session already established", 400, w)
		} else {
			lastLogin := session.Get(ctx, goradd.SessionAuthTime).(int64)
			now := time.Now().Unix()
			if now-lastLogin < LoginRateLimit {
				authWriteError(ctx, fmt.Sprintf("%d", LoginRateLimit-(now-lastLogin)+1), 425, w)
			} else {
				if authTokenLogin(ctx, []byte(msg), w) {
					// if the user was successfully logged in, we mark that so that the same session cannot create another user or log in
					session.Set(ctx, goradd.SessionAuthSuccess, true)
				}
				session.Set(ctx, goradd.SessionAuthTime, time.Now().Unix())
			}
		}
	case OpRevoke:
		if !session.Has(ctx, goradd.SessionAuthTime) {
			// Session was not established. Ask them to say hello first.
			authWriteError(ctx, "say hello", 404, w)
		} else {
			if authRevoke(ctx, []byte(msg), w) {
				// kill the session
				session.Clear(ctx)
				session.Reset(ctx)
			}
		}
	case OpRecover:
		if !session.Has(ctx, goradd.SessionAuthTime) {
			// Session was not established. Ask them to say hello first.
			authWriteError(ctx, "say hello", 404, w)
		} else if session.Has(ctx, goradd.SessionAuthSuccess) {
			// Trying to recover when the person is already logged in. This makes no sense.
			authWriteError(ctx, "session already established", 400, w)
		} else {
			// rate limit recovery attempts
			lastLogin := session.Get(ctx, goradd.SessionAuthTime).(int64)
			now := time.Now().Unix()
			if now-lastLogin < LoginRateLimit {
				authWriteError(ctx, fmt.Sprintf("%d", LoginRateLimit-(now-lastLogin)+1), 425, w)
			} else {
				authRecover(ctx, []byte(msg), w)
			}
		}
	default:
		if op == "" {
			authWriteError(ctx, "No operation specified", 400, w)
		} else {
			authWriteError(ctx, "Invalid operation: "+op, 400, w)
		}
	}
}

func authNewUser(ctx context.Context, msg []byte, w http.ResponseWriter) bool {
	return authService.NewUser(ctx, msg, w)
}

func authLogin(ctx context.Context, msg []byte, w http.ResponseWriter) bool {
	return authService.Login(ctx, msg, w)
}

func authTokenLogin(ctx context.Context, msg []byte, w http.ResponseWriter) bool {
	return authService.TokenLogin(ctx, msg, w)
}

func authRevoke(ctx context.Context, msg []byte, w http.ResponseWriter) bool {
	// its important the authRevoke not write to the response writer unless there is an error, since we need to control that to close the session
	return authService.RevokeToken(ctx, msg, w)
}

func authRecover(ctx context.Context, msg []byte, w http.ResponseWriter) bool {
	return authService.Recover(ctx, msg, w)
}

func authWriteError(ctx context.Context, errorMessage string, errorCode int, w http.ResponseWriter) {
	authService.WriteError(ctx, errorMessage, errorCode, w)
}
