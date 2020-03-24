/*
The messageServer package implements a general purpose messaging platform based on the gorilla
websocket implementation. The platform is similar to other messaging platforms like pubnub,
in that it is channel based.

This means that you can replace it with the distributed messaging platform of your choice as your app grows.

See the web/assets/js/goradd-ws.js for the client side.
*/
package messageServer

import (
	"encoding/json"
	"github.com/goradd/goradd/pkg/html"
	"github.com/goradd/goradd/pkg/log"
)

var Messenger MessengerI

type MessengerI interface {
	JavascriptInit() string
	Send(channel string, message string)
	JavascriptFiles() map[string]html.Attributes
}

func Send(channel string, message interface{}) {
	msg, err := json.Marshal(message)
	if err != nil {
		log.Debug(err)
		return
	}
	if Messenger != nil {
		Messenger.Send(channel, string(msg))
	}
}