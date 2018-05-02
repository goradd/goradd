package main

import (
	"flag"
	"fmt"
	"goradd/config"
	"github.com/spekary/goradd/codegen/generator"
	_ "github.com/spekary/goradd/orm/template"
)

var test = flag.Bool("test", false, "test")

var Options = make(map[string]interface{})


// Create other flags you might care about here

type myHandler struct {}


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

