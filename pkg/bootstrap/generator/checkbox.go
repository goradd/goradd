package generator

import (
	"fmt"
	generator2 "github.com/goradd/goradd/codegen/generator"
	"github.com/goradd/goradd/pkg/config"
	generator3 "github.com/goradd/goradd/pkg/page/control/generator"
)

func init() {
	if !config.Release {
		generator2.RegisterControlGenerator(Checkbox{})
	}
}

// This structure describes the textbox to the connector dialog and code generator
type Checkbox struct {
	generator3.Checkbox // base it on the built-in generator
}

func (d Checkbox) Imports() []generator2.ImportPath {
	return []generator2.ImportPath{
		{"bootstrapctrl", "github.com/goradd/goradd/pkg/bootstrap/control"},
	}
}

func (d Checkbox) GenerateCreator(col *generator2.ColumnType) (s string) {
	s = fmt.Sprintf(
		`bootstrapctrl.CheckboxCreator{
			ID:        %#v,
			ControlOptions: page.ControlOptions{
				IsRequired:      %#v,
				DataConnector: %s{},
			},
		}`, col.ControlID, !col.IsNullable, col.Connector)
	return
}
