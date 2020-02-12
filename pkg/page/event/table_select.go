package event

import (
	"github.com/goradd/goradd/pkg/page"
)

// RowSelected
func RowSelected() *page.Event {
	e := &page.Event{JsEvent: "rowselected"}
	return e
}
