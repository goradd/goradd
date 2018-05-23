package table

import (
	"context"
)

// MapColumn is a table that works with data that is in the form of a slice. The data item itself must be convertable into
// a string, either by normal string conversion symantecs, or using the supplied format string. The format string will
// be applied to a date if the data is a date, or to the string using fmt.Sprintf
type GetterColumn struct {
	ColumnBase
}

type Getter interface {
	Get(string) interface{}
}

type StringGetter interface {
	Get(string) string
}


func NewGetterColumn(index string) *GetterColumn {
	i := GetterColumn{}
	i.Init(index)
	return &i
}

func (c *GetterColumn) Init(index string) {
	c.ColumnBase.Init(c)
	c.SetCellTexter(GetterTexter{Key: index})
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

func (t GetterTexter) CellText (ctx context.Context, col ColumnI, rowNum int, colNum int, data interface{}) string {
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

