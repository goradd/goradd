package session

import (
	"net/http"
	"context"
	"github.com/spekary/goradd/util/types"
)

type sessionContextType string

const sessionContext sessionContextType = "goradd.session"
const sessionResetKey string = "goradd.reset"

var sessionManager ManagerI

// ManagerI is the interface for session managers.
type ManagerI interface {
	// Use wraps the given handler in session management stuff. It must decode the session and put it into the sessionContext key in the context
	// before processing the request, and then encode the session data after the request. See the SCS_Manager for an example.
	Use(http.Handler) http.Handler
}


// SetSessionManager injects the given session manager as the global session manager
func SetSessionManager (m ManagerI) {
	sessionManager = m
}

// Session is the object that is stored in the context. (Actually a pointer to that object).
// You normally do not work with the Session object, but instead should call session.GetInt, session.SetInt, etc.
// Session is exported so that it can be created by custom session managers.
type Session struct {
	*types.SafeMap	// container for session data
}

// NewSession creates a new session object for use by session managers. You should not normally need to call this.
// To reset the session data, call Reset()
func NewSession() *Session {
	return &Session{types.NewSafeMap()}
}

// MarshallBinary serializes the session data for storage.
func (s *Session) MarshalBinary() ([]byte, error) {
	return s.SafeMap.MarshalBinary()
}

// UnmarshallBinary unserializes saved session data
func (s *Session) UnmarshalBinary(data []byte) error {
	return s.SafeMap.UnmarshalBinary(data)
}

// Use injects the session manager into the page management process. It adds to the
// context in the Request object so that later session requests can get to the session information, and also
// wraps the given handler in pre and post processing functions. It should be called from your middleware
// processing stack.
func Use (next http.Handler) http.Handler {
	return sessionManager.Use(next)
}

// getSession returns the session object.
func getSession(ctx context.Context) *Session {
	return ctx.Value(sessionContext).(*Session)
}

// Has returns true if the give key exists in the session store
func Has(ctx context.Context, key string) bool {
	return getSession(ctx).Has(key)
}

// GetInt returns the integer at the given key in the session store. typeOk is false if a value exists, but is not an int.
// If no value exists there, or typeOk is false, 0 is returned.
func GetInt(ctx context.Context, key string) (v int, typeOk bool) {
	return getSession(ctx).GetInt(key)
}

// GetBool returns the boolean at the given key in the session store. typeOk is false if a value exists, but is not a bool.
// If no value exists there, or typeOk is false, false is returned.
func GetBool(ctx context.Context, key string) (v bool, typeOk bool) {
	return getSession(ctx).GetBool(key)
}

// GetString returns the string at the given key in the session store. typeOk is false if a value exists, but is not a string.
// If no value exists there, or typeOk is false, an empty string is returned.
func GetString(ctx context.Context, key string) (v string, typeOk bool) {
	return getSession(ctx).GetString(key)
}

// GetString returns the float64 at the given key in the session store. typeOk is false if a value exists, but is not a float64.
// If no value exists there, or typeOk is false, a zero is returned.
func GetFloat(ctx context.Context, key string) (v float64, typeOk bool) {
	return getSession(ctx).GetFloat(key)
}

// Get returns an interface value stored a the given key. nil is returned if nothing is there.
func Get(ctx context.Context, key string) (v interface{}) {
	return getSession(ctx).Get(key)
}

// Set will put the value at the given key in the session store.
func Set(ctx context.Context, key string, v interface{}) {
	getSession(ctx).Set(key, v)
}

// SetInt will put the int value at the given key in the session store.
func SetInt(ctx context.Context, key string, v int) {
	getSession(ctx).Set(key, v)
}

// SetBool will put the bool value at the given key in the session store.
func SetBool(ctx context.Context, key string, v bool) {
	getSession(ctx).Set(key, v)
}

// SetFloat will put the float64 value at the given key in the session store.
func SetFloat(ctx context.Context, key string, v float64) {
	getSession(ctx).Set(key, v)
}

// SetString will put the string value at the given key in the session store.
func SetString(ctx context.Context, key string, v string) {
	getSession(ctx).Set(key, v)
}

// Remove will remove the value at the given key in the session store.
func Remove(ctx context.Context, key string) {
	getSession(ctx).Remove(key)
}

// Clear removes all of the values from the session store. It does not remove the session token itself (the Cookie for
// example). To remove the token also, call Reset() in addition to Clear()
func Clear(ctx context.Context) {
	getSession(ctx).Clear()
}

// Reset will destroy the old session token. If you also call Clear, or don't have any session data after processing the
// request, it will remove the session token all together after the request. If you do have session data, this will
// cause the session data to be moved to a new session token.
func Reset(ctx context.Context) {
	getSession(ctx).Set(sessionResetKey, true)
}










