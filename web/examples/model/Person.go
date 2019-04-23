package model

// This is the implementation file for the Person ORM object.
// This is where you build the api to your data model for you web application and potentially mobile apps.
// Your edits to this file will be preserved.

import (
	"context"
	"fmt"

	"github.com/goradd/goradd/pkg/orm/query"
)

type Person struct {
	personBase
}

// Create a new Person object and initialize to default values.
func NewPerson(ctx context.Context) *Person {
	o := new(Person)
	o.Initialize(ctx)
	return o
}

// Initialize or re-initialize a Person database object to default values.
func (o *Person) Initialize(ctx context.Context) {
	o.personBase.Initialize()
	// Add your own initializations here
}

// String implements the Stringer interface and returns the default label for the object as it appears in html lists.
// Typically you would change this to whatever was pertinent to your application.
func (o *Person) String() string {
	if o == nil {
		return "" // Possibly - Select One -?
	}
	return fmt.Sprintf("Object id %v", o.PrimaryKey())
}

// QueryPeople returns a new builder that gives you general purpose access to the Person records
// in the database. This is useful for quick queries of the database during development, but eventually you
// should remove this function and move those queries to more specific calls in this file.
func QueryPeople() *PeopleBuilder {
	return queryPeople()
}

// LoadPerson queries for a single Person object by primary key.
// joinOrSelectNodes lets you provide nodes for joining to other tables or selecting specific fields. Table nodes will
// be considered Join nodes, and column nodes will be Select nodes. See Join() and Select() for more info.
// If you need a more elaborate query, use QueryPeople() to start a query builder.
func LoadPerson(ctx context.Context, pk string, joinOrSelectNodes ...query.NodeI) *Person {
	return loadPerson(ctx, pk, joinOrSelectNodes...)
}

// DeletePerson deletes the give record from the database. Note that you can also delete
// loaded Person objects by calling Delete on them.
func DeletePerson(ctx context.Context, pk string) {
	deletePerson(ctx, pk)
}
