//** This file was code generated by GoT. DO NOT EDIT. ***

package orm

import (
	"context"
	"io"

	"github.com/goradd/goradd/web/examples/gen/goradd/model"
	"github.com/goradd/goradd/web/examples/gen/goradd/model/node"
)

func (ctrl *TypesPanel) DrawTemplate(ctx context.Context, _w io.Writer) (err error) {

	if _, err = io.WriteString(_w, `<h1>Type Tables</h1>
<p>
Goradd uses special tables to model enumerated types in the database called Type Tables. A type table is simply
a table that has a minimum of a primary key and a name field. The name-value pairs will be converted to constants in Go,
with the name field becoming the name of the constant in Go,
and the primary key becoming the value.
</p>
<p>
You can refer to these constant values from fields in your database, and that reference will be converted to
the corresponding type in Go. You can also use association tables to create references to a slice of these constant values.
</p>
<p>
Type tables can also have additional fields that you can use as a lookup table in Go. These additional values will become
constant values in Go that you can lookup using the corresponding type value.
</p>
<p>
Type tables give you the ability to not only easily create enumerated types in Go, but also use these values in queries
that you make to the database.
</p>
<p>
Since type tables become constants at compile time, you cannot change the values during program execution. So only use
type tables for values that you know will not change.
</p>
`); err != nil {
		return
	}

	project := model.LoadProject(ctx, "1")

	if _, err = io.WriteString(_w, `<p>
The status of the `); err != nil {
		return
	}

	if _, err = io.WriteString(_w, project.Name()); err != nil {
		return
	}

	if _, err = io.WriteString(_w, ` project is `); err != nil {
		return
	}

	if _, err = io.WriteString(_w, project.StatusType().String()); err != nil {
		return
	}

	if _, err = io.WriteString(_w, `.
</p>
`); err != nil {
		return
	}

	person := model.LoadPerson(ctx, "3", node.Person().PersonTypes()) // Note that we are joining the type table because it has a many-many association

	if _, err = io.WriteString(_w, `<p>
The types of `); err != nil {
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

	if _, err = io.WriteString(_w, ` are
`); err != nil {
		return
	}

	for _i, _j := range person.PersonTypes() {
		_ = _j

		if _, err = io.WriteString(_w, _j.String()); err != nil {
			return
		}

		if _i < len(person.PersonTypes())-1 {
			if _, err = io.WriteString(_w, ", "); err != nil {
				return
			}
		}
	}
	if _, err = io.WriteString(_w, `
</p>

`); err != nil {
		return
	}

	return
}
