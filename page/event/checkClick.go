package event

import "github.com/spekary/goradd/page"

// CheckboxColumnClick retuns an event that will detect a click on a checkbox column in a table, and set up the return
// parameters to return:
//  row: the index of the clicked row
//  col: the index of the clicked column
//  checked: the checked state of the checkbox after the click is processed
//  id: the id of the cell clicked

func CheckboxColumnClick() page.EventI {
	e := &page.Event{
		JsEvent: "click",
	}

	e.ActionValue(`{"row": $j(this).closest("tr")[0].rowIndex, "col": $j(this).closest("th,td")[0].cellIndex, "checked":this.checked, "id":this.id}`)
	e.Selector(`input[gr-checkcol]`)
	return e
}
