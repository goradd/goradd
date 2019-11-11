package generator

import (
	"fmt"
	"github.com/goradd/goradd/codegen/generator"
	"github.com/goradd/goradd/pkg/config"
	"github.com/goradd/goradd/pkg/orm/query"
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

func (d SelectList) Imports() []generator.ImportPath {
	return []generator.ImportPath{
		{Alias: "goraddctrl", Path:"github.com/goradd/goradd/pkg/page/control"},
	}
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
`goraddctrl.SelectListCreator{
	ID:           %#v,
	DataProvider: p,
	ControlOptions: page.ControlOptions{
		IsRequired:      %#v,
		DataConnector: %s{},
	},
}`, col.ControlID, !col.IsNullable, col.Connector)
	return
}


func (d SelectList) GenerateRefresh(col *generator.ColumnType) string {
	return `ctrl.SetValue(val)`
}

func (d SelectList) GenerateUpdate(col *generator.ColumnType) string {
	var s string
	switch col.ColumnType {
	case query.ColTypeInteger:
		s = `val,_ := strconv.Atoi(ctrl.StringValue())`
	case query.ColTypeUnsigned:
		s = `val,_ := strconv.ParseUint(ctrl.StringValue(), 10, 0)`
	case query.ColTypeInteger64:
		s = `val,_ := strconv.ParseInt(ctrl.StringValue(), 10, 64)`
	case query.ColTypeUnsigned64:
		s = `val,_ := strconv.ParseUint(ctrl.StringValue(), 10, 64)`
	default:
		s = `val := ctrl.StringValue()`
	}

	return s
}

func (d SelectList) GenerateProvider(col *generator.ColumnType) string {
	if col.ForeignKey.IsType {
		return fmt.Sprintf(`return model.%sI()`, col.ForeignKey.GoTypePlural)
	} else {
		return fmt.Sprintf(`return model.Query%s(ctx).LoadI(ctx)`, col.ForeignKey.GoTypePlural)
	}
}
