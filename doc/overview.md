# How to GoRADD

## Introduction

GoRADD is a monolithic web development framework designed to quickly get you
from idea to working web application, and then to let you change the
application over time with as little rewrite as possible.

Instead of having to build separate server and client applications, and having 
to stitch together a bunch of different technologies just to get something 
running, GoRADD lets you start with your data model, gets you a working app
quickly, and then build from there. As you learn more about your app, your
audience, and your requirements, you can build and branch out to include them
as your needs grow. Your app grows in incremental steps, and GoRADD helps you on the way.

GoRADD is great for beginners and experienced developers alike.

This guide will walk you through the steps of building and maintaining your
application.

## Install

If you have not yet installed your application, see the [Quickstart] document.

## Database Configuration

Start your application by modeling your data in a SQL database like MySQL.
SQL databases are great for modeling relationships, they are easy to restructure,
and they have some built-in integrity checks. When your application grows to the
point of needing the kind of benefits that NoSQL databases provide, and your 
data model is more firmly established, you can transition to using a NoSQL
database.

To start, create a SQL database and then edit the goradd-project/config/db.go
file and enter the database credentials that will allow the application to
access the database.

See the article [Structuring the Database](#) for details, but the basic idea is 
to create tables with
foreign keys that link to other tables. These form relationships that will be
reflected by the code generated object-relational model (ORM) data access code, 
and the generated forms.
The ORM is go code that lets you query the database, update and delete records
in a way that is independent of the underlying database. If you are using a SQL
database, you likely will not need to write any SQL code to interact with the 
database (though you can if you need to).

Depending on how you structure these foreign keys, you can create one-to-one,
one-to-many, and many-to-many relationships.

You can use more than one database in your application, but GoRADD currently does
not support relationships between tables in different databases.

## Code Generation

After structuring your data, run the code generator by executing the following
from the goradd-project directory:

`go generate codegen/cmd/build.go`

This will do the following:
1. Delete previously generated templates,
2. Read the templates and generate corresponding go code, placing this code in the goradd-project/tmp directory,
3. Run the resulting go code to generate the forms and ORM and place them in the
`goradd-project/gen` directory

Each time you change something in the database, you should run the code generator.


## The Gen Directory

The `gen` directory is where the code generator places the generated forms,
related objects, and ORM objects. Some of the files are replaced each time
you run the code generator, and some are stub files that are generated once,
and then not touched. Some files are meant to be edited in-place, some not
changed by you at all (since they get replaced each time you codegen), 
and some are designed to be copied and moved to another location.

The gen directory is organized as follows:

* gen
  * (database name)
    * connector
    * form
    * model
      * node
    * panel
      * inactive_templates

<dl>
  <dt>model</dt>
  <dd>The model directory contains .go code to access the database using the ORM.</dd>
  <dt>form</dt>
  <dd>Forms represent the top level object in a page and the enclosure for the rest of the controls on the page.</dd>
  <dt>panel</dt>
  <dd>Panels are div objects that encapsulate most of the generated controls in a form.</dd>
  <dt>inactive_templates</dt>
  <dd>Contains default templates that you can activate by moving them to the directory one level above them.</dd>
</dl>

See [Anatomy of the Generated Framework](#) for more details on the roles of the generated
files to create the default application and how you can modify them.

To see the forms in action, start your application by navigating to the
goradd-project directory, and then from there execute:

`go run goradd-project/main`

After that, navigate to:

`http://localhost:8000/goradd/forms`

##  Setting up a Goradd Form

The first time you view your forms, they will not look all that impressive. However,
you will notice that you have the ability to list the records in the
database and perform Create, Update, and Delete operations on each record 
(also known a (CrUD)) using the generated forms.

The generated list forms list all the records 
TODO: Talk about list forms and edit forms

The form object represents the html `<form>` tag object in the page. It 
encloses your controls, and so is the place where you create and initialzie
the top-level controls on the page. 

The forms are located in the gen directory, and any changes you make will be
lost as they will be overwritten the next time you do code generation.
To make permanent changes to a form file, you should move the file to
a different directory, and then make sure that directory gets imported
into your project. A good location for the form file is the 
goradd-project/web/form directory, though any imported directory will work.
Be sure to also move the matching template file (ending in .got)
to the new location.

After moving the file, you should change the *Path* const at the top of
the file to the path you would like the user to use to get to the form.
At any point you can restart your application and test your changes to 
make sure they worked.

If you change the .got template file, (perhaps
to add additional html), you will need to rebuild the template to reflect
the changes. To rebuild the template, run `go generate` on the build.go
file.

You do not have to use the code generated forms. You can start with an
empty form, add controls and initialize them with data from the database
directly. The advantage of following the directions here is that as you
change the structure of your database over time, it will be easier to 
maintain the code that relies on this structure.

##  Customizing the Models

The `model` directory contains "base" files that are regenerated every
time you code-generate, and implementation files that are stub files 
that you can edit. The implementation file embeds the base file, so
all the base file functions are available through the implementation file,
and they are overridable by the implementation file.

You should not move these files to different directories,
they are designed to remain in the model directory.

While you should not change the "base" file, feel free to change the
implementation file. Good candidates for functions in the implementation file
would be functions that represent specialized queries, overrides of
base functions to do additional database validation before records are saved,
calculations based on database fields, etc. These functions would be the
part of your business logic that specifically relates to your data.
 
For example, to change how database objects are displayed in the list view, you
should edit the `String()` function located in each of the implementation
files. You should change it to whatever combination of data in a record
would correctly represent the record.

For example, if you had a "Name" field in the Person table, your String function
in the model/Person.go file might look like this:

```
func (o *Person) String() string {
    return o.Name()
}
```

Whenever the implementation file refers to database fields, it should 
always use the accessor functions, and not the local variables of the
base file.

The other common way to customize the models is to change the database
itself and then run the code generator. The changes will automatically get
picked up in the "base" file.

## Customizing the Look of the Panels
By default, a panel control will just print out its contained
child controls in the order they were added to the control. Sometimes
this works fine, but often you will want to add additional html to
the controls and their surroundings.

To customize the output of a panel, you need to give it a template.
In the panel/inactive_templates directory you will find *got* templates
that contain the controls that will be printed by each panel. To activate
the template and prepare it for editing, move it from the inactive_templates directory
to the directory above it. After you edit it, you also need
to run `go generate build.go` on the build.go file to generate the template.

The templates that you move to the panel directory are not touched by
the code generator. This means that if you add a field to a database, and
then run the code generator, you will not see that field automatically
appear in the form in the browser. You will need to add code to the
template file to draw the new field.

The panel directory contains "base" files and implementation files, like
the model directory. Feel free to edit the implementation file, but do
not make changes to the "base" file. The implementation file contains
a number of commented-out sample functions to give you guidance on
what kinds of changes you might make there. Some possible changes
can be made in the template file as well.

For example, if you had a field called "username" in the Person table, 
and you wanted to change the label on the field, 
as well as give some instructions
to the user, you could do that in the panel/PersonEditPanel.go file
with the following code:
```
func (p *PersonEditPanel) CreateControls() {
    p.PersonEditPanelBase.CreateControls()
    
    p.UsernameTextbox.SetLabel ("User Name")
    p.UsernameTextbox.SetInstructions ("Enter a user name that is at least 8 characters long.")
}

```

You could also do the same thing in the template file using the
following code:
```
{{draw p.UsernameTextbox.SetLabel("User Name").SetInstructions("Enter a user name that is at least 8 characters long.") }}

```

or

```
{{go
    p.UsernameTextbox.SetLabel ("User Name")
    p.UsernameTextbox.SetInstructions ("Enter a user name that is at least 8 characters long.")
}}
{{draw p.UsernameTextbox}}

```

## Custom Controls and Customizing How Data Is Saved

The connectors provide the link between the edit panels and the database.
They move data back and forth between the database and the controls in
the form. The edit panel calls the connector to create its controls,
and also passes off a variety of management responsibilities to the
conector.

If you want to implement a custom control to manage a particular field
or group of fields in the database, or if you want to change how a control
saves its data, you likely should do this in the connector file.

Like the other files in the framework, there is a "base" connector
which you should not edit, and an implementation file. Place your
changes in the implementation file.

One common example of when you might want to change the connector is
if you are implementing a password textbox. By default, the framework
will save whatever the user types into a control, but this would be
bad practice for passwords. Instead, you should save a hash of the 
password in the database so that there is no risk of the password
ever being seen. Here is an example of code in the connector that
will do this for you:

In your connector file:

```
package connector

import (
    "your/model"
	"github.com/goradd/goradd/pkg/page/control"
	"github.com/goradd/goradd/pkg/page"
	"context"
)

// This is the connector override. Feel free to edit.

const DefaultPassword = "******"

type Person struct {
	personBase
	PasswordTextbox *control.Textbox
}

func NewPersonConnector(parent page.ControlI) *Person {
	c := new(Person)
	c.ParentControl = parent
	return c
}


func (c *Person) NewPasswordTextbox(id string) *bootstrap.Textbox {
	var ctrl *control.Textbox
	ctrl = control.NewTextbox(c.ParentControl, id)
	ctrl.SetLabel("Password")
	ctrl.SetMaxLength(300)
	ctrl.SetIsRequired(true)
	ctrl.SetType(control.TextboxTypePassword)

	c.PasswordTextbox = ctrl
	return ctrl
}

// Load will associate the controls with data from the given model.Email object and load the controls with data.
// Generally call this after creating the controls. Otherwise, call Refresh if you Load before creating the controls.
// If you pass a new object, it will prepare the controls for creating a new record in the database.
func (c *Person) Load(ctx context.Context, modelObj *model.Person) {
	if modelObj == nil {
		modelObj = model.NewPerson(ctx)
	}
	c.Person = modelObj
	if modelObj.PrimaryKey() == "" {
		c.EditMode = false
	} else {
		c.EditMode = true
	}
	c.Refresh()
}

// Save takes the data from the controls and saves it in the database.
func (c *Person) Save(ctx context.Context) {
	c.Update()
	c.Person.Save(ctx)
}

func (c *Person) Refresh() {
	c.personBase.Refresh()

	if c.PasswordTextbox != nil {
		if c.Person.PwhashIsValid() && c.Person.Pwhash() != "" {
			c.PasswordTextbox.SetText(DefaultPassword)
		} else {
			c.PasswordTextbox.SetText("")
		}
	}
}

func (c *Person) Update() {
	c.personBase.Update()

	if c.PasswordTextbox != nil {
		text := c.PasswordTextbox.Text()
		if text != "" && text != DefaultPassword {
			c.Person.SetPassword(text)
		}
	}
}
```

And then in your model/Person.go file:

```
func (o *Person) SetPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	o.SetPwhash(string(bytes))
	return nil
}

func (o *Person) VerifyPassword(password string) bool {
	if !o.PwhashIsValid() { // we didn't query for the password hash in the last query
		return false
	}
	err := bcrypt.CompareHashAndPassword([]byte(o.Pwhash()), []byte(password))
	if err != nil {
		return false
	}
	return true
}

```

## Changing the Generated Code

You can change the type of control that is associated with a particular type
of database field. See the comment in the goradd-project/codegen/cmd/codegen.go file
for an example of how to use bootstrap controls.

You can also change the template files themselves. Like most of goradd, the templates
use a layered architecture so you can change specific files without having to
change everything, and so that your changes will not be overwritten when the
framework is updated. See the goradd-project/codegen/templates/readme.txt file
for more information.

## Styling and JavaScript

By default, GoRADD uses the goradd.css file to provide basic styling. To
provide additional styling, create your own .css files and put them in 
the goradd-project/web/assets/css directory.

To have the css file load for a particular form, you would create
an AddRelatedFiles function and add it to your form, like so:
```
func (f *MyForm) AddReleatedFiles() {
	f.FormBase.AddRelatedFiles()

	f.AddStyleSheetFile(path.Join(config.AssetPrefix, "project", css","styles.css"), nil)
}
``` 

Similarly, call the AddJavaScriptFile function in the AddRelatedFiles function
to add a javascript file to a form.

To add a css file or JavaScript file to all of your forms, add the corresponding call
to the AddRelatedFiles function in the goradd-project/control/form_base.go file.

## Wrap Up
For the most part, developing a GoRADD application is iterating all of the above.
As your requirements change, you change the database, generate the code,
then add your business logic and customizations. For basic business
websites, that will be enough. You can do lots more though, and the
following topics may help you explore more:

[Creating Custom Controls including Javascript Controls]
[Using CSS Frameworks including Bootstrap]
[Deployment]
