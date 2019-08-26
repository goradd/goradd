package template_source

// This file will generate the templates and install them in the directory above. Feel free to edit.

//go:generate got -t tpl.got -i -I github.com/goradd/goradd/pkg/page/template_source/macros.inc.got:github.com/goradd/goradd/pkg/bootstrap/control/template_source/macros.inc.got -o ..
