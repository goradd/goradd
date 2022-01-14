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
// The pattern will be behind the ApiPrefix path.
func RegisterPattern(pattern string, handler http.HandlerFunc) {
	http2.RegisterPathHandler(path.Join(config.ApiPrefix, pattern), http.StripPrefix(config.ApiPrefix, handler))
}

// RegisterAppPattern associates the given URL path with the given handler.
// The pattern will be behind the ApiPrefix path and so will benefit from Session management and the
// rest of the handlers.
func RegisterAppPattern(pattern string, handler http.HandlerFunc) {
	http2.RegisterAppPathHandler(path.Join(config.ApiPrefix, pattern), http.StripPrefix(config.ApiPrefix, handler))
}


