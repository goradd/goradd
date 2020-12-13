package ws

import (
	"fmt"
	"github.com/goradd/goradd/pkg/config"
	"github.com/goradd/goradd/pkg/goradd"
	"github.com/goradd/goradd/pkg/html"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/sys"
	"log"
	"net/http"
	"path"
	"path/filepath"
)

type WsMessenger struct {
	port int
	tlsPort int
	hub *WebSocketHub
}

func (m *WsMessenger) Start(mux *http.ServeMux, wsPort int, tlsCertFile string, tlsKeyFile string, tlsPort int) {
	m.port = wsPort
	m.tlsPort = tlsPort
	m.makeHub()

	if wsPort != 0 {
		go func() {
			log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", wsPort), mux))
		}()
	}

	if tlsPort != 0 {
		// Here we confirm that the CertFile and KeyFile exist. If they don't, ListenAndServe just exits with an open error
		// and you will not know why.
		if !sys.PathExists(tlsCertFile) {
			log.Fatalf("WebSocketTLSCertFile does not exist: %s", tlsCertFile)
		}

		if !sys.PathExists(tlsKeyFile) {
			log.Fatalf("WebSocketTLSKeyFile does not exist: %s", tlsKeyFile)
		}

		go func() {
			log.Fatal(http.ListenAndServeTLS(fmt.Sprintf(":%d", tlsPort), tlsCertFile, tlsKeyFile, mux))
		}()
	}
}

func (m *WsMessenger) makeHub() {
	m.hub = NewWebSocketHub()
	go m.hub.run()
}

func (m *WsMessenger) JavascriptInit() string {
	return fmt.Sprintf("goradd.initMessagingClient(%d, %d);\n", m.port, m.tlsPort)
}

func (m *WsMessenger) JavascriptFiles() map[string]html.Attributes {
	ret := make (map[string]html.Attributes)
	var p string
	if config.Release {
		// Note that this is going to get the file out of the location we copied it to for deployment.
		// This must be coordinated with that copy operation, which is located by default in the goradd-project/build/makeAssets.sh file
		p = path.Join(config.AssetPrefix, "messenger","js", "goradd-ws.js")
	} else {
		cur := sys.SourceDirectory()
		p = filepath.Join(cur, "assets", "js", "goradd-ws.js")
	}
	ret[p] = nil
	return ret
}


func (m *WsMessenger) Send(channel string, message string) {
	if m.hub != nil {
		m.hub.send <- clientMessage{channel, message}
	}
}

// WebSocketHandler handles web socket requests to send messages to clients.
// It gets the client id from the context in the request. You should intercept
// the request, authorize the client, then insert the client ID into the context of the
// Request
func (m *WsMessenger) WebSocketHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientID := r.Context().Value(goradd.WebSocketContext).(string)
		serveWs(m.hub, w, r, clientID)
	})
}

func WsAssets() string {
	wsAssets := sys.SourceDirectory()
	wsAssets = path.Join(wsAssets, "assets")
	return wsAssets
}

func init() {
	page.RegisterAssetDirectory(WsAssets(), config.AssetPrefix+"messenger")
}
