package main

// This file executes the complete codegen process by:
// 1) Removing old template files
// 2) Generating new template files from the template source
// 3) Building and then running the codegen app

// Generate the templates
//go:generate got -t got -o goradd-project/tmp/template -I goradd-project/codegen/templates/orm -d github.com/goradd/goradd/codegen/templates/orm -i
//go:generate got -t got -o goradd-project/tmp/template -I goradd-project/codegen/templates/page -d github.com/goradd/goradd/codegen/templates/page -i

// Run the code generator
//go:generate go run .

// Build the resulting templates
//go:generate got -t tpl.got -i -I github.com/goradd/goradd/pkg/page/macros.inc.got -d goradd-project/gen/*/*
