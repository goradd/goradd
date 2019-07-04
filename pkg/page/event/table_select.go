package event

import (
	"github.com/goradd/goradd/pkg/javascript"
	"github.com/goradd/goradd/pkg/page"
)

// RowSelected
func RowSelected() page.EventI {
	e := &page.Event{JsEvent: "rowselected"}
	e.ActionValue(javascript.JsCode("ui")) // the data id of the row
	return e
}
