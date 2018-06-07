package generator

import (
	"bytes"
	"github.com/spekary/goradd/orm/db"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type Codegen struct {
	Tables     map[string]map[string]*db.TableDescription
	TypeTables map[string]map[string]*db.TypeTableDescription
}

type TableType struct {
	db.TableDescription
	Columns    []*ColumnType
	PrimaryKey *ColumnType
}

type TypeTableType struct {
	db.TypeTableDescription

	// Filled in by analyezer
	Constants map[uint]string
}

type ForeignKeyType struct {
	//DbKey string	// We don't support cross database foreign keys yet. Someday maybe.
	TableName string
	ColName   string
}

type ColumnType struct {
	db.ColumnDescription

	// Filled in by analyzer
	ForeignKey *ForeignKeyType
}

func Generate() {

	codegen := Codegen{
		Tables:     map[string]map[string]*db.TableDescription{},
		TypeTables: map[string]map[string]*db.TypeTableDescription{},
	}

	// Map object names to tables, making sure there are no duplicates
	for key, database := range db.GetDatabases() {
		codegen.Tables[key] = make(map[string]*db.TableDescription)
		codegen.TypeTables[key] = make(map[string]*db.TypeTableDescription)
		dd := database.Describe()
		for _, typeTable := range dd.TypeTables {
			if _, ok := codegen.TypeTables[key][typeTable.GoName]; ok {
				log.Println("Error: type table " + typeTable.GoName + " is defined more than once.")
			} else {
				codegen.TypeTables[key][typeTable.GoName] = typeTable
			}
		}
		for _, table := range dd.Tables {
			if _, ok := codegen.Tables[key][table.GoName]; ok {
				log.Println("Error:  table " + table.GoName + " is defined more than once.")
			} else if !table.IsAssociation {
				codegen.Tables[key][table.GoName] = table
			}
		}
	}

	buf := new(bytes.Buffer)

	// Generate the templates.
	for key, database := range db.GetDatabases() {
		dd := database.Describe()
		for _, typeTable := range dd.TypeTables {
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

		for _, table := range dd.Tables {
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
			}
		}

	}

}
