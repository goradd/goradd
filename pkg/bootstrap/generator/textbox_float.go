package generator

import (
	"fmt"
	generator2 "github.com/goradd/goradd/codegen/generator"
	"github.com/goradd/goradd/pkg/orm/db"
	generator3 "github.com/goradd/goradd/pkg/page/control/generator"
)

func init() {
	generator2.RegisterControlGenerator(FloatTextbox{}, "github.com/goradd/goradd/pkg/bootstrap/control/Float")
}

// This structure describes the textbox to the connector dialog and code generator
type FloatTextbox struct {
	generator3.FloatTextbox // base it on the built-in generator
}

func (d FloatTextbox) GenerateCreator(ref interface{}, desc *generator2.ControlDescription) (s string) {
	col := ref.(*db.Column)

	s = fmt.Sprintf(
		`%s.FloatTextboxCreator{
			ID:        p.ID() + "-%s",
			ControlOptions: page.ControlOptions{
				IsRequired:      %#v,
				DataConnector: %s{},
			},
		}`, desc.Package, desc.ControlID, !col.IsNullable, desc.Connector)
	return
}
