# Databases
The database structure is central to Goradd. Goradd is designed to generate code from the database
structure, and then try to respond gracefully to changes in that structure. Your database will
strongly influence how data is presented and manipulated in a goradd application, and you will spend
much of your initial design process in working out the details of your database structure.

Goradd borrows from SQL terminology to think of databases in terms of tables and columns. You can
Each row in the spreadsheet is a collection of information, and each column in that row is a field 
of data of a particular type.

## Tables
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
