package event

const ClickEvent = "click"

// Click is an event that responds to the javascript "click" event.
func Click() *Event {
	return &Event{jsEvent: "click"}
}

// DoubleClick is an event that responds to the javascript "dblclick" event.
func DoubleClick() *Event {
	return &Event{jsEvent: "dblclick"}
}

// ContextMenu returns an event that responds to a context menu click, which is typically done by right-clicking on a two
// mouse button, option-clicking or two-finger clicking on a Mac, or tap and hold on a touch device.
func ContextMenu() *Event {
	return &Event{jsEvent: "contextmenu"}
}
