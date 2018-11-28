package generator

import (
	"fmt"
	"github.com/spekary/gengen/maps"
	"github.com/spekary/goradd/codegen/generator"
	"github.com/spekary/goradd/pkg/orm/query"
	"github.com/spekary/goradd/pkg/page"
	"github.com/spekary/goradd/pkg/page/control"
	"goradd-project/config"
	"goradd-project/config/codegen"
)

func init() {
	if !config.Release {
		generator.RegisterControlGenerator(FloatTextbox{})
	}
}

// This structure describes the FloatTextbox to the connector dialog and code generator
type FloatTextbox struct {

}

func (d FloatTextbox) Type() string {
	return "FloatTextbox"
}

func (d FloatTextbox) NewFunc() string {
	return "NewFloatTextbox"
}

func (d FloatTextbox) Import() string {
	return "github.com/spekary/goradd/pkg/page/control"
}

func (d FloatTextbox) SupportsColumn(col *generator.ColumnType) bool {
	if col.GoType == query.ColTypeFloat || col.GoType == query.ColTypeDouble {
		return true
	}
	return false
}

func (d FloatTextbox) GenerateCreate(namespace string, col *generator.ColumnType) (s string) {
	s = fmt.Sprintf(
`	ctrl = %s.NewFloatTextbox(c.ParentControl, id)
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

func (d FloatTextbox) GenerateGet(ctrlName string, objName string, col *generator.ColumnType) (s string) {
	if col.GoType == query.ColTypeFloat {
		s = fmt.Sprintf(`c.%s.SetFloat32(c.%s.%s())`, ctrlName, objName, col.GoName)
	} else {
		s = fmt.Sprintf(`c.%s.SetFloat64(c.%s.%s())`, ctrlName, objName, col.GoName)
	}
	return
}

func (d FloatTextbox) GeneratePut(ctrlName string, objName string, col *generator.ColumnType) (s string) {
	if col.GoType == query.ColTypeFloat {
		s = fmt.Sprintf(`c.%s.Set%s(c.%s.Float32())`, objName,  col.GoName, ctrlName)
	} else {
		s = fmt.Sprintf(`c.%s.Set%s(c.%s.Float64())`, objName,  col.GoName, ctrlName)
	}
	return
}


func (d FloatTextbox) ConnectorParams() *maps.SliceMap {
	paramControls := page.ControlConnectorParams()
	paramSet := maps.NewSliceMap()

	// TODO: Get the regular Textbox's parameters too
	paramSet.Set("ColumnCount", generator.ConnectorParam {
		"Column Count",
		"Width of field by the number of characters.",
		generator.ControlTypeInteger,
		`{{var}}.SetColumnCount{{val}}`,
		func(c page.ControlI, val interface{}) {
			c.(*control.FloatTextbox).SetColumnCount(val.(int))
		}})


	paramControls.Set("FloatTextbox", paramSet)

	return paramControls
}

