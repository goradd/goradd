// The main package runs a web server that works as an aid to installation, updating and building your application.
// Most of the code is in the buildtools directory
package main

import (
	"github.com/goradd/goradd/internal/goraddtool"
	"log"
)

// Create other flags you might care about here

func main() {
	rootCmd := goraddtool.MakeRootCommand()
	err := rootCmd.Execute()

	if err != nil {
		log.Fatal(err)
	}

}

