package event

const DialogButtonEvent = "gr-dlgbtn"

// DialogButton returns an event that detects clicking on a dialog's button.
func DialogButton() *Event {
	return NewEvent(DialogButtonEvent)
}

const DialogClosedEvent = "gr-dlgclosed"

// DialogClosed indicates that a dialog has closed. This is a good time to do any required cleanup.
func DialogClosed() *Event {
	return NewEvent(DialogClosedEvent)
}

const TableSortEvent = "gr-sort"

// TableSort is a custom event for responding to a table sort event
func TableSort() *Event {
	return NewEvent(TableSortEvent)
}

// RowSelected is an event that indicates a row was selected on a SelectTable.
func RowSelected() *Event {
	return NewEvent("gr-rowselected")
}
