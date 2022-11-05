package pgsql

type Options struct {
	// TypeTableSuffix is the suffix in the name of a table that tells GoRADD to treat
	// the table as a type table. Defaults to "_type" if not set.
	TypeTableSuffix string
	// AssociationTableSuffix is the suffix in the name of a table that tells GoRADD to
	// treat the table as an association table. Defaults to "_assn".
	AssociationTableSuffix string
	// ForeignKeySuffix is the suffix to strip off the ends of names of foreign keys when converting
	// them to internal names. For example, if the suffix is "_id", and a column named
	// manager_id a "project" table is a foreign key to a "person" table, then GoRADD will
	// create "Person" objects with the name "Manager" inside the "Project" object.
	// A suffix is required since it will also create a "ManagerID" member variable, and
	// without the suffix the two values will have the same name.
	// The default is "_id".
	ForeignKeySuffix string
	// UseQualifiedNames will force goradd to prepend schema names in front of table names.
	// This is required to clear up ambiguity when different schemas have the same table name.
	// Postgres documentation discourages the use of repeated table names in different schemas.
	// If you know you do not have repeats, you can leave this option false.
	UseQualifiedNames bool
	// Schemas lets you specify which specific schemas you want to select.
	// If no schemas are specified, then all available schemas will be used.
	// Note that you cannot specify a subset of schemas, ones that do not have overlapping
	// table names, to avoid setting useQualifiedNames to true. The Postgres default search
	// path will be used to find unqualified names during a query, and it might choose
	// schemas not in the schema list you specify.
	Schemas []string
}

// DefaultOptions returns default database analysis options for MySQL databases.
func DefaultOptions() Options {
	return Options{
		TypeTableSuffix:        "_type",
		AssociationTableSuffix: "_assn",
		ForeignKeySuffix:       "_id",
		UseQualifiedNames:      false,
		Schemas:                nil,
	}
}
