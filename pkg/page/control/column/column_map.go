package column

import (
	"context"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/control"
	"reflect"
)

// MapColumn is a table that works with data that is in the form of a map. The data item itself must be convertible into
// a string, either by normal string conversion semantics, or using the supplied format string.
type MapColumn struct {
	control.ColumnBase
	key interface{}
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
	c.key = index
}

func (c *MapColumn) CellData(_ context.Context, rowNum int, colNum int, data interface{}) interface{} {
	vKey := reflect.ValueOf(c.key)
	vMap := reflect.ValueOf(data)
	vValue := vMap.MapIndex(vKey)

	if vValue.IsValid() {
		i := vValue.Interface()
		return i
	}
	return ""
}

func (c *MapColumn) Serialize(e page.Encoder) {
	c.ColumnBase.Serialize(e)
	if err := e.Encode(&c.key); err != nil {
		panic(err)
	}
}

func (c *MapColumn) Deserialize(dec page.Decoder) {
	c.ColumnBase.Deserialize(dec)
	if err := dec.Decode(&c.key); err != nil {
		panic(err)
	}
}



// MapColumnCreator creates a column that treats each row of data as a map of data.
// The index can be any valid map index, and the value must be a standard kind of value that
// can be converted to a string.
type MapColumnCreator struct {
	// ID will assign the given id to the column. If you do not specify it, an id will be given it by the framework.
	ID string
	// Index is the key to use to get to the map data
	Index interface{}
	// Title is the title of the column that appears in the header
	Title string
	// Sortable makes the column display sort arrows in the header
	Sortable bool
	control.ColumnOptions
}

func (c MapColumnCreator) Create(ctx context.Context, parent control.TableI) control.ColumnI {
	col := NewMapColumn(c.Index)
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
	control.RegisterColumn(MapColumn{})
}
