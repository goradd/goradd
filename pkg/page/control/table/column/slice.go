package column

import (
	"context"
	"github.com/goradd/goradd/pkg/page/control/table"
	"reflect"

	"github.com/goradd/goradd/pkg/page"
)

// SliceColumn is a table that works with data that is in the form of a slice. The data item itself must be convertible into
// a string, either by normal string conversion semantics, or using the supplied format string.
type SliceColumn struct {
	table.ColumnBase
	index int
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
	c.index = index
}

func (c *SliceColumn) CellData(_ context.Context, _ int, _ int, data interface{}) interface{} {
	vSlice := reflect.ValueOf(data)
	if vSlice.Kind() != reflect.Slice {
		panic("data must be a slice.")
	}
	return vSlice.Index(c.index).Interface()
}

func (c *SliceColumn) Serialize(e page.Encoder) {
	c.ColumnBase.Serialize(e)
	if err := e.Encode(c.index); err != nil {
		panic(err)
	}
}

func (c *SliceColumn) Deserialize(dec page.Decoder) {
	c.ColumnBase.Deserialize(dec)
	if err := dec.Decode(&c.index); err != nil {
		panic(err)
	}
}

// SliceColumnCreator creates a column that treats each row as a slice of data.
type SliceColumnCreator struct {
	// ID will assign the given id to the column. If you do not specify it, an id will be given it by the framework.
	ID string
	// Index is the slice index that will be used to get to the data in the column
	Index int
	// Title is the title of the column and will appear in the header
	Title string
	// Sortable makes the column display sort arrows in the header
	// Deprecated: Use SortDirection instead
	Sortable bool
	// SortDirection sets the initial sorting direction of the column, and will make the column sortable
	// By default, the column is not sortable.
	SortDirection table.SortDirection
	// IsHtml indicates that the texter is producing HTML rather than text that should be escaped.
	table.ColumnOptions
}

func (c SliceColumnCreator) Create(ctx context.Context, parent table.TableI) table.ColumnI {
	col := NewSliceColumn(c.Index)
	if c.ID != "" {
		col.SetID(c.ID)
	}
	col.SetTitle(c.Title)
	if c.Sortable {
		col.SetSortable()
	}
	if c.SortDirection != table.NotSortable {
		col.SetSortDirection(c.SortDirection)
	}
	col.ApplyOptions(ctx, parent, c.ColumnOptions)
	return col
}

func init() {
	table.RegisterColumn(SliceColumn{})
}
