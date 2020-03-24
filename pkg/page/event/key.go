package event

import "github.com/goradd/goradd/pkg/page"

func KeyDown() *page.Event {
	return &page.Event{JsEvent: "keydown"}
}

func KeyUp() *page.Event {
	return &page.Event{JsEvent: "keyup"}
}

func KeyPress() *page.Event {
	return &page.Event{JsEvent: "keypress"}
}

func Backspace() *page.Event {
	return KeyDown().Condition("event.keyCode == 8")
}

func EnterKey() *page.Event {
	return KeyDown().Condition("event.keyCode == 13")
}

func EscapeKey() *page.Event {
	return KeyDown().Condition("event.keyCode == 27")
}

func TabKey() *page.Event {
	return KeyDown().Condition("event.keyCode == 9")
}

func UpArrow() *page.Event {
	return KeyDown().Condition("event.keyCode == 38")
}

func DownArrow() *page.Event {
	return KeyDown().Condition("event.keyCode == 40")
}
