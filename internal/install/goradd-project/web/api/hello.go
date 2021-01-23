package rest

import (
	"fmt"
	"github.com/goradd/goradd/pkg/datetime"
	"github.com/goradd/goradd/pkg/api"
	"net/http"
)

const HelloPath = "/hello"

func init() {
	api.RegisterPattern(HelloPath, HelloHandler)
}

// HelloHandler sets up the initial communication with the client, establishing
// a session. It responds with the current unix time so that the client knows
// the current world time from the server, rather than relying on the client's
// current time setting, which might be wrong.
func HelloHandler(w http.ResponseWriter, r *http.Request) {
	dt := datetime.Now().Unix()
	w.Header().Add("Content-Type","application/json")
	_,_ = fmt.Fprintf(w, `{"dt":%d}`, dt)

	return
}