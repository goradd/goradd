package main

// This file executes the complete codegen process by:
// 1) Removing old files
// 2) Generating new template files from the template source
// 3) Building and then running the codegen app

// Remove old files
//go:generate gofile remove goradd-tmp/template/*.tpl.go
//go:generate gofile remove goradd-project/gen/*/connector/*.base.go
//go:generate gofile remove goradd-project/gen/*/form/*.go
//go:generate gofile remove goradd-project/gen/*/form/template_source/*.tpl.got
//go:generate gofile remove goradd-project/gen/*/model/*.base.go
//go:generate gofile remove goradd-project/gen/*/model/node/*.go
//go:generate gofile remove goradd-project/gen/*/panel/*.base.go
//go:generate gofile remove goradd-project/gen/*/panel/template_source/inactive/*

// Generate the templates
//go:generate got -t got -o goradd-tmp/template -I goradd-project/codegen/templates/orm -d github.com/goradd/goradd/codegen/templates/orm
//go:generate got -t got -o goradd-tmp/template -I goradd-project/codegen/templates/page -d github.com/goradd/goradd/codegen/templates/page

// Run the code generator
//go:generate go run codegen.go

// Build the resulting templates
//go:generate gofile generate goradd-project/gen/*/form/template_source/build.go
//go:generate gofile generate goradd-project/gen/*/panel/template_source/build.go

// Build the templates that were moved to the form directory
//go:generate gofile generate goradd-project/web/form/template_source/build.go
