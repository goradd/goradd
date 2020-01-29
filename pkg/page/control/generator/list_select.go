package generator

import (
	"fmt"
	"github.com/goradd/goradd/codegen/generator"
	"github.com/goradd/goradd/pkg/orm/db"
	"github.com/goradd/goradd/pkg/orm/query"
)

func init() {
	generator.RegisterControlGenerator(SelectList{}, "github.com/goradd/goradd/pkg/page/control/SelectList")
}

// This structure describes the SelectList to the connector dialog and code generator
type SelectList struct {
}


func (d SelectList) SupportsColumn(ref interface{}) bool {
	if col,ok := ref.(*db.Column); ok && col.ForeignKey != nil {
		return true
	}
	return false
}

func (d SelectList) GenerateCreator(ref interface{}, desc *generator.ControlDescription) (s string) {
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


func (d SelectList) GenerateRefresh(ref interface{}, desc *generator.ControlDescription) string {
	return `ctrl.SetValue(val)`
}

func (d SelectList) GenerateUpdate(ref interface{}, desc *generator.ControlDescription) string {
	col := ref.(*db.Column)
	s1 :=  `
sv := ctrl.StringValue()
`
	var s string
	switch col.ColumnType {
	case query.ColTypeInteger:
		s = `val,_ = strconv.Atoi(sv)`
	case query.ColTypeUnsigned:
		s = `val,_ = strconv.ParseUint(sv, 10, 0)`
	case query.ColTypeInteger64:
		s = `val,_ = strconv.ParseInt(sv, 10, 64)`
	case query.ColTypeUnsigned64:
		s = `val,_ = strconv.ParseUint(sv, 10, 64)`
	default:
		s = `val = sv`
	}

	if col.IsNullable {
		s = fmt.Sprintf(
`
var val interface{}
if sv == "" {
	val = nil
} else {
	%s
}`, s)
	} else {
	s =	fmt.Sprintf(`
var val string
%s
`, s)
	}

	return s1 + s
}

func (d SelectList) GenerateProvider(ref interface{}, desc *generator.ControlDescription) string {
	col := ref.(*db.Column)
	if col.ForeignKey.IsType {
		return fmt.Sprintf(`return model.%sI()`, col.ForeignKey.GoTypePlural)
	} else {
		return fmt.Sprintf(`return model.Query%s(ctx).LoadI(ctx)`, col.ForeignKey.GoTypePlural)
	}
}
