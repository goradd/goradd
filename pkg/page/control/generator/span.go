package generator

import (
	"fmt"
	"github.com/goradd/goradd/codegen/generator"
	"github.com/goradd/goradd/pkg/config"
)

func init() {
	if !config.Release {
		generator.RegisterControlGenerator(Span{})
	}
}

// This structure describes the Span to the connector dialog and code generator
type Span struct {
}

func (d Span) Type() string {
	return "Span"
}

func (d Span) NewFunc() string {
	return "NewSpan"
}

func (d Span) Imports() []string {
	return []string{"github.com/goradd/goradd/pkg/page/control", "fmt"}
}

func (d Span) SupportsColumn(col *generator.ColumnType) bool {
	return true
}

func (d Span) GenerateCreator(col *generator.ColumnType) (s string) {
	s = fmt.Sprintf(
`control.SpanCreator{
	ID:        %#v,
	ControlOptions: page.ControlOptions{
		Disabled:	   %#v,
		Required:      %#v,
		DataConnector: %s{},
	},
}`, col.ControlID, col.IsPk, !col.IsNullable, col.Connector)
	return
}


func (d Span) GenerateRefresh(col *generator.ColumnType) (s string) {
	return `ctrl.SetText(val)`
}

func (d Span) GenerateUpdate(col *generator.ColumnType) (s string) {
	return ""
}

