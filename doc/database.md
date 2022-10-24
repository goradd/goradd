# Databases
The database structure is central to Goradd. The GoRADD code generator generates Go objects based
on the structure of SQL databaase code (called the ORM). Whenever you change the structure of the database, 
you should run the code generator to make sure the Go code reflects the structure of the database.

You will spend much of your initial design process in working out the details of your database structure.

Goradd borrows from SQL terminology to think of databases in terms of tables and columns.
Each row in the spreadsheet is a record, and each column in that row is a field 
of data.

## Tables
Tables (with a few exceptions mentioned below) become Go struct objects in GoRADD. Tables contain columns,
and these columns become the member variables in the Go struct.

The name of the table in the database is used by default to name the object type in Go code. You can 
change this value by setting the goName option in the database. See [Options](#options) below.

Tables can contain indexes to columns. All unique indexes become LoadBy* functions
in Go. For example, a unique index on the "id" field will create a LoadById() function
in the Go object.

## Columns
Columns define the fields of data that make up a row. You will want to define the following in a column.

### Name
The name of the column. Names should be lower_snake_case, as in words should be separated by underscores.
The primary reason for this is for cross-database compatibility, as some databases do not distinguish
between upper and lower case column names.

Column names will be converted to CamelCase inside of the Goradd ORM.

For example, a column name of *first_name* in the the *person* table would appear as Person().FirstName()
in Goradd.

### Type
A column in a database will generally have a certain type of data. Integers, strings, dates and times
are some examples. Go also has its own types of data, like strings, float32, float64, time.Time, etc.
Each database will have some mechanism that will map a database's column to a Goradd data type.
In some cases, the mapping is obvious, like a Varchar or Char type in a SQL database mapping to a 
string type in Go, and the mapping is automatically handled by Goradd. 
However, in other cases the mapping will need to be done in some kind of explicit manner. 

### Nullable
A column that is nullable will be initialized to a nil value by default in Goradd, and a NULL value
in the database. A NULL can represent whatever you want, but generally it means *no value*.

If a column is not Nullable, it also indicates to the 
UI generator that the data for that column is
required. By default, the UI will generate a control that will check that the user has
entered a value during validation of the control's form. You can override this behavior if needed
by calling SetIsRequired(false) on the control in Goradd.

### Default Value
A column's default value is what it will be set to after you create a new record object in Goradd. You
do not need to explicitly set such a column to use it, since it will already have the default value.

If you do not give a column a default value, and a column is Nullable, its default value will be nil.

However, if you do not give a column a default value, and the column is not Nullable, the ORM will
require that you explicity set a value for that column before you can insert the record into the database,
and will panic if you try to insert a record whose required columns have not been set.

## Options
Tables and Columns have a few configurable options to customize code generation. In SQL code, you place
these options in the Comment field of the table or Column as JSON. For example, putting the following in
the Comment field of a table will change the name that Go uses for the table to MyObject:
```json
{"goName":"MyObject"}
```

### Table Options
<dl>
  <dt><strong>literalName</strong></dt>
  <dd>The public word used in panels and forms when referring to this type of object. For example: person, employee, user.
      This should be lowercase. If you specify this, you should also specify a literalPlural</dd>
  <dt><strong>literalPlural</strong></dt>
  <dd>The plural word used to publicly refer to this type of object. For example: people, employees, users.
      This should be lowercase.</dd>
  <dt><strong>goName</strong></dt>
  <dd>The internal name used when referring to the object in Go code.</dd>
  <dt><strong>goPlural</strong></dt>
  <dd>The internal plural name used when referring to the object in Go code.</dd>
  <dt><strong>stringerColumn</strong></dt>
  <dd>The name of the field the object should use in its String() function to name the object.
      If there is a "name" field, it will automatically be used. You can specify a different database
      field here. For example: {"stringer":"title"} will look for a title column in the database and use
      its value when calling the String() function on the object.</dd>

</dl>

### Column Options
<dl>
  <dt><strong>goName</strong></dt>
  <dd>The internal name used when referring to the field in Go code.</dd>
  <dt><strong>goPlural</strong></dt>
  <dd>The internal plural name used when referring to the field in Go code.</dd>
  <dt><strong>min</strong></dt>
  <dd>The minimum value allowed for numeric fields.</dd>
  <dt><strong>max</strong></dt>
  <dd>The maximum value allowed for numeric fields.</dd>
</dl>