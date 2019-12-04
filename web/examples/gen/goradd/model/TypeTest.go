package model

// This is the implementation file for the TypeTest ORM object.
// This is where you build the api to your data model for you web application and potentially mobile apps.
// Your edits to this file will be preserved.

import (
	"context"
	"fmt"
)

type TypeTest struct {
	typeTestBase
}

// Create a new TypeTest object and initialize to default values.
func NewTypeTest() *TypeTest {
	o := new(TypeTest)
	o.Initialize()
	return o
}

// Initialize or re-initialize a TypeTest database object to default values.
func (o *TypeTest) Initialize() {
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
// in the database. Its here to give public access to the query builder, but you can remove it if you do not need it.
func QueryTypeTests(ctx context.Context) *TypeTestsBuilder {
	return queryTypeTests(ctx)
}

// queryTypeTests creates a new builder and is the central spot where all queries are directed.
// You can modify this function to enforce restrictions on queries, for example to make sure the user is authorized to
// access the data.
func queryTypeTests(ctx context.Context) *TypeTestsBuilder {
	return newTypeTestBuilder()
}

// DeleteTypeTest deletes the given record from the database. Note that you can also delete
// loaded TypeTest objects by calling Delete on them.
func DeleteTypeTest(ctx context.Context, pk string) {
	deleteTypeTest(ctx, pk)
}
