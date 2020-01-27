package generator

import (
	"github.com/goradd/goradd/pkg/orm/db"
	strings2 "github.com/goradd/goradd/pkg/strings"
	"github.com/knq/snaker"
	"path"
	"strings"
)

type Importer interface {
	Imports() []string
}

// matchColumnsWithControls maps controls to control descriptions, and returns the imports required by the
// control descriptions
func matchColumnsWithControls(t *db.Table, descriptions map[interface{}]*ControlDescription, importAliases map[string]string) {
	for _, col := range t.Columns {
		controlPath := ControlPath(col)

		if controlPath != "" {
			var imports []string
			typ := path.Base(controlPath)

			generator := GetControlGenerator(controlPath)
			if generator != nil {
				if imp,ok := generator.(Importer); ok {
					imports = imp.Imports()
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
				Path: controlPath,
				Imports: imports,
				ControlType: typ,
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

// ControlPath returns the type of control for a column. It gets this first from the database description, and
// if there is no ControlPath indicated, then from the registered DefaultControlTypeFunc function.
func ControlPath(ref interface{}) string {
	// See if the description has a specific control, which should be a path to the control
	var controlPath string
	switch col := ref.(type) {
	case *db.ReverseReference:
		// need to set this up, getting default from the reverse column
	case *db.ManyManyReference:
		// need to set this up, getting default from the many-many table if it exists
	case *db.Column:
		if i,ok := col.Options["controlPath"]; ok { // a module based control path
			controlPath,ok = i.(string)
			if !ok {
				panic("controlPath must be a string")
			}
			return controlPath // if empty, we want to return empty, because that turns off a control
		}
	}

	if controlPath == "" {
		controlPath = DefaultControlTypeFunc(ref)
	}

	return controlPath
}

func matchReverseReferencesWithControls(t *db.Table, descriptions map[interface{}]*ControlDescription, importAliases map[string]string) {
	for _, rr := range t.ReverseReferences {
		controlPath := ControlPath(rr)

		if controlPath != "" {
			var imports []string
			typ := path.Base(controlPath)

			generator := GetControlGenerator(controlPath)
			if generator != nil {
				if imp,ok := generator.(Importer); ok {
					imports = imp.Imports()
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
				Path: controlPath,
				Imports: imports,
				ControlType: typ,
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

func matchManyManyReferencesWithControls(t *db.Table, descriptions map[interface{}]*ControlDescription, importAliases map[string]string) {
	for _, mm := range t.ManyManyReferences {
		controlPath := ControlPath(mm)
		if controlPath != "" {
			var imports []string
			typ := path.Base(controlPath)

			generator := GetControlGenerator(controlPath)
			if generator != nil {
				if imp,ok := generator.(Importer); ok {
					imports = imp.Imports()
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
				Path: controlPath,
				Imports: imports,
				ControlType: typ,
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
