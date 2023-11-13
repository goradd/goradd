# Database Configuration


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