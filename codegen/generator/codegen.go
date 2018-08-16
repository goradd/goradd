package generator

import (
	"bytes"
	"github.com/spekary/goradd/orm/db"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"os/exec"
	"github.com/spekary/goradd/codegen/connector"
	"github.com/spekary/goradd/util"
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

type ControlDescription struct {
	Import *ImportType
	ControlType string
	NewControlFunc string
	ControlName string
	Generator connector.Generator
}

type ColumnType struct {
	*db.ColumnDescription
	// Related control information
	ControlDescription
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
				columns, imports := ColumnsWithControls(table)
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
				if util.EndsWith(fileName, ".go") {
					execCommand("goimports -w " + fileName)
				}

				// TODO: If a build.go file exists in the directory we are writing to, run it
			}
		}

		for _, typeTableTemplate := range OneTimeTemplates {
			buf.Reset()
			// the template generator function in each template, by convention
			typeTableTemplate.GenerateOnce(codegen, dd, buf)
			fileName := typeTableTemplate.FileName(key)
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

}

func execCommand(command string) {
	parts := strings.Split(command, " ")
	if len(parts) == 0 {
		return
	}

	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		log.Print(err)
	}
}

