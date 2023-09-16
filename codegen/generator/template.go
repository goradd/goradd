package generator

import (
	"github.com/goradd/goradd/pkg/orm/db"
	"io"
)

type Template struct {
	Overwrite bool
	TargetDir string
	//OutFunc TemplateFunc
}

// TableTemplateI represents the interface for templates that are executed once per regular table.
type TableTemplateI interface {
	GenerateTable(codegen CodeGenerator, dd *db.Model, t TableType, _w io.Writer) (err error)
	FileName(key string, t TableType) string
	Overwrite() bool
}

// EnumTableTemplateI represents the interface for templates that are executed once per enum table.
type EnumTableTemplateI interface {
	GenerateEnumTable(codegen CodeGenerator, dd *db.Model, t EnumTableType, _w io.Writer) (err error)
	FileName(key string, t EnumTableType) string
	Overwrite() bool
}

// DatabaseTemplateI represents the interface for templates that are executed once per database.
type DatabaseTemplateI interface {
	GenerateDatabase(codegen CodeGenerator, dd *db.Model, _w io.Writer) (err error)
	FileName(key string) string
	Overwrite() bool
}

// OneTimeTemplateI represents the interface for templates that are executed just once.
type OneTimeTemplateI interface {
	GenerateOnce(codegen CodeGenerator, databases []db.DatabaseI, _w io.Writer) (err error)
	FileName() string
	Overwrite() bool
}

//type TemplateFunc func(codegen *CodeGenerator, t *TableType, _w io.Writer) (err error)

// TableTemplates is the collection of templates that will be populated by the individual templates found
var TableTemplates []TableTemplateI

func AddTableTemplate(t TableTemplateI) {
	TableTemplates = append(TableTemplates, t)
}

// EnumTableTemplates is the list of templates that are executed once per "enum" table.
var EnumTableTemplates []EnumTableTemplateI

// AddEnumTableTemplate adds a template to the EnumTableTemplates list.
func AddEnumTableTemplate(t EnumTableTemplateI) {
	EnumTableTemplates = append(EnumTableTemplates, t)
}

// DatabaseTemplates is the collection of templates that are executed once per database.
var DatabaseTemplates []DatabaseTemplateI

// AddDatabaseTemplate adds a template to the DatabaseTemplates list.
func AddDatabaseTemplate(t DatabaseTemplateI) {
	DatabaseTemplates = append(DatabaseTemplates, t)
}

// OneTimeTemplates is the collection of templates that are executed once per codegen rum.
var OneTimeTemplates []OneTimeTemplateI

// AddOneTimeTemplate adds a template to the OneTimeTemplates list.
func AddOneTimeTemplate(t OneTimeTemplateI) {
	OneTimeTemplates = append(OneTimeTemplates, t)
}
