package auth

import (
	"context"
	"github.com/goradd/goradd/pkg/goradd"
	"goradd-project/gen/bug/model"
)

const AuthContextKey = goradd.ContextKey("app.auth")

// auth is a private structure that caches information on the currently authorized user.
type authContext struct {
	user *model.User
}

type Context struct {
	auth authContext
}

// PutContext allocates our application specific context object and returns it so we can get to it later
// as the context gets passed around the application. Call this from app.PutContext().
func PutContext(ctx context.Context) context.Context {
	authContext := &Context{}

	return context.WithValue(ctx, AuthContextKey, authContext)
}

func getContext(ctx context.Context) *Context {
	c := ctx.Value(AuthContextKey)
	if c == nil {
		return nil
	}
	return c.(*Context)
}

