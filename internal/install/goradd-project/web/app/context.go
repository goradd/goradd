package app

import (
	"context"
	"github.com/goradd/goradd/pkg/goradd"
	"github.com/goradd/goradd/pkg/page"
)

const LocalContextKey = goradd.ContextKey("app.local")

// LocalContext contains the items you want available to
type LocalContext struct {
	//auth auth
}

type Context struct {
	page.Context // standard context stuff goradd requires

	// Your per-request stuff, like per-request instances or dependency injections
	LocalContext
}

// PutContext allocates our application specific context object and returns it so we can get to it later
// as the context gets passed around the application.
func PutContext(ctx context.Context) context.Context {
	localContext := &LocalContext{}

	return context.WithValue(ctx, LocalContextKey, localContext)
}

func GetContext(ctx context.Context) *LocalContext {
	return ctx.Value(LocalContextKey).(*LocalContext)
}
