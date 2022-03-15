package ws

import (
	"fmt"
	"github.com/goradd/goradd/pkg/config"
	"github.com/goradd/goradd/pkg/goradd"
	"github.com/goradd/html5tag"
	http2 "github.com/goradd/goradd/pkg/http"
	_ "github.com/goradd/goradd/pkg/messageServer/ws/assets"
	"net/http"
	"path"
)

type WsMessenger struct {
	hub *WebSocketHub
}

func (m *WsMessenger) Start() *WebSocketHub {
	m.hub = NewWebSocketHub()
	go m.hub.run()
	return m.hub
}

func (m *WsMessenger) JavascriptInit() string {
	return fmt.Sprintf("goradd.initMessagingClient(%q);\n", http2.MakeLocalPath(config.WebsocketMessengerPrefix))
}

func (m *WsMessenger) JavascriptFiles() map[string]html5tag.Attributes {
	ret := make(map[string]html5tag.Attributes)
	p := path.Join(config.AssetPrefix, "messenger", "js", "goradd-ws.js")
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
