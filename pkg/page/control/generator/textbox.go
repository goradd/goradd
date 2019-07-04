package generator

import (
	"fmt"
	"github.com/goradd/goradd/codegen/generator"
	"github.com/goradd/goradd/pkg/config"
	"github.com/goradd/goradd/pkg/orm/query"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/control"
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

func (d Textbox) Imports() []string {
	return []string{"github.com/goradd/goradd/pkg/page/control"}
}

func (d Textbox) SupportsColumn(col *generator.ColumnType) bool {
	if col.ColumnType == query.ColTypeBytes ||
		col.ColumnType == query.ColTypeString {
		return true
	}
	return false
}

func (d Textbox) GenerateCreate(namespace string, col *generator.ColumnType) (s string) {
	s = fmt.Sprintf(
		`	ctrl = %s.NewTextbox(c.ParentControl, id)
	ctrl.SetLabel(ctrl.T("%s"))
`, namespace, col.DefaultLabel)
	if col.MaxCharLength > 0 {
		s += fmt.Sprintf(`	ctrl.SetMaxLength(%d)	
`, col.MaxCharLength)
	}

	if col.IsPk {
		s += `	ctrl.SetDisabled(true)
`
	} else if !col.IsNullable {
		s += `	ctrl.SetIsRequired(true)
`
	}

	return
}

func (d Textbox) GenerateCreator(col *generator.ColumnType, connector page.DataConnector) page.Creator {
	creator := control.TextboxCreator{
		ID: col.ControlID,
		MaxLength: int(col.MaxCharLength),
	}

	creator.ControlOptions.Disabled = col.IsPk
	creator.ControlOptions.Required = !col.IsNullable
	creator.ControlOptions.DataConnector = connector
	return creator
}


func (d Textbox) GenerateRefresh(col *generator.ColumnType) (s string) {
	return `ctrl.SetText(val)`
}

func (d Textbox) GenerateUpdate(col *generator.ColumnType) (s string) {
	return `val := ctrl.Text()`
}

