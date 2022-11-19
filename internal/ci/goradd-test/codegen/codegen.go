package main

import (
	"github.com/goradd/goradd/codegen/generator"
	_ "github.com/goradd/goradd/pkg/bootstrap/generator"
	//"github.com/goradd/goradd/pkg/orm/db/sql/pgsql"
	_ "github.com/goradd/goradd/pkg/page/control/generator"
	_ "goradd-project/config" // Initialize required variables
	_ "goradd-project/tmp/template"
)

func main() {
	config()
	generator.Generate()
}

func config() {
}
