package column

import (
	"context"
	"fmt"
	"github.com/goradd/goradd/pkg/datetime"
	"github.com/goradd/goradd/pkg/page/control"
	"reflect"
	"time"
)

// SliceColumn is a table that works with data that is in the form of a slice. The data item itself must be convertable into
// a string, either by normal string conversion symantecs, or using the supplied format string.
type SliceColumn struct {
	control.ColumnBase
}

func NewSliceColumn(index int, format ...string) *SliceColumn {
	i := SliceColumn{}
	var f string
	if len(format) > 0 {
		f = format[0]
	}
	i.Init(index, f, "")
	return &i
}

func NewTimeSliceColumn(index int, timeFormat string, format ...string) *SliceColumn {
	i := SliceColumn{}
	var f string
	if len(format) > 0 {
		f = format[0]
	}
	i.Init(index, f, timeFormat)
	return &i
}

func (c *SliceColumn) Init(index int, format string, timeFormat string) {
	c.ColumnBase.Init(c)
	c.SetCellTexter(SliceTexter{Index: index, Format: format, TimeFormat: timeFormat})
}

func (c *SliceColumn) SetFormat(format string) *SliceColumn {
	c.CellTexter().(*SliceTexter).Format = format
	return c
}

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

func (t SliceTexter) CellText(ctx context.Context, col 	control.ColumnI, rowNum int, colNum int, data interface{}) string {
	vSlice := reflect.ValueOf(data)
	if vSlice.Kind() != reflect.Slice {
		panic("data must be a slice.")
	}
	v := vSlice.Index(t.Index)

	return ApplyFormat(v, t.Format, t.TimeFormat)

}

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
