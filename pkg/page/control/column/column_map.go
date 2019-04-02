package column

import (
	"context"
	"github.com/goradd/goradd/pkg/page/control"
	"reflect"
)

// MapColumn is a table that works with data that is in the form of a map. The data item itself must be convertible into
// a string, either by normal string conversion semantics, or using the supplied format string.
type MapColumn struct {
	control.ColumnBase
}

// NewMapColumn creates a new column that will treat the data as a map. It will get the data from given index in the
// map and then attempt to convert it into a string. Call SetFormat to explicitly tell the column how to convert the
// data into a string using the fmt.Sprintf function. If the data is a Date, Time or DateTime type, you MUST call
// SetTimeFormat to describe the date or time format.
func NewMapColumn(index interface{}) *MapColumn {
	i := MapColumn{}
	i.Init(index)
	return &i
}

func (c *MapColumn) Init(index interface{}) {
	c.ColumnBase.Init(c)
	c.SetCellTexter(&MapTexter{Index: index})
}

// SetFormat sets an optional format string for the column, which will be passed to fmt.Sprintf
// to format the data.
func (c *MapColumn) SetFormat(format string) *MapColumn {
	c.CellTexter().(*MapTexter).Format = format
	return c
}

// SetTimeFormat sets the format for Date, Time or DateTime type data. The format will be passed to time.Format
// to produce the text to print for the column.
func (c *MapColumn) SetTimeFormat(format string) *MapColumn {
	c.CellTexter().(*MapTexter).TimeFormat = format
	return c
}

// MapTexter is the default CellTexter for tables. It lets you get items out of maps.
type MapTexter struct {
	// Index is the index into the data that corresponds to this table
	Index interface{}
	// Format is a format string. It will be applied using fmt.Sprintf. If you don't provide a Format string, standard
	// string conversion operations will be used.
	Format string
	// TimeFormat is applied to the data using time.Format. You can have both a Format and TimeFormat, and the Format
	// will be applied using fmt.Sprintf after the TimeFormat is applied using time.Format.
	TimeFormat string
}

func (t MapTexter) CellText(ctx context.Context, col control.ColumnI, rowNum int, colNum int, data interface{}) string {

	vKey := reflect.ValueOf(t.Index)
	vMap := reflect.ValueOf(data)
	vValue := vMap.MapIndex(vKey)

	if vValue.IsValid() {
		i := vValue.Interface()
		return ApplyFormat(i, t.Format, t.TimeFormat)
	}
	return ""
}
