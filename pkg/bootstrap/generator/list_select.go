package generator

import (
	"fmt"
	generator2 "github.com/goradd/goradd/codegen/generator"
	"github.com/goradd/goradd/pkg/orm/db"
	generator3 "github.com/goradd/goradd/pkg/page/control/generator"
)

func init() {
	generator2.RegisterControlGenerator(SelectList{}, "github.com/goradd/goradd/pkg/bootstrap/control/SelectList")
}

// This structure describes the textbox to the connector dialog and code generator
type SelectList struct {
	generator3.SelectList // base it on the built-in generator
}

func (d SelectList) GenerateCreator(ref interface{}, desc *generator2.ControlDescription) (s string) {
	col := ref.(*db.Column)
	s = fmt.Sprintf(
		`%s.SelectListCreator{
	ID:           %#v,
	DataProvider: p,
	ControlOptions: page.ControlOptions{
		IsRequired:      %#v,
		DataConnector: %s{},
	},
}`, desc.Package, desc.ControlID, !col.IsNullable, desc.Connector)
	return
}
