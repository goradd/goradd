package event

import (
	"github.com/goradd/goradd/pkg/javascript"
	"github.com/goradd/goradd/pkg/page"
)

// CheckboxColumnClick retuns an event that will detect a click on a checkbox table in a table, and set up the return
// parameters to return:
//  row: the index of the clicked row
//  col: the index of the clicked table
//  checked: the checked state of the checkbox after the click is processed
//  id: the id of the cell clicked
func CheckboxColumnClick() page.EventI {
	e := &page.Event{
		JsEvent: "click",
	}

	m := map[string]interface{}{
		"row":     javascript.JsCode(`g$(event.target).closest("tr").rowIndex`),
		"col":     javascript.JsCode(`g$(event.target).closest("th,td").cellIndex`),
		"checked": javascript.JsCode(`event.target.checked`),
		"id":      `event.target.id`,
	}

	e.ActionValue(m)
	e.Selector(`input[data-gr-checkcol]`)
	return e
}

// CheckboxColumnActionValues can be used to get the values out of the Event.
type CheckboxColumnActionValues struct {
	Row int 		`json:"row"`
	Column int 		`json:"col"`
	Checked bool 	`json:"checked"`
	Id string 		`json:"id"`
}