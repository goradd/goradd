package generator

import (
	"fmt"
	"github.com/goradd/goradd/codegen/generator"
)

func init() {
	generator.RegisterControlGenerator(Span{}, "github.com/goradd/goradd/pkg/page/control/Span")
}

// Span describes the Span to the connector dialog and code generator
type Span struct {
}

func (d Span) Imports() []string {
	return []string{"fmt"}
}

func (d Span) SupportsColumn(ref interface{}) bool {
	return true
}

func (d Span) GenerateCreator(ref interface{}, desc *generator.ControlDescription) (s string) {
	s = fmt.Sprintf(
		`%s.SpanCreator{
	ID:        p.ID() + "-%s",
	ControlOptions: page.ControlOptions{
		DataConnector: %s{},
	},
}`, desc.Package, desc.ControlID, desc.Connector)
	return
}

func (d Span) GenerateRefresh(ref interface{}, desc *generator.ControlDescription) (s string) {
	return `ctrl.SetText(fmt.Sprint(val))`
}

func (d Span) GenerateUpdate(ref interface{}, desc *generator.ControlDescription) (s string) {
	return ""
}

func (d Span) GenerateModifies(ref interface{}, desc *generator.ControlDescription) (s string) {
	return ""
}
