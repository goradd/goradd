// Portions copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package messageServer

import (
	"github.com/goradd/goradd/pkg/log"
	"net/http"
	"strings"
)

// Hub maintains the set of active clients and broadcasts messages to the clients.
var hub *WebSocketHub

// helper type to synchronize the clients map using a channel
type messageType struct {
	channel string
	message map[string]interface{} // messages will be converted to json objects. Items in the object must be json serializable.
}

type WebSocketHub struct {
	// Registered clients, first by channel, and then by pagestate.
	// This means a form or objects on a form can subscribe to multiple channels, but more
	// then one object on the same form cannot subscribe to the same channel.
	clients map[string]map[string]*Client

	// Inbound messages from the clients.
	//Broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	send chan messageType
}

func NewWebSocketHub() *WebSocketHub {
	return &WebSocketHub{
		//Broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		send:       make(chan messageType),
		clients:    make(map[string]map[string]*Client),
	}
}

func (h *WebSocketHub) run() {
	for {
		select {
		case client := <-h.register:
			log.FrameworkDebugf("Client registering for channels %s pagestate %s", client.channels, client.pagestate)
			var clientsByFormstate map[string]*Client
			var ok bool

			channels := strings.Split(client.channels, ",")
			for _,channel := range channels {
				if clientsByFormstate, ok = h.clients[channel]; !ok {
					clientsByFormstate = make(map[string]*Client)
					h.clients[channel] = clientsByFormstate
				}
				if _, ok := clientsByFormstate[client.pagestate]; !ok {
					clientsByFormstate[client.pagestate] = client
				} else {
					// The page is registering again for a particular channel. Maybe a page refresh? Close the previous channel to prevent a memory leak
					h.unregisterClient(channel, client.pagestate)
					clientsByFormstate[client.pagestate] = client
				}
			}

		case client := <-h.unregister:
			log.FrameworkDebugf("Client unregistering for channels %s pagestate %s", client.channels, client.pagestate)
			h.unregisterClient(client.channels, client.pagestate)

		case msg := <-h.send:
			if clients, ok := h.clients[msg.channel]; ok {
				log.FrameworkDebugf("Sending to channel %s - %v", msg.channel, msg.message)
				for _, client := range clients {
					client.send <- msg.message
				}
			} else {
				log.Errorf("Could not find channel %s", msg.channel)
			}

			/* not broadcasting currently. This might change
			case message := <-h.Broadcast:
				for client := range h.clients {
					select {
					case client.send <- message: // echo
					default:
						close(client.send)
						delete(h.clients, client)
					}
				}
			*/
		}
	}
}

func (h *WebSocketHub) unregisterClient(channels string, pagestate string) {
	sChannels := strings.Split(channels, ",")
	for _,channel := range sChannels {
		if clientsByFormstate, ok := h.clients[channel]; ok {
			if client, ok2 := clientsByFormstate[pagestate]; ok2 {
				close(client.send)
				delete(clientsByFormstate, pagestate)
				if len(clientsByFormstate) == 0 {
					delete(h.clients, channel)
				}
			}
		}
	}
}

func HasChannel(channel string) bool {
	_, ok := hub.clients[channel]
	return ok
}

func SendMessage(channel string, message map[string]interface{}) {
	if message == nil {
		message = make(map[string]interface{})
	}
	message["channel"] = channel // send the channel we are sending to with the message
	hub.send <- messageType{channel, message}
}

func MakeHub() {
	hub = NewWebSocketHub()
	go hub.run()
}

func WebsocketHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
}
