package event

import (
	"github.com/spekary/goradd/javascript"
	"github.com/spekary/goradd/page"
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
		"row":     javascript.JsCode(`$j(this).closest("tr")[0].rowIndex`),
		"col":     javascript.JsCode(`$j(this).closest("th,td")[0].cellIndex`),
		"checked": javascript.JsCode(`this.checked`),
		"id":      `this.id`,
	}

	e.ActionValue(m)
	e.Selector(`input[data-gr-checkcol]`)
	return e
}
