package column

import (
	"context"
	"github.com/spekary/goradd/pkg/page/control"
	"github.com/spekary/goradd/pkg/orm/query"
)

type AliasGetter interface {
	GetAlias(key string) query.AliasValue
}

// AliasColumn is a column that uses the AliasGetter interface to get the alias text out of a database object.
// The data therefore should be a slice of objects that implement the AliasGetter interface. All ORM objects
// are AliasGetters (or should be).
type AliasColumn struct {
	control.ColumnBase
	alias string
}

// NewAliasColumn creates a new table column that gets its text from an alias attached to an ORM object.
func NewAliasColumn(alias string, format ...string) *AliasColumn {
	i := AliasColumn{}
	var f string
	if len(format) > 0 {
		f = format[0]
	}
	i.Init(alias, f, "")
	return &i
}

// NewAliasColumn creates a new table column that gets its text from an alias attached to an ORM object.
// The alias should get a DateTime type of data.
func NewDateAliasColumn(alias string, timeFormat string, format ...string) *AliasColumn {
	i := AliasColumn{}
	var f string
	if len(format) > 0 {
		f = format[0]
	}
	i.Init(alias, f, timeFormat)
	return &i
}

func (c *AliasColumn) Init(alias string, format string, timeFormat string) {
	c.ColumnBase.Init(c)
	c.SetCellTexter(&AliasTexter{Alias: alias, Format: format, TimeFormat: timeFormat})
	c.SetTitle(alias)
}

// SetFormat sets the format string of the node column.
func (c *AliasColumn) SetFormat(format string) *AliasColumn {
	c.CellTexter().(*AliasTexter).Format = format
	return c
}

// SetTimeFormat sets the time format of the string, specifically for a DateTime column.
func (c *AliasColumn) SetTimeFormat(format string) *AliasColumn {
	c.CellTexter().(*AliasTexter).TimeFormat = format
	return c
}

// AliasTexter gets text out of an ORM object with an alias. If the alias does not exist, it will panic.
type AliasTexter struct {
	// Alias is the alias name in the database object that we are interested in.
	Alias string
	// Format is a format string. It will be applied using fmt.Sprintf. If you don't provide a Format string, standard
	// string conversion operations will be used.
	Format string
	// TimeFormat is applied to the data using time.Format. You can have both a Format and TimeFormat, and the Format
	// will be applied using fmt.Sprintf after the TimeFormat is applied using time.Format.
	TimeFormat string
}

func (t AliasTexter) CellText(ctx context.Context, col control.ColumnI, rowNum int, colNum int, data interface{}) string {
	if v,ok := data.(AliasGetter); !ok {
		return ""
	} else {
		a := v.GetAlias(t.Alias)
		if a.IsNil() {
			return ""
		}
		if t.TimeFormat != "" {
			// assume we are trying to get a time
			d := a.DateTime()
			return ApplyFormat(d, t.Format, t.TimeFormat)
		}
		s := v.GetAlias(t.Alias).String()
		return ApplyFormat(s, t.Format, t.TimeFormat)
	}
}
