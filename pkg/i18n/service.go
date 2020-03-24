package i18n

import (
	"fmt"
)

// Predefined domains. Plugins can add their own domains.
const GoraddDomain = "goradd"
const ProjectDomain = "project"

var translators = map[string]TranslatorI{GoraddDomain: NonTranslator{}, ProjectDomain: NonTranslator{}}

// TranslatorI is the interface that translators must fulfill
type TranslatorI interface {
	// Translate returns the translation of the string contained in the translationBuilder
	Translate(b *translationBuilder) string
}

// NonTranslator is the default translator that just passes all strings through unchanged.
type NonTranslator struct {
}

func (n NonTranslator) Translate(b *translationBuilder) string {
	if b.arguments == nil {
		// Just want a passthrough. If b.message has Sprintf format commands, calling fmt.Sprintf with no arguments will err.
		return b.message
	}
	return fmt.Sprintf(b.message, b.arguments...)
}

// RegisterTranslator sets the translation service for the given domain to the given translator
func RegisterTranslator(domain string, t TranslatorI) {
	translators[domain] = t
}
