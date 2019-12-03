package generator

import (
	"github.com/goradd/goradd/pkg/orm/db"
	strings2 "github.com/goradd/goradd/pkg/strings"
	"github.com/knq/snaker"
	"strings"
)

// matchColumnsWithControls maps controls to control descriptions, and returns the imports required by the
// control descriptions
func matchColumnsWithControls(t *db.Table, descriptions map[interface{}]*ControlDescription, aliasToImport map[string]*ImportType) {
	for _, col := range t.Columns {
		typ, newFunc, importName := defaultControlType(col)

		if typ != "" {
			var mainImport *ImportType

			generator := GetControlGenerator(importName, typ)
			if generator != nil {
				for i, importPath := range generator.Imports() {
					var ok bool
					var imp *ImportType
					if imp, ok = aliasToImport[importPath.Alias]; ok {
						if imp.Path != importPath.Path {
							panic("found the same alias with different import path")
						}
					} else {
						imp = &ImportType {
							importPath.Path,
							importPath.Alias,
							i == 0,
						}

						aliasToImport[importPath.Alias] = imp
					}
					if mainImport == nil {
						mainImport = imp
					}
				}
			}

			// TODO: Get this from a database comment if provided
			var defaultLabel string
			var controlName string

			if col.ForeignKey != nil {
				defaultLabel = strings2.Title(col.ForeignKey.GoName)
				controlName = col.ForeignKey.GoName + typ
			} else {
				defaultLabel = strings2.Title(col.DbName)
				controlName = col.GoName + typ
			}

			var defaultID string
			defaultID = strings.Replace(t.DbName, "_", "-", -1) + "-" + strings.Replace(col.DbName, "_", "-", -1)

			cd := ControlDescription{
				Import: mainImport,
				ControlType: typ,
				NewControlFunc: newFunc,
				ControlName: controlName,
				ControlID: defaultID,
				DefaultLabel: defaultLabel,
				Generator: generator,
				Connector:t.GoName + controlName + "Connector",
			}
			descriptions[col] = &cd
		}
	}

	return
}

// defaultControlType returns the default type of control for a column. ControlBase types can be customized in other ways too.
func defaultControlType(ref interface{}) (typ string, createFunc string, importName string) {
	d := DefaultControlTypeFunc(ref)
	return d.Typ, d.CreateFunc, d.ImportName
}

func matchReverseReferencesWithControls(t *db.Table, descriptions map[interface{}]*ControlDescription, aliasToImport map[string]*ImportType) {
	for _, rr := range t.ReverseReferences {
		typ, newFunc, importName := defaultControlType(rr)

		if typ != "" {
			var mainImport *ImportType

			generator := GetControlGenerator(importName, typ)
			if generator != nil {
				for i, importPath := range generator.Imports() {
					var ok bool
					var imp *ImportType
					if imp, ok = aliasToImport[importPath.Alias]; ok {
						if imp.Path != importPath.Path {
							panic("found the same alias with different import path")
						}
					} else {
						imp = &ImportType {
							importPath.Path,
							importPath.Alias,
							i == 0,
						}

						aliasToImport[importPath.Alias] = imp
					}
					if mainImport == nil {
						mainImport = imp
					}
				}
			}

			// TODO: Get this from a database comment if provided
			var defaultLabel string
			var controlName string

			defaultLabel = strings2.Title(rr.GoPlural)
			controlName = rr.GoPlural + typ

			var defaultID string
			defaultID = strings.Replace(t.DbName, "_", "-", -1) + "-" +
				strings.Replace(snaker.CamelToSnake(rr.GoPlural), "_", "-", -1)

			cd := ControlDescription{
				Import: mainImport,
				ControlType: typ,
				NewControlFunc: newFunc,
				ControlName: controlName,
				ControlID: defaultID,
				DefaultLabel: defaultLabel,
				Generator: generator,
				Connector:t.GoName + controlName + "Connector",
			}
			descriptions[rr] = &cd
		}
	}

	return
}

func matchManyManyReferencesWithControls(t *db.Table, descriptions map[interface{}]*ControlDescription, aliasToImport map[string]*ImportType) {
	for _, mm := range t.ManyManyReferences {
		typ, newFunc, importName := defaultControlType(mm)

		if typ != "" {
			var mainImport *ImportType

			generator := GetControlGenerator(importName, typ)
			if generator != nil {
				for i, importPath := range generator.Imports() {
					var ok bool
					var imp *ImportType
					if imp, ok = aliasToImport[importPath.Alias]; ok {
						if imp.Path != importPath.Path {
							panic("found the same alias with different import path")
						}
					} else {
						imp = &ImportType {
							importPath.Path,
							importPath.Alias,
							i == 0,
						}

						aliasToImport[importPath.Alias] = imp
					}
					if mainImport == nil {
						mainImport = imp
					}
				}
			}

			// TODO: Get this from a database comment if provided
			var defaultLabel string
			var controlName string

			defaultLabel = strings2.Title(mm.GoPlural)
			controlName = mm.GoPlural + typ

			var defaultID string
			defaultID = strings.Replace(t.DbName, "_", "-", -1) + "-" +
				strings.Replace(snaker.CamelToSnake(mm.GoPlural), "_", "-", -1)

			cd := ControlDescription{
				Import: mainImport,
				ControlType: typ,
				NewControlFunc: newFunc,
				ControlName: controlName,
				ControlID: defaultID,
				DefaultLabel: defaultLabel,
				Generator: generator,
				Connector:t.GoName + controlName + "Connector",
			}
			descriptions[mm] = &cd
		}
	}

	return
}
