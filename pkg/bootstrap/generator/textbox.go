package generator

import (
	"fmt"
	generator2 "github.com/goradd/goradd/codegen/generator"
	"github.com/goradd/goradd/pkg/orm/db"
	generator3 "github.com/goradd/goradd/pkg/page/control/generator"
)

func init() {
	generator2.RegisterControlGenerator(Textbox{}, "github.com/goradd/goradd/pkg/bootstrap/control/Textbox")
}

// This structure describes the textbox to the connector dialog and code generator
type Textbox struct {
	generator3.Textbox // base it on the built-in generator
}

func (d Textbox) GenerateCreator(ref interface{}, desc *generator2.ControlDescription) (s string) {
	col := ref.(*db.Column)
	s = fmt.Sprintf(
		`%s.TextboxCreator{
			ID:        %#v,
			MaxLength: %d,
			ControlOptions: page.ControlOptions{
				IsRequired:      %#v,
				DataConnector: %s{},
			},
		}`, desc.Import, desc.ControlID, col.MaxCharLength, !col.IsNullable, desc.Connector)
	return
}
