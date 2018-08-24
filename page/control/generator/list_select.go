package generator

import (
	"goradd-project/config"
	"fmt"
	"goradd-project/config/codegen"
	"github.com/spekary/goradd/util/types"
	"github.com/spekary/goradd/page"
	"github.com/gedex/inflector"
	"github.com/spekary/goradd/codegen/generator"
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

func (d SelectList) Import() string {
	// TODO: Add fmt to the import list
	return "github.com/spekary/goradd/page/control"
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
	ctrl.SetLabel("%s")
`, namespace, col.DefaultLabel)

	if codegen.DefaultWrapper != "" {
		s += fmt.Sprintf(`	ctrl.With(page.NewWrapper("%s"))
`, codegen.DefaultWrapper)
	}

	if col.ForeignKey != nil {
		if !col.IsNullable {
			s += `	ctrl.AddItem(ctrl.ParentForm().T("- Select One -"), 0)
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


func (d SelectList) ConnectorParams() *types.OrderedMap {
	paramControls := page.ControlConnectorParams()

	return paramControls
}

