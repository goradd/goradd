package model

import (
	"context"
	"fmt"
)

type Tmp struct {
	tmpBase
}

// Create a new Tmp object and initialize to default values.
func NewTmp(ctx context.Context) *Tmp {
	o := Tmp{}
	o.Initialize(ctx)
	return &o
}

// Initialize or re-initialize a Tmp database object to default values.
func (o *Tmp) Initialize(ctx context.Context) {
	o.tmpBase.Initialize()
	// Add your own initializations here
}

// String implements the Stringer interface and returns the default label for the object as it appears in html lists.
// Typically you would change this to whatever was pertinent to your application.
func (o *Tmp) String() string {
	return fmt.Sprintf("Object id %v", o.PrimaryKey())
}
