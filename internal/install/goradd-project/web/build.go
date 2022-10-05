package web

// Builds the templates that are in the web subdirectories.
// To execute it, run 'go generate build.go'.
// Every time you change a template file, you should run this file.
// Or, alternatively, set up your IDE to run it before you do a build.

//go:generate got -t tpl.got -i -I github.com/goradd/goradd/pkg/page/macros.inc.got -d goradd-project/web/*/*
