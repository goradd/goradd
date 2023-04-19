package generator

import (
	"fmt"
	"github.com/goradd/goradd/codegen/generator"
	"github.com/goradd/goradd/pkg/orm/db"
)

func init() {
	generator.RegisterControlGenerator(CheckboxList{}, "github.com/goradd/goradd/pkg/page/control/list/CheckboxList")
}

// CheckboxList describes the CheckboxList to the connector dialog and code generator
type CheckboxList struct {
}

func (d CheckboxList) NewFunc() string {
	return "NewCheckboxList"
}

func (d CheckboxList) SupportsColumn(ref interface{}) bool {
	if rr, ok := ref.(*db.ReverseReference); ok {
		if rr.IsUnique() {
			return false
		}
		return true
	}
	if _, ok := ref.(*db.ManyManyReference); ok {
		return true
	}

	return false
}

func (d CheckboxList) GenerateCreator(ref interface{}, desc *generator.ControlDescription) (s string) {
	s = fmt.Sprintf(
		`%s.CheckboxListCreator{
	ID:           p.ID() + "-%s",
	DataProvider: p,
	ControlOptions: page.ControlOptions{
		DataConnector: %s{},
	},
}`, desc.Package, desc.ControlID, desc.Connector)
	return
}

func (d CheckboxList) GenerateRefresh(ref interface{}, desc *generator.ControlDescription) string {
	switch ref.(type) {
	case *db.ReverseReference:
		return `
			var values []string
			for _,obj := range objects {
				values = append(values, fmt.Sprint(obj.PrimaryKey()))
			}
			ctrl.SetSelectedValues(values)`
	case *db.ManyManyReference:
		return ``
	}
	return ``
}

func (d CheckboxList) GenerateUpdate(ref interface{}, desc *generator.ControlDescription) string {
	/*	switch col := ref.(type) {
		case *db.ReverseReference:
			return fmt.Sprintf(`
				values := ctrl.SelectedValues()
				model.Unasso
				`,
				col.GoPlural)
		case *db.ManyManyReference:
			return fmt.Sprintf(`
				values := []string
				for _,obj := range model.Load%s(ctx) {
					values = append(values, obj.PrimaryKey())
				}
				ctrl.SetSelectedValues(values)`,
				col.GoPlural)
		}
	*/
	return ``
}

func (d CheckboxList) GenerateModifies(ref interface{}, desc *generator.ControlDescription) string {
	/*	switch ref.(type) {
		case *db.ReverseReference:
			return `
				var values []string
				for _,obj := range objects {
					values = append(values, fmt.Sprint(obj.PrimaryKey()))
				}
				ctrl.SetSelectedValues(values)`
		case *db.ManyManyReference:
			return `false`
		}*/
	return ``
}

func (d CheckboxList) GenerateProvider(ref interface{}, desc *generator.ControlDescription) string {
	switch col := ref.(type) {
	case *db.ReverseReference:
		return fmt.Sprintf(`return model.Query%s(ctx).LoadI()`, col.AssociatedTable.GoPlural)
	case *db.ManyManyReference:
		return fmt.Sprintf(`return model.Query%s(ctx).LoadI()`, col.AssociatedTableName)
	}
	return ``
}
