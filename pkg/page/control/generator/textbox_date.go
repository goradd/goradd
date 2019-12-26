package generator

import (
	"fmt"
	"github.com/goradd/goradd/codegen/generator"
	"github.com/goradd/goradd/pkg/orm/db"
	"github.com/goradd/goradd/pkg/orm/query"
)

func init() {
	generator.RegisterControlGenerator(DateTextbox{}, "github.com/goradd/goradd/pkg/page/control/DateTextbox")
}

// This structure describes the IntegerTextbox to the connector dialog and code generator
type DateTextbox struct {
}


func (d DateTextbox) SupportsColumn(ref interface{}) bool {
	if col,ok := ref.(*db.Column); ok &&
		col.ColumnType == query.ColTypeDateTime {
		return true
	}
	return false
}

func (d DateTextbox) GenerateCreator(ref interface{}, desc *generator.ControlDescription) (s string) {
	col := ref.(*db.Column)
	s = fmt.Sprintf(
`%s.DateTextboxCreator{
	ID:        %#v,
	ControlOptions: page.ControlOptions{
		IsRequired:      %#v,
		DataConnector: %s{},
	},
}`, desc.Import, desc.ControlID, !col.IsNullable, desc.Connector)
	return
}


func (d DateTextbox) GenerateRefresh(ref interface{}, desc *generator.ControlDescription) (s string) {
	return `ctrl.SetValue(val)`
}

func (d DateTextbox) GenerateUpdate(ref interface{}, desc *generator.ControlDescription) (s string) {
	return `val := ctrl.Date()`
}
