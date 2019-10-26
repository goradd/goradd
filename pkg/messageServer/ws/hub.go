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
	pagestate string
	channel string
}

type WebSocketHub struct {
	// Channels that clients have subscribed to. Each channel points to a map of pagestates
	channels map[string]map[string]bool

	// Registered clients, keyed by pagestate
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
			log.FrameworkDebugf("New client registering - pagestate: %s", client.pagestate)

			if _,ok := h.clients[client.pagestate]; ok {
				// The same client is registering again. Unregister first.
				h.unregisterClient(client.pagestate)
			}
			h.clients[client.pagestate] = client

		case client := <-h.unregister:
			log.FrameworkDebugf("Client unregistering pagestate: %s", client.pagestate)
			h.unregisterClient(client.pagestate)

		case msg := <-h.send:
			if pagestates, ok := h.channels[msg.Channel]; ok {
				log.FrameworkDebugf("Sending to channel %s - %v", msg.Channel, msg.Message)

				for pagestate := range pagestates {
					if client,ok2 := h.clients[pagestate]; ok2 {
						client.send <- msg
					}
				}
			} else {
				log.Errorf("Could not find channel %s", msg.Channel)
			}

		case sub := <-h.subscribe:
			h.subscribeChannel(sub.pagestate, sub.channel)

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

func (h *WebSocketHub) unregisterClient(pagestate string) {
	var client,_ = h.clients[pagestate]

	if client == nil {
		return
	}

	for channel := range client.channels {
		h.unsubscribeChannel(pagestate, channel)
	}
	delete (h.clients, pagestate)
}

func (h *WebSocketHub) subscribeChannel(pagestate string, channel string) {
	var client,_ = h.clients[pagestate]

	if client == nil {
		return
	}

	client.channels[channel] = true
	if pagestates, ok := h.channels[channel]; !ok {
		h.channels[channel] = make(map[string]bool)
	} else {
		pagestates[pagestate] = true
	}
}

func (h *WebSocketHub) unsubscribeChannel(pagestate string, channel string) {
	if pagestates, ok := h.channels[channel]; ok {
		delete(pagestates, pagestate)
		if len(pagestates) == 0 {
			delete(h.channels, channel)
		} else {
			h.channels[channel] = pagestates
		}
	}
	if client, ok := h.clients[pagestate]; ok {
		delete(client.channels, channel)
	}
}




