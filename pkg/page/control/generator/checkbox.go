package generator

import (
	"github.com/spekary/gengen/maps"
	"goradd-project/config"
	"fmt"
	"goradd-project/config/codegen"
	"github.com/spekary/goradd/pkg/page"
	"github.com/spekary/goradd/pkg/orm/query"
	"github.com/spekary/goradd/tools/codegen/generator"
)

func init() {
	if !config.Release {
		generator.RegisterControlGenerator(Checkbox{})
	}
}

// This structure describes the Checkbox to the connector dialog and code generator
type Checkbox struct {

}

func (d Checkbox) Type() string {
	return "Checkbox"
}

func (d Checkbox) NewFunc() string {
	return "NewCheckbox"
}

func (d Checkbox) Import() string {
	// TODO: Add fmt to the import list
	return "github.com/spekary/goradd/pkg/page/control"
}

func (d Checkbox) SupportsColumn(col *generator.ColumnType) bool {
	if col.GoType == query.ColTypeBool {
		return true
	}
	return false
}

func (d Checkbox) GenerateCreate(namespace string, col *generator.ColumnType) (s string) {
	s = fmt.Sprintf(
		`	ctrl = %s.NewCheckbox(c.ParentControl, id)
	ctrl.SetLabel("%s")
`, namespace, col.DefaultLabel)

	if codegen.DefaultWrapper != "" {
		s += fmt.Sprintf(`	ctrl.With(page.NewWrapper("%s"))
`, codegen.DefaultWrapper)
	}

	return
}

func (d Checkbox) GenerateGet(ctrlName string, objName string, col *generator.ColumnType) (s string) {
	s = fmt.Sprintf(`c.%s.SetChecked(c.%s.%s())`, ctrlName, objName, col.GoName)
	return
}

func (d Checkbox) GeneratePut(ctrlName string, objName string, col *generator.ColumnType) (s string) {
	s = fmt.Sprintf(`c.%s.Set%s(c.%s.Checked())`, objName,  col.GoName, ctrlName)
	return
}


func (d Checkbox) ConnectorParams() *maps.SliceMap {
	paramControls := page.ControlConnectorParams()

	return paramControls
}

