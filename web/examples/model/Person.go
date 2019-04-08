package model

import (
	"context"
	"fmt"
)

type Person struct {
	personBase
}

// Create a new Person object and initialize to default values.
func NewPerson(ctx context.Context) *Person {
	o := Person{}
	o.Initialize(ctx)
	return &o
}

// Initialize or re-initialize a Person database object to default values.
func (o *Person) Initialize(ctx context.Context) {
	o.personBase.Initialize()
	// Add your own initializations here
}

// String implements the Stringer interface and returns the default label for the object as it appears in html lists.
// Typically you would change this to whatever was pertinent to your application.
func (o *Person) String() string {
	return fmt.Sprintf("Object id %v", o.PrimaryKey())
}
