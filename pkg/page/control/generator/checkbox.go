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
	return "github.com/goradd/goradd/pkg/page/control/Checkbox"
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
		`%s.CheckboxCreator{
			ID:        %#v,
			ControlOptions: page.ControlOptions{
				IsRequired:      %#v,
				DataConnector: %s{},
			},
		}`, desc.Import, desc.ControlID, !col.IsNullable, desc.Connector)
	return
}

func (d Checkbox) GenerateRefresh(ref interface{}, desc *generator.ControlDescription) (s string) {
	return `ctrl.SetChecked(val)`
}

func (d Checkbox) GenerateUpdate(ref interface{}, desc *generator.ControlDescription) (s string) {
	return `val := ctrl.Checked()`
}


