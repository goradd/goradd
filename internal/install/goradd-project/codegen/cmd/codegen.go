package main

import (
	"github.com/goradd/goradd/codegen/generator"
	_ "github.com/goradd/goradd/pkg/bootstrap/generator"
	_ "github.com/goradd/goradd/pkg/page/control/generator"
	_ "goradd-project/config" // Initialize required variables
	_ "goradd-tmp/template"
)

func main() {
	config()
	generator.Generate()
}

func config() {
	// Customize the codegen process here

	// Replace the DefaultControlTypeFunc with one that calls the default, and then alters it
	// This lets you specify what kind of control you want to associate with a particular type of database field.

	// Uncomment the line below to generate bootstrap controls
	//generator2.BootstrapCodegenSetup()

}
