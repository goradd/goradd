package generator

import (
	"goradd-project/config"
	"github.com/spekary/goradd/orm/query"
	"fmt"
	"goradd-project/config/codegen"
	"github.com/spekary/goradd/util/types"
	"github.com/spekary/goradd/page"
	"github.com/spekary/goradd/page/control"
	"github.com/spekary/goradd/codegen/generator"
)

func init() {
	if !config.Release {
		generator.RegisterControlGenerator(IntegerTextbox{})
	}
}

// This structure describes the IntegerTextbox to the connector dialog and code generator
type IntegerTextbox struct {

}

func (d IntegerTextbox) Type() string {
	return "IntegerTextbox"
}

func (d IntegerTextbox) NewFunc() string {
	return "NewIntegerTextbox"
}

func (d IntegerTextbox) Import() string {
	return "github.com/spekary/goradd/page/control"
}

func (d IntegerTextbox) SupportsColumn(col *generator.ColumnType) bool {
	if (col.GoType == query.ColTypeInteger ||
		col.GoType == query.ColTypeInteger64) &&
		!col.IsReference() {
		return true
	}
	return false
}

func (d IntegerTextbox) GenerateCreate(namespace string, col *generator.ColumnType) (s string) {
	s = fmt.Sprintf(
`	ctrl = %s.NewIntegerTextbox(c.ParentControl, id)
	ctrl.SetLabel("%s")
`, namespace, col.DefaultLabel)

	// TODO: Set a maximum value based on database limit

	if codegen.DefaultWrapper != "" {
		s += fmt.Sprintf(`	ctrl.With(page.NewWrapper("%s"))
`, codegen.DefaultWrapper)
	}
	if !col.IsNullable {
		s += `	ctrl.SetIsRequired(true)
`
	}

	return
}

func (d IntegerTextbox) GenerateGet(ctrlName string, objName string, col *generator.ColumnType) (s string) {
	s = fmt.Sprintf(`c.%s.SetInt(int(c.%s.%s()))`, ctrlName, objName, col.GoName)
	return
}

func (d IntegerTextbox) GeneratePut(ctrlName string, objName string, col *generator.ColumnType) (s string) {
	if col.GoType == query.ColTypeInteger64 {
		s = fmt.Sprintf(`c.%s.Set%s(c.%s.Int64())`, objName,  col.GoName, ctrlName)
	} else {
		s = fmt.Sprintf(`c.%s.Set%s(c.%s.Int())`, objName,  col.GoName, ctrlName)
	}
	return
}


func (d IntegerTextbox) ConnectorParams() *types.OrderedMap {
	paramControls := page.ControlConnectorParams()
	paramSet := types.NewOrderedMap()

	// TODO: Get the regular Textbox's parameters too
	paramSet.Set("ColumnCount", generator.ConnectorParam {
		"Column Count",
		"Width of field by the number of characters.",
		generator.ControlTypeInteger,
		`{{var}}.SetColumnCount{{val}}`,
		func(c page.ControlI, val interface{}) {
			c.(*control.IntegerTextbox).SetColumnCount(val.(int))
		}})


	paramControls.Set("IntegerTextbox", paramSet)

	return paramControls
}

