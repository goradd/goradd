package column

import (
	"fmt"
	"context"
	"time"
	"github.com/spekary/goradd/datetime"
	reflect "reflect"
)

// Slice is a column that works with data that is in the form of a slice. The data item itself must be convertable into
// a string, either by normal string conversion symantecs, or using the supplied format string. The format string will
// be applied to a date if the data is a date, or to the string using fmt.Sprintf
type Slice struct {
	ColumnBase
}

func NewSliceColumn(index int) *Slice {
	i := Slice{}
	i.Init(index)
	return &i
}

func (c *Slice) Init(index int) {
	c.ColumnBase.Init(c)
	c.SetCellTexter(SliceTexter{Index: index})
}

// SliceTexter is the default CellTexter for tables. It lets you get items out of slices and maps.
type SliceTexter struct {
	// Index is the index into the data that corresponds to this column
	Index int
	// Format is a format string. It will be applied using fmt.Sprintf. If you don't provide a Format string, standard
	// string conversion operations will be used.
	Format string
	// TimeFormat is applied to the data using time.Format. You can have both a Format and TimeFormat, and the Format
	// will be applied using fmt.Sprintf after the TimeFormat is applied using time.Format.
	TimeFormat string
}

func (t SliceTexter) CellText (ctx context.Context, row int, col int, data interface{}) string {
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