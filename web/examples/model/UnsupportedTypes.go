package model

// This is the implementation file for the UnsupportedTypes ORM object.
// This is where you build the api to your data model for you web application and potentially mobile apps.
// Your edits to this file will be preserved.

import (
	"context"
	"fmt"

	"github.com/goradd/goradd/pkg/orm/query"
)

type UnsupportedTypes struct {
	unsupportedTypesBase
}

// Create a new UnsupportedTypes object and initialize to default values.
func NewUnsupportedTypes(ctx context.Context) *UnsupportedTypes {
	o := new(UnsupportedTypes)
	o.Initialize(ctx)
	return o
}

// Initialize or re-initialize a UnsupportedTypes database object to default values.
func (o *UnsupportedTypes) Initialize(ctx context.Context) {
	o.unsupportedTypesBase.Initialize()
	// Add your own initializations here
}

// String implements the Stringer interface and returns the default label for the object as it appears in html lists.
// Typically you would change this to whatever was pertinent to your application.
func (o *UnsupportedTypes) String() string {
	if o == nil {
		return "" // Possibly - Select One -?
	}
	return fmt.Sprintf("Object id %v", o.PrimaryKey())
}

// QueryUnsupportedTypes returns a new builder that gives you general purpose access to the UnsupportedTypes records
// in the database. This is useful for quick queries of the database during development, but eventually you
// should remove this function and move those queries to more specific calls in this file.
func QueryUnsupportedTypes() *UnsupportedTypesBuilder {
	return queryUnsupportedTypes()
}

// LoadUnsupportedTypes queries for a single UnsupportedTypes object by primary key.
// joinOrSelectNodes lets you provide nodes for joining to other tables or selecting specific fields. Table nodes will
// be considered Join nodes, and column nodes will be Select nodes. See Join() and Select() for more info.
// If you need a more elaborate query, use QueryUnsupportedTypes() to start a query builder.
func LoadUnsupportedTypes(ctx context.Context, pk string, joinOrSelectNodes ...query.NodeI) *UnsupportedTypes {
	return loadUnsupportedTypes(ctx, pk, joinOrSelectNodes...)
}

// LoadUnsupportedTypesByTypeSerial queries for a single UnsupportedTypes object by the given unique index values.
// joinOrSelectNodes lets you provide nodes for joining to other tables or selecting specific fields. Table nodes will
// be considered Join nodes, and column nodes will be Select nodes. See Join() and Select() for more info.
// If you need a more elaborate query, use QueryUnsupportedTypes() to start a query builder.
func LoadUnsupportedTypesByTypeSerial(ctx context.Context, type_serial string, joinOrSelectNodes ...query.NodeI) *UnsupportedTypes {
	return loadUnsupportedTypesByTypeSerial(ctx, type_serial, joinOrSelectNodes...)
}

// DeleteUnsupportedTypes deletes the give record from the database. Note that you can also delete
// loaded UnsupportedTypes objects by calling Delete on them.
func DeleteUnsupportedTypes(ctx context.Context, pk string) {
	deleteUnsupportedTypes(ctx, pk)
}
