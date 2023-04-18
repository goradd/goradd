package generator

import (
	"fmt"
	"github.com/goradd/goradd/codegen/generator"
	"github.com/goradd/goradd/pkg/orm/db"
	"github.com/goradd/goradd/pkg/orm/query"
)

func init() {
	generator.RegisterControlGenerator(EmailTextbox{}, "github.com/goradd/goradd/pkg/page/control/textbox/EmailTextbox")
}

// EmailTextbox describes the EmailTextbox to the connector dialog and code generator
type EmailTextbox struct {
}

func (d EmailTextbox) SupportsColumn(ref interface{}) bool {
	if col, ok := ref.(*db.Column); ok &&
		(col.ColumnType == query.ColTypeBytes ||
			col.ColumnType == query.ColTypeString) {
		return true
	}
	return false
}

func (d EmailTextbox) GenerateCreator(ref interface{}, desc *generator.ControlDescription) (s string) {
	col := ref.(*db.Column)
	s = fmt.Sprintf(
		`%s.EmailTextboxCreator{
			ID:        p.ID() + "-%s",
			MaxLength: %d,
			ControlOptions: page.ControlOptions{
				IsRequired:      %#v,
				DataConnector: %s{},
			},
		}`, desc.Package, desc.ControlID, col.MaxCharLength, !col.IsNullable, desc.Connector)
	return
}

func (d EmailTextbox) GenerateRefresh(ref interface{}, desc *generator.ControlDescription) (s string) {
	return `ctrl.SetText(val)`
}

func (d EmailTextbox) GenerateUpdate(ref interface{}, desc *generator.ControlDescription) (s string) {
	return `val := ctrl.Text()`
}

func (d EmailTextbox) GenerateModifies(ref interface{}, desc *generator.ControlDescription) (s string) {
	return `val != ctrl.Text()`
}
