package generator

import (
	"fmt"
	generator2 "github.com/goradd/goradd/codegen/generator"
	"github.com/goradd/goradd/pkg/config"
	"github.com/goradd/goradd/pkg/orm/db"
	generator3 "github.com/goradd/goradd/pkg/page/control/generator"
)

func init() {
	if !config.Release {
		generator2.RegisterControlGenerator(IntegerTextbox{})
	}
}

// This structure describes the textbox to the connector dialog and code generator
type IntegerTextbox struct {
	generator3.IntegerTextbox // base it on the built-in generator
}

func (d IntegerTextbox) Imports() []generator2.ImportPath {
	return []generator2.ImportPath{
		{"bootstrapctrl", "github.com/goradd/goradd/pkg/bootstrap/control"},
	}
}

func (d IntegerTextbox) GenerateCreator(ref interface{}, desc *generator2.ControlDescription) (s string) {
	col := ref.(*db.Column)
	s = fmt.Sprintf(
		`bootstrapctrl.IntegerTextboxCreator{
	ID:        %#v,
	ControlOptions: page.ControlOptions{
		IsRequired:      %#v,
		DataConnector: %s{},
	},
}`, desc.ControlID, !col.IsNullable, desc.Connector)
	return
}
