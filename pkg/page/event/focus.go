package event

// Blur returns an event that responds to the javascript "blur" event. This is typically fired when a
// control loses focus, though there are bugs in some browsers that prevent the blur event when
// the control loses focus when the entire browser is placed in the background.
func Blur() *Event {
	return NewEvent("blur")
}

// Focus returns an event that responds to the javascript "focus" event. This event is triggered when a control
// receives the focus.
func Focus() *Event {
	return NewEvent("focus")
}

// FocusIn returns an event that responds to the javascript "focusin" event. This is fired when a control, or
// any of its nested controls, gains focus. In other words, the event bubbles.
func FocusIn() *Event {
	return NewEvent("focusin")
}

// FocusOut returns an event that responds to the javascript "focusout" event. This is fired when a control,
// or any of its nested controls, loses focus. In other words, the event bubbles.
func FocusOut() *Event {
	return NewEvent("focusout")
}
