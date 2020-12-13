// Portions copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ws

import (
	"github.com/goradd/goradd/pkg/log"
)

// clientMessage is the information that is passed to the client for each message
type clientMessage struct {
	Channel string `json:"channel"`
	Message string `json:"message"`
}

type subscription struct {
	clientID string
	channel  string
}

type WebSocketHub struct {
	// Channels that clients have subscribed to. Each channel points to a map of client IDs
	channels map[string]map[string]bool

	// Registered clients, keyed by client ID
	clients map[string]*Client

	// Inbound messages from the clients.
	//Broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	send chan clientMessage

	subscribe chan subscription
}

func NewWebSocketHub() *WebSocketHub {
	return &WebSocketHub{
		//Broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		send:       make(chan clientMessage),
		clients:    make(map[string]*Client),
		channels:	make(map[string]map[string]bool),
		subscribe:  make(chan subscription),
	}
}

func (h *WebSocketHub) run() {
	for {
		select {
		case client := <-h.register:
			log.FrameworkDebugf("New client registering - client ID: %s", client.clientID)

			if _,ok := h.clients[client.clientID]; ok {
				// The same client is registering again. Unregister first.
				h.unregisterClient(client.clientID)
			}
			h.clients[client.clientID] = client

		case client := <-h.unregister:
			log.FrameworkDebugf("Client unregistering clientID: %s", client.clientID)
			h.unregisterClient(client.clientID)

		case msg := <-h.send:
			if clientIDs, ok := h.channels[msg.Channel]; ok {
				log.FrameworkDebugf("Sending to channel %s - %v", msg.Channel, msg.Message)

				for clientID := range clientIDs {
					if client,ok2 := h.clients[clientID]; ok2 {
						client.send <- msg
					}
				}
			} else {
				//log.Errorf("Could not find channel %s", msg.Channel)
			}

		case sub := <-h.subscribe:
			log.FrameworkDebugf("Subscribing to channel %s - %v", sub.clientID, sub.channel)
			h.subscribeChannel(sub.clientID, sub.channel)

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

func (h *WebSocketHub) unregisterClient(clientID string) {
	var client,_ = h.clients[clientID]

	if client == nil {
		return
	}

	for channel := range client.channels {
		h.unsubscribeChannel(clientID, channel)
	}
	delete (h.clients, clientID)
}

func (h *WebSocketHub) subscribeChannel(clientID string, channel string) {
	var client,_ = h.clients[clientID]

	if client == nil {
		return
	}

	client.channels[channel] = true
	if clientIDs, ok := h.channels[channel]; !ok {
		clientIDs = make(map[string]bool)
		clientIDs[clientID] = true
		h.channels[channel] = clientIDs
	} else {
		clientIDs[clientID] = true
	}
}

func (h *WebSocketHub) unsubscribeChannel(clientID string, channel string) {
	if clientIDs, ok := h.channels[channel]; ok {
		delete(clientIDs, clientID)
		if len(clientIDs) == 0 {
			delete(h.channels, channel)
		} else {
			h.channels[channel] = clientIDs
		}
	}
	if client, ok := h.clients[clientID]; ok {
		delete(client.channels, channel)
	}
}




