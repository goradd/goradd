package event

import (
	"github.com/spekary/goradd/page"
	"github.com/spekary/goradd/javascript"
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

// DataGridSort is a custom event for responding to a data grid sort event
func DataGridSort() page.EventI {
	e := &page.Event{JsEvent: "goradd.datagrid.sort"}
	e.ActionValue(javascript.JsCode("ui"))
	return e
}

// DialogButton returns an event that detects clicking on a dialog's button.
func DialogButton() page.EventI {
	e := &page.Event{JsEvent: "goradd.datagrid.dialogbtn"}
	e.ActionValue(javascript.JsCode("ui"))
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
