package buildtools

// This file executes the complete codegen process by:
// 1) Removing old template files
// 2) Generating new template files from the template source
// 3) Building and then running the codegen app

// TODO: Put the templates in a loadable library so that we are not building the whole application each time. Not sure that really matters though.

// Generate the templates
//go:generate gofile -r GOPATH/src/goradd-tmp/template/*.tpl.go
//go:generate go generate ../orm/codegen/build.go
//go:generate go generate ../page/codegen/build.go

// Run the code generator
//go:generate go run ../codegen/codegen.go

// Build the resulting templates
//go:generate gofile -g GOPATH/src/goradd-project/form/template_source/build.go
//go:generate gofile -g GOPATH/src/goradd-project/gen/*/panel/template_source/build.go
