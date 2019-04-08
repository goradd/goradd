package model

import (
	"context"
	"fmt"
)

type UnsupportedTypes struct {
	unsupportedTypesBase
}

// Create a new UnsupportedTypes object and initialize to default values.
func NewUnsupportedTypes(ctx context.Context) *UnsupportedTypes {
	o := UnsupportedTypes{}
	o.Initialize(ctx)
	return &o
}

// Initialize or re-initialize a UnsupportedTypes database object to default values.
func (o *UnsupportedTypes) Initialize(ctx context.Context) {
	o.unsupportedTypesBase.Initialize()
	// Add your own initializations here
}

// String implements the Stringer interface and returns the default label for the object as it appears in html lists.
// Typically you would change this to whatever was pertinent to your application.
func (o *UnsupportedTypes) String() string {
	return fmt.Sprintf("Object id %v", o.PrimaryKey())
}
