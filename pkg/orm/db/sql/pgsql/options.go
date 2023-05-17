package pgsql

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
	// UseQualifiedNames will force goradd to prepend schema names in front of the generated object names.
	// This is required to clear up ambiguity when different schemas have the same table name, and
	// so will generate the same object name without the schema.
	// Postgres documentation discourages the use of repeated table names in different schemas.
	// If you know you have repeats, or you just want to force the schema names to appear in
	// the object names, you can leave this option false.
	UseQualifiedNames bool
	// Schemas lets you specify which specific schemas you want to select.
	// If no schemas are specified, then all available schemas will be used.
	Schemas []string
}

// DefaultOptions returns default database analysis options for MySQL databases.
func DefaultOptions() Options {
	return Options{
		EnumTableSuffix:        "_enum",
		AssociationTableSuffix: "_assn",
		ForeignKeySuffix:       "_id",
		UseQualifiedNames:      false,
		Schemas:                nil,
	}
}
