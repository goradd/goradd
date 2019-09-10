package generator

import (
	"fmt"
	"github.com/goradd/gengen/pkg/maps"
	"github.com/goradd/goradd/codegen/generator"
	"github.com/goradd/goradd/pkg/config"
	"github.com/goradd/goradd/pkg/orm/query"
	"github.com/goradd/goradd/pkg/page"
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

func (d Checkbox) Imports() []generator.ImportPath {
	return []generator.ImportPath{
		{Alias: "goraddctrl", Path:"github.com/goradd/goradd/pkg/page/control"},
	}
}

func (d Checkbox) SupportsColumn(col *generator.ColumnType) bool {
	if col.ColumnType == query.ColTypeBool {
		return true
	}
	return false
}

func (d Checkbox) GenerateCreate(namespace string, col *generator.ColumnType) (s string) {
	s = fmt.Sprintf(
		`	ctrl = %s.NewCheckbox(c.ParentControl, id)
	ctrl.SetLabel(ctrl.T("%s"))
`, namespace, col.DefaultLabel)


	return
}

func (d Checkbox) GenerateGet(ctrlName string, objName string, col *generator.ColumnType) (s string) {
	s = fmt.Sprintf(`c.%s.SetChecked(c.%s.%s())`, ctrlName, objName, col.GoName)
	return
}

func (d Checkbox) GeneratePut(ctrlName string, objName string, col *generator.ColumnType) (s string) {
	s = fmt.Sprintf(`c.%s.Set%s(c.%s.Checked())`, objName, col.GoName, ctrlName)
	return
}

func (d Checkbox) ConnectorParams() *maps.SliceMap {
	paramControls := page.ControlConnectorParams()

	return paramControls
}

func (d Checkbox) GenerateCreator(col *generator.ColumnType) (s string) {
	s = fmt.Sprintf(
		`goraddctrl.CheckboxCreator{
			ID:        %#v,
			ControlOptions: page.ControlOptions{
				IsRequired:      %#v,
				DataConnector: %s{},
			},
		}`, col.ControlID, !col.IsNullable, col.Connector)
	return
}

func (d Checkbox) GenerateRefresh(col *generator.ColumnType) (s string) {
	return `ctrl.SetChecked(val)`
}

func (d Checkbox) GenerateUpdate(col *generator.ColumnType) (s string) {
	return `val := ctrl.Checked()`
}


