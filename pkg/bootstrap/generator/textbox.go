package generator

import (
	"fmt"
	generator2 "github.com/goradd/goradd/codegen/generator"
	"github.com/goradd/goradd/pkg/config"
	generator3 "github.com/goradd/goradd/pkg/page/control/generator"
)

func init() {
	if !config.Release {
		generator2.RegisterControlGenerator(Textbox{})
	}
}

// This structure describes the textbox to the connector dialog and code generator
type Textbox struct {
	generator3.Textbox // base it on the built-in generator
}

func (d Textbox) Imports() []generator2.ImportPath {
	return []generator2.ImportPath{
		{"bootstrapctrl", "github.com/goradd/goradd/pkg/bootstrap/control"},
	}
}

func (d Textbox) GenerateCreator(col *generator2.ColumnType) (s string) {
	s = fmt.Sprintf(
		`bootstrapctrl.TextboxCreator{
			ID:        %#v,
			MaxLength: %d,
			ControlOptions: page.ControlOptions{
				IsRequired:      %#v,
				DataConnector: %s{},
			},
		}`, col.ControlID, col.MaxCharLength, !col.IsNullable, col.Connector)
	return
}
