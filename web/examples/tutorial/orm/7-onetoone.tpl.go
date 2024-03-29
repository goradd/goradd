//** This file was code generated by GoT. DO NOT EDIT. ***

package orm

import (
	"context"
	"io"

	"github.com/goradd/goradd/web/examples/gen/goradd/model"
	"github.com/goradd/goradd/web/examples/gen/goradd/model/node"
)

func (ctrl *OneOnePanel) DrawTemplate(ctx context.Context, _w io.Writer) (err error) {

	if _, err = io.WriteString(_w, `<h1>One-to-One Relationships</h1>
<p>
By creating a unique index on a foreign key, you will create a one-to-one relationship.
One-to-one relationships link two records with each record being an extension of the other.
You can use this to model subclasses, extend a record with optional values, or to improve search results, allowing the database
to retrieve only the main record on an initial search, and then the extended record if more detail is desired of
a particular record.
</p>
<p>
In the example below, we load a login, and then examine which person the login belongs to.
`); err != nil {
		return
	}

	login := model.LoadLogin(ctx, "2", node.Login().Person())

	if _, err = io.WriteString(_w, `</p>
<p>
    Person for Login `); err != nil {
		return
	}

	if _, err = io.WriteString(_w, login.Username()); err != nil {
		return
	}

	if _, err = io.WriteString(_w, `: `); err != nil {
		return
	}

	if _, err = io.WriteString(_w, login.Person().FirstName()); err != nil {
		return
	}

	if _, err = io.WriteString(_w, ` `); err != nil {
		return
	}

	if _, err = io.WriteString(_w, login.Person().LastName()); err != nil {
		return
	}

	if _, err = io.WriteString(_w, ` <br>
</p>
<p>
Here, we traverse the relationship in the other direction, loading the person first, and then getting the login.
`); err != nil {
		return
	}

	person := model.LoadPerson(ctx, "3", node.Person().Login())

	if _, err = io.WriteString(_w, `</p>
<p>
    Login for `); err != nil {
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

	if _, err = io.WriteString(_w, `: `); err != nil {
		return
	}

	if _, err = io.WriteString(_w, person.Login().Username()); err != nil {
		return
	}

	if _, err = io.WriteString(_w, `  <br>
</p>

<h2>Creating One-to-One Linked Records</h2>
<p>
In a similar fashion to how one-to-many relationships work, you can create a link between two records by saving one,
getting its id, and then setting the foreign key in the other record to that id. However, its easier to use the Set*
functions for the objects themselves and call Save on the parent object.
</p>
`); err != nil {
		return
	}

	newPerson := model.NewPerson()
	newPerson.SetFirstName("Hu")
	newPerson.SetLastName("Man")

	newLogin := model.NewLogin()
	newLogin.SetUsername("human")

	newPerson.SetLogin(newLogin)
	newPerson.Save(ctx)

	if _, err = io.WriteString(_w, `<p>
    New person `); err != nil {
		return
	}

	if _, err = io.WriteString(_w, newPerson.FirstName()); err != nil {
		return
	}

	if _, err = io.WriteString(_w, ` `); err != nil {
		return
	}

	if _, err = io.WriteString(_w, newPerson.LastName()); err != nil {
		return
	}

	if _, err = io.WriteString(_w, ` has been given a Login ID of `); err != nil {
		return
	}

	if _, err = io.WriteString(_w, newPerson.Login().ID()); err != nil {
		return
	}

	if _, err = io.WriteString(_w, `
</p>
`); err != nil {
		return
	}

	// Delete records created above
	newPerson.Delete(ctx) // newLogin will automatically get deleted because its foreign key constraint is set to CASCADE on Delete

	return
}
