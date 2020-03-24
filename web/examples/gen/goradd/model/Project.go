package model

// This is the implementation file for the Project ORM object.
// This is where you build the api to your data model for you web application and potentially mobile apps.
// Your edits to this file will be preserved.

import (
	"context"
	"encoding/gob"
	"fmt"
)

type Project struct {
	projectBase
}

// Create a new Project object and initialize to default values.
func NewProject() *Project {
	o := new(Project)
	o.Initialize()
	return o
}

// Initialize or re-initialize a Project database object to default values.
func (o *Project) Initialize() {
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
// in the database. Its here to give public access to the query builder, but you can remove it if you do not need it.
func QueryProjects(ctx context.Context) *ProjectsBuilder {
	return queryProjects(ctx)
}

// queryProjects creates a new builder and is the central spot where all queries are directed.
// You can modify this function to enforce restrictions on queries, for example to make sure the user is authorized to
// access the data.
func queryProjects(ctx context.Context) *ProjectsBuilder {
	return newProjectBuilder(ctx)
}

// DeleteProject deletes the given record from the database. Note that you can also delete
// loaded Project objects by calling Delete on them.
func DeleteProject(ctx context.Context, pk string) {
	deleteProject(ctx, pk)
}

func init() {
	gob.RegisterName("ExampleProject", new(Project))
}
