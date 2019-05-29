package db

import "github.com/goradd/gengen/pkg/maps"

type TableDescription struct {
	// DbKey is the key used to find the database in the global database cluster
	DbKey string
	// DbName is the name of the database table or object in the database.
	DbName string
	// LiteralName is the name of the object when describing it to the world. Use the "literalName" option in the comment to override the default. Should be lower case.
	LiteralName string
	// LiteralPlural is the plural name of the object. Use the "literalPlural" option in the comment to override the default. Should be lower case.
	LiteralPlural string
	// GoName is name of the struct when referring to it in go code. Use the "goName" option in the comment to override the default.
	GoName string
	// GoPlural is the name of a collection of these objects when referring to them in go code. Use the "goPlural" option in the comment to override the default.
	GoPlural string
	// LcGoName is the same as GoName, but with first letter lower case.
	LcGoName      string
	// Columns is a list of ColumnDescriptions, one for each column in the table.
	Columns       []*ColumnDescription
	// columnMap is an internal map of the columns
	columnMap     map[string]*ColumnDescription
	// Indexes are the indexes defined in the database. Unique indexes will result in LoadBy* functions.
	Indexes       []IndexDescription
	// Options are key-value pairs of values that can be used to customize how code generation is performed
	Options       maps.SliceMap
	// IsType is true if this is a type table
	IsType        bool
	// IsAssociation is true if this is an association table, which is used to create a many-to-many relationship between two tables.
	IsAssociation bool
	// Comment is the general comment included in the database
	Comment       string

	// The following items are filled in by the analyze process

	// ManyManyReferences describe the many-to-many references pointing to this table
	ManyManyReferences []*ManyManyReference
	// ReverseReferences describes the many-to-one references pointing to this table
	ReverseReferences  []*ReverseReference
	// HasDateTime is true if the table contains a DateTime column.
	HasDateTime bool
	// PrimaryKeyColumn points to the column that contains the primary key of the table.
	PrimaryKeyColumn *ColumnDescription

	// Skip will cause the table to be skipped in code generation
	Skip bool
}


