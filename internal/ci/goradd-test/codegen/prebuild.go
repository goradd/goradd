package main

// Just for preparing the build when testing during development

// Generate the templates
//go:generate gofile remove goradd-project/tmp/template/*.tpl.go
//go:generate gofile remove goradd-project/gen/*/form
//go:generate gofile remove goradd-project/gen/*/model/*
//go:generate gofile remove goradd-project/gen/*/panel
//go:generate gofile remove goradd-project/gen/*/panelbase
