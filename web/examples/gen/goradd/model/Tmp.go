package model

// This is the implementation file for the Tmp ORM object.
// This is where you build the api to your data model for you web application and potentially mobile apps.
// Your edits to this file will be preserved.

import (
	"context"
	"fmt"
)

type Tmp struct {
	tmpBase
}

// Create a new Tmp object and initialize to default values.
func NewTmp() *Tmp {
	o := new(Tmp)
	o.Initialize()
	return o
}

// Initialize or re-initialize a Tmp database object to default values.
func (o *Tmp) Initialize() {
	o.tmpBase.Initialize()
	// Add your own initializations here
}

// String implements the Stringer interface and returns the default label for the object as it appears in html lists.
// Typically you would change this to whatever was pertinent to your application.
func (o *Tmp) String() string {
	if o == nil {
		return "" // Possibly - Select One -?
	}
	return fmt.Sprintf("Object id %v", o.PrimaryKey())
}

// QueryTmps returns a new builder that gives you general purpose access to the Tmp records
// in the database. Its here to give public access to the query builder, but you can remove it if you do not need it.
func QueryTmps(ctx context.Context) *TmpsBuilder {
	return queryTmps(ctx)
}

// queryTmps creates a new builder and is the central spot where all queries are directed.
// You can modify this function to enforce restrictions on queries, for example to make sure the user is authorized to
// access the data.
func queryTmps(ctx context.Context) *TmpsBuilder {
	return newTmpBuilder()
}

// DeleteTmp deletes the given record from the database. Note that you can also delete
// loaded Tmp objects by calling Delete on them.
func DeleteTmp(ctx context.Context, pk string) {
	deleteTmp(ctx, pk)
}
