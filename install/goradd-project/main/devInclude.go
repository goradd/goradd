// +build !release

package main

import (

	_ "goradd-project/gen"	// Code-generated forms
	_ "github.com/goradd/goradd/pkg/bootstrap/examples"	// Bootstrap examples

)

// This file conditionally builds examples, generated forms, and other things
// useful during development, but that you definitely do not want in a release of your program.
