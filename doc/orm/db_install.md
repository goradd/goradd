# Database Installation

GoRADD currently supports the following databases in the ORM:

- MySQL
- Postgres

This guide will give you some options on how to quickly install them, but its
not a comprehensive guide on how to use or configure them. The good news is
that although they are SQL databases, you don't need to know any SQL to use
them with GoRADD.

## Local Installation
To effectively work with GoRADD, you will need to install one of the supported
databases locally, so that you can run the code generator to generate the ORM code.

The supported databases have installers that will let you install the servers directly onto
your local computer. The files that contain the data are then stored in default locations.

However, you should also consider what your deployments will look like. If you are supporting
multiple kinds of deployments, perhaps running different versions of database software,
or on machines which may need to be upgraded at times and so database versions will change,
you should consider running your database servers in virtual machines. 

One common virtual machine strategy is to use Docker and download docker images of
the database and version of the database you would like to use. You then would 
have the ability to quickly switch versions of the database as the need arises, and
to better control where the database files are kept and backed up on your system.

Both databases can communicate via a standard port number, or via a unix socket file if
on a unix type platform.

## Database Tools
Once your database is running, you can issue SQL commands to setup the database. But SQL
commands can be hard to remember. Another way of managing a SQL database is to use a
database tool that lets you graphically set up your database. Some common tools that are free
include:
- [DBeaver](https://dbeaver.io) - A desktop application
- [PHPMyAdmin](https://hub.docker.com/_/phpmyadmin) - PHP scripts that can be installed via docker or locally
- [PHPPgAdmin](https://hub.docker.com/r/dockage/phppgadmin) - PHP scripts that can be installed via docker or locally

Once you have a database installed and running, the next step is to create the database
structure for your application.
