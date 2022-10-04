package event

import (
	"github.com/goradd/goradd/pkg/page"
)

func Change() *page.Event {
	return &page.Event{JsEvent: "change"}
}

// DragDrop returns an event that responds to the javascript drop event
func DragDrop() *page.Event {
	return &page.Event{JsEvent: "drop"}
}

func Input() *page.Event {
	return &page.Event{JsEvent: "input"}
}

func Select() *page.Event {
	return &page.Event{JsEvent: "select"}
}

// TableSort is a custom event for responding to a table sort event
func TableSort() *page.Event {
	e := &page.Event{JsEvent: "grsort"}
	return e
}

const DialogButtonEvent = "gr-dlgbtn"

// DialogButton returns an event that detects clicking on a dialog's button.
func DialogButton() *page.Event {
	e := &page.Event{JsEvent: DialogButtonEvent}
	return e
}

const DialogClosedEvent = "grdlgclosed"

// DialogClosed indicates that a dialog has closed. This is a good time to do any required cleanup.
func DialogClosed() *page.Event {
	e := &page.Event{JsEvent: DialogClosedEvent}
	return e
}

// TimerExpired is used in conjunction with a JsTimer control to detect the expiration of the timer
func TimerExpired() *page.Event {
	return &page.Event{JsEvent: "goradd.timerexpired"}
}

// Event returns an event that responds to the given named javascript event. Use this to create a response to a custom
// javascript event.
func Event(event string) *page.Event {
	return &page.Event{JsEvent: event}
}
