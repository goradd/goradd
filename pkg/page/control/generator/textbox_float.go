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
		generator.RegisterControlGenerator(FloatTextbox{})
	}
}

// This structure describes the FloatTextbox to the connector dialog and code generator
type FloatTextbox struct {
}

func (d FloatTextbox) Type() string {
	return "FloatTextbox"
}

func (d FloatTextbox) Imports() []generator.ImportPath {
	return []generator.ImportPath{
		{Alias: "goraddctrl", Path:"github.com/goradd/goradd/pkg/page/control"},
	}
}

func (d FloatTextbox) SupportsColumn(ref interface{}) bool {
	if col,ok := ref.(*db.Column); ok && (col.ColumnType == query.ColTypeFloat || col.ColumnType == query.ColTypeDouble) {
		return true
	}
	return false
}

func (d FloatTextbox) GenerateCreator(ref interface{}, desc *generator.ControlDescription) (s string) {
	col := ref.(*db.Column)
	s = fmt.Sprintf(
		`goraddctrl.FloatTextboxCreator{
			ID:        %#v,
			ControlOptions: page.ControlOptions{
				IsRequired:      %#v,
				DataConnector: %s{},
			},
		}`, desc.ControlID, !col.IsNullable, desc.Connector)
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
