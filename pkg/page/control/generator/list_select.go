package generator

import (
	"fmt"
	"github.com/gedex/inflector"
	"github.com/goradd/gengen/pkg/maps"
	"github.com/goradd/goradd/codegen/generator"
	"github.com/goradd/goradd/pkg/config"
	"github.com/goradd/goradd/pkg/page"
)

func init() {
	if !config.Release {
		generator.RegisterControlGenerator(SelectList{})
	}
}

// This structure describes the SelectList to the connector dialog and code generator
type SelectList struct {
}

func (d SelectList) Type() string {
	return "SelectList"
}

func (d SelectList) NewFunc() string {
	return "NewSelectList"
}

func (d SelectList) Imports() []string {
	return []string{"github.com/goradd/goradd/pkg/page/control"}
}

// TODO: This has to be changed to support virtual column types like ManyMany and Reverse
func (d SelectList) SupportsColumn(col *generator.ColumnType) bool {
	if col.ForeignKey != nil {
		return true
	}
	return false
}

func (d SelectList) GenerateCreate(namespace string, col *generator.ColumnType) (s string) {
	s = fmt.Sprintf(
		`	ctrl = %s.NewSelectList(c.ParentControl, id)
	ctrl.SetLabel(ctrl.T("%s"))
`, namespace, col.DefaultLabel)

	if generator.DefaultWrapper != "" {
		s += fmt.Sprintf(`	ctrl.With(page.NewWrapper("%s"))
`, generator.DefaultWrapper)
	}

	return
}

func (d SelectList) GenerateGet(ctrlName string, objName string, col *generator.ColumnType) (s string) {
	isType := col.ForeignKey != nil && col.ForeignKey.IsType
	colName := col.GoName
	if isType {
		colName = col.ForeignKey.GoName
	}

	// TODO: Possibly deal with non-type lists, though setting up those lists is usually very customized
	if isType {

		s += fmt.Sprintf(
`	ctrl := c.%s
	ctrl.Clear()
`, ctrlName)
		if col.IsNullable {
			s += `	ctrl.AddItem(ctrl.ΩT("- None -"), nil)
`
		} else {
			// Deal with situation where a selection is required, but has not yet been made.
			s += fmt.Sprintf(
`	if c.%s.%s() == 0 {
		ctrl.AddItem(ctrl.ΩT("- Select One -"), 0)
	}
`, objName, colName)
		}

		s += fmt.Sprintf(`	ctrl.AddListItems(model.%s())
`, inflector.Pluralize(col.ForeignKey.GoType))
	}

	if col.IsNullable {
		s += fmt.Sprintf(
			// Use nil for null database values
`	if c.%[2]s.%[3]sIsNull() {
		c.%[1]s.SetValue(nil)
	} else {
		c.%[1]s.SetValue(c.%[2]s.%[3]s())
	}
`, ctrlName, objName, colName)
	} else {
		// Don't use nil, but just use whatever the object has
		// We are assuming that in the situation where no selection has been made, the database object
		// has a default empty value that matches the value given to the - Select One - list item
		s += fmt.Sprintf(
`	c.%s.SetValue(c.%s.%s())
`,
			ctrlName, objName, colName)
	}
	return
}

func (d SelectList) GeneratePut(ctrlName string, objName string, col *generator.ColumnType) (s string) {
	if col.IsNullable {
		if col.ForeignKey != nil && col.ForeignKey.IsType {
			s = fmt.Sprintf(
`	if v := c.%[3]s.Value(); v == nil {
		c.%[1]s.Set%[2]s(nil)
	} else {
		c.%[1]s.Set%[2]s(v.(model.%[4]s))
	}
`,
			objName, col.ForeignKey.GoName, ctrlName, col.ForeignKey.GoType)
		} else {
			s = fmt.Sprintf(`	c.%s.Set%s(c.%s.Value())`, objName, col.GoName, ctrlName)
		}
	} else {
		if col.ForeignKey != nil && col.ForeignKey.IsType {
			s = fmt.Sprintf(`	c.%s.Set%s(c.%s.Value().(model.%s))`, objName, col.ForeignKey.GoName, ctrlName, col.ForeignKey.GoType)
		} else {
			s = fmt.Sprintf(`	c.%s.Set%s(c.%s.StringValue())`, objName, col.GoName, ctrlName)
		}
	}
	return
}

func (d SelectList) ConnectorParams() *maps.SliceMap {
	paramControls := page.ControlConnectorParams()

	return paramControls
}
