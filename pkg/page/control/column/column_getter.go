package column

import (
	"context"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/control"
)

// GetterColumn is a column that uses the Getter interface to get the text out of columns. The data therefore should be
// a slice of objects that implement the Getter interface.
type GetterColumn struct {
	control.ColumnBase
	key string
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
	c.SetTitle(index)
	c.key = index
}

func (c *GetterColumn) CellData(_ context.Context, rowNum int, colNum int, data interface{}) interface{} {
	switch v := data.(type) {
	case Getter:
		return v.Get(c.key)
	case StringGetter:
		return v.Get(c.key)
	}
	return ""
}

func (c *GetterColumn) Serialize(e page.Encoder) {
	c.ColumnBase.Serialize(e)

	if err := e.Encode(c.key); err != nil {
		panic(err)
	}
}

func (c *GetterColumn) Deserialize(dec page.Decoder) {
	c.ColumnBase.Deserialize(dec)

	if err := dec.Decode(&c.key); err != nil {
		panic(err)
	}
}


// GetterColumnCreator creates a column that uses a Getter to get the text of each cell.
type GetterColumnCreator struct {
	// ID will assign the given id to the column. If you do not specify it, an id will be given it by the framework.
	ID string
	// Index is the value passed to the Get function of each row of the data to get the data for the cell.
	Index string
	// Title is the title that appears in the header of the column
	Title string
	// Format is a format string applied to the data using fmt.Sprintf
	Format string
	// TimeFormat is a format string applied specifically to time data using time.Format
	TimeFormat string
	// Sortable makes the column display sort arrows in the header
	Sortable bool
	control.ColumnOptions
}

func (c GetterColumnCreator) Create(ctx context.Context, parent control.TableI) control.ColumnI {
	col := NewGetterColumn(c.Index)
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
	if c.Sortable {
		col.SetSortable()
	}
	col.ApplyOptions(ctx, parent, c.ColumnOptions)
	return col
}

func init() {
	control.RegisterColumn(GetterColumn{})
}
