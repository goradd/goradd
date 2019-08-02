package generator

import (
	"fmt"
	"github.com/goradd/goradd/codegen/generator"
	"github.com/goradd/goradd/pkg/config"
)

func init() {
	if !config.Release {
		generator.RegisterControlGenerator(SelectList{})
	}
}

// This structure describes the SelectList to the connector dialog and code generator
type SelectList struct {
}

func (d SelectList) Type() string {
	return "SelectList"
}

func (d SelectList) NewFunc() string {
	return "NewSelectList"
}

func (d SelectList) Imports() []string {
	return []string{"github.com/goradd/goradd/pkg/page/control"}
}

// TODO: This has to be changed to support virtual column types like ManyMany and Reverse
func (d SelectList) SupportsColumn(col *generator.ColumnType) bool {
	if col.ForeignKey != nil {
		return true
	}
	return false
}

func (d SelectList) GenerateCreator(col *generator.ColumnType) (s string) {
	s = fmt.Sprintf(
`control.SelectListCreator{
	ID:           %#v,
	DataProvider: p.ID(),
	ControlOptions: page.ControlOptions{
		IsRequired:      %#v,
		DataConnector: %s{},
	},
}`, col.ControlID, !col.IsNullable, col.Connector)
	return
}


func (d SelectList) GenerateRefresh(col *generator.ColumnType) string {
	if col.ForeignKey.IsType {
		return `ctrl.SetValue(int(val))`
	} else {
		return `ctrl.SetValue(string(val))` // should be a string id
	}
}

func (d SelectList) GenerateUpdate(col *generator.ColumnType) string {
	if col.ForeignKey.IsType {
		return `val := ctrl.IntValue()`
	} else {
		return `val := ctrl.StringValue()`
	}
}

func (d SelectList) GenerateProvider(col *generator.ColumnType) string {
	if col.ForeignKey.IsType {
		return fmt.Sprintf(`return model.%sI()`, col.ForeignKey.GoTypePlural)
	} else {
		return fmt.Sprintf(`return model.Query%s(ctx).LoadI(ctx)`, col.ForeignKey.GoTypePlural)
	}
}
