package page

// GoraddDomain is the translation file domain for all strings provided by goradd
const GoraddDomain = "goradd"
// ProjectDomain is the translation file domain for all strings that are part of the application
const ProjectDomain = "project"
// libraries should have their own domain

type Translater interface {
	Translate(in string) string
}

type PageTranslator struct {
	Domain string
	Language string
}

func (t *PageTranslator) SetLanguage(l string) {
	t.Language = l
}

func (t PageTranslator) Translate(in string) (out string) {
	// TODO: Call the application translation service
	// return I18N.Translate(t.domain, t.language, in)
	// first argument is the domain, so that goradd can provide its own translation files, and project can have custom files without overlap
	return in
}