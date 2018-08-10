package app

import (
	"context"
	"github.com/spekary/goradd"
	"github.com/spekary/goradd/page"
	"os"
)

const LocalContextKey = goradd.ContextKey("app.local")

// Change this
type LocalContext struct {
	//userId int	// for example
}

type Context struct {
	page.Context // standard context stuff goradd requires

	// Your per-request stuff, like per-request instances or dependency injections
	LocalContext
}

// Our application was called from the command line
// Here we populate the things of the context that we can know about
func NewCliContext() (ctx *Context) {
	ctx = &Context{}
	ctx.Context.FillApp(os.Args[1:])

	// Do your custom context populations here

	return ctx
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
