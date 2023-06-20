package db

import (
	"github.com/gedex/inflector"
	. "github.com/goradd/goradd/pkg/orm/query"
	"github.com/kenshaw/snaker"
)

// EnumTable describes a enum table, which essentially defines an enumerated type.
// In the SQL world, they are a table with an integer key (starting index 1) and a "name" value, though
// they can have other values associated with them too. Goradd will maintain the
// relationships in SQL, but in a No-SQL situation, it will embed all the ids and values.
type EnumTable struct {
	// DbKey is the key used to find the database in the global database cluster
	DbKey string
	// DbName is the name of the table in the database
	DbName string
	// LiteralName is the english name of the object when describing it to the world. Use the "literalName" option in the comment to override the default.
	LiteralName string
	// LiteralPlural is the plural english name of the object. Use the "literalPlural" option in the comment to override the default.
	LiteralPlural string
	// GoName is the name of the item as a go type name.
	GoName string
	// GoPlural is the plural of the go type
	GoPlural string
	// LcGoName is the lower case version of the go name
	LcGoName string
	// FieldNames are the names of the fields defined in the table. The first field name MUST be the name of the id field, and 2nd MUST be the name of the name field, the others are optional extra fields.
	FieldNames []string
	// FieldTypes are the go column types of the fields, indexed by field name
	FieldTypes map[string]GoColumnType
	// Values are the constant values themselves as defined in the table, mapped to field names in each row.
	Values []map[string]interface{}
	// PkField is the name of the private key field
	PkField string

	// Filled in by analyzer
	Constants map[int]string
}

// FieldGoName returns the go name corresponding to the given field offset
func (tt *EnumTable) FieldGoName(i int) string {
	if i >= len(tt.FieldNames) {
		return ""
	}
	fn := tt.FieldNames[i]
	fn = UpperCaseIdentifier(fn)
	return fn
}

// FieldGoPlural returns the go plural name corresponding to the given field offset
func (tt *EnumTable) FieldGoPlural(i int) string {
	if i >= len(tt.FieldNames) {
		return ""
	}
	fn := tt.FieldNames[i]
	fn = inflector.Pluralize(UpperCaseIdentifier(fn))
	return fn
}

// FieldGoColumnType returns the GoColumnType corresponding to the given field offset
func (tt *EnumTable) FieldGoColumnType(i int) GoColumnType {
	if i >= len(tt.FieldNames) {
		return ColTypeUnknown
	}
	fn := tt.FieldNames[i]
	ft := tt.FieldTypes[fn]
	return ft
}

// FileName returns the default file name corresponding to the enum table.
func (tt *EnumTable) FileName() string {
	return snaker.CamelToSnake(tt.GoName)
}
