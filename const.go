package goradd

// Application wide constants and types

// ContextKey is the type used for keys of the values we are storing in the context object that is passed around
// the application. Goradd private contexts will start with "goradd.". You are free to use this type to save your
// own information in the context, but do not start your key names with "goradd." to avoid potential key collissions.
type ContextKey string

const (
	PageContext    = ContextKey("goradd.page")
	SessionContext = ContextKey("goradd.session")
	SqlContext     = ContextKey("goradd.sql")
)
