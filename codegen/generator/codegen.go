package generator

import (
	"bytes"
	"fmt"
	"github.com/goradd/gofile/pkg/sys"
	"github.com/goradd/goradd/pkg/config"
	"github.com/goradd/goradd/pkg/orm/db"
	"github.com/goradd/goradd/pkg/stringmap"
	"github.com/goradd/goradd/pkg/strings"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
)

// BuildingExamples turns on the building of the templates for the examples code.
// This is specific to automating the build of the examples database code.
// You do not normally need to set this.
var BuildingExamples bool

type CodeGenerator struct {
	// Tables is a map of the tables by database
	Tables map[string]map[string]TableType
	// EnumTables is a map of the enum tables by database
	EnumTables map[string]map[string]EnumTableType

	// importAliasesByPath stores import paths by package name to help correctly manage packages with the same name
	importAliasesByPath map[string]string
}

type TableType struct {
	*db.Table
	controlDescriptions map[interface{}]*ControlDescription
	Imports             []ImportType
}

type EnumTableType struct {
	*db.EnumTable
}

// ImportType represents an import path required for a control. This is analyzed per-table.
type ImportType struct {
	Path  string
	Alias string // blank if not needing an alias
}

// ControlDescription is matched with a Column below and provides additional information regarding
// how information in a column can be used to generate a default control to edit that information.
// It is specifically for code generation.
type ControlDescription struct {
	Path string
	// Package is the package alias to be used when referring to the package the control is in. It is generated on a per-file basis.
	Package string
	// Imports is the list of imported packages that the control uses
	Imports      []string
	ControlType  string
	ControlName  string
	ControlID    string // default id to generate
	DefaultLabel string
	Generator    ControlGenerator
	Connector    string
}

func (cd *ControlDescription) ControlIDConst() string {
	if cd.ControlID == "" {
		return ""
	}
	return strings.KebabToCamel(cd.ControlID) + "Id"
}

func (t *TableType) GetColumnByDbName(name string) *db.Column {
	for _, col := range t.Columns {
		if col.DbName == name {
			return col
		}
	}
	return nil
}

func (t *TableType) ControlDescription(ref interface{}) *ControlDescription {
	return t.controlDescriptions[ref]
}

