package model

// This is the implementation file for the Login ORM object.
// This is where you build the api to your data model for you web application and potentially mobile apps.
// Your edits to this file will be preserved.

import (
	"context"
	"fmt"
)

type Login struct {
	loginBase
}

// Create a new Login object and initialize to default values.
func NewLogin(ctx context.Context) *Login {
	o := new(Login)
	o.Initialize(ctx)
	return o
}

// Initialize or re-initialize a Login database object to default values.
func (o *Login) Initialize(ctx context.Context) {
	o.loginBase.Initialize()
	// Add your own initializations here
}

// String implements the Stringer interface and returns the default label for the object as it appears in html lists.
// Typically you would change this to whatever was pertinent to your application.
func (o *Login) String() string {
	if o == nil {
		return "" // Possibly - Select One -?
	}
	return fmt.Sprintf("Object id %v", o.PrimaryKey())
}

// QueryLogins returns a new builder that gives you general purpose access to the Login records
// in the database. Its here to give public access to the query builder, but you can remove it if you do not need it.
func QueryLogins(ctx context.Context) *LoginsBuilder {
	return queryLogins(ctx)
}

// queryLogins creates a new builder and is the central spot where all queries are directed.
// You can modify this function to enforce restrictions on queries, for example to make sure the user is authorized to
// access the data.
func queryLogins(ctx context.Context) *LoginsBuilder {
	return newLoginBuilder()
}

// DeleteLogin deletes the given record from the database. Note that you can also delete
// loaded Login objects by calling Delete on them.
func DeleteLogin(ctx context.Context, pk string) {
	deleteLogin(ctx, pk)
}
