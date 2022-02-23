package event

import (
	"github.com/goradd/goradd/pkg/javascript"
	"github.com/goradd/goradd/pkg/page"
)

const (
	CellClickDefault     = `{"row": event.goradd.match.parentElement.rowIndex, "col": event.goradd.match.cellIndex}`
	CellClickRowIndex    = `event.goradd.match.parentElement.rowIndex`
	CellClickColumnIndex = `event.goradd.match.cellIndex`
	CellClickCellId      = `event.goradd.match.id`
	CellClickRowId       = `event.goradd.match.parentElement.id`
	CellClickRowValue    = `g$(event.goradd.match).closest("tr").data("value")`
	CellClickColId       = `g$(event.goradd.match).columnId()`
)

// CellClick returns an event to detect clicking on a table cell.
// Lots of things can be determined using this event by changing the return values. When this event fires,
// the javascript environment will have the following local variables defined:
//   - this: The object on to which the event listener was attached.
//   - event: The event object for the click.
//   - event.target - the html object clicked in. If your table cell had other objects in it, this will return the
//                    object clicked inside the cell. This could be important, for example,
//                    if you had form objects inside the cell, and you wanted to behave differently
//                    if a form object was clicked on, verses clicking outside the form object.
//   - event.goradd.match: This will be the cell object clicked on, even if an item inside of the cell was clicked.
//
// Here are some examples of return params you can specify to return data to your action handler:
//   event.goradd.match.id - the cell id
//   event.goradd.match.tagName - the tag for the cell (either th or td)
//   event.goradd.match.cellIndex - the table index that was clicked on, starting on the left with table zero
//   g$(event.goradd.match).data('value') - the "data-value" attribute of the cell (if you specify one). Use this formula for any kind of "data-" attribute.
//   g$(event.goradd.match).parent() - the javascript row object
//   event.goradd.match.parentElement - the html row object
//   event.goradd.match.parentElement.rowIndex - the index of the row clicked, starting with zero at the top (including any header rows).
//   event.goradd.match.parentElement.id the id of the row clicked on
//   g$(event.goradd.match).parent().data("value") - the "data-value" attribute of the row.
//   g$(event.goradd.match).columnId() - the id of the column clicked in
//
// You can put your items in a javascript array, and an array will be returned as the strParameter in the action.
// Or you can put it in a javascript object, and a named array(hash) will be returned.
//
// By default, the cell click does not bubble. Add Bubbles() to the event to get the click
// to bubble up from sub objects.
func CellClick() *page.Event {
	e := page.NewEvent("click").
		Selector("td").
		ActionValue(javascript.JsCode{CellClickDefault})
	return e
}

// HeaderCellClick responds to clicks on header cells (th)
func HeaderCellClick() *page.Event {
	e := page.NewEvent("click").
		Selector("th").
		ActionValue(javascript.JsCode{CellClickDefault})
	return e
}


// RowDataActionValue returns code to use in the ActionValue to to return the data value of the row clicked on.
// The code can be used directly, or in a map or array.
// For example:
//   e := event.CellClick().ActionValue(event.RowDataActionValue("rowVal")).Delay(100)
func RowDataActionValue(key string) javascript.JavaScripter {
	return javascript.JsCode{`g$(this).parent().data("` + key + `")`}
}

// CellDataActionValue sets the ActionValue to javascript that will return the data value of the row clicked on.
// If you are going to use this, call it immediately after you call CellClick, and before any other calls on the event.
// For example:
//   e := event.CellClick().ActionValue(event.CellDataActionValue("cellVal")).Delay(100)
func CellDataActionValue(key string) javascript.JavaScripter {
	return javascript.JsCode{`g$(this).data("` + key + `")`}
}
