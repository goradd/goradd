package generator

import (
	"bytes"
	"fmt"
	"github.com/goradd/gofile/pkg/sys"
	"github.com/goradd/goradd/pkg/orm/db"
	"github.com/goradd/goradd/pkg/stringmap"
	"github.com/goradd/goradd/pkg/strings"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
)


// This is specific to automating the build of the examples database code. You do not normally need to
// set this.
var BuildingExamples bool

type CodeGenerator struct {
	// Tables is a map of the tables by database
	Tables     map[string]map[string]TableType
	// TypeTables is a map of the type tables by database
	TypeTables map[string]map[string]TypeTableType

	importAliasesByPath map[string]string
}

type TableType struct {
	*db.Table
	controlDescriptions map[interface{}]*ControlDescription
	Imports             []ImportType
}

type TypeTableType struct {
	*db.TypeTable
}

// ImportType represents an import path required for a control. This is analyzed per-table.
type ImportType struct {
	Path      string
	Alias     string // blank if not needing an alias
}

// ControlDescription is matched with a Column below and provides additional information regarding
// how information in a column can be used to generate a default control to edit that information.
// It is specifically for code generation.
type ControlDescription struct {
	Path		   string
	// Import is the package alias to be used when referring to the package the control is in. It is generated on a per-file basis.
	Import		   string
	// Imports is the list of imported packages that control uses
	Imports        []string
	ControlType    string
	ControlName    string
	ControlID      string // default id to generate
	DefaultLabel   string
	Generator      ControlGenerator
	Connector	   string
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
		TypeTables: make(map[string]map[string]TypeTableType),
	}

	databases := db.GetDatabases()
	if BuildingExamples {
		databases = []db.DatabaseI{db.GetDatabase("goradd")}
	}

	// Map object names to tables, making sure there are no duplicates
	for _, database := range databases {
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
			} else {
				descriptions := make(map[interface{}]*ControlDescription)
				importAliases := make(map[string]string)
				matchColumnsWithControls(table, descriptions, importAliases)
				matchReverseReferencesWithControls(table, descriptions, importAliases)
				matchManyManyReferencesWithControls(table, descriptions, importAliases)

				var i []ImportType
				for _,k := range stringmap.SortedKeys(importAliases) {
					i = append(i, ImportType{k,importAliases[k]})
				}

				t := TableType{
					table,
					descriptions,
					i,
				}
				codegen.Tables[key][table.GoName] = t
			}
		}

	}

	buf := new(bytes.Buffer)

	// Generate the templates.
	for _, database := range databases {
		dd := database.Describe()
		dbKey := dd.DbKey

		for _, tableKey := range stringmap.SortedKeys(codegen.TypeTables[dbKey]) {
			typeTable := codegen.TypeTables[dbKey][tableKey]
			for _, typeTableTemplate := range TypeTableTemplates {
				buf.Reset()
				// the template generator function in each template, by convention
				typeTableTemplate.GenerateTypeTable(codegen, dd, typeTable, buf)
				fileName := typeTableTemplate.FileName(dbKey, typeTable)
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

		for _, tableKey := range stringmap.SortedKeys(codegen.Tables[dbKey]) {
			table := codegen.Tables[dbKey][tableKey]
			for _, tableTemplate := range TableTemplates {
				buf.Reset()
				tableTemplate.GenerateTable(codegen, dd, table, buf)
				fileName := tableTemplate.FileName(dbKey, table)
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
					curDir,_ := os.Getwd()
					_ = os.Chdir(filepath.Dir(fileName)) // run it from the files directory to pick up the correct go.mod file if there is one
					_, err = sys.ExecuteShellCommand("goimports -w " + filepath.Base(fileName))
					_ = os.Chdir(curDir)
					if err != nil {
						panic("error running goimports: " + string(err.(*exec.ExitError).Stderr)) // perhaps goimports is not installed?
					}
				}
			}
		}

		for _, oneTimeTemplate := range OneTimeTemplates {
			buf.Reset()
			// the template generator function in each template, by convention
			oneTimeTemplate.GenerateOnce(codegen, dd, buf)
			fileName := oneTimeTemplate.FileName(dbKey)
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

// Reset resets the internal information of the code generator. Call this just before generating a file.
func (c *CodeGenerator) ResetImports() {
	c.importAliasesByPath = make(map[string]string)
}

// AddImportPath adds an import path to the import path list. In particular, it will help manage the package aliases
// so the path can be referred to using the correct package name or package alias. Call this on all
// paths used by the file before calling ImportString.
func (c *CodeGenerator) AddImportPaths(paths ...string) {
	for _,p := range paths {
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
	for p,a := range c.importAliasesByPath {
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
	imp,t := path.Split(p)
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



