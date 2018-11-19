package event

import "github.com/spekary/goradd/pkg/page"

func Click() page.EventI {
	return &page.Event{JsEvent: "click"}
}

func DoubleClick() page.EventI {
	return &page.Event{JsEvent: "dblclick"}
}

// ContextMenu returns an event that responds to a context menu click, which is typically done by right clicking on a two
// mouse button, option-clicking or two-finger clicking on a Mac, or tap and hold on a touch device.
func ContextMenu() page.EventI {
	return &page.Event{JsEvent: "contextmenu"}
}
