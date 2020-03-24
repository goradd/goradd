package generator

import (
	"fmt"
	"github.com/goradd/goradd/codegen/generator"
	"github.com/goradd/goradd/pkg/orm/db"
	"github.com/goradd/goradd/pkg/orm/query"
)

func init() {
	generator.RegisterControlGenerator(PhoneTextbox{}, "github.com/goradd/goradd/pkg/page/control/PhoneTextbox")
}

// This structure describes the PhoneTextbox to the connector dialog and code generator
type PhoneTextbox struct {
}

func (d PhoneTextbox) SupportsColumn(ref interface{}) bool {
	if col,ok := ref.(*db.Column); ok &&
		(col.ColumnType == query.ColTypeBytes ||
			col.ColumnType == query.ColTypeString) {
		return true
	}
	return false

}

func (d PhoneTextbox) GenerateCreator(ref interface{}, desc *generator.ControlDescription) (s string) {
	col := ref.(*db.Column)
	s = fmt.Sprintf(
`%s.PhoneTextboxCreator{
	ID:        p.ID() + "-%s",
	ControlOptions: page.ControlOptions{
		IsRequired:      %#v,
		DataConnector: %s{},
	},
}`, desc.Package, desc.ControlID, !col.IsNullable, desc.Connector)
	return
}


func (d PhoneTextbox) GenerateRefresh(ref interface{}, desc *generator.ControlDescription) (s string) {
	return `ctrl.SetText(val)`
}

func (d PhoneTextbox) GenerateUpdate(ref interface{}, desc *generator.ControlDescription) (s string) {
	return `val := ctrl.Text()`
}
