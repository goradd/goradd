package generator

import (
	"fmt"
	"github.com/goradd/goradd/codegen/generator"
	"github.com/goradd/goradd/pkg/orm/db"
	"github.com/goradd/goradd/pkg/orm/query"
)

func init() {
	generator.RegisterControlGenerator(Checkbox{}, "github.com/goradd/goradd/pkg/page/control/CheckboxList")
}

// Checkbox describes the Checkbox to the connector dialog and code generator
type Checkbox struct {
}

func (d Checkbox) SupportsColumn(ref interface{}) bool {
	if col, ok := ref.(*db.Column); ok && col.ColumnType == query.ColTypeBool {
		return true
	}
	return false
}

func (d Checkbox) GenerateCreator(ref interface{}, desc *generator.ControlDescription) (s string) {
	col := ref.(*db.Column)
	s = fmt.Sprintf(
		`%s.CheckboxCreator{
			ID:        p.ID() + "-%s",
			ControlOptions: page.ControlOptions{
				IsRequired:      %#v,
				DataConnector: %s{},
			},
		}`, desc.Package, desc.ControlID, !col.IsNullable, desc.Connector)
	return
}

func (d Checkbox) GenerateRefresh(ref interface{}, desc *generator.ControlDescription) (s string) {
	return `ctrl.SetChecked(val)`
}

func (d Checkbox) GenerateUpdate(ref interface{}, desc *generator.ControlDescription) (s string) {
	return `val := ctrl.Checked()`
}
