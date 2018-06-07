package page

import (
	"io"
)

// This file manages a global template pool, which is a place for named templates. It essentially creates one global template
// and allows other templates to be assigned to the template pool, in a similar way to the built-in template engine
// by creating named templates, you can pre-compile templates and more easily serialize them.

// TemplateExecuter is an interface that will work with either text/template or html/template types of templates
type TemplateExecuter interface {
	Execute(wr io.Writer, data interface{}) error
}

var templatePool map[string]TemplateExecuter

func RegisterTemplate(name string, t TemplateExecuter) {
	templatePool[name] = t
}

func GetTemplate(name string) TemplateExecuter {
	t, _ := templatePool[name]
	return t
}
