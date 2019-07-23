package generator

import (
	"fmt"
	"github.com/goradd/goradd/codegen/generator"
	"github.com/goradd/goradd/pkg/config"
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
	return "IntegerTextbox"
}

func (d IntegerTextbox) NewFunc() string {
	return "NewIntegerTextbox"
}

func (d IntegerTextbox) Imports() []string {
	return []string{"github.com/goradd/goradd/pkg/page/control"}
}

func (d IntegerTextbox) SupportsColumn(col *generator.ColumnType) bool {
	if (col.ColumnType == query.ColTypeInteger ||
		col.ColumnType == query.ColTypeInteger64) &&
		!col.IsReference() {
		return true
	}
	return false
}

func (d IntegerTextbox) GenerateCreator(col *generator.ColumnType) (s string) {
	s = fmt.Sprintf(
`control.IntegerTextboxCreator{
	ID:        %#v,
	MaxLength: %#v,
	ControlOptions: page.ControlOptions{
		Required:      %#v,
		DataConnector: %s{},
	},
}`, col.ControlID, col.MaxCharLength, !col.IsNullable, col.Connector)
	return
}


func (d IntegerTextbox) GenerateRefresh(col *generator.ColumnType) (s string) {
	return `ctrl.SetValue(val)`
}

func (d IntegerTextbox) GenerateUpdate(col *generator.ColumnType) (s string) {
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
