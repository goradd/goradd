// Package api gives you a default pattern for registering an API, like a REST API, to be served
// by the app.
//
// Each file that serves an endpoint of the API should call RegisterPattern or
// RegisterAppPattern from an init() function in that file. You then simply include that file's
// package in an import from your app, and everything will get tied together.
//
// This is a pretty simple handler. Likely you will want to add your own middleware on top of
// the handlers, and you can easily do that by copying this file to your own app and have your
// api files register with that copy. You can then add handlers how you want.
package api

import (
	"github.com/goradd/goradd/pkg/config"
	http2 "github.com/goradd/goradd/pkg/http"
	"net/http"
	"path"
)

// RegisterPattern associates the given URL path with the given handler.
// The pattern will be behind the ApiPrefix path. The handler will be processed before
// the regular app management, and so will not get access to session management. Use this
// if you are using an authorization scheme like oAuth.
//
// The pattern is expected to be the beginning of a url that will end in slash.
// Both the slash ending pattern and non-slash ending pattern will be registered.
// Your handler should be able to handle "" and "/" paths as equivalent.
//
// The handler will have the pattern stripped out.
func RegisterPattern(pattern string, handler http.HandlerFunc) {
	l := len(pattern)
	if l > 0 && pattern[l-1] == '/' {
		l -= 1
	}
	p := path.Join(config.ApiPrefix, pattern[:l])
	http2.RegisterPrefixHandler(p, handler)

	// For speed, register the same handler without a trailing slash.
	http2.RegisterHandler(p, http.StripPrefix(p, handler))
}

// RegisterAppPattern associates the given URL path with the given handler.
// The handler will be behind the App handler and so will benefit from Session management and the
// rest of the handlers.
//
// The pattern is expected to be the beginning of a url that will end in slash.
// Both the slash ending pattern and non-slash ending pattern will be registered.
// Your handler should be able to handle "" and "/" paths as equivalent.
//
// The handler will have the pattern stripped out.
func RegisterAppPattern(pattern string, handler http.HandlerFunc) {
	l := len(pattern)
	if l > 0 && pattern[l-1] == '/' {
		l -= 1
	}
	p := path.Join(config.ApiPrefix, pattern[:l])
	http2.RegisterAppPrefixHandler(p, handler)

	// For speed, register the same handler without a trailing slash.
	http2.RegisterAppHandler(p, http.StripPrefix(p, handler))
}
