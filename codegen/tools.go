// +build tools

package codegen


/*
This file is here to resolve dependencies on binary tools. This allows you to run go mod tidy, and maintain pointers to
these tools in the go.mod file. This is an odd remnant of the modules functionality introduced in Go 1.11.
The tools flag above is intended to never be used.
 */

import (
	_ "github.com/goradd/gengen"
	_ "github.com/goradd/got"
)
