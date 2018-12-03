package main

import (
	"github.com/spekary/goradd/codegen/generator"

	_ "github.com/spekary/goradd/pkg/page/control/generator"
	_ "goradd-tmp/template"

	"goradd-project/config"
)

func main() {
	config.InitDatabases()
	generator.Generate()
}
