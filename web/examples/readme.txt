
model

The model package here has been copied from the model package that gets generated for the
goradd sample database. It has been modified so that the examples code will compile even if the database
is not loaded. Normally you would refer directly to the model package in your goradd-project
directory.

db

This directory contains the sample database code that is used by some of the tutorials and examples.
To make these work, you will need to:

1) Install the database that is in the examples/db directory, and
2) Make the database available to the application by configuring it in goradd-project/config/db.go.

controls

The controls directory contains usage examples of the supported base controls provided by the goradd package.
The goal of the package is to be as comprehensive as possible in showing all the different options that
a control has, and the many different ways to use a control. It also includes brower-based tests that are
part of the Travis continuous-integration tests.

The base controls are foundations on which more elaborate controls can be built. But they also
provide enough support that a basic data-driven website can be built with only these controls.
