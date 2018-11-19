package event

import "github.com/spekary/goradd/pkg/page"

func KeyDown() page.EventI {
	return &page.Event{JsEvent: "keydown"}
}

func KeyUp() page.EventI {
	return &page.Event{JsEvent: "keyup"}
}

func KeyPress() page.EventI {
	return &page.Event{JsEvent: "keypress"}
}

func Backspace() page.EventI {
	return KeyDown().Condition("event.keyCode == 8")
}

func EnterKey() page.EventI {
	return KeyDown().Condition("event.keyCode == 13")
}

func EscapeKey() page.EventI {
	return KeyDown().Condition("event.keyCode == 27")
}

func TabKey() page.EventI {
	return KeyDown().Condition("event.keyCode == 9")
}

func UpArrow() page.EventI {
	return KeyDown().Condition("event.keyCode == 38")
}

func DownArrow() page.EventI {
	return KeyDown().Condition("event.keyCode == 40")
}
