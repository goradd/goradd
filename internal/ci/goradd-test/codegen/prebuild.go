package main

// This file executes the complete codegen process by:
// 1) Removing old template files
// 2) Generating new template files from the template source
// 3) Building and then running the codegen app

// Generate the templates
//go:generate gofile remove goradd-project/tmp/template/*.tpl.go
//go:generate gofile remove goradd-project/gen/*/form
//go:generate gofile remove goradd-project/gen/*/model/*
//go:generate gofile remove goradd-project/gen/*/panel
//go:generate gofile remove goradd-project/gen/*/panelbase

//go:generate got -t got -o goradd-project/tmp/template -I goradd-project/codegen/templates/orm -d github.com/goradd/goradd/codegen/templates/orm
//go:generate got -t got -o goradd-project/tmp/template -I goradd-project/codegen/templates/page -d github.com/goradd/goradd/codegen/templates/page
