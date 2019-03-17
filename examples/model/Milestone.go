package model

import (
	"context"
	"fmt"
)

type Milestone struct {
	milestoneBase
}

// Create a new Milestone object and initialize to default values.
func NewMilestone(ctx context.Context) *Milestone {
	o := Milestone{}
	o.Initialize(ctx)
	return &o
}

// Initialize or re-initialize a Milestone database object to default values.
func (o *Milestone) Initialize(ctx context.Context) {
	o.milestoneBase.Initialize()
	// Add your own initializations here
}

// String implements the Stringer interface and returns the default label for the object as it appears in html lists.
// Typically you would change this to whatever was pertinent to your application.
func (o *Milestone) String() string {
	return fmt.Sprintf("Object id %v", o.PrimaryKey())
}
