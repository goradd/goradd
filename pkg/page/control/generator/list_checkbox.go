package generator

import (
	"fmt"
	"github.com/goradd/goradd/codegen/generator"
	"github.com/goradd/goradd/pkg/config"
	"github.com/goradd/goradd/pkg/orm/db"
	"github.com/goradd/goradd/pkg/orm/query"
)

func init() {
	if !config.Release {
		generator.RegisterControlGenerator(CheckboxList{})
	}
}

// This structure describes the CheckboxList to the connector dialog and code generator
type CheckboxList struct {
}

func (d CheckboxList) Type() string {
	return "CheckboxList"
}

func (d CheckboxList) NewFunc() string {
	return "NewCheckboxList"
}

func (d CheckboxList) Imports() []generator.ImportPath {
	return []generator.ImportPath{
		{Alias: "goraddctrl", Path:"github.com/goradd/goradd/pkg/page/control"},
	}
}

// TODO: This has to be changed to support virtual column types like ManyMany and Reverse
func (d CheckboxList) SupportsColumn(ref interface{}) bool {
	if col,ok := ref.(*db.Column); ok && col.ForeignKey != nil {
		return true
	}
	return false
}

func (d CheckboxList) GenerateCreator(ref interface{}, desc *generator.ControlDescription) (s string) {
	col := ref.(*db.Column)
	s = fmt.Sprintf(
`goraddctrl.CheckboxListCreator{
	ID:           %#v,
	DataProvider: p,
	ControlOptions: page.ControlOptions{
		IsRequired:      %#v,
		DataConnector: %s{},
	},
}`, desc.ControlID, !col.IsNullable, desc.Connector)
	return
}


func (d CheckboxList) GenerateRefresh(ref interface{}, desc *generator.ControlDescription) string {
	return `ctrl.SetValue(val)`
}

func (d CheckboxList) GenerateUpdate(ref interface{}, desc *generator.ControlDescription) string {
	col := ref.(*db.Column)
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

func (d CheckboxList) GenerateProvider(ref interface{}, desc *generator.ControlDescription) string {
	col := ref.(*db.Column)
	if col.ForeignKey.IsType {
		return fmt.Sprintf(`return model.%sI()`, col.ForeignKey.GoTypePlural)
	} else {
		return fmt.Sprintf(`return model.Query%s(ctx).LoadI(ctx)`, col.ForeignKey.GoTypePlural)
	}
}
