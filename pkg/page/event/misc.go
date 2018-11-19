package event

import (
	"github.com/spekary/goradd/pkg/javascript"
	"github.com/spekary/goradd/pkg/page"
)

func Change() page.EventI {
	return &page.Event{JsEvent: "change"}
}

// DragDrop returns an event that responds to the javascript drop event
func DragDrop() page.EventI {
	return &page.Event{JsEvent: "drop"}
}

func Input() page.EventI {
	return &page.Event{JsEvent: "input"}
}

func Select() page.EventI {
	return &page.Event{JsEvent: "select"}
}

// TableSort is a custom event for responding to a table sort event
func TableSort() page.EventI {
	e := &page.Event{JsEvent: "grsort"}
	e.ActionValue(javascript.JsCode("ui")) // this will be the column id
	return e
}

const DialogButtonEvent = "grdlgbtn"

// DialogButton returns an event that detects clicking on a dialog's button.
func DialogButton() page.EventI {
	e := &page.Event{JsEvent: DialogButtonEvent}
	e.ActionValue(javascript.JsCode("ui"))
	return e
}

const DialogClosingEvent = "grdlgclosing"

// DialogClosing indicates that a dialog is about to close. This is a good time to gather up any information that
// you might need.
func DialogClosing() page.EventI {
	e := &page.Event{JsEvent: DialogClosingEvent}
	return e
}

const DialogClosedEvent = "grdlgclosed"

// DialogClosed indicates that a dialog has closed. This is a good time to do any required cleanup.
func DialogClosed() page.EventI {
	e := &page.Event{JsEvent: DialogClosedEvent}
	return e
}


// TimerExpired is used in conjunction with a JsTimer control to detect the expiration of the timer
func TimerExpired() page.EventI {
	return &page.Event{JsEvent: "goradd.timerexpired"}
}

// Custom returns an event that responds to the given javascript event
func Event(event string) page.EventI {
	return &page.Event{JsEvent: event}
}
