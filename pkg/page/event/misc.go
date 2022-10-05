package event

func Change() *Event {
	return &Event{jsEvent: "change"}
}

// DragDrop returns an event that responds to the javascript drop event
func DragDrop() *Event {
	return &Event{jsEvent: "drop"}
}

func Input() *Event {
	return &Event{jsEvent: "input"}
}

func Select() *Event {
	return &Event{jsEvent: "select"}
}

// TableSort is a custom event for responding to a table sort event
func TableSort() *Event {
	e := &Event{jsEvent: "grsort"}
	return e
}

const DialogButtonEvent = "gr-dlgbtn"

// DialogButton returns an event that detects clicking on a dialog's button.
func DialogButton() *Event {
	e := &Event{jsEvent: DialogButtonEvent}
	return e
}

const DialogClosedEvent = "grdlgclosed"

// DialogClosed indicates that a dialog has closed. This is a good time to do any required cleanup.
func DialogClosed() *Event {
	e := &Event{jsEvent: DialogClosedEvent}
	return e
}

// TimerExpired is used in conjunction with a JsTimer control to detect the expiration of the timer
func TimerExpired() *Event {
	return &Event{jsEvent: "goradd.timerexpired"}
}
