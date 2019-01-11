package i18n

import (
	"fmt"
)

// Predefined domains. Plugins can add their own domains.
const GoraddDomain = "goradd"
const ProjectDomain = "project"


var translators = map[string]TranslatorI{GoraddDomain: NonTranslator{}, ProjectDomain: NonTranslator{}}

type TranslatorI interface {
	Translate (b *translationBuilder) string
}

type NonTranslator struct {
}

func (n NonTranslator) Translate (b *translationBuilder) string {
	if b.arguments == nil {
		// Just want a passthrough. If b.message has Sprintf format commands, calling fmt.Sprintf with no arguments will err.
		return b.message
	}
	return fmt.Sprintf(b.message, b.arguments...)
}

// RegisterTranslator sets the translation service for the given domain to the given translator
func RegisterTranslator (domain string, t TranslatorI) {
	translators[domain] = t
}
