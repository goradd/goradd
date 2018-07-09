package main

import (
	"flag"
	"fmt"
	"github.com/spekary/goradd/codegen/generator"

	// TODO: build the templates as plugin libraries so they do not need to be hard-linked
	_ "github.com/spekary/goradd/orm/template"
	_ "github.com/spekary/goradd/page/control/codegen/template"
	"goradd/config"
)

var test = flag.Bool("test", false, "test")

var Options = make(map[string]interface{})

// Create other flags you might care about here

type myHandler struct{}

func main() {
	var err error

	config.InitDatabases()

	if *test {
		// run a test
	} else {
		// Run in command line mode.
		generator.Generate()
	}

	if err != nil {
		fmt.Println(err)
	}
}
