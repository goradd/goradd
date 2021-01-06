// Portions copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ws

import (
	"bytes"
	"encoding/json"
	log2 "github.com/goradd/goradd/pkg/log"
	"github.com/goradd/goradd/pkg/messageServer"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *WebSocketHub

	// The websocket connection
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan clientMessage

	channels map[string]bool

	clientID string // authenticator
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
//
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(c.hub.MaxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(c.hub.PongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(c.hub.PongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		c.handleMessage(message)
		c.conn.SetReadDeadline(time.Now().Add(c.hub.PongWait)) // extend pong
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(c.hub.PingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(c.hub.WriteWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// messages are json objects. We gather them up here into an array.
			var messages = []clientMessage{message}

			// Add queued messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				messages = append(messages, <-c.send)
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			if buf, err2 := json.Marshal(messages); err2 != nil {
				panic(err2)
			} else {
				w.Write(buf)
				log2.FrameworkDebugf("Writepump %s", string(buf))
			}

			if err = w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(c.hub.WriteWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// serveWs handles new websocket requests from clients.
func serveWs(hub *WebSocketHub, w http.ResponseWriter, r *http.Request, clientID string) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log2.Error(err)
		return
	}

	client := &Client{hub: hub, conn: conn, send: make(chan clientMessage, 256), channels:make(map[string]bool), clientID: clientID}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}


type inMessage struct {
	// Subscribe indicates subscribing to a channel
	Subscribe []string `json:"subscribe"`
	// Providing a channel will imply you are sending a message to the channel
	Channel string `json:"channel"`
	Message interface{} `json:"message"`
}

func (c *Client) handleMessage(data []byte) {
	var msg inMessage
	_ = json.Unmarshal(data, &msg)

	if msg.Subscribe != nil {
		for _,channel := range msg.Subscribe {
			s := subscription{
				clientID: c.clientID,
				channel:  channel,
			}
			c.hub.subscribe <- s
		}
	}
	if msg.Channel != "" {
		messageServer.Send(msg.Channel, msg.Message)
	}

}
