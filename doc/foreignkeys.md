# Database Relationships
## Foreign Keys

Foreign Keys are a term that SQL databases use to refer to a pointer from one database table
to another. Foreign Keys are what make relational databases "relational". NoSQL databases
can use them too, but NoSQL databases often manage their relationships by duplicating data
instead of using pointers.

Simply put, a foreign key is a field in a table that contains the primary key of a record
in another table. SQL queries can use these links to look for related information in the
other table, and SQL is specially designed to query these relationships.

GoRADD uses SQL foreign keys to learn about the relationships in your data, and then creates
Go ORM classes to manage these relationships. For a MySQL database, you must use the InnoDB
engine to make this work (which is the default engine for MySQL).

(***Future: GoRADD can manage relationships in NoSQL databases
but since NoSQL databases do not have any kind of internal description of these relationships,
GoRADD needs to learn about your intentions from a separate description file. To ease the
transition, GoRADD can create the description file from a SQL database, that can then be
used with the same data in a NoSQL database.)

## Relationship Types
### One-to-Many Relationships
Typically a foreign key will create a one-to-many relationship. The table with the
foreign key points to one record in another table, but that means that the record that
is being pointed to might have more than one record pointing to it, so it might have
many records linked to it.

GoRADD will create a "Reference" from the object with the foreign key to the object being
referenced so that you can easily get from one object to the other. On the other side of the relationship,
it will create a "Reverse Reference", which means that the object has other objects pointing
to it, and you can use this Reverse Reference to get the collection of objects pointing to
that object.

### One-to-One Relationships
You can link two objects together so that it creates a one-to-one reference. You would
typically do this if you had situations where you wanted to load a subset of an object for
general display, but needed to be able to show all the data if the user drilled down to
the detail. 

In SQL, you define this relationship with a foreign key that has a Unique index on it. 
GoRADD will see the unique index and create single links on both sides of the relationship.

### Many-to-Many Relationships
GoRADD uses an intermediate table that sits between two tables that have a many-to-many
relationship, called an "Association" table. This table has only two fields, each with
a foreign key to the other tables being related. GoRADD will see this table and link
the other tables using functions that allow you to get a collection of objects being pointed to
from either side of the relationship, without needing to specifically query the association table.

### Parent-Child Relationships
If a foreign key points to its own table, it creates a parent-child relationship, which can
be used to create a tree structure of objects. The record with the foreign key will point
to its parent object, and the parent object will use the Reverse Reference created by that
to see all of its child objects.

## Specifying Default Values
If the foreign key is nullable, then you are telling GoRADD that the relationship
between the two tables is optional. The default value for this would typically be NULL,
and GoRADD will create a selection list that will start with an empty value. You 
can specify a default in the database if you wish, and GoRADD will read that default
value and use that default value in new records.

If the foreign key is not nullable, you are telling GoRADD that there MUST be a
relationship between two records when the record with the foreign key is created. 
If you do not specify a default
value in the database, then GoRADD will insert a "- Select One -" option at the top
of the corresponding selection list to indicate that the value has not been selected,
but is required.

