package goraddtool

import (
	"fmt"
	"log"
	"os"
)

var cwd string

func loadCwd() {
	var err error
	if cwd, err = os.Getwd(); err != nil {
		log.Fatal(fmt.Errorf("could not get the working directory: %s", err.Error()))
	}
}
