package event

// KeyDown responds to the javascript "keydown" event.
func KeyDown() *Event {
	return NewEvent("keydown")
}

// KeyUp responds to the javascript "keyup" event.
func KeyUp() *Event {
	return NewEvent("keyup")

}

// KeyPress responds to the javascript "keypress" event.
// Deprecated: this is deprecated by the web standards. Use KeyDown or BeforeInput instead.
func KeyPress() *Event {
	return NewEvent("keypress")
}

// BeforeInput responds to the javascript "beforeinput" event.
// This event is fired before a control is changed by text edits.
func BeforeInput() *Event {
	return NewEvent("beforeinput")
}

// Backspace is a keydown event for the backspace key.
func Backspace() *Event {
	return KeyDown().Condition("event.keyCode == 8")
}

// EnterKey is a keydown event for the enter key.
func EnterKey() *Event {
	return KeyDown().Condition("event.keyCode == 13")
}

// EscapeKey is a keydown event for the escape key.
func EscapeKey() *Event {
	return KeyDown().Condition("event.keyCode == 27")
}

// TabKey is a keydown event for the tab key.
func TabKey() *Event {
	return KeyDown().Condition("event.keyCode == 9")
}

// UpArrow is a keydown event for the up arrow.
func UpArrow() *Event {
	return KeyDown().Condition("event.keyCode == 38")
}

// DownArrow is a keydown event for the down arrow.
func DownArrow() *Event {
	return KeyDown().Condition("event.keyCode == 40")
}
