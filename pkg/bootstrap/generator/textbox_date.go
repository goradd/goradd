package generator

import (
	"fmt"
	generator2 "github.com/goradd/goradd/codegen/generator"
	"github.com/goradd/goradd/pkg/config"
	"github.com/goradd/goradd/pkg/orm/db"
	generator3 "github.com/goradd/goradd/pkg/page/control/generator"
)

func init() {
	generator2.RegisterControlGenerator(DateTextbox{}, "github.com/goradd/goradd/pkg/bootstrap/control/DateTextbox")
}

// This structure describes the textbox to the connector dialog and code generator
type DateTextbox struct {
	generator3.DateTextbox // base it on the built-in generator
}

func (d DateTextbox) GenerateCreator(ref interface{}, desc *generator2.ControlDescription) (s string) {
	col := ref.(*db.Column)
	var format string
	if col.IsDateOnly {
		format = config.DefaultDateEntryFormat
	} else if col.IsTimeOnly {
		format = config.DefaultTimeEntryFormat
	} else {
		format = config.DefaultDateTimeEntryFormat
	}
	s = fmt.Sprintf(
		`%s.DateTextboxCreator{
	ID:        p.ID() + "-%s",
	Formats:    []string{%#v},
	ControlOptions: page.ControlOptions{
		IsRequired:      %#v,
		DataConnector: %s{},
	},
}`, desc.Package, desc.ControlID, format, !col.IsNullable, desc.Connector)
	return
}
