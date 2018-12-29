package main

import (
	"github.com/goradd/goradd/codegen/generator"
	"goradd-project/config/dbconfig"

	_ "github.com/goradd/goradd/pkg/page/control/generator"
	_ "goradd-tmp/template"

)

func main() {
	dbconfig.InitDatabases()
	generator.Generate()
}
