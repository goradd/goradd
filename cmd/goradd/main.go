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

