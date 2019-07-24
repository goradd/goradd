package generator

import (
	"fmt"
	"github.com/goradd/goradd/codegen/generator"
	"github.com/goradd/goradd/pkg/config"
)

func init() {
	if !config.Release {
		generator.RegisterControlGenerator(DateTimeSpan{})
	}
}

// This structure describes the DateTimeSpan to the connector dialog and code generator
type DateTimeSpan struct {
}

func (d DateTimeSpan) Type() string {
	return "DateTimeSpan"
}

func (d DateTimeSpan) NewFunc() string {
	return "NewDateTimeSpan"
}

func (d DateTimeSpan) Imports() []string {
	return []string{"github.com/goradd/goradd/pkg/page/control"}
}

func (d DateTimeSpan) SupportsColumn(col *generator.ColumnType) bool {
	return true
}

func (d DateTimeSpan) GenerateCreator(col *generator.ColumnType) (s string) {
	s = fmt.Sprintf(
`control.DateTimeSpanCreator{
	ID:        %#v,
	ControlOptions: page.ControlOptions{
		Disabled:	   %#v,
		DataConnector: %s{},
	},
}`, col.ControlID, col.IsPk, col.Connector)
	return
}


func (d DateTimeSpan) GenerateRefresh(col *generator.ColumnType) (s string) {
	return `ctrl.SetDateTime(val)`
}

func (d DateTimeSpan) GenerateUpdate(col *generator.ColumnType) (s string) {
	return ""
}

