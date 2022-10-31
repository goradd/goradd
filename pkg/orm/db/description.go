package db

// DatabaseDescription generically describes a database to GoRADD. It is sent to NewDatabase() to create a
// DB object that is used internally by GoRADD to access the database. DatabaseDescription should be able to be
// inferred by reading the structure of SQL databases, or read directly from an import file.
type DatabaseDescription struct {
	// Tables are the tables in the database
	Tables []TableDescription
	// MM are the many-to-many links between tables. In SQL databases, these are actual tables,
	// but in NoSQL, these might be array fields on either side of the relationship.
	MM []ManyManyDescription

	// The prefix for related objects.
	AssociatedObjectPrefix string
}

// TableDescription describes a database object to GoRADD.
type TableDescription struct {
	// Name is the name of the database table or collection.
	Name string
	// Columns is a list of ColumnDescriptions, one for each column in the table.
	// The first columns are the primary keys. Usually there is just one primary key.
	Columns []ColumnDescription
	// Indexes are the indexes defined in the database. Unique indexes will result in LoadBy* functions.
	Indexes []IndexDescription
	// TypeData is the data of the type table if this is a type table. The data structure must match that of the columns.
	TypeData []map[string]interface{}

	// Comment is an optional comment about the table
	Comment string
	// Options are key-value settings that can be used to further describe code generation
	Options map[string]interface{}
}

// ColumnDescription describes a field of a database object to GoRADD.
type ColumnDescription struct {
	// Name is the name of the column in the database. This is blank if this is a "virtual" table for sql tables like an association or virtual attribute query.
	Name string
	// NativeType is the type of the column as described by the database itself.
	NativeType string
	//  GoType is the goradd defined column type
	GoType string
	// SubType has additional information to the type of column that can help control code generation
	// When column type is "time.Time", the column will default to both a date and time format. You can also make it:
	//   date (which means date only)
	//   time (time only)
	//   timestamp (we will track the modification time of the table here)
	//   auto timestamp (the database is automatically updating this timestamp for us)
	SubType string
	// MaxCharLength is the maximum length of characters to allow in the column if a string type column.
	// If the database has the ability to specify this, this will correspond to what is specified.
	// In any case, we will generate code to prevent fields from getting bigger than this. Zero indicates there is
	// no length checking or limiting.
	MaxCharLength uint64
	// DefaultValue is the default value as specified by the database. We will initialize new ORM objects
	// with this value. It will be cast to the corresponding GO type.
	DefaultValue interface{}
	// MaxValue is the maximum value allowed for numeric values. This can be used by UI objects to tell the user what the limits are.
	MaxValue interface{}
	// MinValue is the minimum value allowed for numeric values. This can be used by UI objects to tell the user what the limits are.
	MinValue interface{}
	// IsId is true if this column represents a unique identifier generated by the database
	IsId bool
	// IsPk is true if this is the primary key column. PK's do not necessarily need to be ID columns, and if not, we will need to do our own work to generate unique PKs.
	IsPk bool
	// IsNullable is true if the column can be given a NULL value
	IsNullable bool
	// IsUnique indicates that the field holds unique values
	IsUnique bool
	// ForeignKey is additional information describing a foreign key relationship
	ForeignKey *ForeignKeyDescription
	// Comment is the contents of the comment associated with this column
	Comment string
	// Options are key-value settings that can be used to further describe code generation
	Options map[string]interface{}
}

// ForeignKeyDescription describes a pointer from one database object to another database object.
type ForeignKeyDescription struct {
	//DbKey string	// We don't support cross database foreign keys yet. Someday maybe.
	// ReferencedTable is the name of the table on the other end of the foreign key
	ReferencedTable string
	// ReferencedColumn is the database column name in the linked table that matches this column. Often that is the primary key of the other table.
	ReferencedColumn string
	// UpdateAction indicates how the database will react when the referenced item's id changes.
	UpdateAction string
	// DeleteAction indicates how the database will react when the referenced item is deleted.
	DeleteAction string
	// IsUnique is true if the reference is one-to-one
	IsUnique bool

	// GoName is the name we should use to refer to the related object. Leave blank to get a computed value.
	GoName string
	// ReverseName is the name that the reverse reference should use to refer to the collection of objects pointing to it.
	// Leave blank to get a "ThisAsThat" type default name. The lower-case version of this name will be used as a column name
	// to store the values if using a NoSQL database.
	ReverseName string
}

// IndexDescription gives us information about how columns are indexed.
// If a column has a unique index, it will get a corresponding "LoadBy" function in its table's model.
// Otherwise, it will get a corresponding "LoadSliceBy" function.
type IndexDescription struct {
	// IsUnique indicates whether the index is unique
	IsUnique bool
	// ColumnNames are the columns that are part of the index
	ColumnNames []string
}

// ManyManyDescription describes a many-to-many relationship table that contains a two-way pointer between database objects.
type ManyManyDescription struct {
	// Table1 is the name of the first table that is part of the relationship. The private key of that table will be referred to.
	Table1 string
	// Column1 is the database column name. For SQL databases, this is the name of the column in the assn table. For
	// NoSQL, this is the name of the column that will be used to store the ids of the other side. This is optional for
	// NoSQL, as one will be created based on the table names if left blank.
	Column1 string
	// GoName1 is the singular name of the object that Table2 will use to refer to Table1 objects.
	GoName1 string
	// GoPlural1 is the plural name of the object that Table2 will use to refer to Table1 objects.
	GoPlural1 string

	// Table2 is the name of the second table that is part of the relationship. The private key of that table will be referred to.
	Table2 string
	// Column2 is the database column name. For SQL databases, this is the name of the column in the assn table. For
	// NoSQL, this is the name of the column that will be used to store the ids of the other side. This is optional for
	// NoSQL, as one will be created based on the table names if left blank.
	Column2 string
	// GoName2 is the singular name of the object that Table1 will use to refer to Table2 objects.
	GoName2 string
	// GoPlural2 is the plural name of the object that Table1 will use to refer to Table2 objects.
	GoPlural2 string

	// AssnTableName is the name of the intermediate association table that will be used to create the relationship. This is
	// needed for SQL databases, but not for NoSQL, as NoSQL will create additional array columns on each side of the relationship.
	AssnTableName string
}
