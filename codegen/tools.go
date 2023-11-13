//go:build tools
// +build tools

package codegen

/*
This file is here to resolve dependencies on binary tools. This allows us to run go mod tidy, and maintain pointers to
these tools in the go.mod file so that the correct versions of these tools will be associated with a specific
version of goradd.
*/

import (
	_ "github.com/goradd/got"
	_ "github.com/goradd/moddoc"
)
