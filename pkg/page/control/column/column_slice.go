package column

import (
	"context"
	"fmt"
	"github.com/goradd/goradd/pkg/datetime"
	"github.com/goradd/goradd/pkg/page/control"
	"reflect"
	"time"
)

// SliceColumn is a table that works with data that is in the form of a slice. The data item itself must be convertible into
// a string, either by normal string conversion semantics, or using the supplied format string.
type SliceColumn struct {
	control.ColumnBase
}

// NewSliceColumn creates a new column that treats the supplied row data as a slice. It will use the given numeric
// index to get the data. It will then attempt to convert the data into a string, or you can explicitly tell it
// how to do this by calling SetFormat. If the data is a Date, Time or DateTime type, you MUST call SetTimeFormat
// to describe how to format the date or time.
func NewSliceColumn(index int) *SliceColumn {
	i := SliceColumn{}
	i.Init(index)
	return &i
}

func (c *SliceColumn) Init(index int) {
	c.ColumnBase.Init(c)
	c.SetCellTexter(&SliceTexter{Index: index})
}

// SetFormat sets an optional format string for the column, which will be passed to fmt.Sprintf
// to format the data.
func (c *SliceColumn) SetFormat(format string) *SliceColumn {
	c.CellTexter().(*SliceTexter).Format = format
	return c
}

// SetTimeFormat sets the format for Date, Time or DateTime type data. The format will be passed to time.Format
// to produce the text to print for the column.
func (c *SliceColumn) SetTimeFormat(format string) *SliceColumn {
	c.CellTexter().(*SliceTexter).TimeFormat = format
	return c
}

// SliceTexter is the default CellTexter for tables. It lets you get items out of slices.
type SliceTexter struct {
	// Index is the index into the data that corresponds to this table
	Index int
	// Format is a format string. It will be applied using fmt.Sprintf. If you don't provide a Format string, standard
	// string conversion operations will be used.
	Format string
	// TimeFormat is applied to the data using time.Format. You can have both a Format and TimeFormat, and the Format
	// will be applied using fmt.Sprintf after the TimeFormat is applied using time.Format.
	TimeFormat string
}

func (t SliceTexter) CellText(ctx context.Context, col control.ColumnI, rowNum int, colNum int, data interface{}) string {
	vSlice := reflect.ValueOf(data)
	if vSlice.Kind() != reflect.Slice {
		panic("data must be a slice.")
	}
	v := vSlice.Index(t.Index)

	return ApplyFormat(v, t.Format, t.TimeFormat)

}

// ApplyFormat is used by table columns to apply the given fmt.Sprintf and time.Format strings to the data.
// It is exported to allow custom cell texters to use it.
func ApplyFormat(data interface{}, format string, timeFormat string) string {
	var out string

	switch d := data.(type) {
	case int:
		if format == "" {
			out = fmt.Sprintf("%d", d)
		} else {
			out = fmt.Sprintf(format, d)
		}
	case float64:
		if format == "" {
			out = fmt.Sprintf("%f", d)
		} else {
			out = fmt.Sprintf(format, d)
		}
	case float32:
		if format == "" {
			out = fmt.Sprintf("%f", d)
		} else {
			out = fmt.Sprintf(format, d)
		}

	case time.Time:
		if timeFormat == "" {
			panic("Time format is required for time types")
		}
		out = d.Format(timeFormat)

		if format != "" {
			out = fmt.Sprintf(format)
		}

	case datetime.DateTime:
		if timeFormat == "" {
			panic("Time format is required for time types")
		}
		out = d.Format(timeFormat)

		if format != "" {
			out = fmt.Sprintf(format)
		}
	default:
		if format == "" {
			out = fmt.Sprint(d)
		} else {
			out = fmt.Sprintf(format, d)
		}
	}
	return out
}

// SliceColumnCreator creates a column that treats each row as a slice of data.
type SliceColumnCreator struct {
	// ID will assign the given id to the column. If you do not specify it, an id will be given it by the framework.
	ID string
	// Index is the slice index that will be used to get to the data in the column
	Index int
	// Title is the title of the column and will appear in the header
	Title string
	// Format is a format string applied to the data using fmt.Sprintf
	Format string
	// TimeFormat is a format string applied specifically to time data using time.Format
	TimeFormat string
	control.ColumnOptions
}

func (c SliceColumnCreator) Create(ctx context.Context, parent control.TableI) control.ColumnI {
	col := NewSliceColumn(c.Index)
	if c.ID != "" {
		col.SetID(c.ID)
	}
	col.SetTitle(c.Title)
	if c.Format != "" {
		col.SetFormat(c.Format)
	}
	if c.TimeFormat != "" {
		col.SetTimeFormat(c.TimeFormat)
	}
	col.ApplyOptions(ctx, parent, c.ColumnOptions)
	return col
}

