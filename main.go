package main

// This is the main execution point for the goradd installation and management tool. It is put here to
// make it easier to install for beginners, so that they can just call:
// go install github.com/goradd/goradd and it gets installed
//
// This is NOT the entry point for your web application. That lives in the goradd-project/main directory.
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

