package main

import (
	"github.com/spekary/goradd/codegen/generator"
	"goradd-project/config/dbconfig"

	_ "github.com/spekary/goradd/pkg/page/control/generator"
	_ "goradd-tmp/template"

)

func main() {
	dbconfig.InitDatabases()
	generator.Generate()
}
