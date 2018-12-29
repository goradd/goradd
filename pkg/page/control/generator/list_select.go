package generator

import (
	"fmt"
	"github.com/gedex/inflector"
	"github.com/goradd/gengen/maps"
	"github.com/spekary/goradd/codegen/generator"
	"github.com/spekary/goradd/pkg/config"
	"github.com/spekary/goradd/pkg/page"
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
	return []string{"github.com/spekary/goradd/pkg/page/control"}
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

	if col.ForeignKey != nil {
		if !col.IsNullable {
			s += `	ctrl.AddItem(ctrl.ΩT("- Select One -"), 0)
`
		}
		if col.ForeignKey.IsType {
			s += fmt.Sprintf(`	ctrl.AddListItems(model.%s())
`, inflector.Pluralize(col.ForeignKey.GoType))
		}
	}

	return
}

func (d SelectList) GenerateGet(ctrlName string, objName string, col *generator.ColumnType) (s string) {
	s = fmt.Sprintf(`c.%s.SetValue(c.%s.%s())`, ctrlName, objName, col.ForeignKey.GoName)
	return
}

func (d SelectList) GeneratePut(ctrlName string, objName string, col *generator.ColumnType) (s string) {
	if col.ForeignKey != nil && col.ForeignKey.IsType {
		s = fmt.Sprintf(`c.%s.Set%s(c.%s.Value().(model.%s))`, objName, col.ForeignKey.GoName, ctrlName, col.ForeignKey.GoType)
	} else {
		s = fmt.Sprintf(`c.%s.Set%s(c.%s.StringValue())`, objName, col.GoName, ctrlName)
	}
	return
}


func (d SelectList) ConnectorParams() *maps.SliceMap {
	paramControls := page.ControlConnectorParams()

	return paramControls
}

