package model

// This is the implementation file for the Milestone ORM object.
// This is where you build the api to your data model for you web application and potentially mobile apps.
// Your edits to this file will be preserved.

import (
	"context"
	"encoding/gob"
	"fmt"
)

type Milestone struct {
	milestoneBase
}

// Create a new Milestone object and initialize to default values.
func NewMilestone() *Milestone {
	o := new(Milestone)
	o.Initialize()
	return o
}

// Initialize or re-initialize a Milestone database object to default values.
func (o *Milestone) Initialize() {
	o.milestoneBase.Initialize()
	// Add your own initializations here
}

// String implements the Stringer interface and returns the default label for the object as it appears in html lists.
// Typically you would change this to whatever was pertinent to your application.
func (o *Milestone) String() string {
	if o == nil {
		return "" // Possibly - Select One -?
	}
	return fmt.Sprintf("Object id %v", o.PrimaryKey())
}

// QueryMilestones returns a new builder that gives you general purpose access to the Milestone records
// in the database. Its here to give public access to the query builder, but you can remove it if you do not need it.
func QueryMilestones(ctx context.Context) *MilestonesBuilder {
	return queryMilestones(ctx)
}

// queryMilestones creates a new builder and is the central spot where all queries are directed.
// You can modify this function to enforce restrictions on queries, for example to make sure the user is authorized to
// access the data.
func queryMilestones(ctx context.Context) *MilestonesBuilder {
	return newMilestoneBuilder()
}

// DeleteMilestone deletes the given record from the database. Note that you can also delete
// loaded Milestone objects by calling Delete on them.
func DeleteMilestone(ctx context.Context, pk string) {
	deleteMilestone(ctx, pk)
}

func init() {
	gob.RegisterName("ExampleMilestone", new(Milestone))
}
