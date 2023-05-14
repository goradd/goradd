package db

import "github.com/goradd/goradd/pkg/strings"

// The ManyManyReference structure is used by the templates during the codegen process to describe a many-to-many relationship.
// Underlying the structure is an association table that has two values that are foreign keys pointing
// to the records that are linked. The names of these fields will determine the names of the corresponding accessors
// in each of the model objects. This allows multiple of these many-many relationships to exist
// on the same tables but for different purposes.
type ManyManyReference struct {
	// AssnTableName is the database table that links the two associated tables together.
	AssnTableName string
	// AssnSourceColumnName is the database column in the association table that points at the source table's primary key.
	AssnSourceColumnName string
	// AssnDestColumnName is the database column in the association table that points at the destination table's primary key.
	AssnDestColumnName string
	// DestinationTableName is the database table being linked (the table that we are joining to)
	DestinationTableName string

	// GoName is the name used to refer to an object on the other end of the reference.
	// It is not the same as the object type. For example TeamMember would refer to a Person type.
	// This is derived from the AssnDestColumnName but can be overridden by comments in the column.
	GoName string
	// GoPlural is the name used to refer to the group of objects on the other end of the reference.
	// For example, TeamMembers. This is derived from the AssnDestColumnName but can be overridden by
	// a comment in the table.
	GoPlural string

	// objectType is the type of the model object that this is pointing to.
	objectType string
	// objectTypes is the plural of the type of the model object that this is pointing to.
	objectTypes    string
	primaryKey     string
	primaryKeyType string

	// SupportsForeignKeys indicates that updates and deletes are automatically handled by the database engine.
	// If this is false, the code generator will need to manually update these items.
	SupportsForeignKeys bool

	// IsEnumAssociation is true if this is a many-many relationship with an enum table
	IsEnumAssociation bool

	// MM is the many-many reference on the other end of the relationship that points back to this one.
	MM *ManyManyReference
}

// JsonKey returns the key used when referring to the associated objects in JSON.
func (m *ManyManyReference) JsonKey(dd *Model) string {
	return strings.LcFirst(m.GoPlural)
}

// ObjectType returns the name of the object type the association links to.
func (m *ManyManyReference) ObjectType() string {
	return m.objectType
}

// ObjectTypes returns the plural name of the object type the association links to.
func (m *ManyManyReference) ObjectTypes() string {
	return m.objectTypes
}

// PrimaryKeyType returns the Go type of the primary key of the object the association links to.
func (m *ManyManyReference) PrimaryKeyType() string {
	return m.primaryKeyType
}

// PrimaryKey returns the database field name of the primary key of the object the association links to.
func (m *ManyManyReference) PrimaryKey() string {
	return m.primaryKey
}
