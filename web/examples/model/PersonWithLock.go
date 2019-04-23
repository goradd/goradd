package model

// This is the implementation file for the PersonWithLock ORM object.
// This is where you build the api to your data model for you web application and potentially mobile apps.
// Your edits to this file will be preserved.

import (
	"context"
	"fmt"

	"github.com/goradd/goradd/pkg/orm/query"
)

type PersonWithLock struct {
	personWithLockBase
}

// Create a new PersonWithLock object and initialize to default values.
func NewPersonWithLock(ctx context.Context) *PersonWithLock {
	o := new(PersonWithLock)
	o.Initialize(ctx)
	return o
}

// Initialize or re-initialize a PersonWithLock database object to default values.
func (o *PersonWithLock) Initialize(ctx context.Context) {
	o.personWithLockBase.Initialize()
	// Add your own initializations here
}

// String implements the Stringer interface and returns the default label for the object as it appears in html lists.
// Typically you would change this to whatever was pertinent to your application.
func (o *PersonWithLock) String() string {
	if o == nil {
		return "" // Possibly - Select One -?
	}
	return fmt.Sprintf("Object id %v", o.PrimaryKey())
}

// QueryPersonWithLocks returns a new builder that gives you general purpose access to the PersonWithLock records
// in the database. This is useful for quick queries of the database during development, but eventually you
// should remove this function and move those queries to more specific calls in this file.
func QueryPersonWithLocks() *PersonWithLocksBuilder {
	return queryPersonWithLocks()
}

// LoadPersonWithLock queries for a single PersonWithLock object by primary key.
// joinOrSelectNodes lets you provide nodes for joining to other tables or selecting specific fields. Table nodes will
// be considered Join nodes, and column nodes will be Select nodes. See Join() and Select() for more info.
// If you need a more elaborate query, use QueryPersonWithLocks() to start a query builder.
func LoadPersonWithLock(ctx context.Context, pk string, joinOrSelectNodes ...query.NodeI) *PersonWithLock {
	return loadPersonWithLock(ctx, pk, joinOrSelectNodes...)
}

// DeletePersonWithLock deletes the give record from the database. Note that you can also delete
// loaded PersonWithLock objects by calling Delete on them.
func DeletePersonWithLock(ctx context.Context, pk string) {
	deletePersonWithLock(ctx, pk)
}
