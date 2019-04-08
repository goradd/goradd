package model

import (
	"context"
	"fmt"
)

type Project struct {
	projectBase
}

// Create a new Project object and initialize to default values.
func NewProject(ctx context.Context) *Project {
	o := Project{}
	o.Initialize(ctx)
	return &o
}

// Initialize or re-initialize a Project database object to default values.
func (o *Project) Initialize(ctx context.Context) {
	o.projectBase.Initialize()
	// Add your own initializations here
}

// String implements the Stringer interface and returns the default label for the object as it appears in html lists.
// Typically you would change this to whatever was pertinent to your application.
func (o *Project) String() string {
	return fmt.Sprintf("Object id %v", o.PrimaryKey())
}
