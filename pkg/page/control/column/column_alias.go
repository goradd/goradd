package column

import (
	"context"
	"github.com/goradd/goradd/pkg/orm/query"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/control"
)

type AliasGetter interface {
	GetAlias(key string) query.AliasValue
}

// AliasColumn is a column that uses the AliasGetter interface to get the alias text out of a database object.
// The data therefore should be a slice of objects that implement the AliasGetter interface. All ORM objects
// are AliasGetters (or should be). Call NewAliasColumn to create the column.
type AliasColumn struct {
	control.ColumnBase
	alias string
}

// NewAliasColumn creates a new table column that gets its text from an alias attached to an ORM object.
// If the alias has a Date type, you MUST call SetTimeFormat to set the format of the printed string.
func NewAliasColumn(alias string) *AliasColumn {
	i := AliasColumn{}
	i.Init(alias)
	return &i
}

func (c *AliasColumn) Init(alias string) {
	c.ColumnBase.Init(c)
	c.SetTitle(alias)
	c.alias = alias
}

func GetNode(c *AliasColumn) query.NodeI {
	return query.Alias(c.alias)
}

func (c *AliasColumn) CellData(ctx context.Context, row int, col int, data interface{}) interface{} {
	if v, ok := data.(AliasGetter); !ok {
		return ""
	} else {
		a := v.GetAlias(c.alias)
		if a.IsNil() {
			return ""
		}
		return a
	}
}

func (c *AliasColumn) Serialize(e page.Encoder) (err error) {
	if err = c.ColumnBase.Serialize(e); err != nil {
		return
	}
	if err = e.Encode(c.alias); err != nil {
		return
	}
	return
}

func (c *AliasColumn) Deserialize(dec page.Decoder) (err error) {
	if err = c.ColumnBase.Deserialize(dec); err != nil {
		panic(err)
	}
	if err = dec.Decode(&c.alias); err != nil {
		panic(err)
	}
	return
}

// AliasColumnCreator creates a column that displays the content of a database alias. Each row must be
// an AliasGetter, which by default all the output from database queries provide that.
type AliasColumnCreator struct {
	// ID will assign the given id to the column. If you do not specify it, an id will be given it by the framework.
	ID string
	// Alias is the name of the alias to use when getting data out of the provided database row
	Alias string
	// Title is the static title string to use in the header row
	Title string
	// Sortable makes the column display sort arrows in the header
	Sortable bool
	control.ColumnOptions
}

func (c AliasColumnCreator) Create(ctx context.Context, parent control.TableI) control.ColumnI {
	col := NewAliasColumn(c.Alias)
	if c.ID != "" {
		col.SetID(c.ID)
	}
	col.SetTitle(c.Title)
	if c.Sortable {
		col.SetSortable()
	}
	col.ApplyOptions(ctx, parent, c.ColumnOptions)
	return col
}


func init() {
	control.RegisterColumn(AliasColumn{})
}