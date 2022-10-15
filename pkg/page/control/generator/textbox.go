package generator

import (
	"fmt"
	"github.com/goradd/goradd/codegen/generator"
	"github.com/goradd/goradd/pkg/orm/db"
	"github.com/goradd/goradd/pkg/orm/query"
)

func init() {
	generator.RegisterControlGenerator(Textbox{}, "github.com/goradd/goradd/pkg/page/control/Textbox")
}

// Textbox describes the textbox to the connector dialog and code generator
type Textbox struct {
}

func (d Textbox) SupportsColumn(ref interface{}) bool {
	if col, ok := ref.(*db.Column); ok &&
		(col.ColumnType == query.ColTypeBytes ||
			col.ColumnType == query.ColTypeString) {
		return true
	}
	return false
}

func (d Textbox) GenerateCreator(ref interface{}, desc *generator.ControlDescription) (s string) {
	col := ref.(*db.Column)
	s = fmt.Sprintf(
		`%s.TextboxCreator{
			ID:        p.ID() + "-%s",
			MaxLength: %d,
			ControlOptions: page.ControlOptions{
				IsRequired:      %#v,
				DataConnector: %s{},
			},
		}`, desc.Package, desc.ControlID, col.MaxCharLength, !col.IsNullable, desc.Connector)
	return
}

func (d Textbox) GenerateRefresh(ref interface{}, desc *generator.ControlDescription) (s string) {
	return `ctrl.SetText(val)`
}

func (d Textbox) GenerateUpdate(ref interface{}, desc *generator.ControlDescription) (s string) {
	return `val := ctrl.Text()`
}
