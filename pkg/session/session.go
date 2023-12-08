package session

import (
	"context"
	"encoding/gob"
	http2 "github.com/goradd/goradd/pkg/http"
	"github.com/goradd/maps"
	"net/http"
)

type sessionContextType string

const sessionContext sessionContextType = "goradd.session"
const sessionResetKey string = "goradd.reset"
const timezoneKey string = "goradd.timezone"

type sessionData = maps.SafeMap[string, interface{}]

var sessionManager ManagerI

// ManagerI is the interface for session managers.
type ManagerI interface {
	http2.User
}

// SetSessionManager injects the given session manager as the global session manager
func SetSessionManager(m ManagerI) {
	sessionManager = m
}

func SessionManager() ManagerI {
	return sessionManager
}

// Session is the object that is stored in the context. (Actually a pointer to that object).
// You normally do not work with the Session object, but instead should call session.GetInt, session.SetInt, etc.
// Session is exported so that it can be created by custom session managers.
type Session struct {
	data *sessionData // container for session data
}

// NewSession creates a new session object for use by session managers. You should not normally need to call this.
// To reset the session data, call Reset()
func NewSession() *Session {
	return &Session{new(sessionData)}
}

// MarshalBinary serializes the session data for storage.
func (s *Session) MarshalBinary() ([]byte, error) {
	return s.data.MarshalBinary()
}

// UnmarshalBinary unserializes saved session data
func (s *Session) UnmarshalBinary(data []byte) error {
	if s.data == nil {
		s.data = new(sessionData)
	}
	return s.data.UnmarshalBinary(data)
}

// Use injects the session manager into the page management process. It adds to the
// context in the Request object so that later session requests can get to the session information, and also
// wraps the given handler in pre and post processing functions. It should be called from your middleware
// processing stack.
func Use(next http.Handler) http.Handler {
	return sessionManager.Use(next)
}

// getSession returns the session object.
func getSession(ctx context.Context) *Session {
	return ctx.Value(sessionContext).(*Session)
}

// HasSession returns true if the session exists.
func HasSession(ctx context.Context) bool {
	return ctx.Value(sessionContext) != nil
}

// Has returns true if the given key exists in the session store
func Has(ctx context.Context, key string) bool {
	return getSession(ctx).data.Has(key)
}

// GetInt returns the integer at the given key in the session store. If the key does not exist
// OR if what does exists at that key is not an integer, zero will be returned.
// Call Has() if the zero value has meaning for you and you want to check for existence.
func GetInt(ctx context.Context, key string) (v int) {
	i, ok := getSession(ctx).data.Load(key)
	if ok {
		v, _ = i.(int)
	}
	return
}

// GetBool returns the boolean at the given key in the session store.
// If the key does not exist OR if what does exist at that key is not a boolean,
// false will be returned.
func GetBool(ctx context.Context, key string) (v bool) {
	i, ok := getSession(ctx).data.Load(key)
	if ok {
		v, _ = i.(bool)
	}
	return
}

// GetString returns the string at the given key in the session store.
// If the key does not exist OR if what does exist at that key is not a string,
// an empty string will be returned.
func GetString(ctx context.Context, key string) (v string) {
	i, ok := getSession(ctx).data.Load(key)
	if ok {
		v, _ = i.(string)
	}
	return
}

// GetFloat64 returns the float64 at the given key in the session store.
// If the key does not exist OR if what does exist at that key is not a float,
// 0 will be returned.
func GetFloat64(ctx context.Context, key string) (v float64) {
	i, ok := getSession(ctx).data.Load(key)
	if ok {
		v, _ = i.(float64)
	}
	return
}

// GetFloat32 returns the float32 at the given key in the session store.
// If the key does not exist OR if what does exist at that key is not a float,
// 0 will be returned.
func GetFloat32(ctx context.Context, key string) (v float32) {
	i, ok := getSession(ctx).data.Load(key)
	if ok {
		v, _ = i.(float32)
	}
	return
}

// Get returns an interface value stored at the given key. nil is returned if nothing is there.
func Get(ctx context.Context, key string) (v interface{}) {
	return getSession(ctx).data.Get(key)
}

// Set will put the value at the given key in the session store.
// v must be serializable
func Set(ctx context.Context, key string, v interface{}) {
	getSession(ctx).data.Set(key, v)
}

// SetInt will put the int value at the given key in the session store.
func SetInt(ctx context.Context, key string, v int) {
	getSession(ctx).data.Set(key, v)
}

// SetBool will put the bool value at the given key in the session store.
func SetBool(ctx context.Context, key string, v bool) {
	getSession(ctx).data.Set(key, v)
}

// SetFloat64 will put the float64 value at the given key in the session store.
func SetFloat64(ctx context.Context, key string, v float64) {
	getSession(ctx).data.Set(key, v)
}

// SetFloat32 will put the float32 value at the given key in the session store.
func SetFloat32(ctx context.Context, key string, v float32) {
	getSession(ctx).data.Set(key, v)
}

// SetString will put the string value at the given key in the session store.
func SetString(ctx context.Context, key string, v string) {
	getSession(ctx).data.Set(key, v)
}

// Remove will remove the value at the given key in the session store.
func Remove(ctx context.Context, key string) {
	getSession(ctx).data.Delete(key)
}

// Clear removes all of the values from the session store. It does not remove the session token itself (the Cookie for
// example). To remove the token also, call Reset() in addition to Clear()
func Clear(ctx context.Context) {
	getSession(ctx).data.Clear()
}

// Reset will destroy the old session token. If you also call Clear, or don't have any session data after processing the
// request, it will remove the session token all together after the request. If you do have session data, this will
// cause the session data to be moved to a new session token.
func Reset(ctx context.Context) {
	getSession(ctx).data.Set(sessionResetKey, true)
}

func init() {
	gob.Register(&Session{})
}
