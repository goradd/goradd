package model

// This is the implementation file for the Login ORM object.
// This is where you build the api to your data model for you web application and potentially mobile apps.
// Your edits to this file will be preserved.

import (
	"context"
	"fmt"

	"github.com/goradd/goradd/pkg/orm/query"
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
// in the database. This is useful for quick queries of the database during development, but eventually you
// should remove this function and move those queries to more specific calls in this file.
func QueryLogins() *LoginsBuilder {
	return queryLogins()
}

// LoadLogin queries for a single Login object by primary key.
// joinOrSelectNodes lets you provide nodes for joining to other tables or selecting specific fields. Table nodes will
// be considered Join nodes, and column nodes will be Select nodes. See Join() and Select() for more info.
// If you need a more elaborate query, use QueryLogins() to start a query builder.
func LoadLogin(ctx context.Context, pk string, joinOrSelectNodes ...query.NodeI) *Login {
	return loadLogin(ctx, pk, joinOrSelectNodes...)
}

// LoadLoginByUsername queries for a single Login object by the given unique index values.
// joinOrSelectNodes lets you provide nodes for joining to other tables or selecting specific fields. Table nodes will
// be considered Join nodes, and column nodes will be Select nodes. See Join() and Select() for more info.
// If you need a more elaborate query, use QueryLogins() to start a query builder.
func LoadLoginByUsername(ctx context.Context, username string, joinOrSelectNodes ...query.NodeI) *Login {
	return loadLoginByUsername(ctx, username, joinOrSelectNodes...)
}

// LoadLoginByPersonID queries for a single Login object by the given unique index values.
// joinOrSelectNodes lets you provide nodes for joining to other tables or selecting specific fields. Table nodes will
// be considered Join nodes, and column nodes will be Select nodes. See Join() and Select() for more info.
// If you need a more elaborate query, use QueryLogins() to start a query builder.
func LoadLoginByPersonID(ctx context.Context, person_id string, joinOrSelectNodes ...query.NodeI) *Login {
	return loadLoginByPersonID(ctx, person_id, joinOrSelectNodes...)
}

// DeleteLogin deletes the give record from the database. Note that you can also delete
// loaded Login objects by calling Delete on them.
func DeleteLogin(ctx context.Context, pk string) {
	deleteLogin(ctx, pk)
}
