package temp

import (
	"context"
	reflect "reflect"
	"github.com/spekary/goradd/html"
	html2 "html"
)

// SliceColumn is a table that works with data that is in the form of a slice. The data item itself must be convertable into
// a string, either by normal string conversion symantecs, or using the supplied format string. The format string will
// be applied to a date if the data is a date, or to the string using fmt.Sprintf
type LinkColumn struct {
	ColumnBase
	label interface{}	// could be a string, or if a slice of strings, it will examine the object for a value
	destination interface{}
	getVars interface{}
	tagAttributes *html.Attributes
	asButton bool
}

func NewProxyColumn(label interface{},	// could be a string, or if a slice of strings, it will examine the object for a value
	destination interface{},
	getVars interface{},
	tagAttributes *html.Attributes,
	asButton bool,
) *LinkColumn {
	i := LinkColumn{}
	i.Init(index)
	return &i
}

func (c *LinkColumn) Init(index int) {
	c.ColumnBase.Init(c)
	c.ColumnBase.dontEscape = true
	c.SetCellTexter(SliceTexter{Index: index})
}

// LinkTexter creates proxy links and buttons.
type LinkTexter struct {
	Label interface{}
	Destination interface{}
	GetVarsOrActionValue interface{}
	Attributes *html.Attributes
	// Format is a format string. It will be applied using fmt.Sprintf. If you don't provide a Format string, standard
	// string conversion operations will be used.
	Format string
	// TimeFormat is applied to the target data using time.Format. You can have both a Format and TimeFormat, and the Format
	// will be applied using fmt.Sprintf after the TimeFormat is applied using time.Format.
	TimeFormat string
	AsButton bool
}

func (t LinkTexter) CellText (ctx context.Context, col ColumnI, rowNum int, colNum int, data interface{}) string {
	label := t.getObjectHtmlValue(t.Label, data)
	var getVars interface{}

	if t.GetVarsOrActionValue != nil {
		switch vars := t.GetVarsOrActionValue.(type) {
		case string: // Should already be escaped, since they might be html
			getVars = vars
		case map[string]string: // A map of things to hunt for in the data
			v2 := map[string]string{}
			for key,val := range vars {
				v2[key] = t.getObjectHtmlValue([]string{val}, data)
			}
			vars = v2
		}
	}

	var tagAttributes = html.NewAttributes()
	if t.Attributes

}

func (t LinkTexter) getObjectHtmlValue(spec interface{}, data interface{}) string {
	dataValue := reflect.ValueOf(data)

	switch v := spec.(type) {
	case string:
		return html2.EscapeString(v)
	case int:
		if dataValue.Kind() == reflect.Slice {
			val := dataValue.Index(v)
			return ApplyFormat(val, t.Format, t.TimeFormat)
		}
		return ApplyFormat(v, t.Format, "")
	case []string:
		switch d := data.(type) {
		case Getter:
			v2 := d.Get(v[0])
			if len(v) > 1 {
				return t.getObjectHtmlValue(v[1:], v2)	// Go to the next level
			}
			return ApplyFormat(v2, t.Format, t.TimeFormat)
		case StringGetter:
			v2 := d.Get(v[0])
			if len(v) > 1 {
				panic ("Can't traverse a StringGetter")
			}
			return ApplyFormat(v2, t.Format, t.TimeFormat)

		case map[string]interface{}:
			v2 := d[v[0]]
			if len(v) > 1 {
				return t.getObjectHtmlValue(v[1:], v2)	// Go to the next level
			}
			return ApplyFormat(v2, t.Format, t.TimeFormat)
		case map[string]string:
			v2 := d[v[0]]
			if len(v) > 1 {
				panic ("Can't traverse a string map")
			}
			return ApplyFormat(v2, t.Format, t.TimeFormat)
		}
	default:
		panic("Unknown spec type")
	}

	return ""
}