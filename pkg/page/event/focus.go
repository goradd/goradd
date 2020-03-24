package event

import "github.com/goradd/goradd/pkg/page"

func Blur() *page.Event {
	return &page.Event{JsEvent: "blur"}
}

func Focus() *page.Event {
	return &page.Event{JsEvent: "focus"}
}

func FocusIn() *page.Event {
	return &page.Event{JsEvent: "focusin"}
}

func FocusOut() *page.Event {
	return &page.Event{JsEvent: "focusout"}
}
