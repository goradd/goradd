package i18n

var translator TranslatorI = new(NonTranslator)

type TranslatorI interface {
	Translate (domain string, language string, s string) string
}

type NonTranslator struct {
}

func (n NonTranslator) Translate (domain string, language string, s string) string {
	return s
}

// SetTranslator sets the translation service to the given translator
func SetTranslator (t TranslatorI) {
	translator = t
}

func Translate (domain string, language string, s string) string {
	return translator.Translate(domain, language, s)
}