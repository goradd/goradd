package model

// This is the implementation file for the TypeTest ORM object.
// This is where you build the api to your data model for you web application and potentially mobile apps.
// Your edits to this file will be preserved.

import (
	"context"
	"fmt"

	"github.com/goradd/goradd/pkg/orm/query"
)

type TypeTest struct {
	typeTestBase
}

// Create a new TypeTest object and initialize to default values.
func NewTypeTest(ctx context.Context) *TypeTest {
	o := new(TypeTest)
	o.Initialize(ctx)
	return o
}

// Initialize or re-initialize a TypeTest database object to default values.
func (o *TypeTest) Initialize(ctx context.Context) {
	o.typeTestBase.Initialize()
	// Add your own initializations here
}

// String implements the Stringer interface and returns the default label for the object as it appears in html lists.
// Typically you would change this to whatever was pertinent to your application.
func (o *TypeTest) String() string {
	if o == nil {
		return "" // Possibly - Select One -?
	}
	return fmt.Sprintf("Object id %v", o.PrimaryKey())
}

// QueryTypeTests returns a new builder that gives you general purpose access to the TypeTest records
// in the database. This is useful for quick queries of the database during development, but eventually you
// should remove this function and move those queries to more specific calls in this file.
func QueryTypeTests() *TypeTestsBuilder {
	return queryTypeTests()
}

// LoadTypeTest queries for a single TypeTest object by primary key.
// joinOrSelectNodes lets you provide nodes for joining to other tables or selecting specific fields. Table nodes will
// be considered Join nodes, and column nodes will be Select nodes. See Join() and Select() for more info.
// If you need a more elaborate query, use QueryTypeTests() to start a query builder.
func LoadTypeTest(ctx context.Context, pk string, joinOrSelectNodes ...query.NodeI) *TypeTest {
	return loadTypeTest(ctx, pk, joinOrSelectNodes...)
}

// DeleteTypeTest deletes the give record from the database. Note that you can also delete
// loaded TypeTest objects by calling Delete on them.
func DeleteTypeTest(ctx context.Context, pk string) {
	deleteTypeTest(ctx, pk)
}
