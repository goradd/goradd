package event

import "github.com/goradd/goradd/pkg/page"

func Blur() page.EventI {
	return &page.Event{JsEvent: "blur"}
}

func Focus() page.EventI {
	return &page.Event{JsEvent: "focus"}
}

func FocusIn() page.EventI {
	return &page.Event{JsEvent: "focusin"}
}

func FocusOut() page.EventI {
	return &page.Event{JsEvent: "focusout"}
}
