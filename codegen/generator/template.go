package generator

import (
	"bytes"
	"github.com/spekary/goradd/orm/db"
)

type Template struct {
	Overwrite bool
	TargetDir string
	//OutFunc TemplateFunc
}

type TableTemplateI interface {
	GenerateTable(codegen Codegen, dd *db.DatabaseDescription, t *db.TableDescription, buf *bytes.Buffer)
	FileName(t *db.TableDescription) string
}

type TypeTableTemplateI interface {
	GenerateTypeTable(codegen Codegen, dd *db.DatabaseDescription, t *db.TypeTableDescription, buf *bytes.Buffer)
	FileName(t *db.TypeTableDescription) string
}



//type TemplateFunc func(codegen *Codegen, t *TableType, buf *bytes.Buffer)

// Will be populated by the individual templates found
var TableTemplates []TableTemplateI

func AddTableTemplate(t TableTemplateI) {
	TableTemplates = append(TableTemplates, t)
}

var TypeTableTemplates []TypeTableTemplateI

func AddTypeTableTemplate(t TypeTableTemplateI) {
	TypeTableTemplates = append(TypeTableTemplates, t)
}