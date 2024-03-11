package db

import (
	strings2 "github.com/goradd/goradd/pkg/strings"
	"github.com/kenshaw/snaker"
)

type Table struct {
	// DbKey is the key used to find the database in the global database cluster
	DbKey string
	// DbName is the name of the database table or object in the database.
	DbName string
	// LiteralName is the name of the object when describing it to the world. Use the "literalName" option in the comment to override the default. Should be lower case.
	LiteralName string
	// LiteralPlural is the plural name of the object. Use the "literalPlural" option in the comment to override the default. Should be lower case.
	LiteralPlural string
	// GoName is the name of the struct when referring to it in go code. Use the "goName" option in the comment to override the default.
	GoName string
	// GoPlural is the name of a collection of these objects when referring to them in go code. Use the "goPlural" option in the comment to override the default.
	GoPlural string
	// LcGoName is the same as GoName, but with first letter lower case.
	LcGoName string
	// Columns is a list of ColumnDescriptions, one for each column in the table.
	Columns []*Column
	// columnMap is an internal map of the columns by database name of the column
	columnMap map[string]*Column
	// Indexes are the indexes defined in the database. Unique indexes will result in LoadBy* functions.
	Indexes []Index
	// Options are key-value pairs of values that can be used to customize how code generation is performed
	Options map[string]interface{}
	// isEnum is true if this is a enum table
	isEnum bool
	// Comment is the general comment included in the database
	Comment string

	// The following items are filled in by the importDescription process

	// ManyManyReferences describe the many-to-many references pointing to this table
	ManyManyReferences []*ManyManyReference
	// ReverseReferences describes the many-to-one references pointing to this table
	ReverseReferences []*ReverseReference

	// SupportsForeignKeys determines whether the table uses a storage engine that
	// supports foreign keys. Some examples: Postgres supports foreign keys across all tables,
	// while some MySQL database engines do not support foreign keys, and the engine
	// can be assigned on a table-by-table basis.
	//
	// If the table has foreign key support, the code generator will assume that the
	// database will automatically handle updates and deletes of foreign key values.
	// If not, the code generator will attempt to update the foreign keys as appropriate.
	SupportsForeignKeys bool
}

func (t *Table) PrimaryKeyColumn() *Column {
	if len(t.Columns) == 0 {
		return nil
	}
	if !t.Columns[0].IsPk {
		return nil
	}
	return t.Columns[0]
}

func (t *Table) PrimaryKeyGoType() string {
	return t.PrimaryKeyColumn().ColumnType.GoType()
}

// GetColumn returns a Column given the database name of a column
func (t *Table) GetColumn(name string) *Column {
	return t.columnMap[name]
}

// DefaultHtmlID is the default id of the corresponding form object when used in generated HTML.
func (t *Table) DefaultHtmlID() string {
	return strings2.CamelToKebab(t.GoName)
}

// FileName is the base name of generated file names that correspond to this database table.
// Typically, Go files are lower case snake case by convention.
func (t *Table) FileName() string {
	s := snaker.CamelToSnake(t.GoName)
	if strings2.EndsWith(s, "_test") {
		// Go will ignore files that end with _test. If we somehow create a filename like this,
		// we add an underscore to make sure it is still included in a build.
		s = s + "_"
	}
	return s
}

// HasGetterName returns true if the given name is in use by one of the getters.
// This is used for detecting naming conflicts. Will also return an error string
// to display if there is a conflict.
func (t *Table) HasGetterName(name string) (hasName bool, desc string) {
	for _, c := range t.Columns {
		if c.GoName == name {
			return false, "conflicts with column " + c.GoName
		}
	}

	for _, rr := range t.ReverseReferences {
		if rr.GoName == name {
			return false, "conflicts with reverse reference singular name " + rr.GoName
		}
		if rr.GoPlural == name {
			return false, "conflicts with reverse reference plural name " + rr.GoPlural
		}
	}

	for _, mm := range t.ManyManyReferences {
		if mm.GoName == name {
			return false, "conflicts with many-many singular name " + mm.GoName
		}
		if mm.GoPlural == name {
			return false, "conflicts with many-many plural name " + mm.GoPlural
		}
	}
	return false, ""
}
