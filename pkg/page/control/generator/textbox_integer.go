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
		generator.RegisterControlGenerator(IntegerTextbox{})
	}
}

// This structure describes the IntegerTextbox to the connector dialog and code generator
type IntegerTextbox struct {
}

func (d IntegerTextbox) Type() string {
	return "github.com/goradd/goradd/pkg/page/control/IntegerTextbox"
}

func (d IntegerTextbox) Imports() []generator.ImportPath {
	return []generator.ImportPath{
		{Alias: "goraddctrl", Path:"github.com/goradd/goradd/pkg/page/control"},
	}
}

func (d IntegerTextbox) SupportsColumn(ref interface{}) bool {
	if col,ok := ref.(*db.Column); ok &&
		(col.ColumnType == query.ColTypeInteger ||
		col.ColumnType == query.ColTypeInteger64) &&
		!col.IsReference() {
		return true
	}
	return false
}

func (d IntegerTextbox) GenerateCreator(ref interface{}, desc *generator.ControlDescription) (s string) {
	col := ref.(*db.Column)
	s = fmt.Sprintf(
`%s.IntegerTextboxCreator{
	ID:        %#v,
	ControlOptions: page.ControlOptions{
		IsRequired:      %#v,
		DataConnector: %s{},
	},
}`, desc.Import, desc.ControlID, !col.IsNullable, desc.Connector)
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
