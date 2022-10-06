// GoRADD is a framework for rapid development of websites and API servers.
// See the Readme for details and instructions to get started.
package main

// This is the main execution point for the goradd installation and management tool. It is put here to
// make it easier to install for beginners, so that they can just call:
//   go install github.com/goradd/goradd@latest
// and the goradd tool gets installed
//
// This is NOT the entry point for your web application. That lives in the goradd-project/main directory.
// See the quickstart guide for info on getting started.
import (
	"github.com/goradd/goradd/internal/goraddtool"
	"log"
)

func main() {
	rootCmd := goraddtool.MakeRootCommand()
	err := rootCmd.Execute()

	if err != nil {
		log.Fatal(err)
	}

}
