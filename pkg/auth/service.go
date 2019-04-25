package auth

import (
	"context"
	"net/http"
)

var authService AuthI

// AuthI describes the interface for the authorization service you should implement. Note that rate limiting and
// general hacker protection has already been done for you.
type AuthI interface {
	// NewUser should create a new user in your database with the given credentials. If you are using token based
	// authentication, send back the token you have created and saved in your database. You can send it back in the
	// body of your response, or as a header. You should also save an identifier for the user in the session so that
	// the new user is also logged in. Return true if the attempt was successful, false if not.
	// One reason for an unsuccessful attempt might
	// be too short a user name, or an insecure password. Communicate that information in your response to the client.
	NewUser(ctx context.Context, user string, password string, w http.ResponseWriter) bool
	// Login attempts to log in using the given user name and password. Return true if the login attempt was
	// successful, and false if not. If using tokens, return a saved token on successful login. Also, put the user
	// id in the session. If login was not
	// successful, return that information to the client, likely by returning a 401 response code in the header.
	Login(ctx context.Context, user string, password string, w http.ResponseWriter) bool
	// TokenLogin attempts to log in using the given token. This would be a token that you have previously returned
	// in one of the above methods. Put a user id in the session if successful. Return true if successful and false if not.
	TokenLogin(ctx context.Context, token string, w http.ResponseWriter) bool
	// RevokeToken will revoke the given token and will close the session. You should delete the token from the database.
	RevokeToken(ctx context.Context, token string)
	// Recover should give the user information on how to reset the password. We do not immediately reset anything,
	// since this might come from a malicious attack. Only after the user has successfully completed the recovery
	// process should you remove all tokens associated with this user.
	Recover(ctx context.Context, method string)
}

func RegisterAuthenticationService(a AuthI) {
	authService = a
}

