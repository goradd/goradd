package column

import (
	"context"
	"goradd-project/override/control_base"
	"github.com/spekary/goradd/pkg/page/control/control_base/table"
)

// GetterColumn is a column that uses the Getter interface to get the text out of columns. The data therefore should be
// a slice of objects that implement the Getter interface.
type GetterColumn struct {
	control_base.ColumnBase
}

type Getter interface {
	Get(string) interface{}
}

type StringGetter interface {
	Get(string) string
}

func NewGetterColumn(index string, format ...string) *GetterColumn {
	i := GetterColumn{}
	var f string
	if len(format) > 0 {
		f = format[0]
	}
	i.Init(index, f, "")
	return &i
}

func NewDateGetterColumn(index string, timeFormat string, format ...string) *GetterColumn {
	i := GetterColumn{}
	var f string
	if len(format) > 0 {
		f = format[0]
	}
	i.Init(index, f, timeFormat)
	return &i
}

func (c *GetterColumn) Init(index string, format string, timeFormat string) {
	c.ColumnBase.Init(c)
	c.SetCellTexter(GetterTexter{Key: index, Format: format, TimeFormat: timeFormat})
	c.SetTitle(index)
}

func (c *GetterColumn) SetFormat(format string) *GetterColumn {
	c.CellTexter().(*GetterTexter).Format = format
	return c
}

func (c *GetterColumn) SetTimeFormat(format string) *GetterColumn {
	c.CellTexter().(*GetterTexter).TimeFormat = format
	return c
}

// GetterTexter lets you get items out of map like objects using the Getter interface.
type GetterTexter struct {
	// Key is the key to use when calling the Get function on the object.
	Key string
	// Format is a format string. It will be applied using fmt.Sprintf. If you don't provide a Format string, standard
	// string conversion operations will be used.
	Format string
	// TimeFormat is applied to the data using time.Format. You can have both a Format and TimeFormat, and the Format
	// will be applied using fmt.Sprintf after the TimeFormat is applied using time.Format.
	TimeFormat string
}

func (t GetterTexter) CellText(ctx context.Context, col table.ColumnI, rowNum int, colNum int, data interface{}) string {
	switch v := data.(type) {
	case Getter:
		d := v.Get(t.Key)
		return ApplyFormat(d, t.Format, t.TimeFormat)
	case StringGetter:
		d := v.Get(t.Key)
		return ApplyFormat(d, t.Format, t.TimeFormat)
	}
	return ""
}
