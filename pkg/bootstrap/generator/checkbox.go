package generator

import (
	"fmt"
	generator2 "github.com/goradd/goradd/codegen/generator"
	"github.com/goradd/goradd/pkg/orm/db"
	generator3 "github.com/goradd/goradd/pkg/page/control/generator"
)

func init() {
	generator2.RegisterControlGenerator(Checkbox{}, "github.com/goradd/goradd/pkg/bootstrap/control/Checkbox")
}

// This structure describes the textbox to the connector dialog and code generator
type Checkbox struct {
	generator3.Checkbox // base it on the built-in generator
}

func (d Checkbox) GenerateCreator(ref interface{}, desc *generator2.ControlDescription) (s string) {
	col := ref.(*db.Column)
	s = fmt.Sprintf(
		`%s.CheckboxCreator{
			ID:        %#v,
			ControlOptions: page.ControlOptions{
				IsRequired:      %#v,
				DataConnector: %s{},
			},
		}`, desc.Import, desc.ControlID, !col.IsNullable, desc.Connector)
	return
}
