//go:build !release
// +build !release

package main

import (
	_ "github.com/goradd/goradd/doc"
	_ "github.com/goradd/goradd/pkg/bootstrap/examples/panels" // Bootstrap examples
	_ "github.com/goradd/goradd/test/browsertest"
	_ "github.com/goradd/goradd/web/examples"
	_ "github.com/goradd/goradd/web/welcome"
	_ "goradd-project/gen" // Code-generated forms
	_ "goradd-project/web" // Registers file assets through init calls.
)

// This file conditionally includes examples, generated forms, and other things
// useful during development, but that you definitely do not want in a release of your program.
// If you don't need these things, just delete this file.
