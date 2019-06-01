package event

import (
	"github.com/goradd/goradd/pkg/javascript"
	"github.com/goradd/goradd/pkg/page"
)

// rowSelectedEvent indicates that a row was selected from the SelectTable
type rowSelectedEvent struct {
	page.Event
}

// RowSelected
func RowSelected() *rowSelectedEvent {
	e := &rowSelectedEvent{page.Event{JsEvent: "rowselected"}}
	e.ActionValue(javascript.JsCode("ui")) // the data id of the row
	return e
}
