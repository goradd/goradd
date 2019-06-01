package db

import "github.com/goradd/gengen/pkg/maps"

// The ManyManyReference structure is used by the templates during the codegen process to describe a many-to-many relationship.
type ManyManyReference struct {
	// AssnTableName is the database table creating the association. NoSQL: The originating table. SQL: The association table
	AssnTableName string
	// AssnColumnName is the column creating the association. NoSQL: The table storing the array of ids on the other end. SQL: the column in the association table pointing towards us.
	AssnColumnName string

	// AssociatedTableName is the database table being linked. NoSQL & SQL: The table we are joining to
	AssociatedTableName string
	// AssociatedColumnName is the database column being linked. NoSQL: table point backwards to us. SQL: Column in association table pointing forwards to refTable
	AssociatedColumnName string
	// AssociatedObjectName is the go name of the object created by this reference
	AssociatedObjectName string

	// GoName is the name used to refer to an object on the other end of the reference.
	GoName string
	// GoPlural is the name used to refer to the group of objects on the other end of the reference.
	GoPlural string

	// IsTypeAssociation is true if this is a many-many relationship with a type table
	IsTypeAssociation bool
	// Options are the key-value options taken from the Comments in the association table, if there is one.
	Options maps.SliceMap

	// MM is the many-many reference on the other end of the relationship that points back to this one.
	MM *ManyManyReference
}

func (m *ManyManyReference) ObjName(dd *DatabaseDescription) string {
	return dd.AssociatedObjectPrefix + m.GoPlural
}

func (m *ManyManyReference) JsonKey(dd *DatabaseDescription) string {
	return LowerCaseIdentifier(m.GoPlural)
}
