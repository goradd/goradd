//** This file was code generated by GoT. DO NOT EDIT. ***

package orm

import (
	"context"
	"io"

	"github.com/goradd/goradd/web/examples/gen/goradd/model"
	"github.com/goradd/goradd/web/examples/gen/goradd/model/node"
)

func (ctrl *ManyManyPanel) DrawTemplate(ctx context.Context, _w io.Writer) (err error) {

	if _, err = io.WriteString(_w, `<h1>Many-to-Many Relationships</h1>
<p>
Many-to-many relationships link records where each side of the relationship sees multiple records on the other side.
In the Golang ORM, each side of the relationship will see a slice of records on the other side.
</p>
<p>
Many-to-many relationships are modeled in SQL databases with an intermediate table, called an Association table.
The association table is a table with just two fields, each field being a foreign key pointing to one side of the
relationship. To identify a table as being an association table in SQL, append "_assn" to the end of the name of the table.
</p>
<p>
NoSQL databases store Many-to-many relationships in special fields on each side that is an array of record ids that
point to the records on the other side.
</p>
<p>
In either case, the ORM abstracts out the means of creating the relationship so that you do not have to worry about what
is happening in the database. Simply treat each side as a slice
of objects pointing to the other table, and the Goradd ORM will take care of the rest.
</p>
<p>
In the example below, we are using the team member - project association. Any person can be a team member of many projects,
and any project can have multiple team members.
`); err != nil {
		return
	}

	project := model.LoadProject(ctx, "1", node.Project().TeamMembers())
	person := model.LoadPerson(ctx, "1", node.Person().ProjectsAsTeamMember())

	if _, err = io.WriteString(_w, `</p>
<p>
    Project `); err != nil {
		return
	}

	if _, err = io.WriteString(_w, project.Name()); err != nil {
		return
	}

	if _, err = io.WriteString(_w, ` has team members:
    `); err != nil {
		return
	}

	for _i, _j := range project.TeamMembers() {
		_ = _j

		if _, err = io.WriteString(_w, _j.FirstName()); err != nil {
			return
		}

		if _, err = io.WriteString(_w, ` `); err != nil {
			return
		}

		if _, err = io.WriteString(_w, _j.LastName()); err != nil {
			return
		}

		if _i < len(project.TeamMembers())-1 {
			if _, err = io.WriteString(_w, ", "); err != nil {
				return
			}
		}
	}
	if _, err = io.WriteString(_w, `
    `); err != nil {
		return
	}

	if _, err = io.WriteString(_w, `
</p>
<p>
    Person `); err != nil {
		return
	}

	if _, err = io.WriteString(_w, person.FirstName()); err != nil {
		return
	}

	if _, err = io.WriteString(_w, ` `); err != nil {
		return
	}

	if _, err = io.WriteString(_w, person.LastName()); err != nil {
		return
	}

	if _, err = io.WriteString(_w, ` is a member of these projects:
    `); err != nil {
		return
	}

	for _i, _j := range person.ProjectsAsTeamMember() {
		_ = _j

		if _, err = io.WriteString(_w, _j.Name()); err != nil {
			return
		}

		if _i < len(person.ProjectsAsTeamMember())-1 {
			if _, err = io.WriteString(_w, ", "); err != nil {
				return
			}
		}
	}
	if _, err = io.WriteString(_w, `
</p>
<h2>Creating Many-Many Linked Records</h2>
<p>
Creating many-many linked records is similar to creating linked records in a one-to-many situation. You
simply call the appropriate Set* function to set the slice of items, and then call Save.
</p>
`); err != nil {
		return
	}

	project2 := model.NewProject()
	project2.SetName("NewProject")
	project2.SetNum(100)
	project2.SetStatusType(model.ProjectStatusTypeOpen)

	p1 := model.NewPerson()
	p1.SetFirstName("Me")
	p1.SetLastName("You")
	p2 := model.NewPerson()
	p2.SetFirstName("Him")
	p2.SetLastName("Her")

	project2.SetTeamMembers([]*model.Person{p1, p2})
	project2.Save(ctx)

	project3 := model.LoadProject(ctx, project2.ID(), node.Project().TeamMembers())

	if _, err = io.WriteString(_w, `
<p>
    Project `); err != nil {
		return
	}

	if _, err = io.WriteString(_w, project3.Name()); err != nil {
		return
	}

	if _, err = io.WriteString(_w, ` has team members:
    `); err != nil {
		return
	}

	for _i, _j := range project3.TeamMembers() {
		_ = _j

		if _, err = io.WriteString(_w, _j.FirstName()); err != nil {
			return
		}

		if _, err = io.WriteString(_w, ` `); err != nil {
			return
		}

		if _, err = io.WriteString(_w, _j.LastName()); err != nil {
			return
		}

		if _i < len(project3.TeamMembers())-1 {
			if _, err = io.WriteString(_w, ", "); err != nil {
				return
			}
		}
	}
	if _, err = io.WriteString(_w, `
</p>
<h2>Deleting Many-Many Linked Records</h2>
<p>
Deleting a record will also delete the link between two many-many linked records. However, it will not delete
the record on the other side of the link.
</p>

`); err != nil {
		return
	}

	// Delete the records we created above.
	project3.Delete(ctx)
	p1.Delete(ctx)
	p2.Delete(ctx)

	if _, err = io.WriteString(_w, `
`); err != nil {
		return
	}

	return
}
