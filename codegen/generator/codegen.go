package generator

import (
	"bytes"
	"github.com/spekary/goradd/pkg/orm/db"
	"github.com/spekary/goradd/pkg/strings"
	"github.com/spekary/goradd/pkg/sys"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type Codegen struct {
	Tables     map[string]map[string]TableType	// TODO: Change to ordered maps for consistent codegeneration
	TypeTables map[string]map[string]TypeTableType
}

type TableType struct {
	*db.TableDescription
	Columns    []ColumnType
	Imports    []*ImportType
}

type TypeTableType struct {
	*db.TypeTableDescription
}

// ImportType represents an import path required for a control. This is analyzed per-table.
type ImportType struct {
	Path string
	Namespace string
	Alias string // blank if not needing an alias
}

// ControlDescription is matched with a ColumnDescription below and provides additional information regarding
// how information in a column can be used to generated a default control to edit that information.
type ControlDescription struct {
	Import *ImportType
	ControlType string
	NewControlFunc string
	ControlName string
	ControlID string	// default id to generate
	DefaultLabel string
	Generator ControlGenerator
}

// ColumnType combines a database ColumnDescription with a ControlDescription
type ColumnType struct {
	*db.ColumnDescription
	// Related control information
	ControlDescription
}

func (t *TableType) GetColumnByDbName(name string) *ColumnType {
	for _,col := range t.Columns {
		if col.DbName == name {
			return &col
		}
	}
	return nil
}

func Generate() {

	codegen := Codegen{
		Tables:     make(map[string]map[string]TableType),
		TypeTables: make(map[string]map[string]TypeTableType),
	}

	// Map object names to tables, making sure there are no duplicates
	for _, database := range db.GetDatabases() {
		key := database.Describe().DbKey
		codegen.Tables[key] = make(map[string]TableType)
		codegen.TypeTables[key] = make(map[string]TypeTableType)
		dd := database.Describe()

		// Create wrappers for the tables with extra analysis required for form generation
		for _, typeTable := range dd.TypeTables {
			if _, ok := codegen.TypeTables[key][typeTable.GoName]; ok {
				log.Println("Error: type table " + typeTable.GoName + " is defined more than once.")
			} else {
				tt := TypeTableType{
					typeTable,
				}
				codegen.TypeTables[key][typeTable.GoName] = tt
			}
		}
		for _, table := range dd.Tables {
			if _, ok := codegen.Tables[key][table.GoName]; ok {
				log.Println("Error:  table " + table.GoName + " is defined more than once.")
			} else if !table.IsAssociation {
				columns, imports := columnsWithControls(table)
				t := TableType {
					table,
					columns,
					imports,
				}
				codegen.Tables[key][table.GoName] = t
			}
		}
	}

	buf := new(bytes.Buffer)

	// Generate the templates.
	for _, database := range db.GetDatabases() {
		dd := database.Describe()
		key := dd.DbKey
		for _, typeTable := range codegen.TypeTables[key] {
			for _, typeTableTemplate := range TypeTableTemplates {
				buf.Reset()
				// the template generator function in each template, by convention
				typeTableTemplate.GenerateTypeTable(codegen, dd, typeTable, buf)
				fileName := typeTableTemplate.FileName(key, typeTable)
				path := filepath.Dir(fileName)

				if _, err := os.Stat(fileName); err == nil {
					if !typeTableTemplate.Overwrite() {
						continue
					}
				}

				os.MkdirAll(path, 0777)
				err := ioutil.WriteFile(fileName, buf.Bytes(), 0644)
				if err != nil {
					log.Print(err)
				}
			}
		}

		for _, table := range codegen.Tables[key] {
			if table.IsAssociation || table.Skip {
				continue
			}
			for _, tableTemplate := range TableTemplates {
				buf.Reset()
				tableTemplate.GenerateTable(codegen, dd, table, buf)
				fileName := tableTemplate.FileName(key, table)
				path := filepath.Dir(fileName)

				if _, err := os.Stat(fileName); err == nil {
					if !tableTemplate.Overwrite() {
						continue
					}
				}

				os.MkdirAll(path, 0777)
				err := ioutil.WriteFile(fileName, buf.Bytes(), 0644)
				if err != nil {
					log.Print(err)
				}

				// run imports on all generated go files
				if strings.EndsWith(fileName, ".go") {
					sys.ExecuteShellCommand("goimports -w " + fileName)
				}
			}
		}

		for _, oneTimeTemplate := range OneTimeTemplates {
			buf.Reset()
			// the template generator function in each template, by convention
			oneTimeTemplate.GenerateOnce(codegen, dd, buf)
			fileName := oneTimeTemplate.FileName(key)
			path := filepath.Dir(fileName)

			if _, err := os.Stat(fileName); err == nil {
				if !oneTimeTemplate.Overwrite() {
					continue
				}
			}

			os.MkdirAll(path, 0777)
			err := ioutil.WriteFile(fileName, buf.Bytes(), 0644)
			if err != nil {
				log.Print(err)
			}
		}


	}

}
