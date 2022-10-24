package generator

import (
	"fmt"
	"github.com/goradd/goradd/codegen/generator"
	"github.com/goradd/goradd/pkg/orm/db"
	"github.com/goradd/goradd/pkg/orm/query"
)

func init() {
	generator.RegisterControlGenerator(FloatTextbox{}, "github.com/goradd/goradd/pkg/page/control/textbox/FloatTextbox")
}

// FloatTextbox describes the FloatTextbox to the connector dialog and code generator
type FloatTextbox struct {
}

func (d FloatTextbox) SupportsColumn(ref interface{}) bool {
	if col, ok := ref.(*db.Column); ok && (col.ColumnType == query.ColTypeFloat || col.ColumnType == query.ColTypeDouble) {
		return true
	}
	return false
}

func (d FloatTextbox) GenerateCreator(ref interface{}, desc *generator.ControlDescription) (s string) {
	col := ref.(*db.Column)
	sMinVal := fmt.Sprintf("%v", col.MinValue)
	sMaxVal := fmt.Sprintf("%v", col.MaxValue)

	s = fmt.Sprintf(
		`%s.FloatTextboxCreator{
			ID:        p.ID() + "-%s",
`, desc.Package, desc.ControlID)

	s += `    // Set this with a "min" value in the column comment. For example: {"min":1.00}
    MinValue: &textbox.FloatLimit{
		Value: ` + sMinVal + `,
		InvalidMessage: fmt.Sprintf(p.GT("Must be at least %g"),` + sMinVal + `),
	},
    // Set this with a "max" value in the column comment. For example: {"max":10.00}
	MaxValue: &textbox.FloatLimit{
		Value: ` + sMaxVal + `,
		InvalidMessage: fmt.Sprintf(p.GT("Must be at most %g"), ` + sMaxVal + `),
	},		
`
	s += fmt.Sprintf(`
			ControlOptions: page.ControlOptions{
				IsRequired:      %#v,
				DataConnector: %s{},
			},
		}`, !col.IsNullable, desc.Connector)
	return
}

func (d FloatTextbox) GenerateRefresh(ref interface{}, desc *generator.ControlDescription) (s string) {
	return `ctrl.SetValue(val)`
}

func (d FloatTextbox) GenerateUpdate(ref interface{}, desc *generator.ControlDescription) (s string) {
	col := ref.(*db.Column)
	if col.ColumnType == query.ColTypeFloat {
		return `val := ctrl.Float32()`
	} else {
		return `val := ctrl.Float64()`
	}
}
