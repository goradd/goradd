package event

import (
	"github.com/goradd/goradd/pkg/javascript"
	"github.com/goradd/goradd/pkg/page"
)

const (
	CellClickDefault     = `{"row": this.parentElement.rowIndex, "col": this.cellIndex}`
	CellClickRowIndex    = `this.parentElement.rowIndex`
	CellClickColumnIndex = `this.cellIndex`
	CellClickCellId      = `this.id`
	CellClickRowId       = `$j(this).parent().attr("id")`
	CellClickRowValue    = `$j(this).parent().data("value")`
	CellClickColId       = `$j(this).parent().closest("table").find("thead").find("th")[this.cellIndex].id`
)

// CellClickEvent is a cell click event. Do not create it directly, but rather use CellClick()
type CellClickEvent struct {
	page.Event
}

/**
 * CellClick returns an event to detect clicking on a table cell.
 * Lots of things can be determined using this event by changing the Returns values. When this event fires,
 * the javascript environment will have the following local variables defined:
 * - this: The html object for the cell clicked.
 * - event: The event object for the click.
 *
 * Here are some examples of return params you can specify to return data to your action handler:
 * 	this.id - the cell id
 *  this.tagName - the tag for the cell (either th or td)
 *  this.cellIndex - the table index that was clicked on, starting on the left with table zero
 *  $j(this).data('value') - the "data-value" attribute of the cell (if you specify one). Use this formula for any kind of "data-" attribute.
 *  $j(this).parent() - the jQuery row object
 *  this.parentElement - the html row object
 *  this.parentElement.rowIndex - the index of the row clicked, starting with zero at the top (including any header rows).
 *  $j(this).parent().attr('id') or $j(this).parent()[0].id - the id of the row clicked on
 *  $j(this).parent().data("value") - the "data-value" attribute of the row. Use this formula for any kind of "data-" attribute.
 *  $j(this).parent().closest('table').find('thead').find('th')[this.cellIndex].id - the id of the table clicked in
 *  event.target - the html object clicked in. If your table cell had other objects in it, this will return the
 *    object clicked inside the cell. This could be important, for example, if you had form objects inside the cell,
 *    and you wanted to behave differently if a form object was clicked on, verses clicking outside the form object.
 *
 * You can put your items in a javascript array, and an array will be returned as the strParameter in the action.
 * Or you can put it in a javascript object, and a named array(hash) will be returned.
 */
func CellClick() *CellClickEvent {
	e := &CellClickEvent{page.Event{JsEvent: "click"}}
	e.Selector("th,td").ActionValue(javascript.JsCode(CellClickDefault))
	return e
}

// RowDataActionValue sets the ActionValue to javascript that will return the data value of the row clicked on.
// If you are going to use this, call it immediately after you call CellClick, and before any other calls on the event.
// For example:
//   e := CellClick().RowDataActionValue("rowVal").Delay(100)
func (e *CellClickEvent) RowDataActionValue(key string) *CellClickEvent {
	e.ActionValue(javascript.JsCode(`$j(this).parent().data("` + key + `")`))
	return e
}

// CellDataActionValue sets the Returns value to javascript that will return the data value of the row clicked on.
// If you are going to use this, call it immediately after you call CellClick, and before any other calls on the event.
// For example:
//   e := CellClick().CellDataActionValue("cellVal").Delay(100)
func (e *CellClickEvent) CellDataActionValue(key string) *CellClickEvent {
	e.ActionValue(javascript.JsCode(`$j(this).data("` + key + `")`))
	return e
}
