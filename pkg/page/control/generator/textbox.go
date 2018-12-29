package generator

import (
	"fmt"
	"github.com/goradd/gengen/maps"
	"github.com/spekary/goradd/codegen/generator"
	"github.com/spekary/goradd/pkg/config"
	"github.com/spekary/goradd/pkg/orm/query"
	"github.com/spekary/goradd/pkg/page"
	"github.com/spekary/goradd/pkg/page/control"
)

func init() {
	if !config.Release {
		generator.RegisterControlGenerator(Textbox{})
	}
}

// This structure describes the textbox to the connector dialog and code generator
type Textbox struct {

}

func (d Textbox) Type() string {
	return "Textbox"
}

func (d Textbox) NewFunc() string {
	return "NewTextbox"
}

func (d Textbox) Imports() []string {
	return []string{"github.com/spekary/goradd/pkg/page/control"}
}

func (d Textbox) SupportsColumn(col *generator.ColumnType) bool {
	if col.GoType == query.ColTypeBytes ||
		col.GoType == query.ColTypeString {
		return true
	}
	return false
}

func (d Textbox) GenerateCreate(namespace string, col *generator.ColumnType) (s string) {
	s = fmt.Sprintf(
`	ctrl = %s.NewTextbox(c.ParentControl, id)
	ctrl.SetLabel(ctrl.T("%s"))
`, namespace, col.DefaultLabel)
	if col.MaxCharLength > 0 {
		s += fmt.Sprintf(`	ctrl.SetMaxLength(%d)	
`, col.MaxCharLength)
	}

	if generator.DefaultWrapper != "" {
		s += fmt.Sprintf(`	ctrl.With(page.NewWrapper("%s"))
`, generator.DefaultWrapper)
	}
	if col.IsPk {
		s += `	ctrl.SetDisabled(true)
`
	} else if !col.IsNullable {
		s += `	ctrl.SetIsRequired(true)
`
	}

	return
}

func (d Textbox) GenerateGet(ctrlName string, objName string, col *generator.ColumnType) (s string) {
	s = fmt.Sprintf(`c.%s.SetText(c.%s.%s())`, ctrlName, objName, col.GoName)
	return
}

func (d Textbox) GeneratePut(ctrlName string, objName string, col *generator.ColumnType) (s string) {
	s = fmt.Sprintf(`c.%s.Set%s(c.%s.Text())`, objName,  col.GoName, ctrlName)
	return
}


func (d Textbox) ConnectorParams() *maps.SliceMap {
	paramControls := page.ControlConnectorParams()
	paramSet := maps.NewSliceMap()

	paramSet.Set("ColumnCount", generator.ConnectorParam {
		"Column Count",
		"Width of field by the number of characters.",
		generator.ControlTypeInteger,
		`{{var}}.SetColumnCount{{val}}`,
		func(c page.ControlI, val interface{}) {
			c.(*control.Textbox).SetColumnCount(val.(int))
		}})


	paramControls.Set("Textbox", paramSet)

	return paramControls
}

