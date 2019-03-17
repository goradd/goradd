package model

import (
	"context"
	"fmt"
)

type Address struct {
	addressBase
}

// Create a new Address object and initialize to default values.
func NewAddress(ctx context.Context) *Address {
	o := Address{}
	o.Initialize(ctx)
	return &o
}

// Initialize or re-initialize a Address database object to default values.
func (o *Address) Initialize(ctx context.Context) {
	o.addressBase.Initialize()
	// Add your own initializations here
}

// String implements the Stringer interface and returns the default label for the object as it appears in html lists.
// Typically you would change this to whatever was pertinent to your application.
func (o *Address) String() string {
	return fmt.Sprintf("Object id %v", o.PrimaryKey())
}
