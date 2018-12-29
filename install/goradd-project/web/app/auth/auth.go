package auth

// This is an example auth module that you can include and build out in your app package.
// You will need to provide the Person model and drag the file to the app package above to use it.

import (
	"context"
	"github.com/goradd/goradd/pkg/orm/op"
	"github.com/goradd/goradd/pkg/session"
	"github.com/trustelem/zxcvbn"
	"goradd-project/gen/goradd/model"
	"goradd-project/gen/goradd/model/node"
)

const UserSessionVar = "GoraddUserId"	// Change this to a unique session variable name for your application

// auth is a private structure that caches information on the currently authorized user.
type auth struct {
	user *model.Person
}

// CurrentUser returns the currently logged in user object, or nil if no user is logged in.
// It will cache the user information and store it in the context for future reference in the current request,
// so that the database does not need to be queried each time. It gets the current user out of the sesssion.
// If there is no current user, it returns nil.
func CurrentUser(ctx context.Context) *model.Person {
	l := GetContext(ctx)
	if l == nil {
		panic("LocalContext has not been put in the context")
	}
	if l.auth.user == nil {
		// try to populate the user from the session data
		id, _ := session.GetString(ctx, UserSessionVar)
		if id != "" {
			l.auth.user = model.QueryPeople().
				Where(
					op.And(
						op.Equal(node.Person().ID(), id),
						//op.Equal(node.Person().Active(), 1),
					)).
				Get(ctx)
		}
	}

	return l.auth.user
}

/*
func Authorize(ctx context.Context, permissionLevel model.PermissionType)  {
	u := CurrentUser(ctx)
	grctx := page.GetContext(ctx)
	if grctx.Host == "localhost:8000" {
		//SetCurrentUserID(ctx, "5")
		return
	}
	if u == nil || u.SitePermissionType() > permissionLevel {

			// TODO: Log this attempt at unauthorized access
			page.Redirect("/html/unauthorized.html")
	}
}*/

func SetCurrentUserID(ctx context.Context, id string) {
	l := GetContext(ctx)
	if l == nil {
		panic("LocalContext has not been put in the context")
	}

	if id == "" {
		session.Remove(ctx, UserSessionVar)
		l.auth.user = nil
		return
	}
	l.auth.user = model.LoadPerson(ctx, id)

	if l.auth.user == nil {
		session.Remove(ctx, UserSessionVar)
	} else {
		session.SetString(ctx, UserSessionVar, id)
	}
}

/*
func SetCurrentCompanyID(ctx context.Context, id string) {
	l := GetContext(ctx)
	if l == nil {
		panic("LocalContext has not been put in the context")
	}

	if id == "" {
		session.Remove(ctx, COMPANY_SESSION_VAR)
	} else {
		session.SetString(ctx, COMPANY_SESSION_VAR, id)
	}
}*/



// ValidatePassword enforces the rules about what makes up a good enough password for our system
// We are using a 3rd party library to evaluate password strength. This same library is available in javascript,
// but this is a server side check.
func ValidatePassword(password string, userName string) bool {
	var userInputs = []string{userName}
	r := zxcvbn.PasswordStrength(password, userInputs)
	return r.Score >= 1	// Change this to whatever score you are going to require.
}

// Logout does the steps required to log out the current user
func Logout(ctx context.Context) {
	SetCurrentUserID(ctx, "")
	//SetCurrentCompanyID(ctx, "")
}
