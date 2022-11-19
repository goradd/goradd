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

type TableTemplateI interface {
	GenerateTable(codegen CodeGenerator, dd *db.Model, t TableType, _w io.Writer) (err error)
	FileName(key string, t TableType) string
	Overwrite() bool
}

type TypeTableTemplateI interface {
	GenerateTypeTable(codegen CodeGenerator, dd *db.Model, t TypeTableType, _w io.Writer) (err error)
	FileName(key string, t TypeTableType) string
	Overwrite() bool
}

type OneTimeTemplateI interface {
	GenerateOnce(codegen CodeGenerator, dd *db.Model, _w io.Writer) (err error)
	FileName(key string) string
	Overwrite() bool
}

//type TemplateFunc func(codegen *CodeGenerator, t *TableType, _w io.Writer) (err error)

// Will be populated by the individual templates found
var TableTemplates []TableTemplateI

func AddTableTemplate(t TableTemplateI) {
	TableTemplates = append(TableTemplates, t)
}

var TypeTableTemplates []TypeTableTemplateI

func AddTypeTableTemplate(t TypeTableTemplateI) {
	TypeTableTemplates = append(TypeTableTemplates, t)
}

var OneTimeTemplates []OneTimeTemplateI

func AddOneTimeTemplate(t OneTimeTemplateI) {
	OneTimeTemplates = append(OneTimeTemplates, t)
}