func Generate() {
	codegen := CodeGenerator{
		Tables:     make(map[string]map[string]TableType),
		EnumTables: make(map[string]map[string]EnumTableType),
	}

	databases := db.GetDatabases()

	if len(databases) == 0 {
		log.Println("There are no databases to use for code generation. Setup databases in the db.cfg file.")
	}

	if BuildingExamples {
		databases = []db.DatabaseI{db.GetDatabase("goradd")}
	}

	// map object names to tables, making sure there are no duplicates
	for _, database := range databases {
		if database.Model() == nil {
			panic("Missing model. Did you forget to call Analyze on the database?")
		}
		dbKey := database.Model().DbKey
		codegen.Tables[dbKey] = make(map[string]TableType)
		codegen.EnumTables[dbKey] = make(map[string]EnumTableType)
		dd := database.Model()

		// Create wrappers for the tables with extra analysis required for form generation
		stringmap.Range(dd.EnumTables, func(k string, enumTable *db.EnumTable) bool {
			if _, ok := codegen.EnumTables[dbKey][enumTable.GoName]; ok {
				log.Println("Error: enum table " + enumTable.GoName + " is defined more than once.")
			} else {
				tt := EnumTableType{
					enumTable,
				}
				codegen.EnumTables[dbKey][enumTable.GoName] = tt
			}
			return true
		})

		stringmap.Range(dd.Tables, func(k string, table *db.Table) bool {
			if _, ok := codegen.Tables[dbKey][table.GoName]; ok {
				log.Println("Error:  table " + table.GoName + " is defined more than once.")
			} else {
				descriptions := make(map[interface{}]*ControlDescription)
				importAliases := make(map[string]string)
				matchColumnsWithControls(database, table, descriptions, importAliases)
				matchReverseReferencesWithControls(table, descriptions, importAliases)
				matchManyManyReferencesWithControls(table, descriptions, importAliases)

				var i []ImportType
				for _, k := range stringmap.SortedKeys(importAliases) {
					i = append(i, ImportType{k, importAliases[k]})
				}

				t := TableType{
					table,
					descriptions,
					i,
				}
				codegen.Tables[dbKey][table.GoName] = t
			}
			return true
		})
	}

	buf := new(bytes.Buffer)

	// Generate the templates.
	for _, database := range databases {
		dd := database.Model()
		dbKey := dd.DbKey

		for _, tableKey := range stringmap.SortedKeys(codegen.EnumTables[dbKey]) {
			enumTable := codegen.EnumTables[dbKey][tableKey]
			for _, enumTableTemplate := range EnumTableTemplates {
				buf.Reset()
				// the template generator function in each template, by convention
				enumTableTemplate.GenerateEnumTable(codegen, dd, enumTable, buf)
				fileName := enumTableTemplate.FileName(dbKey, enumTable)
				fp := filepath.Dir(fileName)

				// If the file already exists, and we are not over-writing, skip it
				if _, err := os.Stat(fileName); err == nil {
					if !enumTableTemplate.Overwrite() {
						continue
					}
				}

				if err := os.MkdirAll(fp, 0777); err != nil {
					log.Print(err)
				}
				if err := os.WriteFile(fileName, buf.Bytes(), 0644); err != nil {
					log.Print(err)
				} else {
					if Verbose {
						log.Printf("Writing %s", fileName)
					}
				}
				RunGoImports(fileName)
			}
		}

		for _, tableKey := range stringmap.SortedKeys(codegen.Tables[dbKey]) {
			table := codegen.Tables[dbKey][tableKey]
			if table.PrimaryKeyColumn() == nil {
				log.Println("*** Skipping table " + table.DbName + " since it has no primary key column")
				continue
			}
			for _, tableTemplate := range TableTemplates {
				buf.Reset()
				tableTemplate.GenerateTable(codegen, dd, table, buf)
				fileName := tableTemplate.FileName(dbKey, table)
				fp := filepath.Dir(fileName)

				// If the file already exists, and we are not over-writing, skip it
				if _, err := os.Stat(fileName); err == nil {
					if !tableTemplate.Overwrite() {
						continue
					}
				}

				if err := os.MkdirAll(fp, 0777); err != nil {
					log.Print(err)
				}
				if err := os.WriteFile(fileName, buf.Bytes(), 0644); err != nil {
					log.Print(err)
				} else {
					if Verbose {
						log.Printf("Writing %s", fileName)
					}
				}
				RunGoImports(fileName)
			}
		}

		for _, dbTemplate := range DatabaseTemplates {
			buf.Reset()
			// the template generator function in each template, by convention
			dbTemplate.GenerateDatabase(codegen, dd, buf)
			fileName := dbTemplate.FileName(dbKey)
			fp := filepath.Dir(fileName)

			// If the file already exists, and we are not over-writing, skip it
			if _, err := os.Stat(fileName); err == nil {
				if !dbTemplate.Overwrite() {
					continue
				}
			}

			if err := os.MkdirAll(fp, 0777); err != nil {
				log.Print(err)
			}
			if err := os.WriteFile(fileName, buf.Bytes(), 0644); err != nil {
				log.Print(err)
			} else {
				if Verbose {
					log.Printf("Writing %s", fileName)
				}
			}
			RunGoImports(fileName)
		}
	}

	for _, onceTemplate := range OneTimeTemplates {
		buf.Reset()
		onceTemplate.GenerateOnce(codegen, databases, buf)
		fileName := onceTemplate.FileName()
		fp := filepath.Dir(fileName)
		// If the file already exists, and we are not over-writing, skip it
		if _, err := os.Stat(fileName); err == nil {
			if !onceTemplate.Overwrite() {
				continue
			}
		}

		if err := os.MkdirAll(fp, 0777); err != nil {
			log.Print(err)
		}
		if err := os.WriteFile(fileName, buf.Bytes(), 0644); err != nil {
			log.Print(err)
		} else {
			if Verbose {
				log.Printf("Writing %s", fileName)
			}
		}
		RunGoImports(fileName)
	}
}

