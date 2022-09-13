package api

// This is both an example of how to create handlers for api's, like REST api's, but also is
// an important handler to include if you need access to session variables from your handlers.
//
// To activate this, simply include this file in an import path from your application.
import (
	"fmt"
	"net/http"
	"time"

	"github.com/goradd/goradd/pkg/api"
)

const HelloPath = "/hello"

func init() {
	api.RegisterAppPattern(HelloPath, HelloHandler)
}

// HelloHandler sets up the initial communication with the client, establishing
// a session. It responds with the current unix time so that the client knows
// the current world time from the server, rather than relying on the client's
// current time setting, which might be wrong. The default session handler will
// require that the client supports cookies.
func HelloHandler(w http.ResponseWriter, r *http.Request) {
	dt := time.Now().Unix()
	w.Header().Add("Content-Type", "application/json")
	_, _ = fmt.Fprintf(w, `{"dt":%d}`, dt)

	return
}
