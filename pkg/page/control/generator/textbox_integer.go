package generator

import (
	"fmt"
	"github.com/goradd/gengen/pkg/maps"
	"github.com/goradd/goradd/codegen/generator"
	"github.com/goradd/goradd/pkg/config"
	"github.com/goradd/goradd/pkg/orm/query"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/control"
)

func init() {
	if !config.Release {
		//generator.RegisterControlGenerator(IntegerTextbox{})
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

func (d IntegerTextbox) Imports() []string {
	return []string{"github.com/goradd/goradd/pkg/page/control"}
}

func (d IntegerTextbox) SupportsColumn(col *generator.ColumnType) bool {
	if (col.ColumnType == query.ColTypeInteger ||
		col.ColumnType == query.ColTypeInteger64) &&
		!col.IsReference() {
		return true
	}
	return false
}

func (d IntegerTextbox) GenerateCreate(namespace string, col *generator.ColumnType) (s string) {
	s = fmt.Sprintf(
		`	ctrl = %s.NewIntegerTextbox(c.ParentControl, id)
	ctrl.SetLabel(ctrl.T("%s"))
`, namespace, col.DefaultLabel)

	// TODO: Set a maximum value based on database limit

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
	switch col.ColumnType {
	case query.ColTypeInteger64:
		s = fmt.Sprintf(`c.%s.Set%s(c.%s.Int64())`, objName, col.GoName, ctrlName)
	case query.ColTypeInteger:
		s = fmt.Sprintf(`c.%s.Set%s(c.%s.Int())`, objName, col.GoName, ctrlName)
	case query.ColTypeUnsigned64:
		s = fmt.Sprintf(`c.%s.Set%s(uint64(c.%s.Int64()))`, objName, col.GoName, ctrlName)
	case query.ColTypeUnsigned:
		s = fmt.Sprintf(`c.%s.Set%s(uint(c.%s.Int()))`, objName, col.GoName, ctrlName)

	}
	return
}

func (d IntegerTextbox) ConnectorParams() *maps.SliceMap {
	paramControls := page.ControlConnectorParams()
	paramSet := maps.NewSliceMap()

	// TODO: Get the regular Textbox's parameters too
	paramSet.Set("ColumnCount", generator.ConnectorParam{
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
