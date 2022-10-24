// Package examples contains the example and tutorial code for learning how to use GoRADD.
//
// The gen/goradd/model here has been copied from the model package that gets generated for the
// goradd sample database. It has been placed here so that the examples code will compile even if the database
// is not loaded. Normally you would refer directly to the model package in your goradd-project/gen
// directory.
//
// The db directory has the source code for creating the database that goradd uses for testing
// and for the examples.
package examples

import (
	_ "github.com/goradd/goradd/web/examples/controls/panels"   // panels imports the controls package
	_ "github.com/goradd/goradd/web/examples/tutorial/controls" // imports the orm part of the tutorial
	_ "github.com/goradd/goradd/web/examples/tutorial/orm"      // imports the orm part of the tutorial
)
