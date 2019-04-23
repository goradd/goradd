package model

// This is the implementation file for the Project ORM object.
// This is where you build the api to your data model for you web application and potentially mobile apps.
// Your edits to this file will be preserved.

import (
	"context"
	"fmt"

	"github.com/goradd/goradd/pkg/orm/query"
)

type Project struct {
	projectBase
}

// Create a new Project object and initialize to default values.
func NewProject(ctx context.Context) *Project {
	o := new(Project)
	o.Initialize(ctx)
	return o
}

// Initialize or re-initialize a Project database object to default values.
func (o *Project) Initialize(ctx context.Context) {
	o.projectBase.Initialize()
	// Add your own initializations here
}

// String implements the Stringer interface and returns the default label for the object as it appears in html lists.
// Typically you would change this to whatever was pertinent to your application.
func (o *Project) String() string {
	if o == nil {
		return "" // Possibly - Select One -?
	}
	return fmt.Sprintf("Object id %v", o.PrimaryKey())
}

// QueryProjects returns a new builder that gives you general purpose access to the Project records
// in the database. This is useful for quick queries of the database during development, but eventually you
// should remove this function and move those queries to more specific calls in this file.
func QueryProjects() *ProjectsBuilder {
	return queryProjects()
}

// LoadProject queries for a single Project object by primary key.
// joinOrSelectNodes lets you provide nodes for joining to other tables or selecting specific fields. Table nodes will
// be considered Join nodes, and column nodes will be Select nodes. See Join() and Select() for more info.
// If you need a more elaborate query, use QueryProjects() to start a query builder.
func LoadProject(ctx context.Context, pk string, joinOrSelectNodes ...query.NodeI) *Project {
	return loadProject(ctx, pk, joinOrSelectNodes...)
}

// LoadProjectByNum queries for a single Project object by the given unique index values.
// joinOrSelectNodes lets you provide nodes for joining to other tables or selecting specific fields. Table nodes will
// be considered Join nodes, and column nodes will be Select nodes. See Join() and Select() for more info.
// If you need a more elaborate query, use QueryProjects() to start a query builder.
func LoadProjectByNum(ctx context.Context, num int, joinOrSelectNodes ...query.NodeI) *Project {
	return loadProjectByNum(ctx, num, joinOrSelectNodes...)
}

// DeleteProject deletes the give record from the database. Note that you can also delete
// loaded Project objects by calling Delete on them.
func DeleteProject(ctx context.Context, pk string) {
	deleteProject(ctx, pk)
}
