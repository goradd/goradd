package mysql

type Options struct {
	// EnumTableSuffix is the suffix in the name of a table that tells GoRADD to treat
	// the table as a enum table. Defaults to "_enum" if not set.
	EnumTableSuffix string
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
}

// DefaultOptions returns default database analysis options for MySQL databases.
func DefaultOptions() Options {
	return Options{
		EnumTableSuffix:        "_enum",
		AssociationTableSuffix: "_assn",
		ForeignKeySuffix:       "_id",
	}
}
