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
		generator2.RegisterControlGenerator(PhoneTextbox{})
	}
}

// This structure describes the PhoneTextbox to the connector dialog and code generator
type PhoneTextbox struct {
	generator3.PhoneTextbox // base it on the built-in generator
}

func (d PhoneTextbox) Type() string {
	return "github.com/goradd/goradd/pkg/bootstrap/control/PhoneTextbox"
}

func (d PhoneTextbox) GenerateCreator(ref interface{}, desc *generator2.ControlDescription) (s string) {
	col := ref.(*db.Column)
	s = fmt.Sprintf(
		`%s.PhoneTextboxCreator{
			ID:        %#v,
			ControlOptions: page.ControlOptions{
				IsRequired:      %#v,
				DataConnector: %s{},
			},
		}`, desc.Import, desc.ControlID, !col.IsNullable, desc.Connector)
	return
}
