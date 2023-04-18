package generator

import (
	"fmt"
	"github.com/goradd/goradd/codegen/generator"
	"github.com/goradd/goradd/pkg/orm/db"
	"github.com/goradd/goradd/pkg/orm/query"
)

func init() {
	generator.RegisterControlGenerator(IntegerTextbox{}, "github.com/goradd/goradd/pkg/page/control/textbox/IntegerTextbox")
}

// IntegerTextbox describes the IntegerTextbox to the connector dialog and code generator
type IntegerTextbox struct {
}

func (d IntegerTextbox) SupportsColumn(ref interface{}) bool {
	if col, ok := ref.(*db.Column); ok &&
		(col.ColumnType == query.ColTypeInteger ||
			col.ColumnType == query.ColTypeInteger64) &&
		!col.IsReference() {
		return true
	}
	return false
}

func (d IntegerTextbox) GenerateCreator(ref interface{}, desc *generator.ControlDescription) (s string) {
	col := ref.(*db.Column)
	sMinVal := fmt.Sprintf("%v", col.MinValue)
	sMaxVal := fmt.Sprintf("%v", col.MaxValue)
	s = fmt.Sprintf(
		`%s.IntegerTextboxCreator{
	ID:        p.ID() + "-%s",
`, desc.Package, desc.ControlID)
	s += `    // Set this with a "min" value in the column comment. For example: {"min":100}
    MinValue: &textbox.IntegerLimit{
		Value: ` + sMinVal + `,
		InvalidMessage: fmt.Sprintf(p.GT("Must be at least %d"),` + sMinVal + `),
	},
    // Set this with a "max" value in the column comment. For example: {"max":1000}
	MaxValue: &textbox.IntegerLimit{
		Value: ` + sMaxVal + `,
		InvalidMessage: fmt.Sprintf(p.GT("Must be at most %d"), ` + sMaxVal + `),
	},
`
	s += fmt.Sprintf(`	ControlOptions: page.ControlOptions{
		IsRequired:      %#v,
		DataConnector: %s{},
	},
}`, !col.IsNullable, desc.Connector)
	return
}

func (d IntegerTextbox) GenerateRefresh(ref interface{}, desc *generator.ControlDescription) (s string) {
	return `ctrl.SetValue(val)`
}

func (d IntegerTextbox) GenerateUpdate(ref interface{}, desc *generator.ControlDescription) (s string) {
	col := ref.(*db.Column)
	switch col.ColumnType {
	case query.ColTypeInteger:
		return `val := ctrl.Int()`
	case query.ColTypeInteger64:
		return `val := int64(ctrl.Int())`
	case query.ColTypeUnsigned:
		return `val := uint(ctrl.Int())`
	case query.ColTypeUnsigned64:
		return `val := uint64(ctrl.Int())`
	}
	panic("not a compatible column type")
}

func (d IntegerTextbox) GenerateModifies(ref interface{}, desc *generator.ControlDescription) (s string) {
	return `val != ctrl.Value()`
}
