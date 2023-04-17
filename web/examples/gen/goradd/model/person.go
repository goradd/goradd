package model

// This is the implementation file for the Person ORM object.
// This is where you build the api to your data model for your web application and potentially mobile apps.
// Your edits to this file will be preserved.

import (
	"context"
	"encoding/gob"
	"fmt"
)

type Person struct {
	personBase
}

// NewPerson creates a new Person object and initializes it to default values.
func NewPerson() *Person {
	o := new(Person)
	o.Initialize()
	return o
}

// Initialize will initialize or re-initialize a Person database object to default values.
func (o *Person) Initialize() {
	o.personBase.Initialize()
	// Add your own initializations here
}

// String implements the Stringer interface and returns the default label for the object as it appears in html lists.
// Typically you would change this to whatever was pertinent to your application.
func (o *Person) String() string {
	if o == nil {
		return "" // Possibly - Select One -?
	}
	return fmt.Sprintf("Person %v", o.PrimaryKey())
}

// QueryPeople returns a new builder that gives you general purpose access to the Person records
// in the database. Its here to give public access to the query builder, but you can remove it if you do not need it.
func QueryPeople(ctx context.Context) *PeopleBuilder {
	return queryPeople(ctx)
}

// queryPeople creates a new builder and is the central spot where all queries are directed.
// You can modify this function to enforce restrictions on queries, for example to make sure the user is authorized to
// access the data.
func queryPeople(ctx context.Context) *PeopleBuilder {
	return newPersonBuilder(ctx)
}

// DeletePerson deletes the given record from the database. Note that you can also delete
// loaded Person objects by calling Delete on them.
func DeletePerson(ctx context.Context, pk string) {
	deletePerson(ctx, pk)
}

func init() {
	gob.RegisterName("ExamplePerson", new(Person))
}
