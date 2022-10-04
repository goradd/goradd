package event

// RowSelected is an event that indicates a row was selected on a SelectTable.
func RowSelected() *Event {
	e := &Event{jsEvent: "rowselected"}
	return e
}
