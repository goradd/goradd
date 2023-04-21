package db

// ForeignKeyInfo is additional information to describe what a foreign key points to.
// Cross database foreign keys are not supported.
type ForeignKeyInfo struct {
	// ReferencedTable is the name of the table on the other end of the foreign key
	ReferencedTable string
	// ReferencedColumn is the database column name in the linked table that matches this column name.
	// Often that is the primary key of the other table.
	ReferencedColumn string
	// UpdateAction indicates how the column will react when the other end of the relationship's value changes.
	// Some databases react automatically, but databases or tables that do not support
	// foreign keys will need to be updated manually.
	UpdateAction FKAction
	// DeleteAction indicates how the column will react when the other end of the relationship's record is deleted.
	// Some databases react automatically, but databases or tables that do not support
	// foreign keys will need to be updated manually.
	DeleteAction FKAction
	// GoName is the name we should use to refer to the related object
	GoName string
	// GoType is the type of the related object
	GoType string
	// GoTypePlural is the plural version of the type when referring to groups of related objects
	GoTypePlural string
	// isEnum is true if this is a related type
	IsEnum bool
	// RR is filled in by the analyzer and represents a reverse reference relationship
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
func (fk *ForeignKeyInfo) GoVarName() string {
	return "obj" + fk.GoName
}

// FKAction indicates how the database handles situations when one side of a relationship is deleted or the key
// is changed.
type FKAction int

// The foreign key actions tell us what the database will do automatically if a foreign key object is changed. This allows
// us to do the appropriate thing when we detect in the ORM that a linked object is changing.
const (
	// FKActionNone indicates the database does not support foreign key actions, so we should deal with it ourselves.
	FKActionNone FKAction = iota
	// FKActionSetNull indicates the foreign key will be set to null
	FKActionSetNull
	// FKActionSetDefault indicates the foreign key will be set to a default value
	FKActionSetDefault
	// FKActionCascade indicates the foreign key will automatically update values during an update, and
	// automatically delete related records during a delete.
	FKActionCascade
	// FKActionRestrict indicates the database will prevent the action from happening by likely panicking with an error
	FKActionRestrict
)

func (a FKAction) String() string {
	switch a {
	case FKActionSetNull:
		return "Null"
	case FKActionSetDefault:
		return "Default"
	case FKActionCascade:
		return "Cascade"
	case FKActionRestrict:
		return "Restrict"
	default:
		return "" // None
	}
}

func FKActionFromString(s string) FKAction {
	switch s {
	case "Null":
		return FKActionSetNull
	case "Default":
		return FKActionSetDefault
	case "Cascade":
		return FKActionCascade
	case "Restrict":
		return FKActionRestrict
	default:
		return FKActionNone
	}
}
