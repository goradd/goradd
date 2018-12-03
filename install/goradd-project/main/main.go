package main

import (
	"flag"
	"fmt"
	"goradd-project/web/app"

	// Below is where you import packages that register forms
	_ "goradd-project/web/form" // Your  forms.

	// Custom paths, including additional form directories
	_ "site"
)

var local = flag.String("local", "", "serve as webserver from given port, example: -local 8000")
var useFcgi = flag.Bool("fcgi", false, "serve as fcgi, example: -fcgi")
var assetDir = flag.String("assetDir", "", "The centralized asset directory. Required to run the release version of the app.")
// Create other flags you might care about here

func main() {
	var err error

	a := app.MakeApplication(*assetDir)
	err = a.RunWebServer(*local, *useFcgi)

	if err != nil {
		fmt.Println(err)
	}
}

