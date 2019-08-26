package generator

import (
	"fmt"
	"github.com/goradd/goradd/codegen/generator"
	"github.com/goradd/goradd/pkg/config"
	"github.com/goradd/goradd/pkg/orm/query"
)

func init() {
	if !config.Release {
		generator.RegisterControlGenerator(FloatTextbox{})
	}
}

// This structure describes the FloatTextbox to the connector dialog and code generator
type FloatTextbox struct {
}

func (d FloatTextbox) Type() string {
	return "FloatTextbox"
}

func (d FloatTextbox) Imports() []string {
	return []string{"github.com/goradd/goradd/pkg/page/control"}
}

func (d FloatTextbox) SupportsColumn(col *generator.ColumnType) bool {
	if col.ColumnType == query.ColTypeFloat || col.ColumnType == query.ColTypeDouble {
		return true
	}
	return false
}

func (d FloatTextbox) GenerateCreator(col *generator.ColumnType) (s string) {
	s = fmt.Sprintf(
		`control.IntegerTextboxCreator{
			ID:        %#v,
			MaxLength: %#v,
			ControlOptions: page.ControlOptions{
				IsRequired:      %#v,
				DataConnector: %s{},
			},
		}`, col.ControlID, col.MaxCharLength, !col.IsNullable, col.Connector)
	return
}


func (d FloatTextbox) GenerateRefresh(col *generator.ColumnType) (s string) {
	return `ctrl.SetValue(val)`
}

func (d FloatTextbox) GenerateUpdate(col *generator.ColumnType) (s string) {
	if col.ColumnType == query.ColTypeFloat {
		return `val := ctrl.Float32()`
	} else {
		return `val := ctrl.Float64()`
	}
}
