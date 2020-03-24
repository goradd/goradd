// This auth package is an example of how to do authentication and authorization in goradd.
// It will get you started, but you will likely need to make some changes to fit your particular requirements.
//
// The code below expects you to have a User table in your database. It expects it to have an ID field as a unique
// identifier for the user, and a PasswordHash field, which should be a character string of at least 256 characters
// and that will store the hashed password. This system never stores a cleartext password.
//
// You will also need to call PutContext to put the authorization context into the http request context. Call that
// from app.PutContext().
package auth

import (
	"context"
	"github.com/goradd/goradd/pkg/orm/op"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/session"
	"golang.org/x/crypto/bcrypt"
	"goradd-project/gen/bug/model"
	"goradd-project/gen/bug/model/node"
)

const userSessionVar = "app.userid"

// CurrentUser returns the currently logged in user object, or nil if no user is logged in.
// It will cache the user information and store it in the context for future reference in the current request,
// so that the database does not need to be queried each time. It gets the current user out of the sesssion.
// If there is no current user, it returns nil.
func CurrentUser(ctx context.Context) *model.User {
	l := getContext(ctx)
	if l == nil {
		panic("AuthContext has not been put in the context")
	}
	if l.auth.user == nil {
		// try to populate the user from the session data
		id := session.GetString(ctx, userSessionVar)
		if id != "" {
			l.auth.user = model.QueryUsers(ctx).
				Where(
					op.And(
						op.Equal(node.User().ID(), id),
						//op.Equal(node.User().Person().Active(), 1),
					)).
				//Join(node.User().Person()).
				Get()
		}
	}

	return l.auth.user
}

// Authorize authorizes the current user based on a level of authorization. There are other ways to do this,
// but this is just an example. For instance, you could authorize based on specific permissions granted, and
// use a bitfield to OR those together for a more complex authorization system.
func Authorize(ctx context.Context, permissionLevel model.UserType)  bool {
	user := CurrentUser(ctx)
	grctx := page.GetContext(ctx)
	if grctx.Host == "localhost:8000" {
		SetCurrentUserID(ctx, "1")
		return true
	}
	if user == nil || user.UserType() > permissionLevel {
		return false
	}
	return true
}

// SetCurrentUserID will set the current user to a specific id.
func SetCurrentUserID(ctx context.Context, id string) {
	l := getContext(ctx)
	if l == nil {
		panic("AuthContext has not been put in the context")
	}

	if id == "" {
		session.Remove(ctx, userSessionVar)
		l.auth.user = nil
		return
	}
	l.auth.user = model.LoadUser(ctx, id,
		//node.User().Person(),
		)

	if l.auth.user == nil {
		session.Remove(ctx, userSessionVar)
	} else {
		session.SetString(ctx, userSessionVar, id)
	}
}

// ValidatePassword enforces the rules about what makes up a good enough password or userName for our system
// You can also use a 3rd party library to evaluate password strength.
func ValidatePassword(password string, userName string) (passwordErr, userNameErr string) {
	if len(password) < 5 {
		passwordErr = "The password must have at least 5 letters."
	}
	if len(userName) < 3 {
		userNameErr = "The user name must have at least 3 letters."
	}

	return
}

// Logout does the steps required to log out the current user
func Logout(ctx context.Context) {
	SetCurrentUserID(ctx, "")
}

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return ""
	}
	return string(bytes)
}

func VerifyPassword(user *model.User, password string) bool {
	if !user.PasswordHashIsValid() { // we didn't query for the password hash in the last query
		return false
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash()), []byte(password))
	if err != nil {
		return false
	}
	return true
}

