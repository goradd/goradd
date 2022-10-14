package goradd

// Application wide constants and types

// ContextKey is the type used for keys of the values we are storing in the context object that is passed around
// the application. Goradd private contexts will start with "goradd.". You are free to use this type to save your
// own information in the context, but do not start your key names with "goradd." to avoid potential key collisions.
type ContextKey string

const (
	PageContext      = ContextKey("goradd.page")
	BufferContext    = ContextKey("goradd.buf")
	WebSocketContext = ContextKey("goradd.ws")
)

// Default session values
const (
	SessionLanguage    = "goradd.lang"
	SessionSalt        = "goradd.salt"
	SessionCsrf        = "goradd.csrf"
	SessionAuthTime    = "goradd.auth.time"
	SessionAuthSuccess = "goradd.auth.success"
)
