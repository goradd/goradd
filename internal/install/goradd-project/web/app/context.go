package app

import (
	"context"
	"github.com/goradd/goradd/pkg/goradd"
)

const LocalContextKey = goradd.ContextKey("app.local")

// LocalContext contains the items you want available to
type LocalContext struct {
	//your local context variables here
}

// PutLocalContext allocates our application specific context object and returns it so we can get to it later
// as the context gets passed around the application.
func PutLocalContext(ctx context.Context) context.Context {
	localContext := &LocalContext{}

	return context.WithValue(ctx, LocalContextKey, localContext)
}

func GetLocalContext(ctx context.Context) *LocalContext {
	return ctx.Value(LocalContextKey).(*LocalContext)
}
