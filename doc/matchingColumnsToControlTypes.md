# Matching Columns to Control Types

The code generator will generate an edit panel to provide a starting place for developing an edit screen.
It will match the type of data in each database column with a control type that can be used to edit that
particular type of data, and populate the edit panel with those controls. 
There are a few ways to customize how this is done.

## The DefaultControlTypeFunc function
This function returns the module path to a control that will be used to edit a particular type of
column. It is a global variable, so you can change this function to suit your needs. See the
BootstrapCodegenSetup function for an example of how to change it.

The returned control must be capable of editing the type of data indicated in the column. 
If you are creating your own custom control that you want to return in this function,
you will likely need to create a generator for the control too. More on generators below.

## Specifying the Control Type in the Database Description
If you are using a SQL database, you can specify a specific control type for a column by putting
a JSON variable in the comment attached to a column. The variable name is "controlPath", and the
value should be the module path to the control. For example:
```json
{"controlPath":"github.com/goradd/goradd/pkg/bootstrap/control/EmailTextbox"}
```
in the comment of a database column will force that column to generate a bootstrap EmailTextbox control
to edit that column.

In a NoSQL database, you would edit the database description file to specify this in the Options area
of the column.

The advantage of specifying the control type in the database is that as the database changes, the
control type will stick to the column. If you delete the column, the control type gets deleted too, 
and the control will automatically not be created in the future.

## Specifying the Control Type in the generated Panel
In your goradd-project/gen/*/panel directory, you will see a variety of panels generated. For each
database table, there
are two panels for editing objects in the table, and *EditPanel and a *EditPanelBase,
and two for listing them.

The CreateControls function in the panel is responsible for creating the controls for each column
in the table, and it will call individual control creation functions. You can override any of these
functions to change how individual controls are created. To do so, find the control creation function
in the .base.go file, copy that function to the corresponding non-base file, and then modify it
to suit your needs.

The .base.go file gets recreated every time you perform a code generation. The corresponding non-base
file will not be changed when you code-generate.

## Generators
A generator is responsible for generating code for a particular type of control. It will generate
the control creator from information contained in the database column, and also generate the connector
which transfers information between the control and the database at set times. 

If you are creating your own custom control that you want to use in code generation, you will need
to create a generator and register it with the code generator. Examples can be found
in the goradd.com/goradd/goradd/pkg/page/control/generator package.