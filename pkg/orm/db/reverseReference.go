package db

import (
	"strings"
)

// ReverseReference represents a kind of virtual column that is a result of a foreign-key
// pointing back to this column. This is the "one" side of a one-to-many relationship. Or, if
// the relationship is unique, this creates a one-to-one relationship.
// In SQL, since there is only a one-way foreign key, the side being pointed at does not have any direct
// data in a table indicating the relationship. We create a ReverseReference during data analysis and include
// it with the table description so that the table can know about the relationship and use it when doing queries.
type ReverseReference struct {
	// DbColumn is only used in NoSQL databases, and is the name of a column that will hold the pk(s) of the referring column(s)
	DbColumn string
	// AssociatedTableName is the table on the "many" end that is pointing to the table containing the ReverseReference.
	AssociatedTable *Table
	// AssociatedColumn is the column on the "many" end that is pointing to the table containing the ReverseReference. It is a foreign-key.
	AssociatedColumn *Column
	// AssociatedPkType is the go type of the primary key column of the AssociatedTable
	AssociatedPkType string
	// GoName is the name used to represent an object in the reverse relationship
	GoName string
	// GoPlural is the name used to represent the group of objects in the reverse relationship
	GoPlural string
	// GoType is the type of object in the collection of "many" objects, which corresponds to the name of the struct corresponding to the table
	GoType string
	// GoTypePlural is the plural of the type of object in the collection of "many" objects
	GoTypePlural string
}

func (r *ReverseReference) ObjName(dd *Database) string {
	if r.IsUnique() {
		return dd.AssociatedObjectPrefix + r.GoName
	} else {
		return dd.AssociatedObjectPrefix + r.GoPlural
	}
}

func (r *ReverseReference) MapName() string {
	if r.IsUnique() {
		return "" // no map
	} else {
		return "m" + r.GoPlural
	}
}

// AssociatedGoName returns the name of the column that is pointing back to us. The name returned
// is the Go name that we could use to name the referenced object.
func (r *ReverseReference) AssociatedGoName() string {
	return UpperCaseIdentifier(r.AssociatedColumn.DbName)
}

func (r *ReverseReference) JsonKey(dd *Database) string {
	return LowerCaseIdentifier(strings.TrimSuffix(r.AssociatedColumn.DbName, dd.ForeignKeySuffix))
}

func (r *ReverseReference) IsUnique() bool {
	return r.AssociatedColumn.IsUnique
}

func (r *ReverseReference) IsNullable() bool {
	return r.AssociatedColumn.IsNullable
}

func (r *ReverseReference) AssociatedTableName() string {
	return r.AssociatedTable.DbName
}



