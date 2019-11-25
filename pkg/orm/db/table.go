package db

type Table struct {
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
	LcGoName string
	// Columns is a list of ColumnDescriptions, one for each column in the table.
	Columns []*Column
	// columnMap is an internal map of the columns
	columnMap map[string]*Column
	// Indexes are the indexes defined in the database. Unique indexes will result in LoadBy* functions.
	Indexes []Index
	// Options are key-value pairs of values that can be used to customize how code generation is performed
	Options map[string]interface{}
	// IsType is true if this is a type table
	IsType bool
	// Comment is the general comment included in the database
	Comment string

	// The following items are filled in by the analyze process

	// ManyManyReferences describe the many-to-many references pointing to this table
	ManyManyReferences []*ManyManyReference
	// ReverseReferences describes the many-to-one references pointing to this table
	ReverseReferences []*ReverseReference
	// HasDateTime is true if the table contains a DateTime column.
	HasDateTime bool
}

func (t *Table) PrimaryKeyColumn() *Column {
	return t.Columns[0]
}

func (t *Table) PrimaryKeyGoType() string {
	return t.PrimaryKeyColumn().ColumnType.GoType()
}


// GetColumn returns a Column given the name of a column
func (t *Table) GetColumn(name string) *Column {
	return t.columnMap[name]
}
