package model

// This is the implementation file for the PersonWithLock ORM object.
// This is where you build the api to your data model for your web application and potentially mobile apps.
// Your edits to this file will be preserved.

import (
	"context"
	"encoding/gob"
	"fmt"
)

type PersonWithLock struct {
	personWithLockBase
}

// Create a new PersonWithLock object and initialize to default values.
func NewPersonWithLock() *PersonWithLock {
	o := new(PersonWithLock)
	o.Initialize()
	return o
}

// Initialize or re-initialize a PersonWithLock database object to default values.
func (o *PersonWithLock) Initialize() {
	o.personWithLockBase.Initialize()
	// Add your own initializations here
}

// String implements the Stringer interface and returns the default label for the object as it appears in html lists.
// Typically you would change this to whatever was pertinent to your application.
func (o *PersonWithLock) String() string {
	if o == nil {
		return "" // Possibly - Select One -?
	}
	return fmt.Sprintf("PersonWithLock %v", o.PrimaryKey())
}

// QueryPersonWithLocks returns a new builder that gives you general purpose access to the PersonWithLock records
// in the database. Its here to give public access to the query builder, but you can remove it if you do not need it.
func QueryPersonWithLocks(ctx context.Context) *PersonWithLocksBuilder {
	return queryPersonWithLocks(ctx)
}

// queryPersonWithLocks creates a new builder and is the central spot where all queries are directed.
// You can modify this function to enforce restrictions on queries, for example to make sure the user is authorized to
// access the data.
func queryPersonWithLocks(ctx context.Context) *PersonWithLocksBuilder {
	return newPersonWithLockBuilder(ctx)
}

// DeletePersonWithLock deletes the given record from the database. Note that you can also delete
// loaded PersonWithLock objects by calling Delete on them.
func DeletePersonWithLock(ctx context.Context, pk string) {
	deletePersonWithLock(ctx, pk)
}

func init() {
	gob.RegisterName("ExamplePersonWithLock", new(PersonWithLock))
}
