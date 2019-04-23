package model

// This is the implementation file for the Tmp ORM object.
// This is where you build the api to your data model for you web application and potentially mobile apps.
// Your edits to this file will be preserved.

import (
	"context"
	"fmt"

	"github.com/goradd/goradd/pkg/orm/query"
)

type Tmp struct {
	tmpBase
}

// Create a new Tmp object and initialize to default values.
func NewTmp(ctx context.Context) *Tmp {
	o := new(Tmp)
	o.Initialize(ctx)
	return o
}

// Initialize or re-initialize a Tmp database object to default values.
func (o *Tmp) Initialize(ctx context.Context) {
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
// in the database. This is useful for quick queries of the database during development, but eventually you
// should remove this function and move those queries to more specific calls in this file.
func QueryTmps() *TmpsBuilder {
	return queryTmps()
}

// LoadTmp queries for a single Tmp object by primary key.
// joinOrSelectNodes lets you provide nodes for joining to other tables or selecting specific fields. Table nodes will
// be considered Join nodes, and column nodes will be Select nodes. See Join() and Select() for more info.
// If you need a more elaborate query, use QueryTmps() to start a query builder.
func LoadTmp(ctx context.Context, pk string, joinOrSelectNodes ...query.NodeI) *Tmp {
	return loadTmp(ctx, pk, joinOrSelectNodes...)
}

// LoadTmpByD queries for a single Tmp object by the given unique index values.
// joinOrSelectNodes lets you provide nodes for joining to other tables or selecting specific fields. Table nodes will
// be considered Join nodes, and column nodes will be Select nodes. See Join() and Select() for more info.
// If you need a more elaborate query, use QueryTmps() to start a query builder.
func LoadTmpByD(ctx context.Context, d string, joinOrSelectNodes ...query.NodeI) *Tmp {
	return loadTmpByD(ctx, d, joinOrSelectNodes...)
}

// DeleteTmp deletes the give record from the database. Note that you can also delete
// loaded Tmp objects by calling Delete on them.
func DeleteTmp(ctx context.Context, pk string) {
	deleteTmp(ctx, pk)
}
