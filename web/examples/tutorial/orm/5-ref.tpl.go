//** This file was code generated by GoT. DO NOT EDIT. ***

package orm

import (
	"context"
	"io"

	"github.com/goradd/goradd/web/examples/gen/goradd/model"
	"github.com/goradd/goradd/web/examples/gen/goradd/model/node"
)

func (ctrl *RefPanel) DrawTemplate(ctx context.Context, _w io.Writer) (err error) {

	if _, err = io.WriteString(_w, `<h1>References</h1>
<h2>Foreign Keys</h2>
<p>
    Relational databases let you link records together using record ids, called foreign keys. At its basic level, a foreign key is
    just a field that
    contains a key that identifies a record in another table (or even in the same table). Many databases have a mechanism that lets you further describe
    a foreign key and how it behaves. For example, MySQL calls these CONSTRAINTs. These descriptions help to maintain the
    integrity of the database while modifying inter-related records.
</p>
<p>
    Goradd will detect these relationships in your database and create links to these related objects so that you can get to them easily.
    If you are not using a SQL database, or you are using a SQL database that does not have a CONSTRAINT mechanism,
    you can still get the same behavior by creating a data description file to tell Goradd about these relationships, and Goradd will
    then manage these links.
</p>
<p>
    One important thing to do is decide what should happen if the referenced record
    is deleted. Usually, you will want one of two behaviors:
    <ol>
    <li>Set the reference to NULL, or</li>
    <li>Delete this record</li>
    </ol>
    Goradd will look at what direction you have given in the constraint for the foreign key to determine what to do.
    If the constraint is specified to Set Null on Delete, then it will set the foreign key to NULL when the record
    on the other side of the relationship is deleted.
    If it is directed to Cascade on Delete, it will delete any records pointing to it with a foreign key.
    You can override this behavior, but that is what happens by default.
</p>

<h2>Loading Referenced Records</h2>
<p>
    In the example below, we get the first address record, and then we follow the link to the person that has that address
    by using the LoadPerson function from that Address. That will query the database again for the related address.
</p>
`); err != nil {
		return
	}

	address := model.LoadAddress(ctx, "1")
	person := address.LoadPerson(ctx)

	if _, err = io.WriteString(_w, `<p>
    Address: `); err != nil {
		return
	}

	if _, err = io.WriteString(_w, address.Street()); err != nil {
		return
	}

	if _, err = io.WriteString(_w, `, `); err != nil {
		return
	}

	if _, err = io.WriteString(_w, address.City()); err != nil {
		return
	}

	if _, err = io.WriteString(_w, `<br>
    Person: `); err != nil {
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

	if _, err = io.WriteString(_w, `
</p>
<h2>Pre-loading Referenced Records</h2>
<p>
    In the example above, we made two queries to the database. All SQL databases, and some NoSQL databases, have the ability
    to combine queries like this into one query. In SQL, you use a JOIN statement, and Goradd adopts this terminology
    to indicate that you want to use a foreign key to pre-load related records.
</p>
<p>
    To preload a connection using a Load* function, simply pass in nodes for the tables that you want to preload as an extra
    parameter to the Load* function.
</p>
`); err != nil {
		return
	}

	address = model.LoadAddress(ctx, "2", node.Address().Person())

	if _, err = io.WriteString(_w, `<p>
    Address: `); err != nil {
		return
	}

	if _, err = io.WriteString(_w, address.Street()); err != nil {
		return
	}

	if _, err = io.WriteString(_w, `, `); err != nil {
		return
	}

	if _, err = io.WriteString(_w, address.City()); err != nil {
		return
	}

	if _, err = io.WriteString(_w, `<br>
    Person: `); err != nil {
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

	if _, err = io.WriteString(_w, `
</p>
<p>
    You can pre-load slice queries too.
</p>
<p>
    `); err != nil {
		return
	}

	for _, project := range model.QueryProjects(ctx).
		Join(node.Project().Manager()).
		Load() {

		if _, err = io.WriteString(_w, `
            <div>Project: `); err != nil {
			return
		}

		if _, err = io.WriteString(_w, project.Name()); err != nil {
			return
		}

		if _, err = io.WriteString(_w, `, Manager: `); err != nil {
			return
		}

		if _, err = io.WriteString(_w, project.Manager().FirstName()); err != nil {
			return
		}

		if _, err = io.WriteString(_w, ` `); err != nil {
			return
		}

		if _, err = io.WriteString(_w, project.Manager().LastName()); err != nil {
			return
		}

		if _, err = io.WriteString(_w, `</div>
    `); err != nil {
			return
		}

	}

	if _, err = io.WriteString(_w, `
</p>

`); err != nil {
		return
	}

	return
}
