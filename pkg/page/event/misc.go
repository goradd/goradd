package event

// Change triggers on the javascript change event.
// The change event triggers after a change has been recorded on a control. For text boxes, this occurs
// after focus leaves the text box. Other controls, like select controls, change immediately when a new item
// is selected.
func Change() *Event {
	return NewEvent("change")
}

// DragDrop returns an event that responds to the javascript drop event
func DragDrop() *Event {
	return NewEvent("drop")

}

// Input triggers on the input event. The input event happens when text box type controls have
// been changed at all. This is the event you want to watch if you want to know when a user has typed
// in a text box, or pressed backspace, or cut or pasted into the text box.
func Input() *Event {
	return NewEvent("input")

}

// Select triggers on the select event. The select event happens when text is selected in the control.
func Select() *Event {
	return NewEvent("select")
}
