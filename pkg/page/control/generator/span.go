package generator

import (
	"fmt"
	"github.com/goradd/goradd/codegen/generator"
	"github.com/goradd/goradd/pkg/config"
	"github.com/goradd/goradd/pkg/orm/db"
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

func (d Span) Imports() []generator.ImportPath {
	return []generator.ImportPath{
		{Alias: "goraddctrl", Path:"github.com/goradd/goradd/pkg/page/control"},
		{Alias: "", Path:"fmt"},
	}
}

func (d Span) SupportsColumn(ref interface{}) bool {
	return true
}

func (d Span) GenerateCreator(ref interface{}, desc *generator.ControlDescription) (s string) {
	col := ref.(*db.Column)
	s = fmt.Sprintf(
`goraddctrl.SpanCreator{
	ID:        %#v,
	ControlOptions: page.ControlOptions{
		IsDisabled:	   %#v,
		IsRequired:      %#v,
		DataConnector: %s{},
	},
}`, desc.ControlID, col.IsPk, !col.IsNullable, desc.Connector)
	return
}


func (d Span) GenerateRefresh(ref interface{}, desc *generator.ControlDescription) (s string) {
	return `ctrl.SetText(fmt.Sprintf("%v", val))`
}

func (d Span) GenerateUpdate(ref interface{}, desc *generator.ControlDescription) (s string) {
	return ""
}

