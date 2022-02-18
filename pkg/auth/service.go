package auth

import (
	"context"
	"net/http"
)

var authService AuthI

// AuthI describes the interface for the authorization service you should implement. Note that rate limiting and
// general hacker protection has already been done for you. The msg parameter below is the decoded json sent
// by the client, and that you can use for your own needs.
type AuthI interface {
	// NewUser should create a new user in your database with credentials specified in the message.
	// The message comes from the "msg" form value sent by the client, so the client should put
	// whatever credentials are needed in that form value, and then you can extract the credentials on the server
	// end from the message.
	// If you are doing token based authentication, send back the token you have created and saved in your
	// database as well as a hash of the user's id. You can send it back in the
	// body of your response, or as a header. You should also save an identifier for the user in the session so that
	// the new user is also logged in. Return true if the attempt was successful, false if not.
	// One reason for an unsuccessful attempt might
	// be too short a user name, or an insecure password. Communicate that information in your response to the client.
	NewUser(ctx context.Context, message []byte, w http.ResponseWriter) bool
	// Login attempts to log in using the credentials in the message.
	// The message comes from the "msg" form value sent by the client, so the client should put
	// whatever credentials are needed in that form value, and then you can extract the credentials on the server
	// end from the message.
	// Return true if the login attempt was successful, and false if not.
	// If using tokens, return a saved token and a hash of your user id if login was successful.
	// Also, put the user id in the session. If login was not successful, return false, but also
	// write the reason to the writer so the client can know what happened, and also return an error code,
	// likely by returning a 401 response code in the header.
	Login(ctx context.Context, msg []byte, w http.ResponseWriter) bool
	// TokenLogin should attempt to log in using a token and id hash in the message. This would be a token that you have
	// previously returned in one of the above methods. Puts a user id in the session if successful.
	// Returns true if successful and false if not.
	TokenLogin(ctx context.Context, msg []byte, w http.ResponseWriter) bool
	// RevokeToken should revoke the token in the message. The session will be closed by the service as well.
	// You should delete the token from the database.
	RevokeToken(ctx context.Context, msg []byte, w http.ResponseWriter) bool
	// Recover should give the user information on how to reset the password. We do not immediately reset anything,
	// since this might come from a malicious attack. Only after the user has successfully completed the recovery
	// process should you remove all tokens associated with this user. You could potentially include in the msg
	// instructions on what recovery method to use. Write back instructions for the user on what to do next.
	Recover(ctx context.Context, msg []byte, w http.ResponseWriter) bool
	// WriteError should write the errorMessage to the response writer, together with the code. You should format
	// the message in whatever way your client expects to receive it.
	WriteError(ctx context.Context, errorMessage string, httpCode int, w http.ResponseWriter)
}

func RegisterAuthenticationService(a AuthI) {
	authService = a
}
