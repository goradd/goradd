package page

import "github.com/spekary/goradd/pkg/i18n"

// GoraddDomain is the translation file domain for all strings provided by goradd
const GoraddDomain = "goradd"

// ProjectDomain is the translation file domain for all strings that are part of the application
const ProjectDomain = "project"

// libraries should have their own domain

type Translater interface {
	Translate(in string) string
}

type PageTranslator struct {
	Domain   string
	Language string
}

func (t *PageTranslator) SetLanguage(l string) {
	t.Language = l
}

func (t PageTranslator) Translate(in string) (out string) {
	return i18n.Translate(t.Domain, t.Language, in)
}
