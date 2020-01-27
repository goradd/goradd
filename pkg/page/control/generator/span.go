package generator

import (
	"fmt"
	"github.com/goradd/goradd/codegen/generator"
	"github.com/goradd/goradd/pkg/orm/db"
)

func init() {
	generator.RegisterControlGenerator(Span{}, "github.com/goradd/goradd/pkg/page/control/Span")
}

// This structure describes the Span to the connector dialog and code generator
type Span struct {
}

func (d Span) Imports() []string {
	return []string{"fmt"}
}

func (d Span) SupportsColumn(ref interface{}) bool {
	return true
}

func (d Span) GenerateCreator(ref interface{}, desc *generator.ControlDescription) (s string) {
	col := ref.(*db.Column)
	s = fmt.Sprintf(
`%s.SpanCreator{
	ID:        %#v,
	ControlOptions: page.ControlOptions{
		IsDisabled:	   %#v,
		IsRequired:      %#v,
		DataConnector: %s{},
	},
}`, desc.Import, desc.ControlID, col.IsPk, !col.IsNullable, desc.Connector)
	return
}


func (d Span) GenerateRefresh(ref interface{}, desc *generator.ControlDescription) (s string) {
	return `ctrl.SetText(fmt.Sprint(val))`
}

func (d Span) GenerateUpdate(ref interface{}, desc *generator.ControlDescription) (s string) {
	return ""
}

