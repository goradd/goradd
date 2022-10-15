package generator

import (
	"fmt"
	generator2 "github.com/goradd/goradd/codegen/generator"
	"github.com/goradd/goradd/pkg/orm/db"
	generator3 "github.com/goradd/goradd/pkg/page/control/generator"
)

func init() {
	generator2.RegisterControlGenerator(EmailTextbox{}, "github.com/goradd/goradd/pkg/bootstrap/control/EmailTextbox")
}

// EmailTextbox describes the EmailTextbox to the connector dialog and code generator
type EmailTextbox struct {
	generator3.EmailTextbox // base it on the built-in generator
}

func (d EmailTextbox) GenerateCreator(ref interface{}, desc *generator2.ControlDescription) (s string) {
	col := ref.(*db.Column)
	s = fmt.Sprintf(
		`%s.EmailTextboxCreator{
			ID:        p.ID() + "-%s",
			ControlOptions: page.ControlOptions{
				IsRequired:      %#v,
				DataConnector: %s{},
			},
		}`, desc.Package, desc.ControlID, !col.IsNullable, desc.Connector)
	return
}
