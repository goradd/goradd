# Anatomy of the Generated Framework

The framework in its current form generates two basic types
of objects:
- [ORM](#orm). 
- [Forms](#forms)
## ORM
The ORM, or Object Relational Model, are Go objects that map
to database objects, and that provide a layer of abstraction
between your code and however your database is accessed. Code
is generated for setters and getters, marshallers to JSON and
other formats, and code that commits changes to the database.

See [Databases](database.md) for details on how to structure your
database to achieve particular results in the generated objects.

