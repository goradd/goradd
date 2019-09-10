package generator

import (
	"fmt"
	"github.com/goradd/goradd/codegen/generator"
	"github.com/goradd/goradd/pkg/config"
	"github.com/goradd/goradd/pkg/orm/query"
)

func init() {
	if !config.Release {
		generator.RegisterControlGenerator(Textbox{})
	}
}

// This structure describes the textbox to the connector dialog and code generator
type Textbox struct {
}

func (d Textbox) Type() string {
	return "Textbox"
}

func (d Textbox) Imports() []generator.ImportPath {
	return []generator.ImportPath{
		{Alias: "goraddctrl", Path:"github.com/goradd/goradd/pkg/page/control"},
	}
}

func (d Textbox) SupportsColumn(col *generator.ColumnType) bool {
	if col.ColumnType == query.ColTypeBytes ||
		col.ColumnType == query.ColTypeString {
		return true
	}
	return false
}

func (d Textbox) GenerateCreator(col *generator.ColumnType) (s string) {
	s = fmt.Sprintf(
		`goraddctrl.TextboxCreator{
			ID:        %#v,
			MaxLength: %d,
			ControlOptions: page.ControlOptions{
				IsRequired:      %#v,
				DataConnector: %s{},
			},
		}`, col.ControlID, col.MaxCharLength, !col.IsNullable, col.Connector)
	return
}



func (d Textbox) GenerateRefresh(col *generator.ColumnType) (s string) {
	return `ctrl.SetText(val)`
}

func (d Textbox) GenerateUpdate(col *generator.ColumnType) (s string) {
	return `val := ctrl.Text()`
}

