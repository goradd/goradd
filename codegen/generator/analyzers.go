package generator

import (
	"fmt"
	"github.com/goradd/goradd/pkg/orm/db"
	"strconv"
	"strings"
)

func columnsWithControls(t *db.TableDescription) (columns []ColumnType, imports []*ImportType) {
	var pathToImport = make(map[string]*ImportType)
	var namespaceToImport = make(map[string]*ImportType)

	for _, col := range t.Columns {
		col2 := ColumnType{ColumnDescription: col}

		typ, newFunc, importName := controlType(col)

		if typ != "" {
			var mainImport *ImportType

			generator := GetControlGenerator(importName, typ)
			if generator == nil {
				panic(fmt.Errorf("Generator for control type %s/%s is not defined", importName, typ))
			}
			for i, importPath := range generator.Imports() {
				var ok bool
				var imp *ImportType
				if imp, ok = pathToImport[importPath]; !ok {
					var namespace string
					// add new import path
					items := strings.Split(importPath, `/`)
					lastName := items[len(items)-1]
					var suffix = ""
					var count = 1
					for {
						if _, ok = namespaceToImport[lastName+suffix]; !ok {
							break
						}
						count++
						suffix = strconv.Itoa(count)
					}
					namespace = lastName + suffix

					if suffix == "" {
						imp = &ImportType{
							importPath,
							lastName,
							"",
							i == 0,
						}
					} else {
						imp = &ImportType{
							importPath,
							namespace,
							namespace,
							i == 0,
						}
					}
					imports = append(imports, imp)
					pathToImport[importPath] = imp
					namespaceToImport[namespace] = imp
					if mainImport == nil {
						mainImport = imp
					}
				} else {
					if mainImport == nil {
						mainImport = imp
					}
				}
			}
			defaultLabel := strings.Title(strings.Replace(col.DbName, "_", " ", -1))

			var defaultID string
			if GenerateControlIDs {
				defaultID = strings.Replace(t.DbName, "_", "-", -1) + "-" + strings.Replace(col.DbName, "_", "-", -1)
			}

			col2.ControlDescription = ControlDescription{
				mainImport,
				typ,
				newFunc,
				col.GoName + typ,
				defaultID,
				defaultLabel,
				generator,
			}
		}
		columns = append(columns, col2)
	}

	return
}

// ControlType returns the default type of control for a column. Control types can be customized in other ways too.
func controlType(col *db.ColumnDescription) (typ string, createFunc string, importName string) {
	d := DefaultControlTypeFunc(col)
	return d.Typ, d.CreateFunc, d.ImportName
}
