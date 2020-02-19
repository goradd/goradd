package model

// This is the implementation file for the EmployeeInfo ORM object.
// This is where you build the api to your data model for you web application and potentially mobile apps.
// Your edits to this file will be preserved.

import (
	"context"
	"encoding/gob"
	"fmt"
)

type EmployeeInfo struct {
	employeeInfoBase
}

// Create a new EmployeeInfo object and initialize to default values.
func NewEmployeeInfo() *EmployeeInfo {
	o := new(EmployeeInfo)
	o.Initialize()
	return o
}

// Initialize or re-initialize a EmployeeInfo database object to default values.
func (o *EmployeeInfo) Initialize() {
	o.employeeInfoBase.Initialize()
	// Add your own initializations here
}

// String implements the Stringer interface and returns the default label for the object as it appears in html lists.
// Typically you would change this to whatever was pertinent to your application.
func (o *EmployeeInfo) String() string {
	if o == nil {
		return "" // Possibly - Select One -?
	}
	return fmt.Sprintf("Object id %v", o.PrimaryKey())
}

// QueryEmployeeInfos returns a new builder that gives you general purpose access to the EmployeeInfo records
// in the database. Its here to give public access to the query builder, but you can remove it if you do not need it.
func QueryEmployeeInfos(ctx context.Context) *EmployeeInfosBuilder {
	return queryEmployeeInfos(ctx)
}

// queryEmployeeInfos creates a new builder and is the central spot where all queries are directed.
// You can modify this function to enforce restrictions on queries, for example to make sure the user is authorized to
// access the data.
func queryEmployeeInfos(ctx context.Context) *EmployeeInfosBuilder {
	return newEmployeeInfoBuilder(ctx)
}

// DeleteEmployeeInfo deletes the given record from the database. Note that you can also delete
// loaded EmployeeInfo objects by calling Delete on them.
func DeleteEmployeeInfo(ctx context.Context, pk string) {
	deleteEmployeeInfo(ctx, pk)
}

func init() {
	gob.RegisterName("ExampleEmployeeInfo", new(EmployeeInfo))
}
