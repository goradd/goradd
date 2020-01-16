package model

// This is the implementation file for the Reverse ORM object.
// This is where you build the api to your data model for you web application and potentially mobile apps.
// Your edits to this file will be preserved.

import (
	"context"
	"fmt"
)

type Reverse struct {
	reverseBase
}

// Create a new Reverse object and initialize to default values.
func NewReverse() *Reverse {
	o := new(Reverse)
	o.Initialize()
	return o
}

// Initialize or re-initialize a Reverse database object to default values.
func (o *Reverse) Initialize() {
	o.reverseBase.Initialize()
	// Add your own initializations here
}

// String implements the Stringer interface and returns the default label for the object as it appears in html lists.
// Typically you would change this to whatever was pertinent to your application.
func (o *Reverse) String() string {
	if o == nil {
		return "" // Possibly - Select One -?
	}
	return fmt.Sprintf("Object id %v", o.PrimaryKey())
}

// QueryReverses returns a new builder that gives you general purpose access to the Reverse records
// in the database. Its here to give public access to the query builder, but you can remove it if you do not need it.
func QueryReverses(ctx context.Context) *ReversesBuilder {
	return queryReverses(ctx)
}

// queryReverses creates a new builder and is the central spot where all queries are directed.
// You can modify this function to enforce restrictions on queries, for example to make sure the user is authorized to
// access the data.
func queryReverses(ctx context.Context) *ReversesBuilder {
	return newReverseBuilder()
}

// DeleteReverse deletes the given record from the database. Note that you can also delete
// loaded Reverse objects by calling Delete on them.
func DeleteReverse(ctx context.Context, pk string) {
	deleteReverse(ctx, pk)
}
