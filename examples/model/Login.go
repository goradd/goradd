package model

import (
	"context"
	"fmt"
)

type Login struct {
	loginBase
}

// Create a new Login object and initialize to default values.
func NewLogin(ctx context.Context) *Login {
	o := Login{}
	o.Initialize(ctx)
	return &o
}

// Initialize or re-initialize a Login database object to default values.
func (o *Login) Initialize(ctx context.Context) {
	o.loginBase.Initialize()
	// Add your own initializations here
}

// String implements the Stringer interface and returns the default label for the object as it appears in html lists.
// Typically you would change this to whatever was pertinent to your application.
func (o *Login) String() string {
	return fmt.Sprintf("Object id %v", o.PrimaryKey())
}
