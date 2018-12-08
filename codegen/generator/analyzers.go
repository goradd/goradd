package generator

import (
	"github.com/spekary/goradd/pkg/orm/db"
	"strings"
	"strconv"
)


func columnsWithControls(t *db.TableDescription) (columns []ColumnType, imports []*ImportType) {
	var pathToImport = make(map[string]*ImportType)
	var namespaceToImport = make(map[string]*ImportType)

	for _,col := range t.Columns {
		col2 := ColumnType{ColumnDescription:col}

		typ, newFunc, importPath := controlType(col)

		var namespace string
		var imp *ImportType
		var ok bool
		if typ != "" {
			if imp,ok = pathToImport[importPath]; !ok {
				// add new import path
				items := strings.Split(importPath, `/`)
				lastName := items[len(items)-1]
				var suffix = ""
				var count = 1
				for  {
					if _,ok = namespaceToImport[lastName + suffix]; !ok {
						break
					}
					count ++
					suffix = strconv.Itoa(count)
				}
				namespace = lastName + suffix

				if suffix == "" {
					imp = &ImportType{
						importPath,
						lastName,
						"",
					}
				} else {
					imp = &ImportType{
						importPath,
						namespace,
						namespace,
					}
				}
				imports = append(imports, imp)
				pathToImport[importPath] = imp
				namespaceToImport[namespace] = imp
			}
			defaultLabel := strings.Title(strings.Replace(col.DbName, "_", " ", -1))

			var defaultID string
			if GenerateControlIDs {
				defaultID = strings.Replace(t.DbName, "_", "-", -1) + "-" + strings.Replace(col.DbName, "_", "-", -1)
			}

			col2.ControlDescription = ControlDescription{
				imp,
				typ,
				newFunc,
				col.GoName + typ,
				defaultID,
				defaultLabel,
				GetControlGenerator(importPath, typ),
			}
		}
		columns = append(columns, col2)
	}

	return
}


// ControlType returns the default type of control for a column. Control types can be customized in other ways too.
func  controlType(col *db.ColumnDescription) (typ string, createFunc string, importName string) {
	d := DefaultControlTypeFunc(col)
	return d.Typ, d.CreateFunc, d.ImportName
}