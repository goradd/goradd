// +build unitTest

package main

import (
// my unit test package here
	_ "github.com/goradd/goradd/pkg/page/test"
)

// This file is here just to offer a conditional build that would include tests. Add a -tag unitTest flag to the
// build to get this to compile
