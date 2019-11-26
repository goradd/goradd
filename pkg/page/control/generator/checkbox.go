package generator

import (
	"fmt"
	"github.com/goradd/goradd/codegen/generator"
	"github.com/goradd/goradd/pkg/config"
	"github.com/goradd/goradd/pkg/orm/db"
	"github.com/goradd/goradd/pkg/orm/query"
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

func (d Checkbox) SupportsColumn(ref interface{}) bool {
	if col,ok := ref.(*db.Column); ok && col.ColumnType == query.ColTypeBool {
		return true
	}
	return false
}

func (d Checkbox) GenerateCreator(ref interface{}, desc *generator.ControlDescription) (s string) {
	col := ref.(*db.Column)
	s = fmt.Sprintf(
		`goraddctrl.CheckboxCreator{
			ID:        %#v,
			ControlOptions: page.ControlOptions{
				IsRequired:      %#v,
				DataConnector: %s{},
			},
		}`, desc.ControlID, !col.IsNullable, desc.Connector)
	return
}

func (d Checkbox) GenerateRefresh(ref interface{}, desc *generator.ControlDescription) (s string) {
	return `ctrl.SetChecked(val)`
}

func (d Checkbox) GenerateUpdate(ref interface{}, desc *generator.ControlDescription) (s string) {
	return `val := ctrl.Checked()`
}


