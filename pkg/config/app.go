package config

import "time"

// HSTSTimeout sets the HSTS timeout length in seconds. See HSTSHandler in app_base.go for info
// Set this to -1 to turn off HSTS
var HSTSTimeout int64 = 86400 // one day

// ReadTimeout specifies the time that a client has to complete sending us its request.
// It helps prevent an attack where the client opens a connection and then sends us data really slowly.
// See go's http package, server.go for details.
var ReadTimeout = 5 * time.Second
// ReadTimeout specifies the time that a client has to complete sending us its headers. See go's http package, server.go for details.
// This can be used to control per request read timeouts. If zero, ReadTimeout is used.
var ReadHeaderTimeout = 0 * time.Second
// WriteTimeout is the amount of time our server will wait for our app to finish writing the response. It helps prevent
// an attack where the server makes a request, but then reads the response very slowly.
var WriteTimeout = 10 * time.Second
// IdleTimout is used during keep-alive connections to control how often the client must ping us to keep the connection alive.
// It helps us detect whether the client has gone away so that we can then close the connection.
var IdleTimeout = 180 * time.Second
// AjaxTimeout is the amount of time in milliseconds that we direct the browser to wait until it determines that an ajax
// call timed out. This would mean that the browser has lost the connection to the server. The goradd.js file put up a
// dialog on the screen telling the user to refresh the page to re-establish the connection. This only happens in release
// mode so that you don't have to worry about timeouts when debugging ajax code.
var AjaxTimeout = 10000;