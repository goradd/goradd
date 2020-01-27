package generator

import (
	"fmt"
	"github.com/goradd/goradd/codegen/generator"
	"github.com/goradd/goradd/pkg/orm/db"
	"github.com/goradd/goradd/pkg/orm/query"
)

func init() {
	generator.RegisterControlGenerator(DateTimeSpan{}, "github.com/goradd/goradd/pkg/page/control/DateTimeSpan")
}

// This structure describes the DateTimeSpan to the connector dialog and code generator
type DateTimeSpan struct {
}

func (d DateTimeSpan) NewFunc() string {
	return "NewDateTimeSpan"
}


func (d DateTimeSpan) SupportsColumn(ref interface{}) bool {
	if col,ok := ref.(*db.Column); ok &&
		col.ColumnType == query.ColTypeDateTime {
		return true
	}
	return false
}

func (d DateTimeSpan) GenerateCreator(ref interface{}, desc *generator.ControlDescription) (s string) {
	col := ref.(*db.Column)
	s = fmt.Sprintf(
`%s.DateTimeSpanCreator{
	ID:        %#v,
	ControlOptions: page.ControlOptions{
		IsDisabled:	   %#v,
		DataConnector: %s{},
	},
}`, desc.Package, desc.ControlID, col.IsPk, desc.Connector)
	return
}


func (d DateTimeSpan) GenerateRefresh(ref interface{}, desc *generator.ControlDescription) (s string) {
	return `ctrl.SetDateTime(val)`
}

func (d DateTimeSpan) GenerateUpdate(ref interface{}, desc *generator.ControlDescription) (s string) {
	return ""
}

