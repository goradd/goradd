/*
Package db works with the rest of the orm to interface between a database and the ORM abstraction of
reading and querying a database. It is architected so that both SQL and NoSQL databases can be used with the
abstraction. This allows you to potentially write code that is completely agnostic to what kind of database
you are using. Even if 100% portability is not achievable in your use case, the ORM and database abstraction
layers should be able to handle most of your needs, such that if you ever did have to port to a different
database, it would minimize the amount of custom code you would need to write.

Generally, SQL and NoSQL databases work very differently. However, many SQL databases have recently added
NoSQL capabilities, like storing and searching JSON text. Similarly, NoSQL databases have added features to
enable searching a database through relationships, similar to SQL capabilities. In addition, NoSQL design advice
is often to flatten the database structure as much as possible, so that it looks a whole lot like a SQL database.

The general approach Goradd takes is to generally describe data with key/value pairs. This fits in well with SQL,
as key/value pairs are just table-column/field pairs. NoSQL generally works with key-value pairs anyways.

Relationships between structures are described as relationships, either one-to-many, or many-to-many. By keeping
the description at a higher level, we allow databases to implement those relationships in the way that works
best.

SQL implements one-to-many relationships using foreign keys. In the data description, you will see a
Reference type of relationship, which points from the many to the one, and a ReverseRelationship, which is a kind of
virtual representation of pointing from the one side to the many. ReverseRelationship lists are populated at
the time of a query. In SQL, Many-to-many relationships use an
intermediate table, called an Association Table, that has foreign keys pointing in both directions.

NoSQL implements one-to-many relationships using foreign-keys as well, but both sides of the relationship store
the foreign key, or keys, that point to the other side. This means that a ReverseRelationship will represent an
actual field in the database structure that contains a list of all of the items pointing back at itself. This also
means that when these relationships are created or destroyed, both sides of the relationship need to be updated.
Similarly, NoSQL many-to-many relationships have lists of foreign keys stored on both sides of the relationship.

The other major difference between SQL and NoSQL databases is the built-in capabilities to do aggregate calculations.
In SQL, you generally can create a filtered list of records and ask SQL to sum all the values from a particular field.
Some NoSQL databases can do this, and some cannot. The ones that cannot expect the programmer to do their own filtering
and summing. GoRADD handles this difference by allowing individual GORADD database drivers to be written that add
some aggregate capabilities to a database, and also providing ways for individual developers to simply create their
own custom queries that will be non-portable between databases. In any case, there is always a way to do what you
want to do, just some databases are easier to work with. It depends what you want to do.

*/
package db
