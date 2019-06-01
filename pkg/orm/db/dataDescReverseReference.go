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
	// DbTable is the database table on the "one" end of the relationship. Its the table that the ReverseReference belongs to.
	DbTable string
	// DbColumn is the database column that is referred to from the many end. This is most likely the primary key of DbTable.
	DbColumn string
	// AssociatedTableName is the table on the "many" end that is pointing to the table containing the ReverseReference.
	AssociatedTableName string
	// AssociatedColumnName is the column on the "many" end that is pointing to the table containing the ReverseReference. It is a foreign-key.
	AssociatedColumnName string
	// GoName is the name used to represent an object in the reverse relationship
	GoName string
	// GoPlural is the name used to represent the group of objects in the reverse relationship
	GoPlural string
	// GoType is the type of object in the collection of "many" objects, which corresponds to the name of the struct corresponding to the table
	GoType string
	// GoTypePlural is the plural of the type of object in the collection of "many" objects
	GoTypePlural string
	// IsUnique is true if the ReverseReference is unique. A unique reference creates a one-to-one relationship rather than one-to-many
	IsUnique bool
	//Options maps.SliceMap
}

func (r *ReverseReference) ObjName(dd *DatabaseDescription) string {
	if r.IsUnique {
		return dd.AssociatedObjectPrefix + r.GoName
	} else {
		return dd.AssociatedObjectPrefix + r.GoPlural
	}
}

func (r *ReverseReference) MapName() string {
	if r.IsUnique {
		return "" // no map
	} else {
		return "m" + r.GoPlural
	}
}

// AssociatedGoName returns the name of the column that is pointing back to us. The name returned
// is the Go name that we could use to name the referenced object.
func (r *ReverseReference) AssociatedGoName() string {
	return UpperCaseIdentifier(r.AssociatedColumnName)
}

func (r *ReverseReference) JsonKey(dd *DatabaseDescription) string {
	return LowerCaseIdentifier(strings.TrimSuffix(r.AssociatedColumnName, dd.ForeignKeySuffix))
}
