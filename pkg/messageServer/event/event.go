package event

import "github.com/goradd/goradd/pkg/page"

// Ready returns an event to indicate that the messenger is ready. Messenger implementations
// should send this event to the goradd form once the messenger service has been initialized and is ready to
// receive messages.
func MessengerReady() *page.Event {
	return &page.Event{JsEvent: "messengerReady"}
}
