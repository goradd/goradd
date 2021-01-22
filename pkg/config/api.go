package config

import "net/http"

// ApiPrefix is the url prefix that indicates this is an API call, like a REST or GraphQL call.
// The default is blank and will not serve an API.
// Override this in the goradd_project/config/goradd.go file to turn on the default API handler.
var ApiPrefix = ""

// The ApiManager is a singleton global that manages the deployment of api handlers,
// like REST or GraphQL handlers.
//
// If you specify a value for ApiPrefix, a default ApiManager will be created that uses the ApiPrefix
// to indicate that a route is for an API handler. You can inject your own ApiManager at any time,
// but it should be set up before calling RegisterPattern.
var ApiManager ApiManagerI

type ApiManagerI interface {
	// HandleRequest will route the request to an API handler. It will also
	// detect that the request was, in fact, an API request and return true
	// if so, and false if the request was not an API request.
	HandleRequest(w http.ResponseWriter, r *http.Request) bool

	// RegisterPattern associates the given URL path with the given handler.
	// Whenever a client navigates to the given pattern, the
	// handler will be called. Note that the pattern should NOT include the ApiPathPrefix.
	// If your ApiManager does not need this, just stub it to do nothing.
	RegisterPattern(pattern string, handler http.Handler)
}
