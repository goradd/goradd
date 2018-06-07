package table

import (
	"context"
	"reflect"
)

// MapColumn is a table that works with data that is in the form of a slice. The data item itself must be convertable into
// a string, either by normal string conversion symantecs, or using the supplied format string. The format string will
// be applied to a date if the data is a date, or to the string using fmt.Sprintf
type MapColumn struct {
	ColumnBase
}

func NewMapColumn(index interface{}, format ...string) *MapColumn {
	i := MapColumn{}
	var f string
	if len(format) > 0 {
		f = format[0]
	}
	i.Init(index, f, "")
	return &i
}

func NewTimeMapColumn(index interface{}, timeFormat string, format ...string) *MapColumn {
	i := MapColumn{}
	var f string
	if len(format) > 0 {
		f = format[0]
	}
	i.Init(index, f, timeFormat)
	return &i
}

func (c *MapColumn) Init(index interface{}, format string, timeFormat string) {
	c.ColumnBase.Init(c)
	c.SetCellTexter(MapTexter{Index: index, Format: format, TimeFormat: timeFormat})
}

func (c *MapColumn) SetFormat(format string) *MapColumn {
	c.cellTexter.(*MapTexter).Format = format
	return c
}

func (c *MapColumn) SetTimeFormat(format string) *MapColumn {
	c.cellTexter.(*MapTexter).TimeFormat = format
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

func (t MapTexter) CellText(ctx context.Context, col ColumnI, rowNum int, colNum int, data interface{}) string {

	vKey := reflect.ValueOf(t.Index)
	vMap := reflect.ValueOf(data)
	vValue := vMap.MapIndex(vKey)

	if vValue.IsValid() {
		i := vValue.Interface()
		return ApplyFormat(i, t.Format, t.TimeFormat)
	}
	return ""
}
