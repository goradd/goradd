package model

// This is the implementation file for the UnsupportedTypes ORM object.
// This is where you build the api to your data model for you web application and potentially mobile apps.
// Your edits to this file will be preserved.

import (
	"context"
	"fmt"
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
// in the database. Its here to give public access to the query builder, but you can remove it if you do not need it.
func QueryUnsupportedTypes(ctx context.Context) *UnsupportedTypesBuilder {
	return queryUnsupportedTypes(ctx)
}

// queryUnsupportedTypes creates a new builder and is the central spot where all queries are directed.
// You can modify this function to enforce restrictions on queries, for example to make sure the user is authorized to
// access the data.
func queryUnsupportedTypes(ctx context.Context) *UnsupportedTypesBuilder {
	return newUnsupportedTypesBuilder()
}

// DeleteUnsupportedTypes deletes the given record from the database. Note that you can also delete
// loaded UnsupportedTypes objects by calling Delete on them.
func DeleteUnsupportedTypes(ctx context.Context, pk string) {
	deleteUnsupportedTypes(ctx, pk)
}
