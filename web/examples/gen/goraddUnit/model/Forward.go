package model

// This is the implementation file for the Forward ORM object.
// This is where you build the api to your data model for you web application and potentially mobile apps.
// Your edits to this file will be preserved.

import (
	"context"
	"fmt"
)

type Forward struct {
	forwardBase
}

// Create a new Forward object and initialize to default values.
func NewForward(ctx context.Context) *Forward {
	o := new(Forward)
	o.Initialize(ctx)
	return o
}

// Initialize or re-initialize a Forward database object to default values.
func (o *Forward) Initialize(ctx context.Context) {
	o.forwardBase.Initialize()
	// Add your own initializations here
}

// String implements the Stringer interface and returns the default label for the object as it appears in html lists.
// Typically you would change this to whatever was pertinent to your application.
func (o *Forward) String() string {
	if o == nil {
		return "" // Possibly - Select One -?
	}
	return fmt.Sprintf("Object id %v", o.PrimaryKey())
}

// QueryForwards returns a new builder that gives you general purpose access to the Forward records
// in the database. Its here to give public access to the query builder, but you can remove it if you do not need it.
func QueryForwards(ctx context.Context) *ForwardsBuilder {
	return queryForwards(ctx)
}

// queryForwards creates a new builder and is the central spot where all queries are directed.
// You can modify this function to enforce restrictions on queries, for example to make sure the user is authorized to
// access the data.
func queryForwards(ctx context.Context) *ForwardsBuilder {
	return newForwardBuilder()
}

// DeleteForward deletes the given record from the database. Note that you can also delete
// loaded Forward objects by calling Delete on them.
func DeleteForward(ctx context.Context, pk string) {
	deleteForward(ctx, pk)
}
