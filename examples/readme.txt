
model

The model package here has been copied from the model package that gets generated for the
goradd sample database. It is here so that the examples code will compile even if the database
is not loaded. Normally you would refer directly to the model package in your goradd-project
directory.

db

This directory contains the sample database code in a variety of supported forms. You would install
one of these databases as the "goradd" database to run the examples. You will also need to add
this database to the program's database list. See the code in goradd-project/config/db.go.

controls

The controls directory contains usage examples of the supported base controls provided by the goradd package.
The goal of the package is to be as comprehensive as possible in showing all the different options that
a control has, and the many different ways to use a control. It also includes brower-based tests that are
part of the Travis continuous-integration tests.

The base controls are foundations on which more elaborate controls can be built. But they also
provide enough support that a basic data-driven website can be built with only these controls.
A website built with just the base controls will function even with javascript turned off
at the browser.

tutorial

This is the tutorial website.