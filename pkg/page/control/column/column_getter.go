package column

import (
	"context"
	"github.com/goradd/goradd/pkg/page/control"
)

// GetterColumn is a column that uses the Getter interface to get the text out of columns. The data therefore should be
// a slice of objects that implement the Getter interface.
type GetterColumn struct {
	control.ColumnBase
}

type Getter interface {
	Get(string) interface{}
}

type StringGetter interface {
	Get(string) string
}

// NewGetterColumn creates a new column that will call Get on the column to figure out what data to display.
// If the data is a Date, Time or DateTime type, you MUST also call SetTimeFormat.
// You can also optionally call SetFormat to pass it a fmt.Sprintf string to further format the data before printing.
func NewGetterColumn(index string) *GetterColumn {
	i := GetterColumn{}
	i.Init(index)
	return &i
}

func (c *GetterColumn) Init(index string) {
	c.ColumnBase.Init(c)
	c.SetCellTexter(&GetterTexter{Key: index})
	c.SetTitle(index)
}

// SetFormat sets an optional format string for the column, which will be passed to fmt.Sprintf
// to format the data.
func (c *GetterColumn) SetFormat(format string) *GetterColumn {
	c.CellTexter().(*GetterTexter).Format = format
	return c
}

// SetTimeFormat sets the format for Date, Time or DateTime type data. The format will be passed to time.Format
// to produce the text to print for the column.
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

func (t GetterTexter) CellText(ctx context.Context, col control.ColumnI, rowNum int, colNum int, data interface{}) string {
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
