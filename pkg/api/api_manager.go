package api

import (
	"github.com/goradd/goradd/pkg/config"
	http2 "github.com/goradd/goradd/pkg/http"
	"net/http"
)


func init() {
	// Create a default ApiManager if one does not already exist
	// The default can be over-ridden in the project's config package, since those
	// init functions will be called before this one
	if config.ApiManager == nil && config.ApiPrefix != "" {
		config.ApiManager = NewDefaultApiManager()
	}
}

// The DefaultApiManager type represents the default api manager that is created at bootup.
// The Mux variable can be replaced with the Muxer of your choice.
type DefaultApiManager struct {
	Mux http2.Muxer
}

// NewDefaultApiManager creates a new api manager that uses config.ApiPrefix as a path prefix to determine
// if the given call is a REST call. It also uses Go's standard ServeMux as a muxer to
// direct which handler should handle the request. The Mux is public, so you can replace
// it with your own Muxer if desired.
func NewDefaultApiManager() DefaultApiManager {
	return DefaultApiManager{
		Mux: http.NewServeMux(),
	}
}

func (a DefaultApiManager) RegisterPattern(pattern string, handler http.Handler) {
	a.Mux.Handle(pattern, handler)
}

func (a DefaultApiManager) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.Mux.ServeHTTP(w, r)
}

func (a DefaultApiManager) Use() {
	if config.ApiPrefix != "" && config.ApiManager != nil {
		http2.RegisterAppMuxerHandler(config.ApiPrefix, http2.ErrorHandler(a))
	}
}

// RegisterPattern associates the given URL path with the given handler.
// Call this from an init() function. Whenever a user navigates to the given pattern, the
// handler will be called. Note that the pattern should NOT include the ApiPathPrefix.
func RegisterPattern(pattern string, handler http.HandlerFunc) {
	if config.ApiManager == nil {
		panic ("the ApiManager has not been initialized")
	}
	config.ApiManager.RegisterPattern(pattern, handler)
}

