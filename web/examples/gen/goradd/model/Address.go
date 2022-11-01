package model

// This is the implementation file for the Address ORM object.
// This is where you build the api to your data model for your web application and potentially mobile apps.
// Your edits to this file will be preserved.

import (
	"context"
	"encoding/gob"
	"fmt"
)

type Address struct {
	addressBase
}

// Create a new Address object and initialize to default values.
func NewAddress() *Address {
	o := new(Address)
	o.Initialize()
	return o
}

// Initialize or re-initialize a Address database object to default values.
func (o *Address) Initialize() {
	o.addressBase.Initialize()
	// Add your own initializations here
}

// String implements the Stringer interface and returns the default label for the object as it appears in html lists.
// Typically you would change this to whatever was pertinent to your application.
func (o *Address) String() string {
	if o == nil {
		return "" // Possibly - Select One -?
	}
	return fmt.Sprintf("Address %v", o.PrimaryKey())
}

// QueryAddresses returns a new builder that gives you general purpose access to the Address records
// in the database. Its here to give public access to the query builder, but you can remove it if you do not need it.
func QueryAddresses(ctx context.Context) *AddressesBuilder {
	return queryAddresses(ctx)
}

// queryAddresses creates a new builder and is the central spot where all queries are directed.
// You can modify this function to enforce restrictions on queries, for example to make sure the user is authorized to
// access the data.
func queryAddresses(ctx context.Context) *AddressesBuilder {
	return newAddressBuilder(ctx)
}

// DeleteAddress deletes the given record from the database. Note that you can also delete
// loaded Address objects by calling Delete on them.
func DeleteAddress(ctx context.Context, pk string) {
	deleteAddress(ctx, pk)
}

func init() {
	gob.RegisterName("ExampleAddress", new(Address))
}
