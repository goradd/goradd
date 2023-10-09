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

If you have not yet installed GoRADD, see the [Quickstart] document.

## Database Configuration

Start your development by modeling your data in a SQL database like MySQL or Postgres.
SQL databases are great for modeling relationships, they are easy to restructure,
and they have some built-in integrity checks. When your application grows to the
point of needing the kind of benefits that NoSQL databases provide, and your 
data model is more firmly established, you can transition to using a NoSQL
database.

It is best to install a SQL database system locally on your development computer. 
After you create the SQL database that will host your data, edit the 
*goradd-project/config/db.go*
file and enter the database credentials that will allow the application to
access the database. During development, be sure these credentials have access to 
your database, and access to the additional tables that describe the database. For
example, in MySQL, the credentials should be able to access the "mysql" table, and
in postrgres, the "information_schema" and the "pg_catalog" tables. You can use the
root user during development to accomplish this, and then for deployment, specify
a user that only has the minimal credentials to access the database.

See the article [Structuring the Database](database.md) for details, but the basic idea is 
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
`goradd-project/gen` directory,
4. Build documentation from the generated source and your own source, and place it in the 
`goradd-project/doc` directory.

Each time you change the structure of the database, you should run the code generator.


## The Gen Directory

The `gen` directory is where the code generator places the generated forms,
related objects, and ORM objects. The files are replaced each time
you run the code generator. Some files are meant to be copied to a new location
to use them, and some files are meant to be left in the gen directory and used
from there.

The gen directory is organized as follows:

* gen
  * (database name)
    * form
    * model
      * node
    * panelbase

<dl>
  <dt>model</dt>
  <dd>The model directory contains .go code to access the database using the ORM. It is meant
      to be used in place.</dd>
  <dt>form</dt>
  <dd>Forms represent the top level object in a page and the enclosure for the rest of the controls on the page.
      Forms that you wish to use in your application should be copied to the 
      goradd-project/web/form directory and modified there. This directory also includes 
      panels that help in the process of listing and editing objects (rows) in the database tables.
      To use them, copy them to your goradd-project/web/panels directory.</dd>
  <dd>panelbase contains code that the panel objects rely on. They should be used in place. They
      are regenerated every time the code generator is run.</dd>

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

If you change the .tpl.got template file, (perhaps
to add additional html), you will need to rebuild the template to reflect
the changes. To rebuild the template, run `go generate` on the build.go
file.

##  Customizing the Models

The model code is meant to be used in place.

To customize them, you have a few options:
1) Change the templates used to create them.
2) Some aspects of the model class can be change through comments in the database or data description file.

## Customizing the Look of the Panels
By default, a panel control will just print out its contained
child controls in the order they were added to the control. Sometimes
this works fine, but often you will want to add additional html to
the controls and their surroundings.

To change how the panel draws, do the following:
1) Copy the panel and its accompanying .tpl.go file from the gen/<db>/panel directory and put it in your goradd-project/web/panel directory. 
2) Edit the tpl.go file to control how the panel displays itself and the objects inside of it. 
3) Generate the build.go file to turn the template into go code.
4) Change the form file to import the panel from the new location.

To customize the output of a panel, you need to give it a template.
In the panel/inactive_templates directory you will find *got* templates
that contain the controls that will be printed by each panel. To activate
the template and prepare it for editing, move it from the inactive_templates directory
to the directory above it. After you edit it, you also need
to run `go generate build.go` on the build.go file to generate the template.

The templates that you move to the panel directory are not touched by
the code generator. This means that if you add a field to a database, and
then run the code generator, you will not see that field automatically
appear in the form in the browser. You will need to look at the changes
in the panel file under the gen directory and make those changes to your file.

## Changing the Generated Code

You can change the type of control that is associated with a particular type
of database field. See the comment in the goradd-project/codegen/cmd/codegen.go file
for an example of how to use bootstrap controls.

You can also change the template files themselves. The templates
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