func RunGoImports(fileName string) {
	// run imports on all generated go files
	if strings.EndsWith(fileName, ".go") {
		curDir, _ := os.Getwd()
		_ = os.Chdir(filepath.Dir(fileName)) // run it from the file's directory to pick up the correct go.mod file if there is one
		_, err := sys.ExecuteShellCommand("goimports -w " + filepath.Base(fileName))
		_ = os.Chdir(curDir)
		if err != nil {
			if e, ok := err.(*exec.Error); ok {
				panic("error running goimports: " + e.Error()) // perhaps goimports is not installed?
			} else if e2, ok2 := err.(*exec.ExitError); ok2 {
				// Likely a syntax error in the resulting file
				log.Print(string(e2.Stderr))
			}
		}
	}
}

// ResetImports resets the internal information of the code generator. Call this just before generating a file.
func (c *CodeGenerator) ResetImports() {
	c.importAliasesByPath = make(map[string]string)
}

// AddImportPaths adds an import path to the import path list. In particular, it will help manage the package aliases
// so the path can be referred to using the correct package name or package alias. Call this on all
// paths used by the file before calling ImportString.
func (c *CodeGenerator) AddImportPaths(paths ...string) {
	for _, p := range paths {
		if p == "" {
			return
		}
		if _, ok := c.importAliasesByPath[p]; ok {
			return
		}
		alias := path.Base(p)
		count := 2
		newAlias := alias
	Found:
		for _, a := range c.importAliasesByPath {
			if a == newAlias {
				newAlias = fmt.Sprint(alias, count)
				count++
				continue Found
			}
		}
		c.importAliasesByPath[p] = newAlias
	}
}

// AddObjectPath adds an object path to the import path list. In particular, it will help manage the package list
// so the object can referred to using the correct package name or package alias. Call this on all object
// paths used by the form before calling ImportString.
func (c *CodeGenerator) AddObjectPath(p string) {
	c.AddImportPaths(path.Dir(p))
}

// ImportStrings returns strings to use in an import statement for all of the objects and imports entered
func (c *CodeGenerator) ImportStrings() (ret string) {
	for p, a := range c.importAliasesByPath {
		pkg := path.Base(p)
		if a == pkg {
			ret += fmt.Sprintf("%#v\n", p)
		} else {
			ret += fmt.Sprintf("%s %#v\n", a, p)
		}
	}
	return
}

// ObjectType returns the string that should be used for an object type given its module path
func (c *CodeGenerator) ObjectType(p string) string {
	imp, t := path.Split(p)
	imp = path.Clean(imp)
	if a := c.importAliasesByPath[imp]; a == "" {
		panic("unknown object path: " + p)
	} else {
		return a + "." + t
	}
}

func (c *CodeGenerator) ImportPackage(imp string) string {
	if a := c.importAliasesByPath[imp]; a == "" {
		panic("unknown import path: " + imp)
	} else {
		return a
	}
}

func (c *CodeGenerator) ObjectPackage(imp string) string {
	imp = path.Dir(imp)
	if a := c.importAliasesByPath[imp]; a == "" {
		panic("unknown import path: " + imp)
	} else {
		return a
	}
}

// WrapFormField returns a creator template for a field wrapper with type wrapperType.
//
// child should be the creator template for the control that will be wrapped.
func (c *CodeGenerator) WrapFormField(wrapperType string, label string, forId string, child string) string {
	return fmt.Sprintf(
		`%sCreator{
	ID: p.ID() + "-%s%s",
	For:  p.ID() + "-%s",
	Label: "%s",
	Child: %s,
}
`, wrapperType, forId, config.DefaultFormFieldWrapperIdSuffix, forId, label, child)
}
