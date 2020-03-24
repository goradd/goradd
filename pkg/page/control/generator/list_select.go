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

func (d SelectList) Imports() []string {
	return []string{"strconv"}
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
	ID:           p.ID() + "-%s",
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
	var s string
	if col.IsNullable {
		s = `var val interface{}`
	} else {
		switch col.ColumnType {
		case query.ColTypeInteger:
			s = `var val int`
		case query.ColTypeUnsigned:
			s = `var val uint`
		case query.ColTypeInteger64:
			s = `var val int64`
		case query.ColTypeUnsigned64:
			s = `var val uint64`
		default:
			s = `var val string`
		}
	}
	s += `
	sv := ctrl.StringValue()
`
	var s2 string
	switch col.ColumnType {
	case query.ColTypeInteger:
		s2 = `val,_ = strconv.Atoi(sv)`
	case query.ColTypeUnsigned:
		s2 = `v2,_ := strconv.ParseUint(sv, 10, 0); val = uint(v2)`
	case query.ColTypeInteger64:
		s2 = `val,_ = strconv.ParseInt(sv, 10, 64)`
	case query.ColTypeUnsigned64:
		s2 = `val,_ = strconv.ParseUint(sv, 10, 64)`
	default:
		s2 = `val = sv`
	}

	if col.IsNullable {
		s += fmt.Sprintf(
`
if sv == "" {
	val = nil
} else {
	%s
}`, s2)
	} else {
		s += s2
	}

	return s
}

func (d SelectList) GenerateProvider(ref interface{}, desc *generator.ControlDescription) string {
	col := ref.(*db.Column)
	if col.ForeignKey.IsType {
		return fmt.Sprintf(`return model.All%sI()`, col.ForeignKey.GoTypePlural)
	} else {
		return fmt.Sprintf(`return model.Query%s(ctx).LoadI()`, col.ForeignKey.GoTypePlural)
	}
}
