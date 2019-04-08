package model

import (
	"context"
	"fmt"
)

type PersonWithLock struct {
	personWithLockBase
}

// Create a new PersonWithLock object and initialize to default values.
func NewPersonWithLock(ctx context.Context) *PersonWithLock {
	o := PersonWithLock{}
	o.Initialize(ctx)
	return &o
}

// Initialize or re-initialize a PersonWithLock database object to default values.
func (o *PersonWithLock) Initialize(ctx context.Context) {
	o.personWithLockBase.Initialize()
	// Add your own initializations here
}

// String implements the Stringer interface and returns the default label for the object as it appears in html lists.
// Typically you would change this to whatever was pertinent to your application.
func (o *PersonWithLock) String() string {
	return fmt.Sprintf("Object id %v", o.PrimaryKey())
}
