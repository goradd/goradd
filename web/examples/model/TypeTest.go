package model

import (
	"context"
	"fmt"
)

type TypeTest struct {
	typeTestBase
}

// Create a new TypeTest object and initialize to default values.
func NewTypeTest(ctx context.Context) *TypeTest {
	o := TypeTest{}
	o.Initialize(ctx)
	return &o
}

// Initialize or re-initialize a TypeTest database object to default values.
func (o *TypeTest) Initialize(ctx context.Context) {
	o.typeTestBase.Initialize()
	// Add your own initializations here
}

// String implements the Stringer interface and returns the default label for the object as it appears in html lists.
// Typically you would change this to whatever was pertinent to your application.
func (o *TypeTest) String() string {
	return fmt.Sprintf("Object id %v", o.PrimaryKey())
}
