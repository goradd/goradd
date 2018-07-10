package generator

import (
	"goradd/config"
	"github.com/spekary/goradd/codegen/connector"
	"github.com/spekary/goradd/orm/db"
	"github.com/spekary/goradd/orm/query"
	"fmt"
	"goradd/config/codegen"
	"github.com/spekary/goradd/util/types"
	"github.com/spekary/goradd/page"
	"github.com/spekary/goradd/page/control"
)

func init() {
	if config.Mode == config.AppModeDevelopment {
		connector.RegisterGenerator(Textbox{})
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

func (d Textbox) Import() string {
	return "github.com/spekary/goradd/page/control"
}

func (d Textbox) SupportsColumn(col *db.ColumnDescription) bool {
	if col.GoType == query.ColTypeBytes ||
		col.GoType == query.ColTypeString {
		return true
	}
	return false
}

func (d Textbox) GenerateCreate(namespace string, col *db.ColumnDescription) (s string) {
	s = fmt.Sprintf(
`	ctrl = %s.NewTextbox(c.ParentControl, id)
	ctrl.SetLabel("%s")
`, namespace, col.GoName)
	if col.MaxCharLength > 0 {
		s += fmt.Sprintf(`	ctrl.SetMaxLength(%d)	
`, col.MaxCharLength)
	}

	if codegen.DefaultWrapper != "" {
		s += fmt.Sprintf(`	ctrl.With(page.NewWrapper("%s"))
`, codegen.DefaultWrapper)
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

func (d Textbox) GenerateGet(ctrlName string, objName string, col *db.ColumnDescription) (s string) {
	s = fmt.Sprintf(`c.%s.SetText(c.%s.%s())`, ctrlName, objName, col.GoName)
	return
}

func (d Textbox) GeneratePut(ctrlName string, objName string, col *db.ColumnDescription) (s string) {
	s = fmt.Sprintf(`c.%s.Set%s(c.%s.Text())`, objName,  col.GoName, ctrlName)
	return
}


func (d Textbox) ConnectorParams() *types.OrderedMap {
	paramControls := page.ControlConnectorParams()
	paramSet := types.NewOrderedMap()

	paramSet.Set("ColumnCount", connector.ConnectorParam {
		"Column Count",
		"Width of field by the number of characters.",
		connector.ControlTypeInteger,
		`{{var}}.SetColumnCount{{val}}`,
		func(c page.ControlI, val interface{}) {
			c.(*control.Textbox).SetColumnCount(val.(int))
		}})


	paramControls.Set("Textbox", paramSet)

	return paramControls
}

