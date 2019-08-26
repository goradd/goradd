package db

// ForeignKeyColumn is additional information to describe what a foreign key points to
type ForeignKeyColumn struct {
	//DbKey string	// We don't support cross database foreign keys yet. Someday maybe.
	// TableName is the name of the table on the other end of the foreign key
	TableName string
	// ColumnName is the database column name in the linked table that matches this column name. Often that is the primary key of the other table.
	ColumnName string
	// UpdateAction indicates how the database will react when the other end of the relationship's value changes.
	UpdateAction FKAction
	// DeleteAction indicates how the database will react when the other end of the relationship's record is deleted.
	DeleteAction FKAction
	// GoName is the name we should use to refer to the related object
	GoName string
	// GoType is the type of the related object
	GoType string
	// GoTypePlural is the plural version of the type when referring to groups of related objects
	GoTypePlural string
	// IsType is true if this is a related type
	IsType bool
	// RR is filled in by the analyzer and represent a reverse reference relationship
	RR *ReverseReference
}

/*
// ForeignKeyDescription describes a foreign key relationship between columns in one table and columns in a different table.
// We currently allow the collection of multi-table and cross-database fk data, but we don't currently support them in codegen.
// ForeignKeyDescription is primarily for SQL databases that have specific support for foreign keys, like MySQL InnoDB tables.
type ForeignKeyDescription struct {
	KeyName         string
	Columns         []string
	RelationSchema  string
	RelationTable   string
	relationColumns []string // must be ordered that same as columns
}
*/

// GoVarName returns the name of the go object used to refer to the kind of object the foreign key points to.
func (fk *ForeignKeyColumn) GoVarName() string {
	return "obj" + fk.GoName
}
