package event

import (
	"github.com/goradd/goradd/pkg/page/event"
)

// MessengerReady returns an event to indicate that the messenger is ready. Messenger implementations
// should send this event to the GoRADD form once the messenger service has been initialized and is ready to
// receive messages.
func MessengerReady() *event.Event {
	return event.NewEvent("messengerReady")
}
