package main

// This file executes the complete codegen process by:
// 1) Removing old template files
// 2) Generating new template files from the template source
// 3) Building and the running the codegen app

// TODO: Put the templates in a loadable library so that we are not building the whole application each time. Not sure that really matters though.

// TODO: Create a goradd version of file manipulation tools for use by the build system so we can be cross-platform
// go:generate rm -fv ../../../../goradd-tmp/template/*
// go:generate go generate ../orm/codegen/build.go
//go:generate got -t got -i -o goradd-tmp/template -I "goradd/codegen/orm;github.com/spekary/goradd/orm/codegen"

// go:generate go generate ../page/codegen/build.go


