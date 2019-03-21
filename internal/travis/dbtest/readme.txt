These tests are designed to be run within a context that includes the database in examples/db/mysql.sql and
the code generated from this database. Combined with modules, this makes running these tests a little
complicated. The goradd/test directory has code which will copy out a project directory,
copy this directory, and then run the code generator for that project such that it will
create the code needed for these tests.