package event

import "github.com/goradd/goradd/pkg/page"

const ClickEvent = "click"

func Click() *page.Event {
	return &page.Event{JsEvent: "click"}
}

func DoubleClick() *page.Event {
	return &page.Event{JsEvent: "dblclick"}
}

// ContextMenu returns an event that responds to a context menu click, which is typically done by right clicking on a two
// mouse button, option-clicking or two-finger clicking on a Mac, or tap and hold on a touch device.
func ContextMenu() *page.Event {
	return &page.Event{JsEvent: "contextmenu"}
}
